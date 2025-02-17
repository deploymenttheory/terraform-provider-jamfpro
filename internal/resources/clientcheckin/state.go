package computercheckin

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Client Check-In information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceClientCheckinSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	checkinData := map[string]interface{}{
		"check_in_frequency":                  resp.CheckInFrequency,
		"create_hooks":                        resp.CreateHooks,
		"hook_log":                            resp.HookLog,
		"hook_policies":                       resp.HookPolicies,
		"create_startup_script":               resp.CreateStartupScript,
		"startup_log":                         resp.StartupLog,
		"startup_policies":                    resp.StartupPolicies,
		"startup_ssh":                         resp.StartupSsh,
		"enable_local_configuration_profiles": resp.EnableLocalConfigurationProfiles,
	}

	for key, val := range checkinData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
