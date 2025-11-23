package provider

import (
	"context"

	jamfProDockItem "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/dock_item"
	jamfProStaticComputerGroup "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/static_computer_group"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *frameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		jamfProDockItem.NewDockItemFrameworkResource,
		jamfProStaticComputerGroup.NewStaticComputerGroupFrameworkResource,
	}
}
