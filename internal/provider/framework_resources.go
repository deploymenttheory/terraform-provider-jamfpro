package provider

import (
	"context"

	jamfProCloudDistributionPoint "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/cloud_distribution_point"
	jamfProDockItem "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/dock_item"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *frameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		jamfProCloudDistributionPoint.NewCloudDistributionPointFrameworkResource,
		jamfProDockItem.NewDockItemFrameworkResource,
	}
}
