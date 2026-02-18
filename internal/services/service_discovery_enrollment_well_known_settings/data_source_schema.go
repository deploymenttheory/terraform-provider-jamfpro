package service_discovery_enrollment_well_known_settings

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the datasource.DataSource interface.
var (
	_ datasource.DataSource              = &serviceDiscoveryEnrollmentWellKnownSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceDiscoveryEnrollmentWellKnownSettingsDataSource{}
)

// NewServiceDiscoveryEnrollmentWellKnownSettingsDataSource returns a new instance of the data source.
func NewServiceDiscoveryEnrollmentWellKnownSettingsDataSource() datasource.DataSource {
	return &serviceDiscoveryEnrollmentWellKnownSettingsDataSource{}
}

type serviceDiscoveryEnrollmentWellKnownSettingsDataSource struct {
	client *jamfpro.Client
}

func (d *serviceDiscoveryEnrollmentWellKnownSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_discovery_enrollment_well_known_settings"
}

func (d *serviceDiscoveryEnrollmentWellKnownSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jamfpro.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *jamfpro.Client. Please report this issue to the provider developers.",
		)
		return
	}

	d.client = client
}

func (d *serviceDiscoveryEnrollmentWellKnownSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads Jamf Pro service discovery enrollment well-known settings via the `/api/v1/service-discovery-enrollment/well-known-settings` endpoint. Requires Jamf Pro 11.25 or later",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier for this data source instance.",
				Computed:            true,
			},
			"well_known_settings": schema.ListNestedAttribute{
				MarkdownDescription: "List of service discovery enrollment well-known settings.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"org_name": schema.StringAttribute{
							MarkdownDescription: "Organization name reported by Jamf Pro.",
							Computed:            true,
						},
						"server_uuid": schema.StringAttribute{
							MarkdownDescription: "Jamf Pro server UUID for the organization.",
							Computed:            true,
						},
						"enrollment_type": schema.StringAttribute{
							MarkdownDescription: "Enrollment type for the organization.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
