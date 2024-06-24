// accounts_data_handling.go
package accounts

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customDiffAccounts is the top-level custom diff function.
func customDiffAccounts(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := validateAccessLevelSiteRequirement(ctx, d, meta); err != nil {
		return err
	}

	if err := validatePrivilegesForSiteAccess(ctx, d, meta); err != nil {
		return err
	}

	if err := validateGroupAccessPrivilegeSetRequirement(ctx, d, meta); err != nil {
		return err
	}

	if err := validateCasperAdminUsePrivileges(ctx, d, meta); err != nil {
		return err
	}

	if err := validateCasperAdminSavePrivileges(ctx, d, meta); err != nil {
		return err
	}

	if err := validateReverseDependencyChecks(ctx, d, meta); err != nil {
		return err
	}

	return nil
}

// validateAccessLevelSiteRequirement checks that the 'site' attribute is set when access_level is 'Site Access'.
func validateAccessLevelSiteRequirement(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	accessLevel, ok := d.GetOk("access_level")
	if !ok || accessLevel == nil {
		return nil
	}

	if accessLevel.(string) == "Site Access" {
		if site, ok := d.GetOk("site_id"); ok {
			siteList := site.([]interface{})
			if len(siteList) == 0 || siteList[0] == nil {
				return fmt.Errorf("'site' block must be set when 'access_level' is 'Site Access'")
			}

			siteMap := siteList[0].(map[string]interface{})
			if id, ok := siteMap["id"]; !ok || id == 0 {
				return fmt.Errorf("'site.id' must be set when 'access_level' is 'Site Access'")
			}
		} else {
			return fmt.Errorf("'site' block must be set when 'access_level' is 'Site Access'")
		}
	}

	return nil
}

// validatePrivilegesForSiteAccess checks that certain privileges are not set when access_level is 'Site Access'.
func validatePrivilegesForSiteAccess(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	accessLevel, ok := d.GetOk("access_level")
	if !ok || accessLevel.(string) != "Site Access" {
		return nil
	}

	if jssSettingsPrivileges, ok := d.GetOk("jss_settings_privileges"); ok && len(jssSettingsPrivileges.([]interface{})) > 0 {
		return fmt.Errorf("when 'access_level' is 'Site Access', 'jss_settings_privileges' are not allowed")
	}

	if casperAdminPrivileges, ok := d.GetOk("casper_admin_privileges"); ok && len(casperAdminPrivileges.([]interface{})) > 0 {
		return fmt.Errorf("when 'access_level' is 'Site Access', 'casper_admin_privileges' are not allowed")
	}

	return nil
}

// validateGroupAccessPrivilegeSetRequirement ensures that if access_level is "Group Access", then privilege_set must be "Custom".
func validateGroupAccessPrivilegeSetRequirement(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	accessLevel, accessLevelOk := d.GetOk("access_level")
	privilegeSet, privilegeSetOk := d.GetOk("privilege_set")

	if accessLevelOk && accessLevel.(string) == "Group Access" {
		if !privilegeSetOk || privilegeSet.(string) != "Custom" {
			return fmt.Errorf("when 'access_level' is 'Group Access', 'privilege_set' must be 'Custom'")
		}
	}

	return nil
}

// validateCasperAdminUsePrivileges checks for required privileges when "Use Casper Admin" is selected.
func validateCasperAdminUsePrivileges(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if v, ok := d.GetOk("casper_admin_privileges"); ok {
		casperAdminPrivilegesSet := v.(*schema.Set)
		if casperAdminPrivilegesSet.Contains("Use Casper Admin") {
			requiredReadPrivileges := []string{
				"Read Categories",
				"Read Directory Bindings",
				"Read Dock Items",
				"Read Packages",
				"Read Printers",
				"Read Scripts",
			}

			jssObjectsPrivilegesSet, ok := d.GetOk("jss_objects_privileges")
			if !ok {
				return fmt.Errorf("when 'Use Casper Admin' is selected, the following JSS Object Privileges are required: %v", requiredReadPrivileges)
			}

			jssPrivilegesSet := make(map[string]bool)
			for _, priv := range jssObjectsPrivilegesSet.(*schema.Set).List() {
				jssPrivilegesSet[priv.(string)] = true
			}

			for _, requiredPriv := range requiredReadPrivileges {
				if !jssPrivilegesSet[requiredPriv] {
					return fmt.Errorf("missing required privilege '%s' in 'jss_objects_privileges' when 'Use Casper Admin' is selected in 'casper_admin_privileges'", requiredPriv)
				}
			}
		}
	}

	return nil
}

// validateCasperAdminSavePrivileges checks for required privileges when "Save With Casper Admin" is selected.
func validateCasperAdminSavePrivileges(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if v, ok := d.GetOk("casper_admin_privileges"); ok {
		casperAdminPrivilegesSet := v.(*schema.Set)
		if casperAdminPrivilegesSet.Contains("Save With Casper Admin") {
			requiredCrudPrivileges := []string{
				"Create Categories", "Read Categories", "Update Categories", "Delete Categories",
				"Create Directory Bindings", "Read Directory Bindings", "Update Directory Bindings", "Delete Directory Bindings",
				"Create Dock Items", "Read Dock Items", "Update Dock Items", "Delete Dock Items",
				"Create Packages", "Read Packages", "Update Packages", "Delete Packages",
				"Create Printers", "Read Printers", "Update Printers", "Delete Printers",
				"Create Scripts", "Read Scripts", "Update Scripts", "Delete Scripts",
			}

			jssObjectsPrivilegesSet, ok := d.GetOk("jss_objects_privileges")
			if !ok {
				return fmt.Errorf("when 'Save With Casper Admin' is selected, the following JSS Object Privileges are required: %v", requiredCrudPrivileges)
			}

			jssPrivilegesSet := make(map[string]bool)
			for _, priv := range jssObjectsPrivilegesSet.(*schema.Set).List() {
				jssPrivilegesSet[priv.(string)] = true
			}

			for _, requiredPriv := range requiredCrudPrivileges {
				if !jssPrivilegesSet[requiredPriv] {
					return fmt.Errorf("missing required privilege '%s' in 'jss_objects_privileges' when 'Save With Casper Admin' is selected in 'casper_admin_privileges'", requiredPriv)
				}
			}
		}
	}

	return nil
}

// validateReverseDependencyChecks performs reverse dependency checks for specific privileges.
func validateReverseDependencyChecks(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	jssObjectsPrivileges, ok := d.GetOk("jss_objects_privileges")
	if ok {
		jssObjectsPrivilegesSet := jssObjectsPrivileges.(*schema.Set)
		jssPrivilegesSet := make(map[string]bool)
		for _, priv := range jssObjectsPrivilegesSet.List() {
			jssPrivilegesSet[priv.(string)] = true
		}

		readPrivileges := []string{"Read Categories", "Read Directory Bindings", "Read Dock Items", "Read Packages", "Read Printers", "Read Scripts"}
		allReadPrivilegesPresent := true
		for _, priv := range readPrivileges {
			if !jssPrivilegesSet[priv] {
				allReadPrivilegesPresent = false
				break
			}
		}

		if allReadPrivilegesPresent {
			casperAdminPrivileges, ok := d.GetOk("casper_admin_privileges")
			if ok {
				casperAdminPrivilegesSet := casperAdminPrivileges.(*schema.Set)
				if !casperAdminPrivilegesSet.Contains("Use Casper Admin") {
					return fmt.Errorf("when the following JSS Object Privileges are set %v, 'Use Casper Admin' must be included in 'casper_admin_privileges'", readPrivileges)
				}
			} else {
				return fmt.Errorf("when the following JSS Object Privileges are set %v, 'Use Casper Admin' must be included in 'casper_admin_privileges'", readPrivileges)
			}
		}

		crudPrivileges := []string{
			"Create Categories", "Read Categories", "Update Categories", "Delete Categories",
			"Create Directory Bindings", "Read Directory Bindings", "Update Directory Bindings", "Delete Directory Bindings",
			"Create Dock Items", "Read Dock Items", "Update Dock Items", "Delete Dock Items",
			"Create Packages", "Read Packages", "Update Packages", "Delete Packages",
			"Create Printers", "Read Printers", "Update Printers", "Delete Printers",
			"Create Scripts", "Read Scripts", "Update Scripts", "Delete Scripts",
		}
		allCrudPrivilegesPresent := true
		for _, priv := range crudPrivileges {
			if !jssPrivilegesSet[priv] {
				allCrudPrivilegesPresent = false
				break
			}
		}

		if allCrudPrivilegesPresent {
			casperAdminPrivileges, ok := d.GetOk("casper_admin_privileges")
			if ok {
				casperAdminPrivilegesSet := casperAdminPrivileges.(*schema.Set)
				if !casperAdminPrivilegesSet.Contains("Save With Casper Admin") {
					return fmt.Errorf("when the following JSS Object Privileges are set %v, 'Save With Casper Admin' must be included in 'casper_admin_privileges'", crudPrivileges)
				}
			} else {
				return fmt.Errorf("when the following JSS Object Privileges are set %v, 'Save With Casper Admin' must be included in 'casper_admin_privileges'", crudPrivileges)
			}
		}
	}

	return nil
}
