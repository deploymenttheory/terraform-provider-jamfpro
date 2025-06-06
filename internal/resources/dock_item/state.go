// dockitems_state.go
package dock_item

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Dock Item information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceDockItem) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"id":       strconv.Itoa(resp.ID),
		"name":     resp.Name,
		"type":     resp.Type,
		"path":     resp.Path,
		"contents": resp.Contents,
	}

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
