package smart_computer_group_v2

import (
	"context"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStateEmptyCriteriaKeepTypedNull(t *testing.T) {
	ctx := context.Background()
	site := "-1"
	model := smartComputerGroupV2ResourceModel{
		ID: types.StringValue("42"), Name: types.StringValue("empty"), Description: types.StringNull(), SiteID: types.StringValue(site),
		Criteria: types.ListNull(schemahelpers.ObjectType),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType})},
	}
	diags := state(ctx, &model, "42", &jamfpro.ResourceSmartComputerGroupV2{Name: "empty", SiteId: &site})
	require.False(t, diags.HasError(), "%v", diags)
	require.True(t, model.Criteria.Equal(types.ListNull(schemahelpers.ObjectType)))
	var definition resource.SchemaResponse
	NewSmartComputerGroupV2FrameworkResource().Schema(ctx, resource.SchemaRequest{}, &definition)
	result := tfsdk.State{Schema: definition.Schema}
	require.False(t, result.Set(ctx, &model).HasError(), "empty criteria must retain the schema's object type")
}
