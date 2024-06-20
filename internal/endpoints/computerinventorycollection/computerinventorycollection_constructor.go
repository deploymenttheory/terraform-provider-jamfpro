// computerinventorycollection_object.go
package computerinventorycollection

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProComputerInventoryCollection constructs a ResourceComputerInventoryCollection object from the provided schema data and logs its XML representation.
func constructJamfProComputerInventoryCollection(d *schema.ResourceData) (*jamfpro.ResourceComputerInventoryCollection, error) {

	inventoryCollection := &jamfpro.ResourceComputerInventoryCollection{
		LocalUserAccounts:             d.Get("local_user_accounts").(bool),
		HomeDirectorySizes:            d.Get("home_directory_sizes").(bool),
		HiddenAccounts:                d.Get("hidden_accounts").(bool),
		Printers:                      d.Get("printers").(bool),
		ActiveServices:                d.Get("active_services").(bool),
		MobileDeviceAppPurchasingInfo: d.Get("mobile_device_app_purchasing_info").(bool),
		ComputerLocationInformation:   d.Get("computer_location_information").(bool),
		PackageReceipts:               d.Get("package_receipts").(bool),
		AvailableSoftwareUpdates:      d.Get("available_software_updates").(bool),
		InclueApplications:            d.Get("include_applications").(bool),
		InclueFonts:                   d.Get("include_fonts").(bool),
		IncluePlugins:                 d.Get("include_plugins").(bool),
	}

	// Process applications
	if v, ok := d.GetOk("applications"); ok {
		applications := v.([]interface{})
		for _, application := range applications {
			appMap := application.(map[string]interface{})
			inventoryCollection.Applications = append(inventoryCollection.Applications, jamfpro.ApplicationEntry{
				Application: jamfpro.Application{
					Path:     appMap["path"].(string),
					Platform: appMap["platform"].(string),
				},
			})
		}
	}

	// Process fonts
	if v, ok := d.GetOk("fonts"); ok {
		fonts := v.([]interface{})
		for _, font := range fonts {
			fontMap := font.(map[string]interface{})
			inventoryCollection.Fonts = append(inventoryCollection.Fonts, jamfpro.FontEntry{
				Font: jamfpro.Font{
					Path:     fontMap["path"].(string),
					Platform: fontMap["platform"].(string),
				},
			})
		}
	}

	// Process plugins
	if v, ok := d.GetOk("plugins"); ok {
		plugins := v.([]interface{})
		for _, plugin := range plugins {
			pluginMap := plugin.(map[string]interface{})
			inventoryCollection.Plugins = append(inventoryCollection.Plugins, jamfpro.PluginEntry{
				Plugin: jamfpro.Plugin{
					Path:     pluginMap["path"].(string),
					Platform: pluginMap["platform"].(string),
				},
			})
		}
	}

	// Serialize and pretty-print the inventory collection object as XML for logging
	resourceXML, err := xml.MarshalIndent(inventoryCollection, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Inventory Collection to XML: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Inventory Collection XML:\n%s\n", string(resourceXML))

	return inventoryCollection, nil
}
