// state.go
package managed_software_update

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest ResponseManagedSoftwareUpdatePlan
// information from the Jamf Pro API.
func updateState(d *schema.ResourceData, plan *jamfpro.ResponseManagedSoftwareUpdatePlan) diag.Diagnostics {
	if plan == nil {
		return diag.Errorf("no managed software update plan found in the response")
	}

	d.SetId(plan.PlanUuid)
	if err := d.Set("plan_uuid", plan.PlanUuid); err != nil {
		return diag.FromErr(fmt.Errorf("error setting plan_uuid: %v", err))
	}

	// Set config fields
	config := map[string]interface{}{
		"update_action":                 plan.UpdateAction,
		"version_type":                  plan.VersionType,
		"specific_version":              plan.SpecificVersion,
		"max_deferrals":                 plan.MaxDeferrals,
		"force_install_local_date_time": plan.ForceInstallLocalDateTime,
	}
	if err := d.Set("config", []interface{}{config}); err != nil {
		return diag.FromErr(fmt.Errorf("error setting config: %v", err))
	}

	// Set group or device based on the object type
	if plan.Device.ObjectType == "COMPUTER_GROUP" || plan.Device.ObjectType == "MOBILE_DEVICE_GROUP" {
		group := map[string]interface{}{
			"group_id":    plan.Device.DeviceId,
			"object_type": plan.Device.ObjectType,
		}
		if err := d.Set("group", []interface{}{group}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting group: %v", err))
		}
	} else if plan.Device.ObjectType == "COMPUTER" || plan.Device.ObjectType == "MOBILE_DEVICE" {
		device := map[string]interface{}{
			"device_id":   plan.Device.DeviceId,
			"object_type": plan.Device.ObjectType,
		}
		if err := d.Set("device", []interface{}{device}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting device: %v", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("unknown object type: %s", plan.Device.ObjectType))
	}

	return nil
}
