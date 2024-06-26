// apiroles_data_validation.go
package apiroles

import (
	"fmt"
	"sort"
	"strings"
)

// TODO review this whole file

// List of valid privileges - replace these with actual valid privileges.
var validPrivileges = []string{
	"Allow User to Enroll",
	"Assign Users to Computers",
	"Assign Users to Mobile Devices",
	"Change Password",
	"CLEAR_TEACHER_PROFILE_PRIVILEGE",
	"Create Accounts",
	"Create Advanced Computer Searches",
	"Create Advanced Mobile Device Searches",
	"Create Advanced User Content Searches",
	"Create Advanced User Searches",
	"Create AirPlay Permissions",
	"Create Allowed File Extension",
	"Create API Integrations",
	"Create API Roles",
	"Create Attachment Assignments",
	"Create Buildings",
	"Create Cache",
	"Create Categories",
	"Create Classes",
	"Create Computer Enrollment Invitations",
	"Create Computer Extension Attributes",
	"Create Computer PreStage Enrollments",
	"Create Computers",
	"Create Custom Paths",
	"Create Departments",
	"Create Device Enrollment Program Instances",
	"Create Device Name Patterns",
	"Create Directory Bindings",
	"Create Disk Encryption Configurations",
	"Create Disk Encryption Institutional Configurations",
	"Create Distribution Points",
	"Create Dock Items",
	"Create eBooks",
	"Create Enrollment Customizations",
	"Create Enrollment Profiles",
	"Create File Attachments",
	"Create iOS Configuration Profiles",
	"Create iBeacon",
	"Create Infrastructure Managers",
	"Create Inventory Preload Records",
	"Create Jamf Connect Deployments",
	"Create Jamf Content Distribution Server Files",
	"Create Jamf Protect Deployments",
	"Create JSON Web Token Configuration",
	"Create Keystore",
	"Create LDAP Servers",
	"Create Licensed Software",
	"Create Local Admin Password",
	"Create Mac Applications",
	"Create macOS Configuration Profiles",
	"Create Maintenance Pages",
	"Create Managed Software Updates",
	"Create Mobile Device Applications",
	"Create Mobile Device Enrollment Invitations",
	"Create Mobile Device Extension Attributes",
	"Create Mobile Device Managed App Configurations",
	"Create Mobile Device PreStage Enrollments",
	"Create Mobile Devices",
	"Create Network Integration",
	"Create Network Segments",
	"Create Packages",
	"Create Parent App Settings",
	"Create Patch External Source",
	"Create Patch Management Software Titles",
	"Create Patch Policies",
	"Create Personal Device Configurations",
	"Create Personal Device Profiles",
	"Create Policies",
	"Create Printers",
	"Create Provisioning Profiles",
	"Create Push Certificates",
	"Create Removable MAC Address",
	"Create Restricted Software",
	"Create Scripts",
	"Create Self Service",
	"Create Self Service Bookmarks",
	"Create Self Service Branding Configuration",
	"Create Sites",
	"Create Smart Computer Groups",
	"Create Smart Mobile Device Groups",
	"Create Smart User Groups",
	"Create Software Update Servers",
	"Create Static Computer Groups",
	"Create Static Mobile Device Groups",
	"Create Static User Groups",
	"Create Teacher App Settings",
	"Create User",
	"Create User Extension Attributes",
	"Create VPP Assignment",
	"Create VPP Invitations",
	"Create Webhooks",
	"Delete Accounts",
	"Delete Advanced Computer Searches",
	"Delete Advanced Mobile Device Searches",
	"Delete Advanced User Content Searches",
	"Delete Advanced User Searches",
	"Delete AirPlay Permissions",
	"Delete Allowed File Extension",
	"Delete API Integrations",
	"Delete API Roles",
	"Delete Attachment Assignments",
	"Delete Buildings",
	"Delete Cache",
	"Delete Categories",
	"Delete Classes",
	"Delete Computer Enrollment Invitations",
	"Delete Computer Extension Attributes",
	"Delete Computer PreStage Enrollments",
	"Delete Computers",
	"Delete Custom Paths",
	"Delete Departments",
	"Delete Device Enrollment Program Instances",
	"Delete Device Name Patterns",
	"Delete Directory Bindings",
	"Delete Disk Encryption Configurations",
	"Delete Disk Encryption Institutional Configurations",
	"Delete Distribution Points",
	"Delete Dock Items",
	"Delete eBooks",
	"Delete Enrollment Customizations",
	"Delete Enrollment Profiles",
	"Delete File Attachments",
	"Delete GSX Connection",
	"Delete Infrastructure Managers",
	"Delete iBeacon",
	"Delete iOS Configuration Profiles",
	"Delete Jamf Connect Deployments",
	"Delete Jamf Content Distribution Server Files",
	"Delete Jamf Protect Deployments",
	"Delete JSON Web Token Configuration",
	"Delete Keystores",
	"Delete LDAP Servers",
	"Delete Licensed Software",
	"Delete Limited Access Settings",
	"Delete Local Admin Password",
	"Delete Mac Applications",
	"Delete macOS Configuration Profiles",
	"Delete Maintenance Pages",
	"Delete Managed Software Updates",
	"Delete Mobile Device Applications",
	"Delete Mobile Device Enrollment Invitations",
	"Delete Mobile Device Extension Attributes",
	"Delete Mobile Device Managed App Configurations",
	"Delete Mobile Device PreStage Enrollments",
	"Delete Mobile Devices",
	"Delete Network Integration",
	"Delete Network Segments",
	"Delete Packages",
	"Delete Parent App Settings",
	"Delete Patch External Source",
	"Delete Patch Management Software Titles",
	"Delete Patch Policies",
	"Delete Peripheral Types",
	"Delete Personal Device Configurations",
	"Delete Personal Device Profiles",
	"Delete Policies",
	"Delete Printers",
	"Delete Provisioning Profiles",
	"Delete Push Certificates",
	"Delete Removable MAC Address",
	"Delete Remote Administration",
	"Delete Restricted Software",
	"Delete Retention Policy",
	"Delete Scripts",
	"Delete Self Service",
	"Delete Self Service Bookmarks",
	"Delete Self Service Branding Configuration",
	"Delete Sites",
	"Delete Smart Computer Groups",
	"Delete Smart Mobile Device Groups",
	"Delete Smart User Groups",
	"Delete Software Update Servers",
	"Delete Static Computer Groups",
	"Delete Static Mobile Device Groups",
	"Delete Static User Groups",
	"Delete Teacher App Settings",
	"Delete User",
	"Delete User Extension Attributes",
	"Delete VPP Assignment",
	"Delete VPP Invitations",
	"Delete Webhooks",
	"Dismiss Notifications",
	"Enroll Computers and Mobile Devices",
	"Flush MDM Commands",
	"Flush Policy Logs",
	"Jamf Connect Deployment Retry",
	"Jamf Packages Action",
	"Jamf Protect Deployment Retry",
	"Read Accounts",
	"Read Activation Code",
	"Read Activation Lock Bypass Code",
	"Read Advanced Computer Searches",
	"Read Advanced Mobile Device Searches",
	"Read Advanced User Content Searches",
	"Read Advanced User Searches",
	"Read AirPlay Permissions",
	"Read Allowed File Extension",
	"Read API Integrations",
	"Read API Roles",
	"Read App Request Settings",
	"Read Apple Configurator Enrollment",
	"Read Automatically Renew MDM Profile Settings",
	"Read Buildings",
	"Read Cache",
	"Read Categories",
	"Read Change Management",
	"Read Classes",
	"Read Clustering",
	"Read Cloud Distribution Point",
	"Read Cloud Services Settings",
	"Read Computer Check-In",
	"Read Computer Enrollment Invitations",
	"Read Computer Extension Attributes",
	"Read Computer Inventory Collection",
	"Read Computer Inventory Collection Settings",
	"Read Computer PreStage Enrollments",
	"Read Computer Security",
	"Read Computers",
	"Read Conditional Access",
	"Read Custom Paths",
	"Read Customer Experience Metrics",
	"Read Departments",
	"Read Device Compliance Information",
	"Read Device Enrollment Program Instances",
	"Read Device Name Patterns",
	"Read Directory Bindings",
	"Read Disk Encryption Configurations",
	"Read Disk Encryption Institutional Configurations",
	"Read Distribution Points",
	"Read Dock Items",
	"Read eBooks",
	"Read Education Settings",
	"Read Engage Settings",
	"Read Enrollment Customizations",
	"Read Enrollment Profiles",
	"Read Event Logs",
	"Read File Attachments",
	"Read GSX Connection",
	"Read Infrastructure Managers",
	"Read Inventory Preload Records",
	"Read iBeacon",
	"Read Jamf Connect Deployments",
	"Read Jamf Connect Settings",
	"Read Jamf Protect Deployments",
	"Read Jamf Protect Settings",
	"Read Jamf Content Distribution Server Files",
	"Read JSON Web Token Configuration",
	"Read Keystores",
	"Read Knobs",
	"Read LDAP Servers",
	"Read Licensed Software",
	"Read Limited Access Settings",
	"Read Local Admin Password",
	"Read Local Admin Password Audit History",
	"Read Mac Applications",
	"Read macOS Configuration Profiles",
	"Read Maintenance Pages",
	"Read Managed Software Updates",
	"Read MDM command information in Jamf Pro API",
	"Read Mobile Device App Maintenance Settings",
	"Read Mobile Device Applications",
	"Read Mobile Device Enrollment Invitations",
	"Read Mobile Device Extension Attributes",
	"Read Mobile Device Inventory Collection",
	"Read Mobile Device Managed App Configurations",
	"Read Mobile Device PreStage Enrollments",
	"Read Mobile Device Self Service",
	"Read Mobile Devices",
	"Read Network Integration",
	"Read Network Segments",
	"Read Onboarding Configuration",
	"Read Packages",
	"Read Parent App Settings",
	"Read Password Policy",
	"Read Patch External Source",
	"Read Patch Internal Source",
	"Read Patch Management Settings",
	"Read Patch Management Software Titles",
	"Read Patch Policies",
	"Read Personal Device Configurations",
	"Read Personal Device Profiles",
	"Read PKI",
	"Read Policies",
	"Read Printers",
	"Read Privileges",
	"Read Provisioning Profiles",
	"Read Push Certificates",
	"Read Re-enrollment",
	"Read Remote Administration",
	"Read Removable MAC Address",
	"Read Renewal of the Built-in Certificate Authority",
	"Read Restricted Software",
	"Read Retention Policy",
	"Read SSO Settings",
	"Read Scripts",
	"Read Self Service",
	"Read Self Service Bookmarks",
	"Read Self Service Branding Configuration",
	"Read Sites",
	"Read Smart Computer Groups",
	"Read Smart Mobile Device Groups",
	"Read Smart User Groups",
	"Read SMTP Server",
	"Read Software Update Servers",
	"Read Static Computer Groups",
	"Read Static Mobile Device Groups",
	"Read Static User Groups",
	"Read Teacher App Settings",
	"Read User",
	"Read User Extension Attributes",
	"Read User-Initiated Enrollment",
	"Read VPP Assignment",
	"Read VPP Invitations",
	"Read Volume Purchasing Locations",
	"Read Webhooks",
	"Remove Jamf Parent management capabilities",
	"Remove restrictions set by Jamf Parent",
	"Renewal of the Built-in Certificate Authority",
	"Send Application Attributes Command",
	"Send Blank Pushes to Mobile Devices",
	"Send Command to Renew MDM Profile",
	"Send Computer Bluetooth Command",
	"Send Computer delete User Account Command",
	"Send Computer Remote Command to Download and Install OS X Update",
	"Send Computer Remote Command to Install Package",
	"Send Computer Remote Desktop Command",
	"Send Computer Remote Lock Command",
	"Send Computer Remote Wipe Command",
	"Send Computer Set Activation Lock Command",
	"Send Computer Unmanage Command",
	"Send Computer Unlock User Account Command",
	"Send Declarative Management Command",
	"Send Disable Bootstrap Token Command",
	"Send Email to End Users via JSS",
	"Send Enable Bootstrap Token Command",
	"Send Inventory Requests to Mobile Devices",
	"Send Local Admin Password Command",
	"Send Messages to Self Service Mobile",
	"Send MDM Check In Command",
	"Send Mobile Device Bluetooth Command",
	"Send Mobile Device Disable Data Roaming Command",
	"Send Mobile Device Disable Voice Roaming Command",
	"Send Mobile Device Enable Data Roaming Command",
	"Send Mobile Device Enable Voice Roaming Command",
	"Send Mobile Device Lost Mode Command",
	"Send Mobile Device Managed Settings Command",
	"Send Mobile Device Mirroring Command",
	"Send Mobile Device Personal Hotspot Command",
	"Send Mobile Device Refresh Cellular Plans Command",
	"Send Mobile Device Remote Command to Download and Install iOS Update",
	"Send Mobile Device Remote Lock Command",
	"Send Mobile Device Remote Wipe Command",
	"Send Mobile Device Remove Passcode Command",
	"Send Mobile Device Remove Restrictions Password Command",
	"Send Mobile Device Restart Device Command",
	"Send Mobile Device Set Activation Lock Command",
	"Send Mobile Device Set Device Name Command",
	"Send Mobile Device Set Wallpaper Command",
	"Send Mobile Device Shared Device Configuration Commands",
	"Send Mobile Device Shared iPad Commands",
	"Send Mobile Device Shut Down Command",
	"Send Mobile Device Software Update Recommendation Cadence Command",
	"Send Set Recovery Lock Command",
	"Send Set Timezone Command",
	"Send Software Update Settings Command",
	"Send Update Passcode Lock Grace Period Command",
	"Start Remote Assist Session",
	"Unmanage Mobile Devices",
	"Update Accounts",
	"Update Activation Code",
	"Update Advanced Computer Searches",
	"Update Advanced Mobile Device Searches",
	"Update Advanced User Content Searches",
	"Update Advanced User Searches",
	"Update AirPlay Permissions",
	"Update API Integrations",
	"Update API Roles",
	"Update Apache Tomcat Settings",
	"Update App Request Settings",
	"Update Apple Configurator Enrollment",
	"Update Attachment Assignments",
	"Update Automatically Renew MDM Profile Settings",
	"Update Automatic Mac App Updates Settings",
	"Update Buildings",
	"Update Cache",
	"Update Categories",
	"Update Change Management",
	"Update Classes",
	"Update Clustering",
	"Update Cloud Distribution Point",
	"Update Cloud Services Settings",
	"Update Computer Check-In",
	"Update Computer Enrollment Invitations",
	"Update Computer Extension Attributes",
	"Update Computer Inventory Collection",
	"Update Computer Inventory Collection Settings",
	"Update Computer PreStage Enrollments",
	"Update Computer Security",
	"Update Computers",
	"Update Conditional Access",
	"Update Custom Paths",
	"Update Customer Experience Metrics",
	"Update Departments",
	"Update Device Enrollment Program Instances",
	"Update Device Name Patterns",
	"Update Directory Bindings",
	"Update Disk Encryption Configurations",
	"Update Disk Encryption Institutional Configurations",
	"Update Distribution Points",
	"Update Dock Items",
	"Update eBooks",
	"Update Education Settings",
	"Update Engage Settings",
	"Update Enrollment Customizations",
	"Update Enrollment Profiles",
	"Update File Attachments",
	"Update GSX Connection",
	"Update Infrastructure Managers",
	"Update iOS Configuration Profiles",
	"Update iBeacon",
	"Update Inventory Preload Records",
	"Update Jamf Connect Deployments",
	"Update Jamf Connect Settings",
	"Update Jamf Protect Deployments",
	"Update Jamf Protect Settings",
	"Update JSON Web Token Configuration",
	"Update JSS URL",
	"Update Keystores",
	"Update Knobs",
	"Update LDAP Servers",
	"Update Licensed Software",
	"Update Limited Access Settings",
	"Update Local Admin Password Settings",
	"Update Mac Applications",
	"Update macOS Configuration Profiles",
	"Update Maintenance Pages",
	"Update Managed Software Updates",
	"Update MDM command information in Jamf Pro API",
	"Update Mobile Device App Maintenance Settings",
	"Update Mobile Device Applications",
	"Update Mobile Device Enrollment Invitations",
	"Update Mobile Device Extension Attributes",
	"Update Mobile Device Inventory Collection",
	"Update Mobile Device Managed App Configurations",
	"Update Mobile Device PreStage Enrollments",
	"Update Mobile Device Self Service",
	"Update Mobile Devices",
	"Update Network Integration",
	"Update Network Segments",
	"Update Onboarding Configuration",
	"Update Packages",
	"Update Parent App Settings",
	"Update Password Policy",
	"Update Patch External Source",
	"Update Patch Management Settings",
	"Update Patch Management Software Titles",
	"Update Patch Policies",
	"Update Peripheral Types",
	"Update Personal Device Configurations",
	"Update Personal Device Profiles",
	"Update PKI",
	"Update Policies",
	"Update Printers",
	"Update Privileges",
	"Update Provisioning Profiles",
	"Update Push Certificates",
	"Update Re-enrollment",
	"Update Remote Administration",
	"Update Removable MAC Address",
	"Update Renewal of the Built-in Certificate Authority",
	"Update Restricted Software",
	"Update Retention Policy",
	"Update SSO Settings",
	"Update Scripts",
	"Update Self Service",
	"Update Self Service Bookmarks",
	"Update Self Service Branding Configuration",
	"Update Sites",
	"Update Smart Computer Groups",
	"Update Smart Mobile Device Groups",
	"Update Smart User Groups",
	"Update SMTP Server",
	"Update Software Update Servers",
	"Update Static Computer Groups",
	"Update Static Mobile Device Groups",
	"Update Static User Groups",
	"Update Teacher App Settings",
	"Update User",
	"Update User Extension Attributes",
	"Update User-Initiated Enrollment",
	"Update VPP Assignment",
	"Update VPP Invitations",
	"Update Volume Purchasing Locations",
	"Update Webhooks",
	"View Activation Lock Bypass Code",
	"View Disk Encryption Recovery Key",
	"View Event Logs",
	"View JSS Information",
	"View License Serial Numbers",
	"View Local Admin Password",
	"View Local Admin Password Audit History",
	"View MDM command information in Jamf Pro API",
	"View Mobile Device Lost Mode Location",
	"View Recovery Lock",
}

// validateResourceApiRolesDataFields checks if a given privilege is in the list of valid privileges
// and groups privileges by category.
func validateResourceApiRolesDataFields(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	categories := make(map[string][]string)
	var nonCrudPrivileges []string
	var MDMCommands []string

	for _, priv := range validPrivileges {
		if v == priv {
			return
		}

		// Split the privilege into operation and category
		parts := strings.SplitN(priv, " ", 2)
		if len(parts) == 2 {
			operation, category := parts[0], parts[1]

			// Group CRUD privileges by category
			if operation == "Create" || operation == "Read" || operation == "Update" || operation == "Delete" {
				categories[category] = append(categories[category], priv)
			} else if operation == "Send" {
				MDMCommands = append(MDMCommands, priv)
			} else {
				nonCrudPrivileges = append(nonCrudPrivileges, priv)
			}
		} else {
			nonCrudPrivileges = append(nonCrudPrivileges, priv)
		}
	}

	var formattedPrivileges strings.Builder

	// Sort categories for consistent ordering
	sortedCategories := make([]string, 0, len(categories))
	for category := range categories {
		sortedCategories = append(sortedCategories, category)
	}
	sort.Strings(sortedCategories)

	for _, category := range sortedCategories {
		privileges := categories[category]
		// Adding a spacer with the category name
		formattedPrivileges.WriteString(fmt.Sprintf("---- Privilege Set: %s ----\n", category))
		formattedPrivileges.WriteString(fmt.Sprintf("    %s\n", strings.Join(privileges, "\n    ")))
		formattedPrivileges.WriteString("---- End ----\n\n")
	}

	if len(MDMCommands) > 0 {
		// Adding a spacer for Send MDM Commands
		formattedPrivileges.WriteString("---- MDM Commands ----\n")
		formattedPrivileges.WriteString(fmt.Sprintf("    %s\n", strings.Join(MDMCommands, "\n    ")))
		formattedPrivileges.WriteString("---- End ----\n\n")
	}

	if len(nonCrudPrivileges) > 0 {
		// Adding a spacer for non-CRUD privileges
		formattedPrivileges.WriteString("---- Other Jamf Pro Operations ----\n")
		formattedPrivileges.WriteString(fmt.Sprintf("    %s\n", strings.Join(nonCrudPrivileges, "\n    ")))
		formattedPrivileges.WriteString("---- End ----\n\n")
	}

	errs = append(errs, fmt.Errorf("%q contains an invalid privilege: %s; must be one of:\n%s", key, v, formattedPrivileges.String()))
	return
}
