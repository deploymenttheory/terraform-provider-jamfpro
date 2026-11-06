package smart_computer_group_v3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *smartComputerGroupV3DataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var data smartComputerGroupV3DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Smart Computer Group V3", "The provider has not supplied a Jamf Pro SDK client.")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	id := data.ID.ValueString()
	if id == "" {
		item, err := d.client.GetSmartComputerGroupByNameV3WithContext(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error Finding Smart Computer Group V3", err.Error())
			return
		}
		if item == nil {
			resp.Diagnostics.AddError("Error Finding Smart Computer Group V3", errMissingResponse.Error())
			return
		}
		id = item.ID
	}
	group, err := d.client.GetSmartComputerGroupByIDV3WithContext(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Smart Computer Group V3", err.Error())
		return
	}
	if group == nil {
		resp.Diagnostics.AddError("Error Reading Smart Computer Group V3", errMissingResponse.Error())
		return
	}
	model := smartComputerGroupV3ResourceModel{ID: types.StringValue(id)}
	resp.Diagnostics.Append(state(ctx, &model, group)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID, data.Name, data.Description = model.ID, model.Name, model.Description
	data.SiteID, data.Criteria = model.SiteID, model.Criteria
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
