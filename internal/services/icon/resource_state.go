// icons_state.go
package icon

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Icon information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResponseIconUpload) diag.Diagnostics {
	var diags diag.Diagnostics

	iconFilePath := ""
	iconFileWebSource := ""
	iconFileBase64 := ""

	if val := d.Get("icon_file_path").(string); val != "" {
		iconFilePath = val
	}
	if val := d.Get("icon_file_web_source").(string); val != "" {
		iconFileWebSource = val
	}
	if val := d.Get("icon_file_base64").(string); val != "" {
		iconFileBase64 = val
	}

	iconData := map[string]any{
		"name":                 resp.Name,
		"url":                  resp.URL,
		"icon_file_path":       iconFilePath,
		"icon_file_web_source": iconFileWebSource,
		"icon_file_base64":     iconFileBase64,
	}

	for key, val := range iconData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
