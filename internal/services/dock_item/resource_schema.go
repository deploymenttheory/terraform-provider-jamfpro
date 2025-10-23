package dock_item

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	commonschema "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	ResourceName  = "jamfpro_dock_item_framework"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &dockItemFrameworkResource{}
	_ resource.ResourceWithConfigure   = &dockItemFrameworkResource{}
	_ resource.ResourceWithImportState = &dockItemFrameworkResource{}
)

// NewDockItemFrameworkResource is a helper function to simplify the provider implementation.
func NewDockItemFrameworkResource() resource.Resource {
	return &dockItemFrameworkResource{}
}

// dockItemFrameworkResource defines the resource implementation.
type dockItemFrameworkResource struct {
	client *jamfpro.Client
}

func (r *dockItemFrameworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dock_item"
}

func (r *dockItemFrameworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dockItemFrameworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *dockItemFrameworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Jamf Pro Dock Item with the `/api/v1/dock-items` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the dock item.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the dock item.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the dock item. Must be one of: `App`, `File`, `Folder`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("App", "File", "Folder"),
				},
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The path of the dock item. e.g `file://localhost/Applications/iTunes.app`",
				Required:            true,
			},
			"contents": schema.StringAttribute{
				MarkdownDescription: "Contents of the dock item.",
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
