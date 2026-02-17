package provider

import (
	"context"

	jamfProAdcsSettings "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/adcs_settings"
	jamfProCloudDistributionPoint "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/cloud_distribution_point"
	jamfProDockItem "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/dock_item"
	jamfProServiceDiscoveryEnrollmentWellKnownSettings "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/service_discovery_enrollment_well_known_settings"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *frameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		jamfProAdcsSettings.NewAdcsSettingsFrameworkResource,
		jamfProCloudDistributionPoint.NewCloudDistributionPointFrameworkResource,
		jamfProDockItem.NewDockItemFrameworkResource,
		jamfProServiceDiscoveryEnrollmentWellKnownSettings.NewServiceDiscoveryEnrollmentWellKnownSettingsFrameworkResource,
	}
}
