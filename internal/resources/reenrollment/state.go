package reenrollment

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Re-enrollment information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceReenrollmentSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	reenrollmentData := map[string]interface{}{
		"flush_location_information":         resp.FlushLocationInformation,
		"flush_location_information_history": resp.FlushLocationInformationHistory,
		"flush_policy_history":               resp.FlushPolicyHistory,
		"flush_extension_attributes":         resp.FlushExtensionAttributes,
		"flush_software_update_plans":        resp.FlushSoftwareUpdatePlans,
		"flush_mdm_queue":                    resp.FlushMdmQueue,
	}

	for key, val := range reenrollmentData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
