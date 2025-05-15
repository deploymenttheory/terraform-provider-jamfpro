package computerinventorycollectionsettings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// pathType defines the mapping between schema keys and API scopes
type pathType struct {
	key   string
	scope string
}

// pathTypes defines all available custom path types
var pathTypes = []pathType{
	{"application_paths", "APP"},
	{"font_paths", "FONT"},
	{"plugin_paths", "PLUGIN"},
}

// construct builds the request for UpdateComputerInventoryCollectionSettings
func construct(d *schema.ResourceData) (*jamfpro.ResourceComputerInventoryCollectionSettings, error) {
	resource := &jamfpro.ResourceComputerInventoryCollectionSettings{}

	if v, ok := d.GetOk("computer_inventory_collection_preferences"); ok {
		prefsList := v.([]interface{})
		if len(prefsList) > 0 {
			prefsMap := prefsList[0].(map[string]interface{})

			resource.ComputerInventoryCollectionPreferences = jamfpro.ComputerInventoryCollectionSettingsSubsetPreferences{
				MonitorApplicationUsage:       prefsMap["monitor_application_usage"].(bool),
				IncludeFonts:                  prefsMap["include_fonts"].(bool),
				IncludePlugins:                prefsMap["include_plugins"].(bool),
				IncludePackages:               prefsMap["include_packages"].(bool),
				IncludeSoftwareUpdates:        prefsMap["include_software_updates"].(bool),
				IncludeSoftwareId:             prefsMap["include_software_id"].(bool),
				IncludeAccounts:               prefsMap["include_accounts"].(bool),
				CalculateSizes:                prefsMap["calculate_sizes"].(bool),
				IncludeHiddenAccounts:         prefsMap["include_hidden_accounts"].(bool),
				IncludePrinters:               prefsMap["include_printers"].(bool),
				IncludeServices:               prefsMap["include_services"].(bool),
				CollectSyncedMobileDeviceInfo: prefsMap["collect_synced_mobile_device_info"].(bool),
				UpdateLdapInfoOnComputerInventorySubmissions: prefsMap["update_ldap_info_on_computer_inventory_submissions"].(bool),
				MonitorBeacons:               prefsMap["monitor_beacons"].(bool),
				AllowChangingUserAndLocation: prefsMap["allow_changing_user_and_location"].(bool),
				UseUnixUserPaths:             prefsMap["use_unix_user_paths"].(bool),
				CollectUnmanagedCertificates: prefsMap["collect_unmanaged_certificates"].(bool),
			}
		}
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return nil, fmt.Errorf("failed to marshal Computer Inventory Collection Settings to JSON: %v", err)
	}
	log.Printf("[DEBUG] Constructed Computer Inventory Collection Settings resource:\n%s", string(resourceJSON))

	return resource, nil
}

// constructCustomPaths builds individual path requests for create requests
func constructCustomPaths(d *schema.ResourceData) ([]jamfpro.ResourceComputerInventoryCollectionSettingsCustomPath, error) {
	var paths []jamfpro.ResourceComputerInventoryCollectionSettingsCustomPath

	for _, pt := range pathTypes {
		if v, ok := d.GetOk(pt.key); ok {
			pathSet := v.(*schema.Set)
			for _, p := range pathSet.List() {
				pathMap := p.(map[string]interface{})
				path := pathMap["path"].(string)

				customPath := jamfpro.ResourceComputerInventoryCollectionSettingsCustomPath{
					Scope: pt.scope,
					Path:  path,
				}

				pathJSON, err := json.MarshalIndent(customPath, "", "  ")
				if err != nil {
					//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
					return nil, fmt.Errorf("failed to marshal custom path to JSON: %v", err)
				}
				log.Printf("[DEBUG] Constructing API request for custom path:\n%s", string(pathJSON))

				paths = append(paths, customPath)
			}
		}
	}

	return paths, nil
}

// constructPathUpdates determines which paths need to be added or removed
func constructPathUpdates(d *schema.ResourceData) (pathsToAdd []jamfpro.ResourceComputerInventoryCollectionSettingsCustomPath, pathIDsToRemove []string) {
	for _, pt := range pathTypes {
		if d.HasChange(pt.key) {
			old, new := d.GetChange(pt.key)
			oldSet := old.(*schema.Set)
			newSet := new.(*schema.Set)

			// Find paths to delete (in old but not in new)
			for _, oldPath := range oldSet.List() {
				oldPathMap := oldPath.(map[string]interface{})
				oldPathStr := oldPathMap["path"].(string)
				oldPathID := oldPathMap["id"].(string)

				log.Printf("[DEBUG] Checking old path for removal:\n  Path: %s\n  ID: %s", oldPathStr, oldPathID)

				if !containsPath(newSet.List(), oldPathStr) {
					log.Printf("[DEBUG] Path marked for removal:\n  ID: %s\n  Path: %s", oldPathID, oldPathStr)
					pathIDsToRemove = append(pathIDsToRemove, oldPathID)
				}
			}

			// Find paths to add (in new but not in old)
			for _, newPath := range newSet.List() {
				newPathMap := newPath.(map[string]interface{})
				newPathStr := newPathMap["path"].(string)

				log.Printf("[DEBUG] Checking new path for addition:\n  Path: %s", newPathStr)

				if !containsPath(oldSet.List(), newPathStr) {
					customPath := jamfpro.ResourceComputerInventoryCollectionSettingsCustomPath{
						Scope: pt.scope,
						Path:  newPathStr,
					}

					pathJSON, err := json.MarshalIndent(customPath, "", "  ")
					if err == nil {
						log.Printf("[DEBUG] Path marked for addition:\n%s", string(pathJSON))
					}

					pathsToAdd = append(pathsToAdd, customPath)
				}
			}
		}
	}

	if len(pathsToAdd) > 0 || len(pathIDsToRemove) > 0 {
		log.Printf("[DEBUG] Total paths to add: %d, Total IDs to remove: %d",
			len(pathsToAdd), len(pathIDsToRemove))
	}

	return pathsToAdd, pathIDsToRemove
}

// Helper function to check if a path exists in a list of path maps
func containsPath(list []interface{}, searchPath string) bool {
	for _, item := range list {
		pathMap := item.(map[string]interface{})
		if pathMap["path"].(string) == searchPath {
			return true
		}
	}
	return false
}
