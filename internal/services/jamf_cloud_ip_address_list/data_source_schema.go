package jamf_cloud_ip_address_list

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const jamfCloudIPListURL = "https://engineering.jamf.com/tc-docs-public-ip-lists/public-ips.json"

// Ensure the implementation satisfies the datasource.DataSource interface.
var _ datasource.DataSource = &JamfCloudIPAddressListDataSource{}

// JamfCloudIPAddressListDataSource defines the data source for fetching Jamf Cloud public IP addresses.
type JamfCloudIPAddressListDataSource struct{}

// NewJamfCloudIPAddressListDataSource returns a new instance of the data source.
func NewJamfCloudIPAddressListDataSource() datasource.DataSource {
	return &JamfCloudIPAddressListDataSource{}
}

// Metadata sets the data source type name.
func (d *JamfCloudIPAddressListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jamf_cloud_ip_address_list"
}

// Configure is a no-op as this data source fetches from a public URL and does not require provider configuration.
func (d *JamfCloudIPAddressListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

// Schema defines the data source schema.
func (d *JamfCloudIPAddressListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the public IP address list from Jamf Cloud for use in firewall rules and network configurations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this data source instance.",
			},
			"service_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter results by service name (e.g., `jamf_pro_cloud`, `jamf_cloud_services`, `jamf_cloud_distribution_service`).",
			},
			"provider_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter results by cloud provider (e.g., `aws`, `azure`).",
			},
			"traffic_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter results by traffic direction (e.g., `inbound`, `outbound`).",
			},
			"region_filter": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Filter results by region. Known regions: AWS - `us-all-regions`, `us-stateramp`, " +
					"`us-gov`, `eu-central-1`, `eu-west-2`, `ap-southeast-2`, `ap-northeast-1`. Azure - `centralus`, `germanywestcentral`.",
			},
			"publish_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publish date of the IP address list.",
			},
			"public_ips": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of public IP entries from Jamf Cloud.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Jamf service name (e.g., `jamf_pro_cloud`, `jamf_cloud_services`, `jamf_cloud_distribution_service`).",
						},
						"provider": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The cloud provider (e.g., `aws`, `azure`).",
						},
						"traffic": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The traffic direction (`inbound` or `outbound`).",
						},
						"region": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The cloud region (e.g., `us-all-regions`, `eu-central-1`).",
						},
						"ip_prefixes": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of IP prefixes in CIDR notation.",
						},
						"fqdns": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of fully qualified domain names (FQDNs).",
						},
					},
				},
			},
		},
	}
}
