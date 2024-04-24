// common/jamfprivileges/validate.go
// This package contains shared / common resource functions
package jamfprivileges

import "fmt"

// ValidateJSSObjectsPrivileges ensures that each privilege in the list is valid.
func ValidateJSSObjectsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	// Check if the value is in the list of valid privileges
	for _, validPrivilege := range ValidJSSObjectsPrivileges {
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

	// Check if the value is in the list of valid JSS Settings Privileges
	for _, validPrivilege := range ValidJSSSettingsPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	// Add an error if the value is not found in the valid privileges list
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

	// Check if the value is in the list of valid JSS Actions Privileges
	for _, validPrivilege := range ValidJSSActionsPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	// Add an error if the value is not found in the valid privileges list
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

	// Check if the value is in the list of valid Casper Admin Privileges
	for _, validPrivilege := range ValidCasperAdminPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	// Add an error if the value is not found in the valid privileges list
	errs = append(errs, fmt.Errorf("invalid value '%s' for %s: not a recognized Casper Admin Privilege", v, key))
	return warns, errs
}
