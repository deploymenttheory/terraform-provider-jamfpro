package static_computer_group

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// staticComputerGroupResourceModel describes the resource data model.
type staticComputerGroupResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	Name                types.String   `tfsdk:"name"`
	Description         types.String   `tfsdk:"description"`
	AssignedComputerIDs types.Set      `tfsdk:"assigned_computer_ids"`
	SiteID              types.String   `tfsdk:"site_id"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}
