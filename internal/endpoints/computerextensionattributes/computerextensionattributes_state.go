// computerextensionattributes_state.go
package computerextensionattributes

import (
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Extension Attribute information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceComputerExtensionAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO review this logic ASAP.
	// Update the Terraform state with the fetched data
	resourceData := map[string]interface{}{
		"name":              resp.Name,
		"enabled":           resp.Enabled,
		"description":       resp.Description,
		"data_type":         strings.ToLower(resp.DataType),
		"inventory_display": resp.InventoryDisplay,
		"recon_display":     resp.ReconDisplay,
		"input_type": []interface{}{
			map[string]interface{}{
				"type":     resp.InputType.Type,
				"platform": resp.InputType.Platform,
				"script":   resp.InputType.Script,
				"choices":  resp.InputType.Choices,
			},
		},
	}

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags

}
