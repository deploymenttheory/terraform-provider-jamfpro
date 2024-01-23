// accounts_data_handling.go
package accounts

import (
	"context"
	"fmt"

	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customDiffAccounts is a custom diff function for the Jamf Pro Account resource.
// This function is used during the Terraform plan phase to apply custom validation rules
// that are not covered by the basic schema validation.
func customDiffAccounts(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Declare variables outside the validation scenarios
	var casperAdminPrivileges interface{}
	var ok bool

	// scenario 01
	accessLevel, ok := d.GetOk("access_level")
	if !ok || accessLevel == nil {
		// If access_level is not set, no further checks required
		return nil
	}

	// Enforce that the 'site' attribute must be set if access_level is 'Site Access'
	if accessLevel.(string) == "Site Access" {
		if _, ok := d.GetOk("site"); !ok {
			return fmt.Errorf("'site' must be set when 'access_level' is 'Site Access'")
		}
	}

	// scenario 02 validation when "Use Casper Admin" is present in casper_admin_privileges
	casperAdminPrivileges, ok = d.GetOk("casper_admin_privileges")
	if ok {
		for _, privilege := range casperAdminPrivileges.([]interface{}) {
			if privilege.(string) == "Use Casper Admin" {
				requiredReadPrivileges := []string{
					"Read Categories",
					"Read Directory Bindings",
					"Read Dock Items",
					"Read Packages",
					"Read Printers",
					"Read Scripts",
				}

				jssObjectsPrivileges, ok := d.GetOk("jss_objects_privileges")
				if !ok {
					return fmt.Errorf("when 'Use Casper Admin' is selected, the following JSS Object Privileges are required: %v", requiredReadPrivileges)
				}

				jssPrivilegesSet := make(map[string]bool)
				for _, priv := range jssObjectsPrivileges.([]interface{}) {
					jssPrivilegesSet[priv.(string)] = true
				}

				for _, requiredPriv := range requiredReadPrivileges {
					if !jssPrivilegesSet[requiredPriv] {
						return fmt.Errorf("missing required privilege '%s' in 'jss_objects_privileges' when 'Use Casper Admin' is selected in 'casper_admin_privileges'", requiredPriv)
					}
				}
			}
		}
	}

	// Scenario 03: Validation when "Save With Casper Admin" is present in casper_admin_privileges
	casperAdminPrivileges, ok = d.GetOk("casper_admin_privileges")
	if ok {
		for _, privilege := range casperAdminPrivileges.([]interface{}) {
			if privilege.(string) == "Save With Casper Admin" {
				requiredCrudPrivileges := []string{
					"Create Categories", "Read Categories", "Update Categories", "Delete Categories",
					"Create Directory Bindings", "Read Directory Bindings", "Update Directory Bindings", "Delete Directory Bindings",
					"Create Dock Items", "Read Dock Items", "Update Dock Items", "Delete Dock Items",
					"Create Packages", "Read Packages", "Update Packages", "Delete Packages",
					"Create Printers", "Read Printers", "Update Printers", "Delete Printers",
					"Create Scripts", "Read Scripts", "Update Scripts", "Delete Scripts",
				}

				jssObjectsPrivileges, ok := d.GetOk("jss_objects_privileges")
				if !ok {
					return fmt.Errorf("when 'Save With Casper Admin' is selected, the following JSS Object Privileges are required: %v", requiredCrudPrivileges)
				}

				jssPrivilegesSet := make(map[string]bool)
				for _, priv := range jssObjectsPrivileges.([]interface{}) {
					jssPrivilegesSet[priv.(string)] = true
				}

				for _, requiredPriv := range requiredCrudPrivileges {
					if !jssPrivilegesSet[requiredPriv] {
						return fmt.Errorf("missing required privilege '%s' in 'jss_objects_privileges' when 'Save With Casper Admin' is selected in 'casper_admin_privileges'", requiredPriv)
					}
				}
			}
		}
	}

	// Additional scenarios for reverse dependency checks for scenarios  3 and 4
	jssObjectsPrivileges, ok := d.GetOk("jss_objects_privileges")
	if ok {
		jssPrivilegesSet := make(map[string]bool)
		for _, priv := range jssObjectsPrivileges.([]interface{}) {
			jssPrivilegesSet[priv.(string)] = true
		}

		// Reverse Dependency Check 1: If specific read privileges are set, ensure "Use Casper Admin" is present
		allReadPrivilegesPresent := true
		readPrivileges := []string{"Read Categories", "Read Directory Bindings", "Read Dock Items", "Read Packages", "Read Printers", "Read Scripts"}
		for _, priv := range readPrivileges {
			if !jssPrivilegesSet[priv] {
				allReadPrivilegesPresent = false
				break
			}
		}
		if allReadPrivilegesPresent {
			casperAdminPrivileges, ok = d.GetOk("casper_admin_privileges")
			if !ok || !util.GetStringFromSlice(casperAdminPrivileges.([]interface{}), "Use Casper Admin") {
				return fmt.Errorf("when the following JSS Object Privileges are set %v, 'Use Casper Admin' must be included in 'casper_admin_privileges'", readPrivileges)
			}
		}

		// Reverse Dependency Check 2:
		//If specific CRUD privileges are set, ensure "Save With Casper Admin" is present
		allCrudPrivilegesPresent := true
		crudPrivileges := []string{
			"Create Categories", "Read Categories", "Update Categories", "Delete Categories",
			"Create Directory Bindings", "Read Directory Bindings", "Update Directory Bindings", "Delete Directory Bindings",
			"Create Dock Items", "Read Dock Items", "Update Dock Items", "Delete Dock Items",
			"Create Packages", "Read Packages", "Update Packages", "Delete Packages",
			"Create Printers", "Read Printers", "Update Printers", "Delete Printers",
			"Create Scripts", "Read Scripts", "Update Scripts", "Delete Scripts",
		}
		for _, priv := range crudPrivileges {
			if !jssPrivilegesSet[priv] {
				allCrudPrivilegesPresent = false
				break
			}
		}
		if allCrudPrivilegesPresent {
			casperAdminPrivileges, ok = d.GetOk("casper_admin_privileges")
			if !ok || !util.GetStringFromSlice(casperAdminPrivileges.([]interface{}), "Save With Casper Admin") {
				return fmt.Errorf("when the following JSS Object Privileges are set %v, 'Save With Casper Admin' must be included in 'casper_admin_privileges'", crudPrivileges)
			}
		}
	}

	return nil
}
