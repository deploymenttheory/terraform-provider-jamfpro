// jamf_protect_state.go
package jamfprotect

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the state of the Jamf Protect settings resource in Terraform
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceJamfProtectIntegrationSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := map[string]interface{}{
		"id":              resp.ID,
		"api_client_id":   resp.APIClientID,
		"api_client_name": resp.APIClientName,
		"registration_id": resp.RegistrationID,
		"protect_url":     resp.ProtectURL,
		"last_sync_time":  resp.LastSyncTime,
		"sync_status":     resp.SyncStatus,
		"auto_install":    resp.AutoInstall,
	}

	for key, val := range settings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// updateStateFromRegisterResponse updates the state from the initial registration response
func updateStateFromRegisterResponse(d *schema.ResourceData, resp *jamfpro.ResourceJamfProtectRegisterResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := map[string]interface{}{
		"id":              resp.ID,
		"api_client_id":   resp.APIClientID,
		"api_client_name": resp.APIClientName,
		"registration_id": resp.RegistrationID,
		"protect_url":     resp.ProtectURL,
		"last_sync_time":  resp.LastSyncTime,
		"sync_status":     resp.SyncStatus,
		"auto_install":    resp.AutoInstall,
	}

	for key, val := range settings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
