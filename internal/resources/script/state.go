// scripts_state.go
package script

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Script information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceScript) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"id":              resp.ID,
		"name":            resp.Name,
		"info":            resp.Info,
		"notes":           resp.Notes,
		"os_requirements": resp.OSRequirements,
		"category_id":     resp.CategoryId,
		"priority":        resp.Priority,
		"script_contents": resp.ScriptContents,
		"parameter4":      resp.Parameter4,
		"parameter5":      resp.Parameter5,
		"parameter6":      resp.Parameter6,
		"parameter7":      resp.Parameter7,
		"parameter8":      resp.Parameter8,
		"parameter9":      resp.Parameter9,
		"parameter10":     resp.Parameter10,
		"parameter11":     resp.Parameter11,
	}

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags

}
