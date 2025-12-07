package smart_mobile_device_group

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	commonschema "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &smartMobileDeviceGroupFrameworkDataSource{}

// smartMobileDeviceGroupFrameworkDataSource defines the data source implementation.
type smartMobileDeviceGroupFrameworkDataSource struct {
	client *jamfpro.Client
}

// NewSmartMobileDeviceGroupFrameworkDataSource creates a new instance of the data source.
func NewSmartMobileDeviceGroupFrameworkDataSource() datasource.DataSource {
	return &smartMobileDeviceGroupFrameworkDataSource{}
}

// Metadata returns the data source type name.
func (d *smartMobileDeviceGroupFrameworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smart_mobile_device_group"
}

// Configure adds the provider configured client to the data source.
func (d *smartMobileDeviceGroupFrameworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jamfpro.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Schema defines the schema for the data source.
func (d *smartMobileDeviceGroupFrameworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for retrieving a Jamf Pro Smart Mobile Device Group using the `/api/v1/mobile-device-groups/smart-groups` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the smart mobile device group.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the smart mobile device group.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the smart mobile device group.",
			},
			"site_id": commonschema.SiteID(ctx),
		},
		Blocks: map[string]schema.Block{
			"criteria": commonschema.CriteriaDataSource(ctx),
		},
	}
}
