package service_discovery_enrollment_well_known_settings

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

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource{}
	_ resource.ResourceWithConfigure   = &serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource{}
	_ resource.ResourceWithImportState = &serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource{}
)

// NewServiceDiscoveryEnrollmentWellKnownSettingsFrameworkResource returns the framework resource implementation.
func NewServiceDiscoveryEnrollmentWellKnownSettingsFrameworkResource() resource.Resource {
	return &serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource{}
}

type serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource struct {
	client *jamfpro.Client
}

func (r *serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_discovery_enrollment_well_known_settings"
}

func (r *serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jamfpro.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *jamfpro.Client. Please report this issue to the provider developers.",
		)
		return
	}

	r.client = client
}

func (r *serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Jamf Pro service discovery enrollment well-known settings via the `/api/v1/service-discovery-enrollment/well-known-settings` endpoint. Requires Jamf Pro 11.25 or later.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier for this singleton configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"well_known_settings": schema.ListNestedBlock{
				MarkdownDescription: "List of service discovery enrollment well-known settings.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"org_name": schema.StringAttribute{
							MarkdownDescription: "Organization name reported by Jamf Pro.",
							Computed:            true,
						},
						"server_uuid": schema.StringAttribute{
							MarkdownDescription: "Automated Device Enrollment Device Management Service UUID for the organization.",
							Required:            true,
						},
						"enrollment_type": schema.StringAttribute{
							MarkdownDescription: "Enrollment type for the organization.",
							Required:            true,
							Validators:          []validator.String{stringvalidator.OneOf("mdm-adde", "mdm-byod", "none")},
						},
					},
				},
			},
		},
	}
}
