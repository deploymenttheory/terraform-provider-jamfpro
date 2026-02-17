package smart_computer_group_v2

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// smartComputerGroupV2ResourceModel describes the resource data model.
type smartComputerGroupV2ResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	SiteID      types.String   `tfsdk:"site_id"`
	Criteria    types.List     `tfsdk:"criteria"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

// smartComputerGroupV2CriteriaDataModel describes the criteria data model.
type smartComputerGroupV2CriteriaDataModel struct {
	Name         types.String `tfsdk:"name"`
	Priority     types.Int32  `tfsdk:"priority"`
	AndOr        types.String `tfsdk:"and_or"`
	SearchType   types.String `tfsdk:"search_type"`
	Value        types.String `tfsdk:"value"`
	OpeningParen types.Bool   `tfsdk:"opening_paren"`
	ClosingParen types.Bool   `tfsdk:"closing_paren"`
}
