package smart_computer_group_v2

import (
	"context"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/stretchr/testify/require"
)

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
