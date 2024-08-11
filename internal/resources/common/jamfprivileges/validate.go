// common/jamfprivileges/validate.go
// This package contains shared / common resource functions
package jamfprivileges

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	LatestVersion  = "11.7.1"
	NMinus1Version = "11.7.0"
	NMinus2Version = "11.6.0"
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

// ValidateCasperAdminPrivileges checks if each value in the casper_admin_privileges field
// matches a value from a predefined list of valid Casper Admin Privileges.
func ValidateCasperAdminPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	// Note: As there's no JSON file specified for Casper Admin Privileges,
	// we're keeping the original logic here. If you want to add a JSON file
	// for this as well, you can modify this function similar to the others.
	for _, validPrivilege := range ValidCasperAdminPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	errs = append(errs, fmt.Errorf("invalid value '%s' for %s: not a recognized Casper Admin Privilege", v, key))
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
