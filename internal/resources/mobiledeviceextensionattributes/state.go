package mobiledeviceextensionattributes

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Mobile Device Extension Attribute information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMobileExtensionAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", strconv.Itoa(resp.ID)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_type", resp.DataType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inventory_display", resp.InventoryDisplay); err != nil {
		return diag.FromErr(err)
	}

	// Handle the nested input_type structure
	inputType := []map[string]interface{}{
		{
			"type":          resp.InputType.Type,
			"popup_choices": resp.InputType.PopupChoices.Choice,
		},
	}
	if err := d.Set("input_type", inputType); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
