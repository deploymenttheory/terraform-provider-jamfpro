package access_management_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Access Management settings information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceAccessManagementSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	accessManagementSettingsConfig := map[string]interface{}{
		"automated_device_enrollment_server_uuid": resp.AutomatedDeviceEnrollmentServerUuid,
	}

	for key, val := range accessManagementSettingsConfig {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
