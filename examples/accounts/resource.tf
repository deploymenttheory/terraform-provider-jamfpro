// example of account with custom privileges and full access
resource "jamfpro_account" "jamf_pro_account_001" {
  name                   = "tf-localtest-account-custom-privileges-full-access"
  directory_user         = false
  full_name              = "micky mouse"
  password               = "mySecretThing10" // Password must be at least 10 characters long. password not stored in state
  email                  = "exampleEmailthing@domain.com"
  enabled                = "Enabled"
  force_password_change  = true
  access_level           = "Full Access" // Full Access / Site Access / Group Access
  privilege_set          = "Custom" // "Administrator", "Auditor", "Enrollment Only", "Custom"

  jss_objects_privileges   = [
    "Create Categories",
    "Read Categories",
    "Update Categories",
    "Delete Categories",
    "Create Directory Bindings",
    "Read Directory Bindings",
    "Update Directory Bindings",
    "Delete Directory Bindings",
    "Create Dock Items",
    "Read Dock Items",
    "Update Dock Items",
    "Delete Dock Items",
    "Create Packages",
    "Read Packages",
    "Update Packages",
    "Delete Packages",
    "Create Printers",
    "Read Printers",
    "Update Printers",
    "Delete Printers",
    "Create Scripts",
    "Read Scripts",
    "Update Scripts",
    "Delete Scripts",
	]
  jss_settings_privileges  = ["Read SSO Settings", "Update User-Initiated Enrollment"]
  jss_actions_privileges   = ["Flush Policy Logs"]
  casper_admin_privileges  = ["Use Casper Admin", "Save With Casper Admin"]

}

# // example of account with group privileges and group access
resource "jamfpro_account" "jamf_pro_account_002" {
  name                   = "tf-localtest-account-group-privileges-group-access"
  directory_user         = false
  full_name              = "jonny bravo"
  password               = "mySecretThing10" // Password must be at least 10 characters long. password not stored in state
  email                  = "exampleEmailthing@domain.com"
  enabled                = "Enabled"
  force_password_change  = true
  access_level           = "Group Access" // Full Access / Site Access / Group Access
  privilege_set          = "Custom" // Must be "Custom" for group access
}

// example of account with group privileges and site access
resource "jamfpro_account" "jamf_pro_account_003" {
  name                   = "tf-localtest-account-custom-privileges-site-access"
  directory_user         = false
  full_name              = "donald duck"
  password               = "securePassword" // Password must be at least 10 characters long. password not stored in state
  email                  = "example@domain.com"
  enabled                = "Enabled"
  force_password_change  = true
  access_level           = "Site Access" // Full Access / Site Access / Group Access
  privilege_set          = "Custom"

  site_id = 1

  jss_objects_privileges   = ["Create Advanced Computer Searches",
	"Read Advanced Computer Searches",
	"Update Advanced Computer Searches",
	"Delete Advanced Computer Searches",
	"Create Advanced Mobile Device Searches",
	"Read Advanced Mobile Device Searches",
	"Update Advanced Mobile Device Searches",
	"Delete Advanced Mobile Device Searches",
	]

  jss_actions_privileges   = ["Allow User to Enroll", "Assign Users to Computers", "Assign Users to Mobile Devices", "Change Password"]

}