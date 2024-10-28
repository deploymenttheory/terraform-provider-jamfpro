// common/jamfprivileges/validate.go
// This package is used to validate the privileges fields in the Jamf Pro configuration.
// Whilst there is an api to return the api priviledge set, they are not categorised in
// the same way as the Jamf Pro UI. There it's required to export the priviledge set from
// the Jamf Pro UI and use this to validate the fields. There is a JSON file for each of
// the priviledge types that are used to validate the fields and these are exported using
// a recipe from the scripts folder in the provider. Based upon the jamf pro sdk recipes
// implementation.

// The json files are updated when a new version of Jamf Pro is released and the provider
// via a pipeline, which will update the provider and the JSON files. The JSON files are
// stored in the privileges folder in the jamfprivileges package. The JSON files are named
// jss_objects_privileges.json, jss_settings_privileges.json and jss_actions_privileges.json.
// Broken down by jamf pro version.

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
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"
)

//go:embed privileges/*/*.json
var privilegesFS embed.FS

const (
	LatestVersion  = "11.10.1"
	NMinus1Version = "11.7.1"
	NMinus2Version = "11.7.0"
)

var validVersions = []string{LatestVersion, NMinus1Version, NMinus2Version}

// PrivilegeSupport tracks which versions support a privilege
type PrivilegeSupport struct {
	Name             string
	SupportedVersion string
}

// loadJSONPrivilegesWithVersion loads privileges from JSON files and tracks which version they're from
func loadJSONPrivilegesWithVersion(filename string) (map[string]PrivilegeSupport, error) {
	privilegeMap := make(map[string]PrivilegeSupport)

	for _, version := range validVersions {
		filePath := path.Join("privileges", version, filename)
		data, err := privilegesFS.ReadFile(filePath)
		if err != nil {
			continue
		}

		var privileges []string
		if err := json.Unmarshal(data, &privileges); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON from %s: %w", filePath, err)
		}

		// Add or update privileges with their version information
		for _, priv := range privileges {
			// Only store the oldest version where the privilege appears
			if _, exists := privilegeMap[priv]; !exists {
				privilegeMap[priv] = PrivilegeSupport{
					Name:             priv,
					SupportedVersion: version,
				}
			}
		}
	}

	if len(privilegeMap) == 0 {
		return nil, fmt.Errorf("no valid privileges found for any version")
	}

	return privilegeMap, nil
}

// formatPrivilegesError creates a detailed error message with version information
func formatPrivilegesError(invalidValue, key string, privileges map[string]PrivilegeSupport) error {
	// Sort privileges by version and then alphabetically
	var sortedPrivileges []PrivilegeSupport
	for _, p := range privileges {
		sortedPrivileges = append(sortedPrivileges, p)
	}
	sort.Slice(sortedPrivileges, func(i, j int) bool {
		if sortedPrivileges[i].SupportedVersion != sortedPrivileges[j].SupportedVersion {
			return sortedPrivileges[i].SupportedVersion > sortedPrivileges[j].SupportedVersion
		}
		return sortedPrivileges[i].Name < sortedPrivileges[j].Name
	})

	// Build the error message
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Invalid value '%s' for %s.\n\n", invalidValue, key))
	sb.WriteString("Valid privileges by version:\n\n")

	currentVersion := ""
	for _, priv := range sortedPrivileges {
		if currentVersion != priv.SupportedVersion {
			currentVersion = priv.SupportedVersion
			sb.WriteString(fmt.Sprintf("\nJamf Pro %s:\n", currentVersion))
		}
		sb.WriteString(fmt.Sprintf("- %s\n", priv.Name))
	}

	return fmt.Errorf(sb.String())
}

// ValidateJSSObjectsPrivileges ensures that each privilege in the list is valid.
func ValidateJSSObjectsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	privileges, err := loadJSONPrivilegesWithVersion("jss_objects_privileges.json")
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load JSS Object Privileges: %v", err))
		return warns, errs
	}

	if privilege, exists := privileges[v]; exists {
		if privilege.SupportedVersion != LatestVersion {
			warns = append(warns, fmt.Sprintf("Privilege '%s' is supported in Jamf Pro %s but may not be available in newer versions",
				v, privilege.SupportedVersion))
		}
		return warns, errs
	}

	errs = append(errs, formatPrivilegesError(v, key, privileges))
	return warns, errs
}

// ValidateJSSSettingsPrivileges checks if each value in the jss_settings_privileges field
func ValidateJSSSettingsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	privileges, err := loadJSONPrivilegesWithVersion("jss_settings_privileges.json")
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load JSS Settings Privileges: %w", err))
		return warns, errs
	}

	if privilege, exists := privileges[v]; exists {
		if privilege.SupportedVersion != LatestVersion {
			warns = append(warns, fmt.Sprintf("Privilege '%s' is supported in Jamf Pro %s but may not be available in newer versions",
				v, privilege.SupportedVersion))
		}
		return warns, errs
	}

	errs = append(errs, formatPrivilegesError(v, key, privileges))
	return warns, errs
}

// ValidateJSSActionsPrivileges checks if each value in the jss_actions_privileges field
func ValidateJSSActionsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	privileges, err := loadJSONPrivilegesWithVersion("jss_actions_privileges.json")
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to load JSS Actions Privileges: %w", err))
		return warns, errs
	}

	if privilege, exists := privileges[v]; exists {
		if privilege.SupportedVersion != LatestVersion {
			warns = append(warns, fmt.Sprintf("Privilege '%s' is supported in Jamf Pro %s but may not be available in newer versions",
				v, privilege.SupportedVersion))
		}
		return warns, errs
	}

	errs = append(errs, formatPrivilegesError(v, key, privileges))
	return warns, errs
}
