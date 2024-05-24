// computerextensionattributes_state.go
package computerextensionattributes

import (
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Extension Attribute information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceComputerExtensionAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	resourceData := map[string]interface{}{
		"name":              resource.Name,
		"enabled":           resource.Enabled,
		"description":       resource.Description,
		"data_type":         strings.ToLower(resource.DataType),
		"inventory_display": resource.InventoryDisplay,
		"recon_display":     resource.ReconDisplay,
		"input_type": []interface{}{
			map[string]interface{}{
				"type":     resource.InputType.Type,
				"platform": resource.InputType.Platform,
				"script":   resource.InputType.Script,
				"choices":  resource.InputType.Choices,
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
