package smart_computer_group_v3

import (
	"context"
	"testing"

	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestValidateConfigCriteria(t *testing.T) {
	for _, tc := range []struct {
		name      string
		change    func(*criterionModel)
		wantError bool
	}{
		{"valid criterion", func(*criterionModel) {}, false},
		{"noncontiguous priority", func(c *criterionModel) { c.Priority = types.Int64Value(3) }, true},
		{"unknown priority deferred", func(c *criterionModel) { c.Priority = types.Int64Unknown() }, false},
		{"unmatched open", func(c *criterionModel) { c.OpeningParen = types.BoolValue(true) }, true},
		{"unmatched close", func(c *criterionModel) { c.ClosingParen = types.BoolValue(true) }, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, definition := testResource(t, &fakeGroupSDK{})
			model := testModel(t)
			criteria, diags := schemahelpers.Expand[criterionModel](context.Background(), model.Criteria)
			require.False(t, diags.HasError())
			tc.change(&criteria[0])
			model.Criteria, diags = schemahelpers.Flatten(context.Background(), criteria)
			require.False(t, diags.HasError())
			plan := testPlan(t, definition, model)
			response := resource.ValidateConfigResponse{}
			r.ValidateConfig(context.Background(), resource.ValidateConfigRequest{
				Config: tfsdk.Config{Schema: definition, Raw: plan.Raw},
			}, &response)
			require.Equal(t, tc.wantError, response.Diagnostics.HasError(), "%v", response.Diagnostics)
		})
	}
}
