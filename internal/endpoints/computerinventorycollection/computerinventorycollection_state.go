// computerinventorycollection_state.go
package computerinventorycollection

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Inventory Collection information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceComputerInventoryCollection) diag.Diagnostics {

	var diags diag.Diagnostics

	// Map the configuration fields from the API response to a structured map
	inventoryCollectionData := map[string]interface{}{
		"local_user_accounts":               resource.LocalUserAccounts,
		"home_directory_sizes":              resource.HomeDirectorySizes,
		"hidden_accounts":                   resource.HiddenAccounts,
		"printers":                          resource.Printers,
		"active_services":                   resource.ActiveServices,
		"mobile_device_app_purchasing_info": resource.MobileDeviceAppPurchasingInfo,
		"computer_location_information":     resource.ComputerLocationInformation,
		"package_receipts":                  resource.PackageReceipts,
		"available_software_updates":        resource.AvailableSoftwareUpdates,
		"include_applications":              resource.InclueApplications,
		"include_fonts":                     resource.InclueFonts,
		"include_plugins":                   resource.IncluePlugins,
	}

	// Set the structured map in the Terraform state
	for key, val := range inventoryCollectionData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Process applications
	if err := d.Set("applications", flattenApplications(resource.Applications)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Process fonts
	if err := d.Set("fonts", flattenFonts(resource.Fonts)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Process plugins
	if err := d.Set("plugins", flattenPlugins(resource.Plugins)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// flattenApplications flattens the applications list for setting in Terraform state
func flattenApplications(applications []jamfpro.ApplicationEntry) []interface{} {
	var result []interface{}
	for _, app := range applications {
		appMap := map[string]interface{}{
			"path":     app.Application.Path,
			"platform": app.Application.Platform,
		}
		result = append(result, appMap)
	}
	return result
}

// flattenFonts flattens the fonts list for setting in Terraform state
func flattenFonts(fonts []jamfpro.FontEntry) []interface{} {
	var result []interface{}
	for _, font := range fonts {
		fontMap := map[string]interface{}{
			"path":     font.Font.Path,
			"platform": font.Font.Platform,
		}
		result = append(result, fontMap)
	}
	return result
}

// flattenPlugins flattens the plugins list for setting in Terraform state
func flattenPlugins(plugins []jamfpro.PluginEntry) []interface{} {
	var result []interface{}
	for _, plugin := range plugins {
		pluginMap := map[string]interface{}{
			"path":     plugin.Plugin.Path,
			"platform": plugin.Plugin.Platform,
		}
		result = append(result, pluginMap)
	}
	return result
}
