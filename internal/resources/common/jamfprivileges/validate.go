// common/jamfprivileges/validate.go
// This package is used to validate the privileges fields in the Jamf Pro configuration.
// Whilst there is an api to return the api priviledge set, they are not categorised in
// the same way as the Jamf Pro UI. There it's required to export the priviledge set from
// the Jamf Pro UI and use this to validate the fields. There is a JSON file for each of
// the priviledge types that are used to validate the fields and these are exported using
// a recipe from the Jamf Pro SDK,
// 'recipes/privileges/Export_All_Available_Privileges_By_Type_To_Json'.

// To update the JSON files, export the priviledge using this recipe and then copy the
// JSON file to the appropriate version folder in the privileges folder within this provider.
// The version folders are named after the Jamf Pro version, e.g. 11.7.1, 11.7.0, 11.6.0 etc.
// The intention is to keep only the latest 3 versions of the JSON files in the provider. So
// when a new version is added, the oldest version should be removed.

// The folder name for the privileges folder is 'privileges' and this should be in the same
// directory as the validate.go file. Beneath this folder are the version folders and the JSON
// files for the priviledge types.

// jamfprivileges
// ├── validate.go
// ├── 11.7.1
// │   ├── jss_actions_privileges.json
// │   ├── jss_objects_privileges.json
// │   └── jss_settings_privileges.json
// ├── 11.7.0
// │   ├── jss_actions_privileges.json
// │   ├── jss_objects_privileges.json
// │   └── jss_settings_privileges.json

package jamfprivileges

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	LatestVersion  = "11.10.1"
	NMinus1Version = "11.7.1"
	NMinus2Version = "11.7.0"
)

var validVersions = []string{LatestVersion, NMinus1Version, NMinus2Version}

// ValidateJSSObjectsPrivileges ensures that each privilege in the list is valid.
func ValidateJSSObjectsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	validPrivileges, err := loadJSONPrivileges("jss_objects_privileges.json")
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load JSS Object Privileges: %v", err))
		return warns, errs
	}

	for _, validPrivilege := range validPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	errs = append(errs, fmt.Errorf("invalid value '%s' for %s: not a recognized JSS Object Privilege", v, key))
	return warns, errs
}

// ValidateJSSSettingsPrivileges checks if each value in the jss_settings_privileges field
// matches a value from a predefined list of valid JSS Setting Privileges.
func ValidateJSSSettingsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	validPrivileges, err := loadJSONPrivileges("jss_settings_privileges.json")
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load JSS Settings Privileges: %w", err))
		return warns, errs
	}

	for _, validPrivilege := range validPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	errs = append(errs, fmt.Errorf("invalid value '%s' for %s: not a recognized JSS Setting Privilege", v, key))
	return warns, errs
}

// ValidateJSSActionsPrivileges checks if each value in the jss_actions_privileges field
// matches a value from a predefined list of valid JSS Actions Privileges.
func ValidateJSSActionsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	validPrivileges, err := loadJSONPrivileges("jss_actions_privileges.json")
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load JSS Actions Privileges: %w", err))
		return warns, errs
	}

	for _, validPrivilege := range validPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	errs = append(errs, fmt.Errorf("invalid value '%s' for %s: not a recognized JSS Action Privilege", v, key))
	return warns, errs
}

// loadJSONPrivileges loads privileges from a JSON file in the appropriate version folder
func loadJSONPrivileges(filename string) ([]string, error) {
	var privileges []string

	privilegesFolder, err := findPrivilegesFolder()
	if err != nil {
		return nil, fmt.Errorf("error finding privileges folder: %w", err)
	}

	for _, version := range validVersions {
		filePath := filepath.Join(privilegesFolder, version, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Try the next version
			}
			return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
		}

		err = json.Unmarshal(data, &privileges)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON from %s: %w", filePath, err)
		}

		return privileges, nil
	}

	return nil, fmt.Errorf("no valid privileges file found for any version in %s", privilegesFolder)
}

// findPrivilegesFolder searches for the "privileges" folder starting from the current file's directory
// and moving up the directory tree.
func findPrivilegesFolder() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get the current file path")
	}

	dir := filepath.Dir(filename)
	for {
		privilegesPath := filepath.Join(dir, "privileges")
		if _, err := os.Stat(privilegesPath); err == nil {
			return privilegesPath, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("privileges folder not found")
}
