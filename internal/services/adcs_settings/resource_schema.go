package adcs_settings

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	commonschema "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	ResourceName  = "jamfpro_adcs_settings_framework"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &adcsSettingsFrameworkResource{}
	_ resource.ResourceWithConfigure   = &adcsSettingsFrameworkResource{}
	_ resource.ResourceWithImportState = &adcsSettingsFrameworkResource{}
)

func NewAdcsSettingsFrameworkResource() resource.Resource {
	return &adcsSettingsFrameworkResource{}
}

type adcsSettingsFrameworkResource struct {
	client *jamfpro.Client
}

func (r *adcsSettingsFrameworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adcs_settings"
}

func (r *adcsSettingsFrameworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *adcsSettingsFrameworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *adcsSettingsFrameworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Jamf Pro AD CS connector settings via the `/api/v1/pki/adcs-settings` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the AD CS settings record.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Friendly name shown in Jamf Pro.",
				Required:            true,
			},
			"ca_name": schema.StringAttribute{
				MarkdownDescription: "The Certificate Authority (CA) common name configured for the connector.",
				Required:            true,
			},
			"fqdn": schema.StringAttribute{
				MarkdownDescription: "FQDN of the Windows Server hosting the AD CS connector.",
				Required:            true,
			},
			"adcs_url": schema.StringAttribute{
				MarkdownDescription: "URL Jamf Pro uses to reach the AD CS connector (for example `https://servername.domain.tld/certsrv`).",
				Optional:            true,
				Computed:            true,
			},
			"api_client_id": schema.StringAttribute{
				MarkdownDescription: "Jamf Pro API Client ID used when the connector communicates outbound.",
				Required:            true,
			},
			"revocation_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable certificate revocation checking when issuing certificates.",
				Required:            true,
			},
			"outbound": schema.BoolAttribute{
				MarkdownDescription: "Set to true when the connector operates in outbound mode.",
				Required:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"connector_last_check_in_timestamp": schema.StringAttribute{
				MarkdownDescription: "Timestamp of the most recent connector heartbeat reported by Jamf Pro.",
				Computed:            true,
			},
			"server_certificate_filename": schema.StringAttribute{
				MarkdownDescription: "Filename that will be stored in Jamf Pro for the uploaded server certificate package.",
				Optional:            true,
				Computed:            true,
			},
			"server_certificate_data": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded PKCS#12 data for the server certificate used by the connector.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("server_certificate_filename")),
				},
			},
			"server_certificate_password": schema.StringAttribute{
				MarkdownDescription: "Password required to decrypt the server certificate package, if applicable.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("server_certificate_data")),
				},
			},
			"server_certificate_serial_number": schema.StringAttribute{
				MarkdownDescription: "Serial number reported by Jamf Pro for the current server certificate.",
				Computed:            true,
			},
			"server_certificate_subject": schema.StringAttribute{
				MarkdownDescription: "Subject of the installed server certificate.",
				Computed:            true,
			},
			"server_certificate_issuer": schema.StringAttribute{
				MarkdownDescription: "Issuer of the installed server certificate.",
				Computed:            true,
			},
			"server_certificate_expiration_date": schema.StringAttribute{
				MarkdownDescription: "Expiration timestamp for the server certificate stored in Jamf Pro.",
				Computed:            true,
			},
			"client_certificate_filename": schema.StringAttribute{
				MarkdownDescription: "Filename that will be stored in Jamf Pro for the uploaded client certificate package.",
				Optional:            true,
				Computed:            true,
			},
			"client_certificate_data": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded PKCS#12 data for the client certificate used by the connector.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("client_certificate_filename")),
				},
			},
			"client_certificate_password": schema.StringAttribute{
				MarkdownDescription: "Password required to decrypt the client certificate package, if applicable.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("client_certificate_data")),
				},
			},
			"client_certificate_serial_number": schema.StringAttribute{
				MarkdownDescription: "Serial number reported by Jamf Pro for the current client certificate.",
				Computed:            true,
			},
			"client_certificate_subject": schema.StringAttribute{
				MarkdownDescription: "Subject of the installed client certificate.",
				Computed:            true,
			},
			"client_certificate_issuer": schema.StringAttribute{
				MarkdownDescription: "Issuer of the installed client certificate.",
				Computed:            true,
			},
			"client_certificate_expiration_date": schema.StringAttribute{
				MarkdownDescription: "Expiration timestamp for the client certificate stored in Jamf Pro.",
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
