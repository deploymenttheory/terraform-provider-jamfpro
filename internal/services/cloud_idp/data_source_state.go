package cloud_idp

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Cloud Identity Provider information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceCloudIdentityProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"id":            resp.ID,
		"display_name":  resp.DisplayName,
		"enabled":       resp.Enabled,
		"provider_name": resp.ProviderName,
	}

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
