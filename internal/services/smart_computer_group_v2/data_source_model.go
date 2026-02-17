package smart_computer_group_v2

import "github.com/hashicorp/terraform-plugin-framework/types"

// smartComputerGroupV2DataSourceModel describes the data source data model.
type smartComputerGroupV2DataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SiteID      types.String `tfsdk:"site_id"`
	Criteria    types.List   `tfsdk:"criteria"`
}
