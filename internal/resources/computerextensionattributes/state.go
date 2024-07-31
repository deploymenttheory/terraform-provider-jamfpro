// computerextensionattributes_state.go
package computerextensionattributes

import (
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Computer Extension Attribute information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceComputerExtensionAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	d.Set("name", resp.Name)
	d.Set("enabled", resp.Enabled)
	d.Set("description", resp.Description)
	d.Set("data_type", strings.ToLower(resp.DataType))
	d.Set("inventory_display", resp.InventoryDisplay)
	d.Set("recon_display", resp.ReconDisplay)
	d.Set("input_type", resp.InputType)
	d.Set("input_popup", resp.InputType.Choices)
	d.Set("input_script", resp.InputType.Script)

	return diags

}
