// common/jamfprivileges/validate.go
// This package is used to validate the privileges fields in the Jamf Pro user accounts
// configuration.

// Whilst there is an api to return a complete api role priviledge set, they are not returned
// in a useful categorised format that's required for comparision for user accounts. Therefore
// there is a seperate maintainence pipeline that exports the priviledge
// sets to json files for user accounts.

// The maintainence pipelines runs once a week and will update the priviledges JSON files
// when there is a new version of Jamf Pro is found. The JSON files are named
// jss_objects_privileges.json, jss_settings_privileges.json and jss_actions_privileges.json
// and the provider maintains 3 jamf pro version permission sets.

// jamfprivileges
// ├── validate.go
// ├── 11.7.1
// │   ├── jss_actions_privileges.json
// │   ├── jss_objects_privileges.json
// │   └── jss_settings_privileges.json
// ├── 11.10.1
// │   ├── jss_actions_privileges.json
// │   ├── jss_objects_privileges.json
// │   └── jss_settings_privileges.json
// ├── 11.10.2
// │   ├── jss_actions_privileges.json
// │   ├── jss_objects_privileges.json
// │   └── jss_settings_privileges.json

package jamfprivileges

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"
)

//go:embed privileges/*/*.json
var privilegesFS embed.FS

const (
	LatestVersion  = "11.11.1"
	NMinus1Version = "11.10.2"
	NMinus2Version = "11.10.1"
)

// PrivilegeSupport tracks which versions support a privilege
type PrivilegeSupport struct {
	Name             string
	SupportedVersion string
}

// loadJSONPrivilegesWithVersion loads privileges from JSON files and tracks which version they're from
func loadJSONPrivilegesWithVersion(filename string) (map[string]PrivilegeSupport, error) {
	privilegeMap := make(map[string]PrivilegeSupport)
	latestPrivileges := make(map[string]bool)
	foundAny := false

	latestFilePath := path.Join("privileges", LatestVersion, filename)
	data, err := privilegesFS.ReadFile(latestFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading latest privileges from %s: %w", latestFilePath, err)
	}

	var privileges []string
	if err := json.Unmarshal(data, &privileges); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON from %s: %w", latestFilePath, err)
	}

	for _, priv := range privileges {
		latestPrivileges[priv] = true
		privilegeMap[priv] = PrivilegeSupport{
			Name:             priv,
			SupportedVersion: LatestVersion,
		}
	}
	foundAny = true

	for _, version := range []string{NMinus1Version, NMinus2Version} {
		filePath := path.Join("privileges", version, filename)
		data, err := privilegesFS.ReadFile(filePath)
		if err != nil {
			continue
		}

		var oldPrivileges []string
		if err := json.Unmarshal(data, &oldPrivileges); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON from %s: %w", filePath, err)
		}

		// Only store privileges that don't exist in the latest version
		for _, priv := range oldPrivileges {
			if !latestPrivileges[priv] {
				privilegeMap[priv] = PrivilegeSupport{
					Name:             priv,
					SupportedVersion: version,
				}
			}
		}
	}

	if !foundAny {
		return nil, fmt.Errorf("no valid privileges found for any version")
	}

	return privilegeMap, nil
}

// formatPrivilegesError creates a detailed error message with version information
func formatPrivilegesError(invalidValue, key string, privileges map[string]PrivilegeSupport) error {
	var b strings.Builder

	b.WriteString("Invalid value '")
	b.WriteString(invalidValue)
	b.WriteString("' for ")
	b.WriteString(key)
	b.WriteString(".\n\n")
	b.WriteString("This provider supports Jamf Pro versions ")
	b.WriteString(LatestVersion)
	b.WriteString(" (Latest), ")
	b.WriteString(NMinus1Version)
	b.WriteString(" (N-1), and ")
	b.WriteString(NMinus2Version)
	b.WriteString(" (N-2)\n\n")

	// Group privileges by version and sort
	versionPrivs := make(map[string][]string)
	for _, priv := range privileges {
		versionPrivs[priv.SupportedVersion] = append(versionPrivs[priv.SupportedVersion], priv.Name)
	}

	for version := range versionPrivs {
		sort.Strings(versionPrivs[version])
	}

	b.WriteString("Current privileges (")
	b.WriteString(LatestVersion)
	b.WriteString("):\n")
	for _, priv := range versionPrivs[LatestVersion] {
		b.WriteString("- ")
		b.WriteString(priv)
		b.WriteString("\n")
	}

	// Check for additional privileges in older versions
	hasOlderPrivs := len(versionPrivs[NMinus1Version]) > 0 || len(versionPrivs[NMinus2Version]) > 0
	if hasOlderPrivs {
		b.WriteString("\nAdditional (deprecated) privileges available in older Jamf Pro versions supported by this provider:\n")

		if len(versionPrivs[NMinus1Version]) > 0 {
			b.WriteString("\nJamf Pro ")
			b.WriteString(NMinus1Version)
			b.WriteString(" (N-1):\n")
			for _, priv := range versionPrivs[NMinus1Version] {
				b.WriteString("- ")
				b.WriteString(priv)
				b.WriteString("\n")
			}
		}

		if len(versionPrivs[NMinus2Version]) > 0 {
			b.WriteString("\nJamf Pro ")
			b.WriteString(NMinus2Version)
			b.WriteString(" (N-2):\n")
			for _, priv := range versionPrivs[NMinus2Version] {
				b.WriteString("- ")
				b.WriteString(priv)
				b.WriteString("\n")
			}
		}

		b.WriteString("\nNote: Privileges from older versions may not be available in the latest Jamf Pro version.\n")
	}

	return errors.New(b.String())
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
