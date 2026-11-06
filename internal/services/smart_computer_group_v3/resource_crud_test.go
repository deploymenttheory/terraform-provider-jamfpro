package smart_computer_group_v3

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	apiresponse "github.com/deploymenttheory/go-api-http-client/response"
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
	group       *jamfpro.ResourceSmartComputerGroupV3
	stale       []*jamfpro.ResourceSmartComputerGroupV3
	modernErr   error
	patchErr    error
	readErr     error
	deleteErr   error
	events      []string
	classicBody *jamfpro.ResourceComputerGroup
}

var errTestPatchPermission = errors.New("description permission denied")

func (f *fakeGroupSDK) GetSmartComputerGroupByIDV3WithContext(
	_ context.Context, _ string,
) (*jamfpro.ResourceSmartComputerGroupV3, error) {
	f.events = append(f.events, "read")
	if len(f.stale) > 0 {
		got := f.stale[0]
		f.stale = f.stale[1:]
		return got, nil
	}
	return f.group, f.readErr
}

func (f *fakeGroupSDK) CreateSmartComputerGroupV3WithContext(
	_ context.Context, input jamfpro.ResourceSmartComputerGroupV3,
) (*jamfpro.ResponseSmartComputerGroupCreateV3, error) {
	f.events = append(f.events, "v3 create")
	if f.modernErr != nil {
		return nil, f.modernErr
	}
	f.group = &input
	return &jamfpro.ResponseSmartComputerGroupCreateV3{ID: "42"}, nil
}

func (f *fakeGroupSDK) UpdateSmartComputerGroupByIDV3WithContext(
	_ context.Context, _ string, input jamfpro.ResourceSmartComputerGroupV3,
) (*jamfpro.ResourceSmartComputerGroupV3, error) {
	f.events = append(f.events, "v3 update")
	if f.modernErr != nil {
		return nil, f.modernErr
	}
	f.group = &input
	return &input, nil
}

func (f *fakeGroupSDK) DeleteSmartComputerGroupByIDV3WithContext(_ context.Context, _ string) error {
	f.events = append(f.events, "v3 delete")
	if f.deleteErr == nil {
		f.readErr = &apiresponse.APIError{StatusCode: 404}
	}
	return f.deleteErr
}

func (f *fakeGroupSDK) CreateComputerGroupWithContext(
	_ context.Context, input *jamfpro.ResourceComputerGroup,
) (*jamfpro.ResponseComputerGroupreatedAndUpdated, error) {
	f.events = append(f.events, "classic create")
	f.saveClassic(input)
	return &jamfpro.ResponseComputerGroupreatedAndUpdated{ID: 42}, nil
}

func (f *fakeGroupSDK) UpdateComputerGroupByIDWithContext(
	_ context.Context, _ string, input *jamfpro.ResourceComputerGroup,
) (*jamfpro.ResponseComputerGroupreatedAndUpdated, error) {
	f.events = append(f.events, "classic update")
	f.saveClassic(input)
	return &jamfpro.ResponseComputerGroupreatedAndUpdated{ID: 42}, nil
}

func (f *fakeGroupSDK) saveClassic(input *jamfpro.ResourceComputerGroup) {
	f.classicBody = input
	description := ""
	if f.group != nil {
		description = f.group.Description
	}
	site := strconv.Itoa(input.Site.ID)
	f.group = &jamfpro.ResourceSmartComputerGroupV3{Name: input.Name, Description: description, SiteID: &site}
	for _, criterion := range *input.Criteria.Criterion {
		f.group.Criteria = append(f.group.Criteria, jamfpro.SubsetComputerSmartGroupCriteriaV3(criterion))
	}
}

func (f *fakeGroupSDK) GetComputerSmartGroupPlatformIDByJamfProIDWithContext(_ context.Context, _ string) (string, error) {
	f.events = append(f.events, "platform ID")
	return "platform-uuid", nil
}

func (f *fakeGroupSDK) UpdateGroupDescriptionByIDV2WithContext(_ context.Context, id, description string) error {
	f.events = append(f.events, "description patch")
	if id != "platform-uuid" {
		panic("description must use the Platform ID")
	}
	if f.patchErr != nil {
		return f.patchErr
	}
	f.group.Description = description
	return nil
}

func testResource(t *testing.T, sdk *fakeGroupSDK) (*smartComputerGroupV3FrameworkResource, schema.Schema) {
	t.Helper()
	r := NewSmartComputerGroupV3FrameworkResource().(*smartComputerGroupV3FrameworkResource)
	r.client, r.pollInterval, r.deleteConfirmationTimeout = sdk, time.Nanosecond, 100*time.Millisecond
	var response resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &response)
	return r, response.Schema
}

func testModel(t *testing.T) smartComputerGroupV3ResourceModel {
	t.Helper()
	criteria, diags := schemahelpers.Flatten(context.Background(), []criterionModel{{
		Name: types.StringValue("Serial Number"), Priority: types.Int64Value(0),
		AndOr: types.StringValue("and"), SearchType: types.StringValue("is"),
		Value: types.StringValue("tf-testing-no-match"), OpeningParen: types.BoolValue(false), ClosingParen: types.BoolValue(false),
	}})
	require.False(t, diags.HasError(), "%v", diags)
	timeoutTypes := map[string]attr.Type{
		"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType,
	}
	timeoutValues := map[string]attr.Value{
		"create": types.StringValue("5m"), "read": types.StringNull(),
		"update": types.StringValue("4m"), "delete": types.StringNull(),
	}
	return smartComputerGroupV3ResourceModel{
		ID: types.StringValue("42"), Name: types.StringValue("tf-testing-smart-v3"),
		Description: types.StringValue("description"), SiteID: types.StringValue("-1"), Criteria: criteria,
		Timeouts: timeouts.Value{Object: types.ObjectValueMust(timeoutTypes, timeoutValues)},
	}
}

func testPlan(t *testing.T, definition schema.Schema, data smartComputerGroupV3ResourceModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: definition}
	diags := plan.Set(context.Background(), &data)
	require.False(t, diags.HasError(), "%v", diags)
	return plan
}

func refusal() error {
	return &apiresponse.APIError{StatusCode: 400, Errors: []apiresponse.Errors{{
		Code: "INVALID_FIELD", Description: "The operator has is not valid for extension attribute Notes",
	}}}
}

func TestCreateRecordsGroupAndTimeouts(t *testing.T) {
	sdk := &fakeGroupSDK{}
	r, definition := testResource(t, sdk)
	plan := testModel(t)
	response := resource.CreateResponse{State: tfsdk.State{Schema: definition}}
	r.Create(context.Background(), resource.CreateRequest{Plan: testPlan(t, definition, plan)}, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	var got smartComputerGroupV3ResourceModel
	require.False(t, response.State.Get(context.Background(), &got).HasError())
	require.Equal(t, "42", got.ID.ValueString())
	require.True(t, plan.Timeouts.Equal(got.Timeouts))
	require.Equal(t, []string{"v3 create", "read"}, sdk.events)
}

func TestUpdateFallbackWaitsBeforeDescriptionPatch(t *testing.T) {
	prior := testModel(t)
	old, diags := constructResource(context.Background(), &prior)
	require.False(t, diags.HasError())
	sdk := &fakeGroupSDK{group: old, stale: []*jamfpro.ResourceSmartComputerGroupV3{old}, modernErr: refusal()}
	r, definition := testResource(t, sdk)
	plan := prior
	plan.Name, plan.Description = types.StringValue("renamed"), types.StringValue("")
	before := testPlan(t, definition, prior)
	response := resource.UpdateResponse{State: tfsdk.State{Schema: definition, Raw: before.Raw}}
	r.Update(context.Background(), resource.UpdateRequest{
		Plan: testPlan(t, definition, plan), State: tfsdk.State{Schema: definition, Raw: before.Raw},
	}, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.Equal(t, []string{"v3 update", "classic update", "read", "read", "platform ID", "description patch", "read"}, sdk.events)
	require.Equal(t, "renamed", sdk.classicBody.Name)
	require.Equal(t, -1, sdk.classicBody.Site.ID)
	require.Len(t, *sdk.classicBody.Criteria.Criterion, 1)
	require.Nil(t, sdk.classicBody.Computers)
	var got smartComputerGroupV3ResourceModel
	require.False(t, response.State.Get(context.Background(), &got).HasError())
	require.Equal(t, "", got.Description.ValueString())
	require.True(t, plan.Timeouts.Equal(got.Timeouts))
}

func TestCreateFallbackPatchFailureRetainsID(t *testing.T) {
	sdk := &fakeGroupSDK{modernErr: refusal(), patchErr: errTestPatchPermission}
	r, definition := testResource(t, sdk)
	response := resource.CreateResponse{State: tfsdk.State{Schema: definition}}
	r.Create(context.Background(), resource.CreateRequest{Plan: testPlan(t, definition, testModel(t))}, &response)
	require.True(t, response.Diagnostics.HasError())
	var got smartComputerGroupV3ResourceModel
	require.False(t, response.State.Get(context.Background(), &got).HasError())
	require.Equal(t, "42", got.ID.ValueString())
	require.Equal(t, "", got.Description.ValueString())
}

func TestReadMissingGroupAndPermissionErrors(t *testing.T) {
	for _, code := range []int{404, 403} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			sdk := &fakeGroupSDK{readErr: &apiresponse.APIError{StatusCode: code}}
			r, definition := testResource(t, sdk)
			before := testPlan(t, definition, testModel(t))
			initial := tfsdk.State{Schema: definition, Raw: before.Raw}
			response := resource.ReadResponse{State: initial}
			r.Read(context.Background(), resource.ReadRequest{State: initial}, &response)
			require.Equal(t, code == 404, response.State.Raw.IsNull())
			require.Equal(t, code == 403, response.Diagnostics.HasError())
		})
	}
}

func TestDeleteAlreadyMissingGroup(t *testing.T) {
	sdk := &fakeGroupSDK{deleteErr: &apiresponse.APIError{StatusCode: 404}}
	r, definition := testResource(t, sdk)
	before := testPlan(t, definition, testModel(t))
	initial := tfsdk.State{Schema: definition, Raw: before.Raw}
	response := resource.DeleteResponse{State: initial}
	r.Delete(context.Background(), resource.DeleteRequest{State: initial}, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.True(t, response.State.Raw.IsNull())
}
