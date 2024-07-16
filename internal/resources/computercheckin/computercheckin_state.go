// computercheckin_state.go
package computercheckin

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Check-In information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceComputerCheckin) diag.Diagnostics {
	var diags diag.Diagnostics

	checkinData := map[string]interface{}{
		"check_in_frequency":                       resp.CheckInFrequency,
		"create_startup_script":                    resp.CreateStartupScript,
		"log_startup_event":                        resp.LogStartupEvent,
		"check_for_policies_at_startup":            resp.CheckForPoliciesAtStartup,
		"apply_computer_level_managed_preferences": resp.ApplyComputerLevelManagedPrefs,
		"ensure_ssh_is_enabled":                    resp.EnsureSSHIsEnabled,
		"create_login_logout_hooks":                resp.CreateLoginLogoutHooks,
		"log_username":                             resp.LogUsername,
		"check_for_policies_at_login_logout":       resp.CheckForPoliciesAtLoginLogout,
		"apply_user_level_managed_preferences":     resp.ApplyUserLevelManagedPreferences,
		"hide_restore_partition":                   resp.HideRestorePartition,
		"perform_login_actions_in_background":      resp.PerformLoginActionsInBackground,
	}

	for key, val := range checkinData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
