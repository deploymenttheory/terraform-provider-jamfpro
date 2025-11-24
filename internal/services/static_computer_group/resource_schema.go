package static_computer_group

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	commonschema "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ResourceName  = "jamfpro_static_computer_group"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &staticComputerGroupFrameworkResource{}
	_ resource.ResourceWithConfigure    = &staticComputerGroupFrameworkResource{}
	_ resource.ResourceWithImportState  = &staticComputerGroupFrameworkResource{}
	_ resource.ResourceWithUpgradeState = &staticComputerGroupFrameworkResource{}
)

// NewStaticComputerGroupFrameworkResource is a helper function to simplify the provider implementation.
func NewStaticComputerGroupFrameworkResource() resource.Resource {
	return &staticComputerGroupFrameworkResource{}
}

// staticComputerGroupFrameworkResource defines the resource implementation.
type staticComputerGroupFrameworkResource struct {
	client *jamfpro.Client
}

func (r *staticComputerGroupFrameworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_computer_group"
}

func (r *staticComputerGroupFrameworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *staticComputerGroupFrameworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *staticComputerGroupFrameworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:             1,
		MarkdownDescription: "Manages a Jamf Pro Static Computer Group using the `/api/v2/computer-groups/static-groups` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the static computer group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the static computer group.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the static computer group.",
				Optional:            true,
			},
			"assigned_computer_ids": schema.SetAttribute{
				MarkdownDescription: "Set of computer IDs assigned to this static computer group. Note: This value cannot be read back from the API.",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id":  commonschema.SiteID(ctx),
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
