package provider

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/jamf_cloud_ip_address_list"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_computer_group"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/smart_mobile_device_group"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (p *frameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		jamf_cloud_ip_address_list.NewJamfCloudIPAddressListDataSource,
		smart_computer_group.NewSmartComputerGroupFrameworkDataSource,
		smart_mobile_device_group.NewSmartMobileDeviceGroupFrameworkDataSource,
	}
}
