package static_computer_group

import "github.com/hashicorp/terraform-plugin-framework/types"

// staticComputerGroupDataSourceModel describes the data source data model.
type staticComputerGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SiteID      types.String `tfsdk:"site_id"`
}
