package cloud_distribution_point

import (
	"context"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead retrieves the Cloud Distribution Point configuration from Jamf Pro.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	config, err := client.GetCloudDistributionPoint()
	if err != nil {
		return diag.FromErr(err)
	}

	uploadCap, err := client.GetCloudDistributionPointUploadCapability()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("has_connection_succeeded", config.HasConnectionSucceeded); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("message", config.Message); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inventory_id", config.InventoryId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cdn_type", config.CdnType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("master", config.Master); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("username", config.Username); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("directory", config.Directory); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cdn_url", config.CdnUrl); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("upload_url", config.UploadUrl); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("download_url", config.DownloadUrl); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secondary_auth_required", config.SecondaryAuthRequired); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secondary_auth_status_code", config.SecondaryAuthStatusCode); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secondary_auth_time_to_live", config.SecondaryAuthTimeToLive); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("require_signed_urls", config.RequireSignedUrls); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("key_pair_id", config.KeyPairId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("expiration_seconds", config.ExpirationSeconds); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("principal_distribution_technology", uploadCap.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("direct_upload_capable", uploadCap.Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().UTC().String())

	return diags
}
