package provider

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/jamf_cloud_ip_address_list"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (p *frameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		jamf_cloud_ip_address_list.NewJamfCloudIPAddressListDataSource,
	}
}
