package dock_item

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dockItemFrameworkResourceModel describes the resource data model.
type dockItemFrameworkResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Name     types.String   `tfsdk:"name"`
	Type     types.String   `tfsdk:"type"`
	Path     types.String   `tfsdk:"path"`
	Contents types.String   `tfsdk:"contents"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
