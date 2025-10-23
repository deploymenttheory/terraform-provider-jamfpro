package jamf_protect

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the current settings of the Jamf Protect integration.
func updateState(d *schema.ResourceData, resp *jamfpro.ResponseJamfProtectSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := map[string]any{
		"client_id":       resp.APIClientID,
		"protect_url":     resp.ProtectURL,
		"auto_install":    resp.AutoInstall,
		"sync_status":     resp.SyncStatus,
		"api_client_name": resp.APIClientName,
		"last_sync_time":  resp.LastSyncTime,
		"registration_id": resp.RegistrationID,
	}

	if d.HasChange("password") {
		settings["password"] = d.Get("password").(string)
	}

	for key, val := range settings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
