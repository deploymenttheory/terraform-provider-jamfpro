package cloud_distribution_point

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ResourceName  = "jamfpro_cloud_distribution_point_framework"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	cloudDistributionPointSingletonID = "jamfpro_cloud_distribution_point_singleton"

	cdnTypeNone      = "NONE"
	cdnTypeJamfCloud = "JAMF_CLOUD"
	cdnTypeRackspace = "RACKSPACE_CLOUD_FILES"
	cdnTypeAmazonS3  = "AMAZON_S3"
	cdnTypeAkamai    = "AKAMAI"
)

// cloudDistributionPointResourceModel models the framework resource data.
type cloudDistributionPointResourceModel struct {
	ID                              types.String   `tfsdk:"id"`
	InventoryID                     types.String   `tfsdk:"inventory_id"`
	CdnType                         types.String   `tfsdk:"cdn_type"`
	Master                          types.Bool     `tfsdk:"master"`
	Username                        types.String   `tfsdk:"username"`
	Password                        types.String   `tfsdk:"password"`
	Directory                       types.String   `tfsdk:"directory"`
	CdnURL                          types.String   `tfsdk:"cdn_url"`
	UploadURL                       types.String   `tfsdk:"upload_url"`
	DownloadURL                     types.String   `tfsdk:"download_url"`
	SecondaryAuthRequired           types.Bool     `tfsdk:"secondary_auth_required"`
	SecondaryAuthStatusCode         types.Int64    `tfsdk:"secondary_auth_status_code"`
	SecondaryAuthTimeToLive         types.Int64    `tfsdk:"secondary_auth_time_to_live"`
	RequireSignedUrls               types.Bool     `tfsdk:"require_signed_urls"`
	KeyPairID                       types.String   `tfsdk:"key_pair_id"`
	ExpirationSeconds               types.Int64    `tfsdk:"expiration_seconds"`
	PrivateKey                      types.String   `tfsdk:"private_key"`
	HasConnectionSucceeded          types.Bool     `tfsdk:"has_connection_succeeded"`
	Message                         types.String   `tfsdk:"message"`
	PrincipalDistributionTechnology types.Bool     `tfsdk:"principal_distribution_technology"`
	DirectUploadCapable             types.Bool     `tfsdk:"direct_upload_capable"`
	Timeouts                        timeouts.Value `tfsdk:"timeouts"`
}
