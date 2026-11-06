package smart_computer_group_v3

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestStatePreservesIDTimeoutsAndCriterionOrder(t *testing.T) {
	model := testModel(t)
	timeouts := model.Timeouts
	group := &jamfpro.ResourceSmartComputerGroupV3{
		Name: "read", Criteria: []jamfpro.SubsetComputerSmartGroupCriteriaV3{
			{Name: "second", Priority: 1, AndOr: "OR"}, {Name: "first", Priority: 0, AndOr: "AND"},
		},
	}
	diags := state(context.Background(), &model, group)
	require.False(t, diags.HasError())
	require.Equal(t, "42", model.ID.ValueString())
	require.True(t, timeouts.Equal(model.Timeouts))
	require.Equal(t, "-1", model.SiteID.ValueString())
	criteria, diags := schemahelpers.Expand[criterionModel](context.Background(), model.Criteria)
	require.False(t, diags.HasError())
	require.Equal(t, "first", criteria[0].Name.ValueString())
	require.Equal(t, "and", criteria[0].AndOr.ValueString())
	require.Equal(t, "or", criteria[1].AndOr.ValueString())
}

func TestStateEmptyCriteriaKeepsTypedNull(t *testing.T) {
	model := smartComputerGroupV3ResourceModel{Criteria: types.ListNull(schemahelpers.ObjectType)}
	diags := state(context.Background(), &model, &jamfpro.ResourceSmartComputerGroupV3{Name: "empty"})
	require.False(t, diags.HasError())
	require.True(t, model.Criteria.IsNull())
	require.Equal(t, schemahelpers.ObjectType, model.Criteria.ElementType(context.Background()))
}
