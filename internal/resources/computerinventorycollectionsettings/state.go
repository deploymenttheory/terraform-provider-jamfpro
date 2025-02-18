package computerinventorycollectionsettings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Computer Inventory Collection information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceComputerInventoryCollectionSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	// Ensure the base resource is correctly set as the singleton
	d.SetId("jamfpro_computer_inventory_collection_settings_singleton")

	// Build the preferences map
	preferences := []interface{}{
		map[string]interface{}{
			"monitor_application_usage":         resp.ComputerInventoryCollectionPreferences.MonitorApplicationUsage,
			"include_fonts":                     resp.ComputerInventoryCollectionPreferences.IncludeFonts,
			"include_plugins":                   resp.ComputerInventoryCollectionPreferences.IncludePlugins,
			"include_packages":                  resp.ComputerInventoryCollectionPreferences.IncludePackages,
			"include_software_updates":          resp.ComputerInventoryCollectionPreferences.IncludeSoftwareUpdates,
			"include_software_id":               resp.ComputerInventoryCollectionPreferences.IncludeSoftwareId,
			"include_accounts":                  resp.ComputerInventoryCollectionPreferences.IncludeAccounts,
			"calculate_sizes":                   resp.ComputerInventoryCollectionPreferences.CalculateSizes,
			"include_hidden_accounts":           resp.ComputerInventoryCollectionPreferences.IncludeHiddenAccounts,
			"include_printers":                  resp.ComputerInventoryCollectionPreferences.IncludePrinters,
			"include_services":                  resp.ComputerInventoryCollectionPreferences.IncludeServices,
			"collect_synced_mobile_device_info": resp.ComputerInventoryCollectionPreferences.CollectSyncedMobileDeviceInfo,
			"update_ldap_info_on_computer_inventory_submissions": resp.ComputerInventoryCollectionPreferences.UpdateLdapInfoOnComputerInventorySubmissions,
			"monitor_beacons":                  resp.ComputerInventoryCollectionPreferences.MonitorBeacons,
			"allow_changing_user_and_location": resp.ComputerInventoryCollectionPreferences.AllowChangingUserAndLocation,
			"use_unix_user_paths":              resp.ComputerInventoryCollectionPreferences.UseUnixUserPaths,
			"collect_unmanaged_certificates":   resp.ComputerInventoryCollectionPreferences.CollectUnmanagedCertificates,
		},
	}

	if err := d.Set("computer_inventory_collection_preferences", preferences); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("application_paths", flattenPaths(resp.ApplicationPaths)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("font_paths", flattenPaths(resp.FontPaths)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("plugin_paths", flattenPaths(resp.PluginPaths)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// flattenPaths flattens path items for setting in Terraform state
func flattenPaths(paths []jamfpro.ComputerInventoryCollectionSettingsSubsetPathItem) *schema.Set {
	pathSet := schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}), []interface{}{})

	for _, path := range paths {
		pathMap := map[string]interface{}{
			"path": path.Path,
			"id":   path.ID,
		}
		pathSet.Add(pathMap)
	}
	return pathSet
}
