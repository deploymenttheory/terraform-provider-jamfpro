// activationcode_state.go
package activationcode

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Check-In information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceComputerCheckin) diag.Diagnostics {

	var diags diag.Diagnostics

	// Map the configuration fields from the API response to a structured map
	checkinData := map[string]interface{}{
		"check_in_frequency":                       resource.CheckInFrequency,
		"create_startup_script":                    resource.CreateStartupScript,
		"log_startup_event":                        resource.LogStartupEvent,
		"check_for_policies_at_startup":            resource.CheckForPoliciesAtStartup,
		"apply_computer_level_managed_preferences": resource.ApplyComputerLevelManagedPrefs,
		"ensure_ssh_is_enabled":                    resource.EnsureSSHIsEnabled,
		"create_login_logout_hooks":                resource.CreateLoginLogoutHooks,
		"log_username":                             resource.LogUsername,
		"check_for_policies_at_login_logout":       resource.CheckForPoliciesAtLoginLogout,
		"apply_user_level_managed_preferences":     resource.ApplyUserLevelManagedPreferences,
		"hide_restore_partition":                   resource.HideRestorePartition,
		"perform_login_actions_in_background":      resource.PerformLoginActionsInBackground,
	}

	// Set the structured map in the Terraform state
	for key, val := range checkinData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}
