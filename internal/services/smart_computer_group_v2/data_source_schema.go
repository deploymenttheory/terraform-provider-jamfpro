package smart_computer_group_v2

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	commonschema "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &smartComputerGroupV2FrameworkDataSource{}

// smartComputerGroupV2FrameworkDataSource defines the data source implementation.
type smartComputerGroupV2FrameworkDataSource struct {
	client *jamfpro.Client
}

// NewSmartComputerGroupV2FrameworkDataSource creates a new instance of the data source.
func NewSmartComputerGroupV2FrameworkDataSource() datasource.DataSource {
	return &smartComputerGroupV2FrameworkDataSource{}
}

// Metadata returns the data source type name.
func (d *smartComputerGroupV2FrameworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smart_computer_group_v2"
}

// Configure adds the provider configured client to the data source.
func (d *smartComputerGroupV2FrameworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *smartComputerGroupV2FrameworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for retrieving a Jamf Pro Smart Computer Group using the `/api/v2/computer-groups/smart-groups` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the smart computer group.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the smart computer group.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the smart computer group.",
			},
			"site_id": commonschema.SiteID(ctx),
		},
		Blocks: map[string]schema.Block{
			"criteria": commonschema.CriteriaDataSource(ctx),
		},
	}
}
