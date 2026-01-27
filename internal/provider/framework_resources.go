package provider

import (
	"context"

	jamfProAdcsSettings "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/adcs_settings"
	jamfProDockItem "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/dock_item"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *frameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		jamfProAdcsSettings.NewAdcsSettingsFrameworkResource,
		jamfProDockItem.NewDockItemFrameworkResource,
	}
}
