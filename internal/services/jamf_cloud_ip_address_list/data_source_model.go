package jamf_cloud_ip_address_list

import "github.com/hashicorp/terraform-plugin-framework/types"

// JamfCloudIPAddressListDataSourceModel describes the Terraform state model.
type JamfCloudIPAddressListDataSourceModel struct {
	ID             types.String         `tfsdk:"id"`
	ServiceFilter  types.String         `tfsdk:"service_filter"`
	ProviderFilter types.String         `tfsdk:"provider_filter"`
	TrafficFilter  types.String         `tfsdk:"traffic_filter"`
	RegionFilter   types.String         `tfsdk:"region_filter"`
	PublishDate    types.String         `tfsdk:"publish_date"`
	PublicIPs      []PublicIPEntryModel `tfsdk:"public_ips"`
}

// PublicIPEntryModel describes the nested public IP entry in Terraform state.
type PublicIPEntryModel struct {
	Service    types.String   `tfsdk:"service"`
	Provider   types.String   `tfsdk:"provider"`
	Traffic    types.String   `tfsdk:"traffic"`
	Region     types.String   `tfsdk:"region"`
	IPPrefixes []types.String `tfsdk:"ip_prefixes"`
	FQDNs      []types.String `tfsdk:"fqdns"`
}

// ResponseJamfCloudIPAddressList represents the API response from the Jamf Cloud IP list endpoint.
type ResponseJamfCloudIPAddressList struct {
	PublishDate string                  `json:"publish_date"`
	PublicIPs   []ResponsePublicIPEntry `json:"public_ips"`
}

// ResponsePublicIPEntry represents a single entry in the API response.
type ResponsePublicIPEntry struct {
	Service    string   `json:"service"`
	Provider   string   `json:"provider"`
	Traffic    string   `json:"traffic"`
	Region     string   `json:"region"`
	IPPrefixes []string `json:"ip_prefixes,omitempty"`
	FQDNs      []string `json:"fqdns,omitempty"`
}
