package smart_computer_group

import "github.com/hashicorp/terraform-plugin-framework/types"

// smartComputerGroupDataSourceModel describes the data source data model.
type smartComputerGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SiteID      types.String `tfsdk:"site_id"`
	Criteria    types.List   `tfsdk:"criteria"`
}
