package smart_computer_group_v2

import (
	"context"
	apiresponse "github.com/deploymenttheory/go-api-http-client/response"

	"strconv"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

type fakeGroupSDK struct {
	group     *jamfpro.ResourceSmartComputerGroupV2
	stale     []*jamfpro.ResourceSmartComputerGroupV2
	readErr   error
	deleteErr error
	reads     int
	writes    int
	deadline  bool
}

func (f *fakeGroupSDK) GetSmartComputerGroupByIDV2WithContext(ctx context.Context, _ string) (*jamfpro.ResourceSmartComputerGroupV2, error) {
	f.reads++
	_, f.deadline = ctx.Deadline()
	if len(f.stale) > 0 {
		group := f.stale[0]
		f.stale = f.stale[1:]
		return group, nil
	}
	return f.group, f.readErr
}

func (f *fakeGroupSDK) CreateSmartComputerGroupV2WithContext(ctx context.Context, input jamfpro.ResourceSmartComputerGroupV2) (*jamfpro.ResponseSmartComputerGroupCreateV2, error) {
	f.writes++
	_, f.deadline = ctx.Deadline()
	f.group = &input
	return &jamfpro.ResponseSmartComputerGroupCreateV2{ID: "42"}, nil
}

func (f *fakeGroupSDK) UpdateSmartComputerGroupByIDV2WithContext(ctx context.Context, _ string, input jamfpro.ResourceSmartComputerGroupV2) (*jamfpro.ResourceSmartComputerGroupV2, error) {
	f.writes++
	_, f.deadline = ctx.Deadline()
	f.group = &input
	return &input, nil
}

func (f *fakeGroupSDK) DeleteSmartComputerGroupByIDV2WithContext(ctx context.Context, _ string) error {
	_, f.deadline = ctx.Deadline()
	return f.deleteErr
}

func testResource(t *testing.T, sdk *fakeGroupSDK) (*smartComputerGroupV2FrameworkResource, schema.Schema) {
	t.Helper()
	r := NewSmartComputerGroupV2FrameworkResource().(*smartComputerGroupV2FrameworkResource)
	r.client, r.pollInterval = sdk, time.Nanosecond
	var response resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &response)
	return r, response.Schema
}

func testModel(t *testing.T, value string) smartComputerGroupV2ResourceModel {
	t.Helper()
	criteria, diags := schemahelpers.Flatten(context.Background(), []smartComputerGroupV2CriteriaDataModel{{
		Name: types.StringValue("Packages Installed by Jamf Pro"), Priority: types.Int64Value(0),
		AndOr: types.StringValue("and"), SearchType: types.StringValue("does not have"), Value: types.StringValue(value),
		OpeningParen: types.BoolValue(false), ClosingParen: types.BoolValue(false),
	}})
	require.False(t, diags.HasError(), "%v", diags)
	return smartComputerGroupV2ResourceModel{
		ID: types.StringValue("42"), Name: types.StringValue("tf-testing-smart-v2"),
		Description: types.StringNull(), SiteID: types.StringValue("-1"), Criteria: criteria,
		Timeouts: timeouts.Value{Object: types.ObjectNull(timeoutTypes())},
	}
}

func timeoutTypes() map[string]attr.Type {
	return map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType}
}

func testPlan(t *testing.T, definition schema.Schema, data smartComputerGroupV2ResourceModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: definition}
	require.False(t, plan.Set(context.Background(), &data).HasError())
	return plan
}

func TestUpdatePreservesPlannedTimeouts(t *testing.T) {
	for _, removing := range []bool{false, true} {
		t.Run(strconv.FormatBool(removing), func(t *testing.T) {
			prior := testModel(t, "old.pkg")
			plan := testModel(t, "new.pkg")
			configured := timeouts.Value{Object: types.ObjectValueMust(timeoutTypes(), map[string]attr.Value{
				"create": types.StringValue("300s"), "read": types.StringValue("300s"),
				"update": types.StringValue("300s"), "delete": types.StringNull(),
			})}
			if removing {
				prior.Timeouts = configured
			} else {
				plan.Timeouts = configured
			}
			sdk := &fakeGroupSDK{}
			r, definition := testResource(t, sdk)
			before := testPlan(t, definition, prior)
			initial := tfsdk.State{Schema: definition, Raw: before.Raw}
			response := resource.UpdateResponse{State: initial}
			r.Update(context.Background(), resource.UpdateRequest{Plan: testPlan(t, definition, plan), State: initial}, &response)
			require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
			var got smartComputerGroupV2ResourceModel
			require.False(t, response.State.Get(context.Background(), &got).HasError())
			require.True(t, plan.Timeouts.Equal(got.Timeouts), "Update returned prior timeouts instead of planned values")
			require.True(t, sdk.deadline)
		})
	}
}

func TestReadMissingGroupAndPermissionErrors(t *testing.T) {
	for _, code := range []int{404, 410, 403} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			sdk := &fakeGroupSDK{readErr: &apiresponse.APIError{StatusCode: code}}
			r, definition := testResource(t, sdk)
			before := testPlan(t, definition, testModel(t, "example.pkg"))
			initial := tfsdk.State{Schema: definition, Raw: before.Raw}
			response := resource.ReadResponse{State: initial}
			r.Read(context.Background(), resource.ReadRequest{State: initial}, &response)
			require.Equal(t, code != 403, response.State.Raw.IsNull())
			require.Equal(t, code == 403, response.Diagnostics.HasError())
		})
	}
}
func TestDeleteAlreadyMissingGroup(t *testing.T) {
	sdk := &fakeGroupSDK{deleteErr: &apiresponse.APIError{StatusCode: 404}}
	r, definition := testResource(t, sdk)
	before := testPlan(t, definition, testModel(t, "example.pkg"))
	initial := tfsdk.State{Schema: definition, Raw: before.Raw}
	response := resource.DeleteResponse{State: initial}
	r.Delete(context.Background(), resource.DeleteRequest{State: initial}, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.True(t, response.State.Raw.IsNull())
}

func TestUpdateWaitsForNewPackageCriterion(t *testing.T) {
	prior := testModel(t, "example_pkg-1.0.pkg")
	old, diags := constructResource(context.Background(), &prior)
	require.False(t, diags.HasError())
	sdk := &fakeGroupSDK{stale: []*jamfpro.ResourceSmartComputerGroupV2{old}}
	r, definition := testResource(t, sdk)
	before := testPlan(t, definition, prior)
	initial := tfsdk.State{Schema: definition, Raw: before.Raw}
	response := resource.UpdateResponse{State: initial}
	plan := testModel(t, "example_pkg-2.0.pkg")
	r.Update(context.Background(), resource.UpdateRequest{Plan: testPlan(t, definition, plan), State: initial}, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	var got smartComputerGroupV2ResourceModel
	require.False(t, response.State.Get(context.Background(), &got).HasError())
	require.True(t, plan.Criteria.Equal(got.Criteria), "Update returned stale package criteria")
	require.Equal(t, 2, sdk.reads)
	require.Equal(t, 1, sdk.writes, "read-back convergence must not replay the write")
}
func TestUpdateStaleReadHonorsDeadlineAndKeepsObservedState(t *testing.T) {
	prior := testModel(t, "old.pkg")
	old, diags := constructResource(context.Background(), &prior)
	require.False(t, diags.HasError())
	sdk := &fakeGroupSDK{stale: make([]*jamfpro.ResourceSmartComputerGroupV2, 1000)}
	for i := range sdk.stale {
		sdk.stale[i] = old
	}
	r, definition := testResource(t, sdk)
	before := testPlan(t, definition, prior)
	initial := tfsdk.State{Schema: definition, Raw: before.Raw}
	response := resource.UpdateResponse{State: initial}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	r.Update(ctx, resource.UpdateRequest{Plan: testPlan(t, definition, testModel(t, "new.pkg")), State: initial}, &response)
	require.True(t, response.Diagnostics.HasError(), "stale state must not be reported as a successful update")
	var got smartComputerGroupV2ResourceModel
	require.False(t, response.State.Get(context.Background(), &got).HasError())
	require.True(t, prior.Criteria.Equal(got.Criteria))
	require.Equal(t, "42", got.ID.ValueString())
	require.Equal(t, 1, sdk.writes)
}
