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

var validJSSObjectsPrivileges = []string{
	"Create Advanced Computer Searches",
	"Read Advanced Computer Searches",
	"Update Advanced Computer Searches",
	"Delete Advanced Computer Searches",
	"Create Advanced Mobile Device Searches",
	"Read Advanced Mobile Device Searches",
	"Update Advanced Mobile Device Searches",
	"Delete Advanced Mobile Device Searches",
	"Create Advanced User Searches",
	"Read Advanced User Searches",
	"Update Advanced User Searches",
	"Delete Advanced User Searches",
	"Create Advanced User Content Searches",
	"Read Advanced User Content Searches",
	"Update Advanced User Content Searches",
	"Delete Advanced User Content Searches",
	"Create AirPlay Permissions",
	"Read AirPlay Permissions",
	"Update AirPlay Permissions",
	"Delete AirPlay Permissions",
	"Create Allowed File Extension",
	"Read Allowed File Extension",
	"Delete Allowed File Extension",
	"Create API Integrations",
	"Read API Integrations",
	"Update API Integrations",
	"Delete API Integrations",
	"Create API Roles",
	"Read API Roles",
	"Update API Roles",
	"Delete API Roles",
	"Create Attachment Assignments",
	"Read Attachment Assignments",
	"Update Attachment Assignments",
	"Delete Attachment Assignments",
	"Create Device Enrollment Program Instances",
	"Read Device Enrollment Program Instances",
	"Update Device Enrollment Program Instances",
	"Delete Device Enrollment Program Instances",
	"Create Buildings",
	"Read Buildings",
	"Update Buildings",
	"Delete Buildings",
	"Create Categories",
	"Read Categories",
	"Update Categories",
	"Delete Categories",
	"Create Classes",
	"Read Classes",
	"Update Classes",
	"Delete Classes",
	"Create Computer Enrollment Invitations",
	"Read Computer Enrollment Invitations",
	"Update Computer Enrollment Invitations",
	"Delete Computer Enrollment Invitations",
	"Create Computer Extension Attributes",
	"Read Computer Extension Attributes",
	"Update Computer Extension Attributes",
	"Delete Computer Extension Attributes",
	"Create Custom Paths",
	"Read Custom Paths",
	"Update Custom Paths",
	"Delete Custom Paths",
	"Create Computer PreStage Enrollments",
	"Read Computer PreStage Enrollments",
	"Update Computer PreStage Enrollments",
	"Delete Computer PreStage Enrollments",
	"Create Computers",
	"Read Computers",
	"Update Computers",
	"Delete Computers",
	"Create Departments",
	"Read Departments",
	"Update Departments",
	"Delete Departments",
	"Create Mobile Device Extension Attributes",
	"Read Mobile Device Extension Attributes",
	"Update Mobile Device Extension Attributes",
	"Delete Mobile Device Extension Attributes",
	"Create Device Name Patterns",
	"Read Device Name Patterns",
	"Update Device Name Patterns",
	"Delete Device Name Patterns",
	"Create Directory Bindings",
	"Read Directory Bindings",
	"Update Directory Bindings",
	"Delete Directory Bindings",
	"Create Disk Encryption Configurations",
	"Read Disk Encryption Configurations",
	"Update Disk Encryption Configurations",
	"Delete Disk Encryption Configurations",
	"Create Disk Encryption Institutional Configurations",
	"Read Disk Encryption Institutional Configurations",
	"Update Disk Encryption Institutional Configurations",
	"Delete Disk Encryption Institutional Configurations",
	"Create Dock Items",
	"Read Dock Items",
	"Update Dock Items",
	"Delete Dock Items",
	"Create eBooks",
	"Read eBooks",
	"Update eBooks",
	"Delete eBooks",
	"Create Enrollment Customizations",
	"Read Enrollment Customizations",
	"Update Enrollment Customizations",
	"Delete Enrollment Customizations",
	"Create Enrollment Profiles",
	"Read Enrollment Profiles",
	"Update Enrollment Profiles",
	"Delete Enrollment Profiles",
	"Create Patch External Source",
	"Read Patch External Source",
	"Update Patch External Source",
	"Delete Patch External Source",
	"Create File Attachments",
	"Read File Attachments",
	"Update File Attachments",
	"Delete File Attachments",
	"Create Distribution Points",
	"Read Distribution Points",
	"Update Distribution Points",
	"Delete Distribution Points",
	"Create Push Certificates",
	"Read Push Certificates",
	"Update Push Certificates",
	"Delete Push Certificates",
	"Create iBeacon",
	"Read iBeacon",
	"Update iBeacon",
	"Delete iBeacon",
	"Create Infrastructure Managers",
	"Read Infrastructure Managers",
	"Update Infrastructure Managers",
	"Delete Infrastructure Managers",
	"Create Inventory Preload Records",
	"Read Inventory Preload Records",
	"Update Inventory Preload Records",
	"Delete Inventory Preload Records",
	"Create VPP Invitations",
	"Read VPP Invitations",
	"Update VPP Invitations",
	"Delete VPP Invitations",
	"Create Jamf Connect Deployments",
	"Read Jamf Connect Deployments",
	"Update Jamf Connect Deployments",
	"Delete Jamf Connect Deployments",
	"Create Jamf Content Distribution Server Files",
	"Read Jamf Content Distribution Server Files",
	"Delete Jamf Content Distribution Server Files",
	"Create Jamf Protect Deployments",
	"Read Jamf Protect Deployments",
	"Update Jamf Protect Deployments",
	"Delete Jamf Protect Deployments",
	"Create JSON Web Token Configuration",
	"Read JSON Web Token Configuration",
	"Update JSON Web Token Configuration",
	"Delete JSON Web Token Configuration",
	"Create Keystore",
	"Read Keystores",
	"Update Keystores",
	"Delete Keystores",
	"Create LDAP Servers",
	"Read LDAP Servers",
	"Update LDAP Servers",
	"Delete LDAP Servers",
	"Create Licensed Software",
	"Read Licensed Software",
	"Update Licensed Software",
	"Delete Licensed Software",
	"Create Mac Applications",
	"Read Mac Applications",
	"Update Mac Applications",
	"Delete Mac Applications",
	"Create macOS Configuration Profiles",
	"Read macOS Configuration Profiles",
	"Update macOS Configuration Profiles",
	"Delete macOS Configuration Profiles",
	"Create Maintenance Pages",
	"Read Maintenance Pages",
	"Update Maintenance Pages",
	"Delete Maintenance Pages",
	"Create Managed Software Updates",
	"Read Managed Software Updates",
	"Update Managed Software Updates",
	"Delete Managed Software Updates",
	"Create Mobile Device Applications",
	"Read Mobile Device Applications",
	"Update Mobile Device Applications",
	"Delete Mobile Device Applications",
	"Create iOS Configuration Profiles",
	"Read iOS Configuration Profiles",
	"Update iOS Configuration Profiles",
	"Delete iOS Configuration Profiles",
	"Create Mobile Device Enrollment Invitations",
	"Read Mobile Device Enrollment Invitations",
	"Update Mobile Device Enrollment Invitations",
	"Delete Mobile Device Enrollment Invitations",
	"Create Mobile Device Managed App Configurations",
	"Read Mobile Device Managed App Configurations",
	"Update Mobile Device Managed App Configurations",
	"Delete Mobile Device Managed App Configurations",
	"Create Mobile Device PreStage Enrollments",
	"Read Mobile Device PreStage Enrollments",
	"Update Mobile Device PreStage Enrollments",
	"Delete Mobile Device PreStage Enrollments",
	"Create Mobile Devices",
	"Read Mobile Devices",
	"Update Mobile Devices",
	"Delete Mobile Devices",
	"Create Network Integration",
	"Read Network Integration",
	"Update Network Integration",
	"Delete Network Integration",
	"Create Network Segments",
	"Read Network Segments",
	"Update Network Segments",
	"Delete Network Segments",
	"Create Packages",
	"Read Packages",
	"Update Packages",
	"Delete Packages",
	"Create Patch Management Software Titles",
	"Read Patch Management Software Titles",
	"Update Patch Management Software Titles",
	"Delete Patch Management Software Titles",
	"Create Patch Policies",
	"Read Patch Policies",
	"Update Patch Policies",
	"Delete Patch Policies",
	"Create Peripheral Types",
	"Read Peripheral Types",
	"Update Peripheral Types",
	"Delete Peripheral Types",
	"Create Personal Device Configurations",
	"Read Personal Device Configurations",
	"Update Personal Device Configurations",
	"Delete Personal Device Configurations",
	"Create Personal Device Profiles",
	"Read Personal Device Profiles",
	"Update Personal Device Profiles",
	"Delete Personal Device Profiles",
	"Create Policies",
	"Read Policies",
	"Update Policies",
	"Delete Policies",
	"Create Printers",
	"Read Printers",
	"Update Printers",
	"Delete Printers",
	"Create Provisioning Profiles",
	"Read Provisioning Profiles",
	"Update Provisioning Profiles",
	"Delete Provisioning Profiles",
	"Create Push Certificates",
	"Read Push Certificates",
	"Update Push Certificates",
	"Delete Push Certificates",
	"Create Remote Administration",
	"Read Remote Administration",
	"Update Remote Administration",
	"Delete Remote Administration",
	"Create Removable MAC Address",
	"Read Removable MAC Address",
	"Update Removable MAC Address",
	"Delete Removable MAC Address",
	"Create Restricted Software",
	"Read Restricted Software",
	"Update Restricted Software",
	"Delete Restricted Software",
	"Create Scripts",
	"Read Scripts",
	"Update Scripts",
	"Delete Scripts",
	"Create Self Service Bookmarks",
	"Read Self Service Bookmarks",
	"Update Self Service Bookmarks",
	"Delete Self Service Bookmarks",
	"Create Self Service Branding Configuration",
	"Read Self Service Branding Configuration",
	"Update Self Service Branding Configuration",
	"Delete Self Service Branding Configuration",
	"Create Sites",
	"Read Sites",
	"Update Sites",
	"Delete Sites",
	"Create Smart Computer Groups",
	"Read Smart Computer Groups",
	"Update Smart Computer Groups",
	"Delete Smart Computer Groups",
	"Create Smart Mobile Device Groups",
	"Read Smart Mobile Device Groups",
	"Update Smart Mobile Device Groups",
	"Delete Smart Mobile Device Groups",
	"Create Smart User Groups",
	"Read Smart User Groups",
	"Update Smart User Groups",
	"Delete Smart User Groups",
	"Create Software Update Servers",
	"Read Software Update Servers",
	"Update Software Update Servers",
	"Delete Software Update Servers",
	"Create Static Computer Groups",
	"Read Static Computer Groups",
	"Update Static Computer Groups",
	"Delete Static Computer Groups",
	"Create Static Mobile Device Groups",
	"Read Static Mobile Device Groups",
	"Update Static Mobile Device Groups",
	"Delete Static Mobile Device Groups",
	"Create Static User Groups",
	"Read Static User Groups",
	"Update Static User Groups",
	"Delete Static User Groups",
	"Create Accounts",
	"Read Accounts",
	"Update Accounts",
	"Delete Accounts",
	"Create User Extension Attributes",
	"Read User Extension Attributes",
	"Update User Extension Attributes",
	"Delete User Extension Attributes",
	"Create User",
	"Read User",
	"Update User",
	"Delete User",
	"Create VPP Assignment",
	"Read VPP Assignment",
	"Update VPP Assignment",
	"Delete VPP Assignment",
	"Create Volume Purchasing Locations",
	"Read Volume Purchasing Locations",
	"Update Volume Purchasing Locations",
	"Delete Volume Purchasing Locations",
	"Create Webhooks",
	"Read Webhooks",
	"Update Webhooks",
	"Delete Webhooks",
}

// validateJSSObjectsPrivileges ensures that each privilege in the list is valid.
func validateJSSObjectsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	// Check if the value is in the list of valid privileges
	for _, validPrivilege := range validJSSObjectsPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	errs = append(errs, fmt.Errorf("invalid value '%s' for %s: not a recognized JSS Object Privilege", v, key))
	return warns, errs
}

var validJSSSettingsPrivileges = []string{
	"Read Activation Code",
	"Update Activation Code",
	"Read Apache Tomcat Settings",
	"Update Apache Tomcat Settings",
	"Read Apple Configurator Enrollment",
	"Update Apple Configurator Enrollment",
	"Read Education Settings",
	"Update Education Settings",
	"Read Mobile Device App Maintenance Settings",
	"Update Mobile Device App Maintenance Settings",
	"Read Automatic Mac App Updates Settings",
	"Update Automatic Mac App Updates Settings",
	"Read Automatically Renew MDM Profile Settings",
	"Update Automatically Renew MDM Profile Settings",
	"Read Cache",
	"Update Cache",
	"Read Change Management",
	"Update Change Management",
	"Read Computer Check-In",
	"Update Computer Check-In",
	"Read Cloud Distribution Point",
	"Update Cloud Distribution Point",
	"Read Cloud Services Settings",
	"Update Cloud Services Settings",
	"Read Clustering",
	"Update Clustering",
	"Read Computer Check-In",
	"Update Computer Check-In",
	"Read Computer Inventory Collection",
	"Update Computer Inventory Collection",
	"Read Computer Inventory Collection Settings",
	"Update Computer Inventory Collection Settings",
	"Read Conditional Access",
	"Update Conditional Access",
	"Read Customer Experience Metrics",
	"Update Customer Experience Metrics",
	"Read Device Compliance Information",
	"Read Mobile Device Inventory Collection",
	"Update Mobile Device Inventory Collection",
	"Read Engage Settings",
	"Update Engage Settings",
	"Read GSX Connection",
	"Update GSX Connection",
	"Read Patch Internal Source",
	"Read Jamf Connect Settings",
	"Update Jamf Connect Settings",
	"Read Parent App Settings",
	"Update Parent App Settings",
	"Read Jamf Protect Settings",
	"Update Jamf Protect Settings",
	"Read JSS URL",
	"Update JSS URL",
	"Read Teacher App Settings",
	"Update Teacher App Settings",
	"Read Limited Access Settings",
	"Update Limited Access Settings",
	"Read Retention Policy",
	"Update Retention Policy",
	"Read Onboarding Configuration",
	"Update Onboarding Configuration",
	"Read Password Policy",
	"Update Password Policy",
	"Read Patch Management Settings",
	"Update Patch Management Settings",
	"Read PKI",
	"Update PKI",
	"Read Re-enrollment",
	"Update Re-enrollment",
	"Read Computer Security",
	"Update Computer Security",
	"Read App Request Settings",
	"Update App Request Settings",
	"Read Mobile Device Self Service",
	"Update Mobile Device Self Service",
	"Read Self Service",
	"Update Self Service",
	"Read SMTP Server",
	"Update SMTP Server",
	"Read SSO Settings",
	"Update SSO Settings",
	"Read User-Initiated Enrollment",
	"Update User-Initiated Enrollment",
}

// validateJSSSettingsPrivileges checks if each value in the jss_settings_privileges field
// matches a value from a predefined list of valid JSS Setting Privileges.
func validateJSSSettingsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	// Check if the value is in the list of valid JSS Settings Privileges
	for _, validPrivilege := range validJSSSettingsPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	// Add an error if the value is not found in the valid privileges list
	errs = append(errs, fmt.Errorf("invalid value '%s' for %s: not a recognized JSS Setting Privilege", v, key))
	return warns, errs
}

var validJSSActionsPrivileges = []string{
	"Allow User to Enroll",
	"Assign Users to Computers",
	"Assign Users to Mobile Devices",
	"Change Password",
	"Dismiss Notifications",
	"Enroll Computers and Mobile Devices",
	"Flush MDM Commands",
	"Flush Policy Logs",
	"Jamf Packages Action",
	"Remove restrictions set by Jamf Parent",
	"CLEAR_TEACHER_PROFILE_PRIVILEGE",
	"Renewal of the Built-in Certificate Authority",
	"Jamf Connect Deployment Retry",
	"Jamf Protect Deployment Retry",
	"Send Application Attributes Command",
	"Send Blank Pushes to Mobile Devices",
	"Send Computer Remote Desktop Command",
	"Send Computer Set Activation Lock Command",
	"Send Declarative Management Command",
	"Send Disable Bootstrap Token Command",
	"Send Email to End Users via JSS",
	"Send Enable Bootstrap Token Command",
	"Send Inventory Requests to Mobile Devices",
	"Send Local Admin Password Command",
	"Send Messages to Self Service Mobile",
	"Send Mobile Device Bluetooth Command",
	"Send Mobile Device Diagnostics and Usage Reporting and App Analytics Commands",
	"Send Mobile Device Disable Data Roaming Command",
	"Send Mobile Device Disable Voice Roaming Command",
	"Send Mobile Device Enable Data Roaming Command",
	"Send Mobile Device Enable Voice Roaming Command",
	"Send Mobile Device Lost Mode Command",
	"Send Mobile Device Managed Settings Command",
	"Send Mobile Device Mirroring Command",
	"Send Mobile Device Personal Hotspot Command",
	"Send Mobile Device Remote Command to Download and Install iOS Update",
	"Send Mobile Device Remote Lock Command",
	"Send Mobile Device Remote Wipe Command",
	"Send Mobile Device Remove Passcode Command",
	"Send Mobile Device Remove Restrictions Password Command",
	"Send Mobile Device Set Activation Lock Command",
	"Send Mobile Device Set Wallpaper Command",
	"Send Mobile Device Software Update Recommendation Cadence Command",
	"Send Set Timezone Command",
	"Send Software Update Settings Command",
	"Unmanage Mobile Devices",
	"Update Local Admin Password Settings",
	"View Event Logs",
	"View JSS Information",
	"View License Serial Numbers",
	"View Local Admin Password",
	"View Local Admin Password Audit History",
	"View MDM command information in Jamf Pro API",
}

// validateJSSActionsPrivileges checks if each value in the jss_actions_privileges field
// matches a value from a predefined list of valid JSS Actions Privileges.
func validateJSSActionsPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	// Check if the value is in the list of valid JSS Actions Privileges
	for _, validPrivilege := range validJSSActionsPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	// Add an error if the value is not found in the valid privileges list
	errs = append(errs, fmt.Errorf("invalid value '%s' for %s: not a recognized JSS Action Privilege", v, key))
	return warns, errs
}

var validCasperAdminPrivileges = []string{
	"Use Casper Admin",
	"Save With Casper Admin",
}

// validateCasperAdminPrivileges checks if each value in the casper_admin_privileges field
// matches a value from a predefined list of valid Casper Admin Privileges.
func validateCasperAdminPrivileges(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return warns, errs
	}

	// Check if the value is in the list of valid Casper Admin Privileges
	for _, validPrivilege := range validCasperAdminPrivileges {
		if v == validPrivilege {
			return warns, errs
		}
	}

	// Add an error if the value is not found in the valid privileges list
	errs = append(errs, fmt.Errorf("invalid value '%s' for %s: not a recognized Casper Admin Privilege", v, key))
	return warns, errs
}
