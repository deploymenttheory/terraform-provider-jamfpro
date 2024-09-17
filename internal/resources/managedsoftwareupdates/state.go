// state.go
package managedsoftwareupdates

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResponseManagedSoftwareUpdatePlan
// information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResponseManagedSoftwareUpdatePlan) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"plan_uuid": resp.PlanUuid,
	}

	// Set group information if present
	if resp.Device.DeviceId != "" {
		group := []map[string]interface{}{
			{
				"group_id":    resp.Device.DeviceId, // Using DeviceId as GroupId
				"object_type": resp.Device.ObjectType,
			},
		}
		if err := d.Set("group", group); err != nil {
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

	// Set other top-level attributes
	for k, v := range resourceData {
		if err := d.Set(k, v); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
