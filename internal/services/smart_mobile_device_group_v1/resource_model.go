package smart_mobile_device_group_v1

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// smartMobileDeviceGroupV1ResourceModel describes the resource data model.
type smartMobileDeviceGroupV1ResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	SiteID      types.String   `tfsdk:"site_id"`
	Criteria    types.List     `tfsdk:"criteria"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

// smartMobileDeviceGroupV1CriteriaDataModel describes the criteria data model.
type smartMobileDeviceGroupV1CriteriaDataModel struct {
	Name         types.String `tfsdk:"name"`
	Priority     types.Int32  `tfsdk:"priority"`
	AndOr        types.String `tfsdk:"and_or"`
	SearchType   types.String `tfsdk:"search_type"`
	Value        types.String `tfsdk:"value"`
	OpeningParen types.Bool   `tfsdk:"opening_paren"`
	ClosingParen types.Bool   `tfsdk:"closing_paren"`
}
