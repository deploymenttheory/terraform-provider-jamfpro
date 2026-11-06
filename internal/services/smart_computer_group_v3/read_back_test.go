package smart_computer_group_v3

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateWaitsForNewCriterion(t *testing.T) {
	prior := testModel(t)
	old, diags := constructResource(context.Background(), &prior)
	require.False(t, diags.HasError())
	sdk := &fakeGroupSDK{stale: []*jamfpro.ResourceSmartComputerGroupV3{old}}
	r, definition := testResource(t, sdk)
	// The write uses its operation deadline, not the delete-confirmation budget.
	r.deleteConfirmationTimeout = time.Nanosecond
	plan := testModel(t)
	var criteria []criterionModel
	require.False(t, plan.Criteria.ElementsAs(context.Background(), &criteria, false).HasError())
	criteria[0].Value = types.StringValue("new-criterion-value")
	plan.Criteria, diags = schemahelpers.Flatten(context.Background(), criteria)
	require.False(t, diags.HasError())
	before := testPlan(t, definition, prior)
	initial := tfsdk.State{Schema: definition, Raw: before.Raw}
	response := resource.UpdateResponse{State: initial}
	r.Update(context.Background(), resource.UpdateRequest{Plan: testPlan(t, definition, plan), State: initial}, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	var got smartComputerGroupV3ResourceModel
	require.False(t, response.State.Get(context.Background(), &got).HasError())
	require.True(t, plan.Criteria.Equal(got.Criteria))
	require.Equal(t, "42", got.ID.ValueString())
	require.Equal(t, []string{"v3 update", "read", "read"}, sdk.events)
}

func TestUpdateStaleReadHonorsDeadlineAndKeepsObservedState(t *testing.T) {
	prior := testModel(t)
	old, diags := constructResource(context.Background(), &prior)
	require.False(t, diags.HasError())
	sdk := &fakeGroupSDK{stale: make([]*jamfpro.ResourceSmartComputerGroupV3, 1000)}
	for i := range sdk.stale {
		sdk.stale[i] = old
	}
	r, definition := testResource(t, sdk)
	plan := testModel(t)
	plan.Name = types.StringValue("new-name")
	before := testPlan(t, definition, prior)
	initial := tfsdk.State{Schema: definition, Raw: before.Raw}
	response := resource.UpdateResponse{State: initial}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	r.Update(ctx, resource.UpdateRequest{Plan: testPlan(t, definition, plan), State: initial}, &response)
	require.True(t, response.Diagnostics.HasError())
	var got smartComputerGroupV3ResourceModel
	require.False(t, response.State.Get(context.Background(), &got).HasError())
	require.Equal(t, "42", got.ID.ValueString())
	require.True(t, prior.Name.Equal(got.Name))
	writes := 0
	for _, event := range sdk.events {
		if event == "v3 update" {
			writes++
		}
	}
	require.Equal(t, 1, writes)
}
