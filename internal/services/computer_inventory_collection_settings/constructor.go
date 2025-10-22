package computer_inventory_collection_settings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errMarshalSettings   = "failed to marshal Computer Inventory Collection Settings to JSON"
	errMarshalCustomPath = "failed to marshal custom path to JSON"
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
		prefsList := v.([]any)
		if len(prefsList) > 0 {
			prefsMap := prefsList[0].(map[string]any)

			resource.ComputerInventoryCollectionPreferences = jamfpro.ComputerInventoryCollectionSettingsSubsetPreferences{
				MonitorApplicationUsage:                      prefsMap["monitor_application_usage"].(bool),
				IncludePackages:                              prefsMap["include_packages"].(bool),
				IncludeSoftwareUpdates:                       prefsMap["include_software_updates"].(bool),
				IncludeSoftwareId:                            prefsMap["include_software_id"].(bool),
				IncludeAccounts:                              prefsMap["include_accounts"].(bool),
				CalculateSizes:                               prefsMap["calculate_sizes"].(bool),
				IncludeHiddenAccounts:                        prefsMap["include_hidden_accounts"].(bool),
				IncludePrinters:                              prefsMap["include_printers"].(bool),
				IncludeServices:                              prefsMap["include_services"].(bool),
				CollectSyncedMobileDeviceInfo:                prefsMap["collect_synced_mobile_device_info"].(bool),
				UpdateLdapInfoOnComputerInventorySubmissions: prefsMap["update_ldap_info_on_computer_inventory_submissions"].(bool),
				MonitorBeacons:                               prefsMap["monitor_beacons"].(bool),
				AllowChangingUserAndLocation:                 prefsMap["allow_changing_user_and_location"].(bool),
				UseUnixUserPaths:                             prefsMap["use_unix_user_paths"].(bool),
				CollectUnmanagedCertificates:                 prefsMap["collect_unmanaged_certificates"].(bool),
			}
		}
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMarshalSettings, err)
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
				pathMap := p.(map[string]any)
				path := pathMap["path"].(string)

				customPath := jamfpro.ResourceComputerInventoryCollectionSettingsCustomPath{
					Scope: pt.scope,
					Path:  path,
				}

				pathJSON, err := json.MarshalIndent(customPath, "", "  ")
				if err != nil {
					return nil, fmt.Errorf("%s: %w", errMarshalCustomPath, err)
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

			for _, oldPath := range oldSet.List() {
				oldPathMap := oldPath.(map[string]any)
				oldPathStr := oldPathMap["path"].(string)
				oldPathID := oldPathMap["id"].(string)

				log.Printf("[DEBUG] Checking old path for removal:\n  Path: %s\n  ID: %s", oldPathStr, oldPathID)

				if !containsPath(newSet.List(), oldPathStr) {
					if oldPathID == "-1" {
						log.Printf("[DEBUG] Skipping built-in path for removal (ID: %s, Path: %s)", oldPathID, oldPathStr)
						continue
					}
					log.Printf("[DEBUG] Path marked for removal:\n  ID: %s\n  Path: %s", oldPathID, oldPathStr)
					pathIDsToRemove = append(pathIDsToRemove, oldPathID)
				}
			}

			for _, newPath := range newSet.List() {
				newPathMap := newPath.(map[string]any)
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
func containsPath(list []any, searchPath string) bool {
	for _, item := range list {
		pathMap := item.(map[string]any)
		if pathMap["path"].(string) == searchPath {
			return true
		}
	}
	return false
}
