package static_computer_group

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &staticComputerGroupFrameworkDataSource{}

// staticComputerGroupFrameworkDataSource defines the data source implementation.
type staticComputerGroupFrameworkDataSource struct {
	client *jamfpro.Client
}

// NewStaticComputerGroupFrameworkDataSource creates a new instance of the data source.
func NewStaticComputerGroupFrameworkDataSource() datasource.DataSource {
	return &staticComputerGroupFrameworkDataSource{}
}

// Metadata returns the data source type name.
func (d *staticComputerGroupFrameworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_computer_group"
}

// Configure adds the provider configured client to the data source.
func (d *staticComputerGroupFrameworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *staticComputerGroupFrameworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for retrieving a static computer group from Jamf Pro.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the static computer group.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the static computer group.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the static computer group.",
			},
			"site_id": schema.StringAttribute{
				Computed:    true,
				Description: "The site ID for the group.",
			},
		},
	}
}
