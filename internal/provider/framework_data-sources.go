package provider

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/guid_list_sharder"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/jamf_cloud_ip_address_list"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/service_discovery_enrollment_well_known_settings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_computer_group_v2"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_mobile_device_group_v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (p *frameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		jamf_cloud_ip_address_list.NewJamfCloudIPAddressListDataSource,
		guid_list_sharder.NewGuidListSharderDataSource,
		smart_computer_group_v2.NewSmartComputerGroupV2FrameworkDataSource,
		smart_mobile_device_group_v1.NewSmartMobileDeviceGroupV1FrameworkDataSource,
		service_discovery_enrollment_well_known_settings.NewServiceDiscoveryEnrollmentWellKnownSettingsDataSource,
	}
}
