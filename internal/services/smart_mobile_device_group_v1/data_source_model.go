package smart_mobile_device_group_v1

import "github.com/hashicorp/terraform-plugin-framework/types"

// smartMobileDeviceGroupV1DataSourceModel describes the data source data model.
type smartMobileDeviceGroupV1DataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SiteID      types.String `tfsdk:"site_id"`
	Criteria    types.List   `tfsdk:"criteria"`
}
