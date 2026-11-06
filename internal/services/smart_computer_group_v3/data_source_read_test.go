package smart_computer_group_v3

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func (f *fakeGroupSDK) GetSmartComputerGroupByNameV3WithContext(
	_ context.Context, _ string,
) (*jamfpro.ResponseSmartComputerGroupListItemV3, error) {
	f.events = append(f.events, "lookup name")
	return &jamfpro.ResponseSmartComputerGroupListItemV3{ID: "42", Name: f.group.Name}, nil
}

func TestDataSourceReadsEmptyCriteriaByIDAndName(t *testing.T) {
	for _, byName := range []bool{false, true} {
		name := "by ID"
		if byName {
			name = "by name"
		}
		t.Run(name, func(t *testing.T) {
			sdk := &fakeGroupSDK{group: &jamfpro.ResourceSmartComputerGroupV3{Name: "empty"}}
			d := &smartComputerGroupV3DataSource{client: sdk}
			var definition datasource.SchemaResponse
			d.Schema(context.Background(), datasource.SchemaRequest{}, &definition)
			model := smartComputerGroupV3DataSourceModel{
				ID: types.StringValue("42"), Name: types.StringNull(),
				Description: types.StringNull(), SiteID: types.StringNull(),
				Criteria: types.ListNull(schemahelpers.ObjectType),
			}
			if byName {
				model.ID, model.Name = types.StringNull(), types.StringValue("empty")
			}
			config := tfsdk.State{Schema: definition.Schema}
			require.False(t, config.Set(context.Background(), &model).HasError())
			response := datasource.ReadResponse{State: tfsdk.State{Schema: definition.Schema}}
			d.Read(context.Background(), datasource.ReadRequest{
				Config: tfsdk.Config{Schema: definition.Schema, Raw: config.Raw},
			}, &response)
			require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
			var got smartComputerGroupV3DataSourceModel
			require.False(t, response.State.Get(context.Background(), &got).HasError())
			require.Equal(t, "42", got.ID.ValueString())
			require.Equal(t, "empty", got.Name.ValueString())
			require.True(t, got.Criteria.IsNull())
		})
	}
}
