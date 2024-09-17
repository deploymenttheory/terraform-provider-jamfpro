// state.go
package managedsoftwareupdates

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResourceManagedSoftwareUpdatePlan
// information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceManagedSoftwareUpdatePlan) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"plan_uuid":                     resp.PlanUuid,
		"recipe_id":                     resp.RecipeId,
		"update_action":                 resp.UpdateAction,
		"version_type":                  resp.VersionType,
		"specific_version":              resp.SpecificVersion,
		"build_version":                 resp.BuildVersion,
		"max_deferrals":                 resp.MaxDeferrals,
		"force_install_local_date_time": resp.ForceInstallLocalDateTime,
	}

	// Set device information
	if resp.Device.DeviceId != "" {
		device := []map[string]interface{}{
			{
				"device_id":   resp.Device.DeviceId,
				"object_type": resp.Device.ObjectType,
			},
		}
		if err := d.Set("devices", device); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set status information
	if resp.Status.State != "" {
		status := []map[string]interface{}{
			{
				"state":         resp.Status.State,
				"error_reasons": resp.Status.ErrorReasons,
			},
		}
		if err := d.Set("status", status); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set config information
	config := []map[string]interface{}{
		{
			"update_action":                 resp.UpdateAction,
			"version_type":                  resp.VersionType,
			"specific_version":              resp.SpecificVersion,
			"build_version":                 resp.BuildVersion,
			"max_deferrals":                 resp.MaxDeferrals,
			"force_install_local_date_time": resp.ForceInstallLocalDateTime,
		},
	}
	if err := d.Set("config", config); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	for k, v := range resourceData {
		if err := d.Set(k, v); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
