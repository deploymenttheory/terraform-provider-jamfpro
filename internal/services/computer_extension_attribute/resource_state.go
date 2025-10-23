package computer_extension_attribute

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Computer Extension Attribute information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceComputerExtensionAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", resp.ID); err != nil {
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
	if err := d.Set("enabled", resp.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inventory_display_type", resp.InventoryDisplayType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("input_type", resp.InputType); err != nil {
		return diag.FromErr(err)
	}

	// Handle input type specific fields
	switch resp.InputType {
	case "SCRIPT":
		normalizedScript := normalizeScript(resp.ScriptContents)
		if err := d.Set("script_contents", normalizedScript); err != nil {
			return diag.FromErr(err)
		}
	case "POPUP":
		if err := d.Set("popup_menu_choices", resp.PopupMenuChoices); err != nil {
			return diag.FromErr(err)
		}
	case "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING":
		if err := d.Set("ldap_attribute_mapping", resp.LDAPAttributeMapping); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("ldap_extension_attribute_allowed", resp.LDAPExtensionAttributeAllowed != nil && *resp.LDAPExtensionAttributeAllowed); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
