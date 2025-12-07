package smart_computer_group

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	commonschema "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

const (
	ResourceName  = "jamfpro_smart_computer_group"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &smartComputerGroupFrameworkResource{}
	_ resource.ResourceWithConfigure    = &smartComputerGroupFrameworkResource{}
	_ resource.ResourceWithImportState  = &smartComputerGroupFrameworkResource{}
	_ resource.ResourceWithUpgradeState = &smartComputerGroupFrameworkResource{}
)

// NewSmartComputerGroupFrameworkResource is a helper function to simplify the provider implementation.
func NewSmartComputerGroupFrameworkResource() resource.Resource {
	return &smartComputerGroupFrameworkResource{}
}

// smartComputerGroupFrameworkResource defines the resource implementation.
type smartComputerGroupFrameworkResource struct {
	client *jamfpro.Client
}

func (r *smartComputerGroupFrameworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smart_computer_group"
}

func (r *smartComputerGroupFrameworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jamfpro.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *jamfpro.Client, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	r.client = client
}

func (r *smartComputerGroupFrameworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *smartComputerGroupFrameworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:             1,
		MarkdownDescription: "Manages a Jamf Pro Smart Computer Group using the `/api/v2/computer-groups/smart-groups` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the smart computer group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the smart computer group.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the smart computer group.",
				Optional:            true,
			},
			"site_id":  commonschema.SiteID(ctx),
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"criteria": commonschema.CriteriaResource(ctx),
		},
	}
}
