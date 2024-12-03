// icons_state.go
package icons

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Icon information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResponseIconUpload) diag.Diagnostics {
	var diags diag.Diagnostics

	iconData := map[string]interface{}{
		"name":                 resp.Name,
		"url":                  resp.URL,
		"icon_file_path":       d.Get("icon_file_path").(string),
		"icon_file_web_source": d.Get("icon_file_web_source").(string),
	}

	for key, val := range iconData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
