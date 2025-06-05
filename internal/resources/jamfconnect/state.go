package jamfconnect

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Computer Check-In information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceJamfConnectConfigProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	profileData := map[string]interface{}{
		"config_profile_uuid":  resp.UUID,
		"profile_id":           resp.ProfileID,
		"profile_name":         resp.ProfileName,
		"scope_description":    resp.ScopeDescription,
		"site_id":              resp.SiteID,
		"jamf_connect_version": resp.Version,
		"auto_deployment_type": resp.AutoDeploymentType,
	}

	for key, val := range profileData {

		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting %s: %v", key, err))...)
		}
	}

	return diags
}
