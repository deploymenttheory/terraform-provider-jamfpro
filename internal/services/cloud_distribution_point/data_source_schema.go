package cloud_distribution_point

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProCloudDistributionPoint returns a Terraform data source for Jamf Pro Cloud Distribution Point.
func DataSourceJamfProCloudDistributionPoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"has_connection_succeeded": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the connection to the Cloud Distribution Point was successful. If true, the connection was successful. If false, the connection failed.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A message detailing the result of the connection test. This could be a success message or an error message if the connection failed.",
			},
			"inventory_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier (inventoryId) that links the Cloud Distribution Point to the inventory data. This ID associates the distribution point with specific software or files in the inventory.",
			},
			"cdn_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the content delivery network (CDN) used to distribute content for the cloud distribution point.",
			},
			"master": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Use as principal distribution point. Use as the authoritative source for all files.",
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The username or access key used for authenticating with the selected CDN. Required for Rackspace Cloud Files, " +
					"Amazon Web Services (AWS Access Key ID), or Akamai authentication.",
			},
			"directory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The directory or path for content delivery in Akamai. Required when cdnType is set to Akamai.",
			},
			"cdn_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CDN URL for the Cloud Distribution Point. Format varies by provider (Rackspace, AWS, Akamai).",
			},
			"upload_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL used to upload files to Akamai's NetStorage. Required when cdnType is Akamai.",
			},
			"download_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL used to access and download content from Akamai's EdgeSuite. Required when cdnType is Akamai.",
			},
			"secondary_auth_required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable Remote Authentication. Authorize requests for files stored on the distribution point. Required when cdnType is Akamai.",
			},
			"secondary_auth_status_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Configure the HTTP response code returned by Jamf Pro during remote authentication. Required when cdnType is Akamai and secondaryAuthRequired is true.",
			},
			"secondary_auth_time_to_live": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of seconds before the authorization token expires. Required when cdnType is Akamai and secondaryAuthRequired is true.",
			},
			"require_signed_urls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Amazon Sign URL. Restricts access to requests that use a signed URL. Required when cdnType is Amazon Web Services.",
			},
			"key_pair_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The CloudFront Access Key ID used to generate signed URLs for secure access to content. " +
					"Required when cdnType is Amazon Web Services and requireSignedUrls is true.",
			},
			"expiration_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of seconds before the signed URL expires. Required when cdnType is Amazon Web Services and requireSignedUrls is true.",
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The CloudFront Private Key file required when cdnType is Amazon Web Services and requireSignedUrls is enabled. " +
					"Used for signing URLs for restricted access to CloudFront-distributed content. Supports .pem or .der formats.",
			},
			"principal_distribution_technology": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if this is the principal distribution technology.",
			},
			"direct_upload_capable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the distribution point supports direct file uploads.",
			},
		},
	}
}
