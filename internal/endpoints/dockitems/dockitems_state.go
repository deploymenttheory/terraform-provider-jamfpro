// dockitems_state.go
package dockitems

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Dock Item information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceDockItem) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check if dockItem data exists and update the Terraform state
	if resource != nil {
		resourceData := map[string]interface{}{
			"id":       strconv.Itoa(resource.ID),
			"name":     resource.Name,
			"type":     resource.Type,
			"path":     resource.Path,
			"contents": resource.Contents,
		}

		// Set each attribute in the Terraform state, checking for errors
		var diags diag.Diagnostics
		for key, val := range resourceData {
			if err := d.Set(key, val); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
		return diags
	}
	return diags
}
