package smart_computer_group_v3

import "github.com/hashicorp/terraform-plugin-framework/types"

type smartComputerGroupV3DataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SiteID      types.String `tfsdk:"site_id"`
	Criteria    types.List   `tfsdk:"criteria"`
}
