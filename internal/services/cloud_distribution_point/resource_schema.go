package cloud_distribution_point

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
	_ resource.Resource                     = &cloudDistributionPointFrameworkResource{}
	_ resource.ResourceWithConfigure        = &cloudDistributionPointFrameworkResource{}
	_ resource.ResourceWithConfigValidators = &cloudDistributionPointFrameworkResource{}
	_ resource.ResourceWithImportState      = &cloudDistributionPointFrameworkResource{}
)

// NewCloudDistributionPointFrameworkResource returns the framework resource implementation.
func NewCloudDistributionPointFrameworkResource() resource.Resource {
	return &cloudDistributionPointFrameworkResource{}
}

// cloudDistributionPointFrameworkResource defines the framework resource implementation.
type cloudDistributionPointFrameworkResource struct {
	client *jamfpro.Client
}

func (r *cloudDistributionPointFrameworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_distribution_point"
}

func (r *cloudDistributionPointFrameworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudDistributionPointFrameworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *cloudDistributionPointFrameworkResource) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		cloudDistributionPointConfigValidator{},
	}
}

func (r *cloudDistributionPointFrameworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the Jamf Pro Cloud Distribution Point using the `/api/v1/cloud-distribution-point` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier reported by Jamf Pro for the Cloud Distribution Point.",
				Computed:            true,
			},
			"inventory_id": schema.StringAttribute{
				MarkdownDescription: "Inventory ID associated with the Cloud Distribution Point.",
				Computed:            true,
			},
			"cdn_type": schema.StringAttribute{
				MarkdownDescription: "CDN backing the Cloud Distribution Point. Allowed values: `" + cdnTypeJamfCloud + "`, `" + cdnTypeRackspace + "`, `" + cdnTypeAmazonS3 + "`, `" + cdnTypeAkamai + "`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(cdnTypeNone, cdnTypeJamfCloud, cdnTypeRackspace, cdnTypeAmazonS3, cdnTypeAkamai),
				},
			},
			"master": schema.BoolAttribute{
				MarkdownDescription: "Set to true to make this the principal distribution point.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username or access key for Rackspace, Amazon S3, or Akamai CDNs. Required when `cdn_type` is `" + cdnTypeRackspace + "`, `" + cdnTypeAmazonS3 + "`, or `" + cdnTypeAkamai + "` and must be omitted for other CDN types.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password or secret key for Rackspace, Amazon S3, or Akamai CDNs. Required when `cdn_type` is `" + cdnTypeRackspace + "`, `" + cdnTypeAmazonS3 + "`, or `" + cdnTypeAkamai + "` and must be omitted for other CDN types.",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"directory": schema.StringAttribute{
				MarkdownDescription: "Directory or storage path for Akamai NetStorage. Only valid and required when `cdn_type` is `" + cdnTypeAkamai + "`.",
				Optional:            true,
			},
			"cdn_url": schema.StringAttribute{
				MarkdownDescription: "CDN URL reported by Jamf Pro after configuration.",
				Computed:            true,
			},
			"upload_url": schema.StringAttribute{
				MarkdownDescription: "Upload endpoint for Akamai NetStorage. Only valid and required when `cdn_type` is `" + cdnTypeAkamai + "`.",
				Optional:            true,
			},
			"download_url": schema.StringAttribute{
				MarkdownDescription: "Download endpoint for Akamai EdgeSuite. Only valid and required when `cdn_type` is `" + cdnTypeAkamai + "`.",
				Optional:            true,
			},
			"secondary_auth_required": schema.BoolAttribute{
				MarkdownDescription: "Enable remote authentication for Akamai deliveries. Only valid when `cdn_type` is `" + cdnTypeAkamai + "` and must be explicitly set for that CDN.",
				Computed:            true,
				Optional:            true,
			},
			"secondary_auth_status_code": schema.Int64Attribute{
				MarkdownDescription: "HTTP response code returned during Akamai secondary authentication. Only valid when `cdn_type` is `" + cdnTypeAkamai + "` and `secondary_auth_required` is true.",
				Computed:            true,
				Optional:            true,
			},
			"secondary_auth_time_to_live": schema.Int64Attribute{
				MarkdownDescription: "Lifetime of Akamai secondary authentication tokens in seconds. Only valid when `cdn_type` is `" + cdnTypeAkamai + "` and `secondary_auth_required` is true.",
				Computed:            true,
				Optional:            true,
			},
			"require_signed_urls": schema.BoolAttribute{
				MarkdownDescription: "Enable CloudFront signed URLs for Amazon S3 distributions. Only valid when `cdn_type` is `" + cdnTypeAmazonS3 + "` and must be provided (true or false) for that CDN.",
				Computed:            true,
				Optional:            true,
			},
			"key_pair_id": schema.StringAttribute{
				MarkdownDescription: "CloudFront key pair ID used to generate signed URLs. Only valid when `cdn_type` is `" + cdnTypeAmazonS3 + "` and `require_signed_urls` is true.",
				Optional:            true,
			},
			"expiration_seconds": schema.Int64Attribute{
				MarkdownDescription: "Lifetime of CloudFront signed URLs in seconds. Only valid when `cdn_type` is `" + cdnTypeAmazonS3 + "` and `require_signed_urls` is true.",
				Computed:            true,
				Optional:            true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "CloudFront private key content used when signed URLs are enabled. Only valid when `cdn_type` is `" + cdnTypeAmazonS3 + "` and `require_signed_urls` is true.",
				Optional:            true,
				Sensitive:           true,
			},
			"has_connection_succeeded": schema.BoolAttribute{
				MarkdownDescription: "Result of the most recent Cloud Distribution Point connectivity test.",
				Computed:            true,
			},
			"message": schema.StringAttribute{
				MarkdownDescription: "Additional context provided by Jamf Pro for the last connection test.",
				Computed:            true,
			},
			"principal_distribution_technology": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the configured CDN is the principal distribution technology.",
				Computed:            true,
			},
			"direct_upload_capable": schema.BoolAttribute{
				MarkdownDescription: "Reports whether direct uploads are supported for the current CDN configuration.",
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
