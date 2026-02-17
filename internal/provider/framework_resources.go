package provider

import (
	"context"

	jamfProAdcsSettings "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/adcs_settings"
	jamfProCloudDistributionPoint "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/cloud_distribution_point"
	jamfProDockItem "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/dock_item"
	jamfProSmartComputerGroupV2 "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_computer_group_v2"
	jamfProSmartMobileDeviceGroupV1 "github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_mobile_device_group_v1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *frameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		jamfProAdcsSettings.NewAdcsSettingsFrameworkResource,
		jamfProCloudDistributionPoint.NewCloudDistributionPointFrameworkResource,
		jamfProDockItem.NewDockItemFrameworkResource,
		jamfProSmartComputerGroupV2.NewSmartComputerGroupV2FrameworkResource,
		jamfProSmartMobileDeviceGroupV1.NewSmartMobileDeviceGroupV1FrameworkResource,
	}
}
