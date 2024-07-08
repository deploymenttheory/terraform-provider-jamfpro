// computerinventorycollection_state.go
package computerinventorycollection

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Computer Inventory Collection information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceComputerInventoryCollection) diag.Diagnostics {
	var diags diag.Diagnostics

	inventoryCollectionData := map[string]interface{}{
		"local_user_accounts":               resp.LocalUserAccounts,
		"home_directory_sizes":              resp.HomeDirectorySizes,
		"hidden_accounts":                   resp.HiddenAccounts,
		"printers":                          resp.Printers,
		"active_services":                   resp.ActiveServices,
		"mobile_device_app_purchasing_info": resp.MobileDeviceAppPurchasingInfo,
		"computer_location_information":     resp.ComputerLocationInformation,
		"package_receipts":                  resp.PackageReceipts,
		"available_software_updates":        resp.AvailableSoftwareUpdates,
		"include_applications":              resp.InclueApplications,
		"include_fonts":                     resp.InclueFonts,
		"include_plugins":                   resp.IncluePlugins,
	}

	for key, val := range inventoryCollectionData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if err := d.Set("applications", flattenApplications(resp.Applications)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("fonts", flattenFonts(resp.Fonts)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("plugins", flattenPlugins(resp.Plugins)); err != nil {
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
