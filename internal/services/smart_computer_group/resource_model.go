package smart_computer_group

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// smartComputerGroupResourceModel describes the resource data model.
type smartComputerGroupResourceModel struct {
	ID          types.String                          `tfsdk:"id"`
	Name        types.String                          `tfsdk:"name"`
	Description types.String                          `tfsdk:"description"`
	SiteID      types.String                          `tfsdk:"site_id"`
	Criteria    []smartComputerGroupCriteriaDataModel `tfsdk:"criteria"`
	Timeouts    timeouts.Value                        `tfsdk:"timeouts"`
}

// smartComputerGroupCriteriaDataModel describes the criteria data model.
type smartComputerGroupCriteriaDataModel struct {
	Name         types.String `tfsdk:"name"`
	Priority     types.Int64  `tfsdk:"priority"`
	AndOr        types.String `tfsdk:"and_or"`
	SearchType   types.String `tfsdk:"search_type"`
	Value        types.String `tfsdk:"value"`
	OpeningParen types.Bool   `tfsdk:"opening_paren"`
	ClosingParen types.Bool   `tfsdk:"closing_paren"`
}
