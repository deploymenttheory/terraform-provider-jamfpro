# data "local_file" "site_and_computer_ids" {
#   filename = "${path.module}/../../data_sources/site_and_computer_ids.json"
# }

# resource "jamfpro_account" "account_min" {
#   name         = "tf-testing-${var.testing_id}-min-${random_id.rng.hex}"
#   enabled      = "Enabled"
#   access_level = "Full Access"
# }

# // Disabled account
# resource "jamfpro_account" "account_disabled" {
#   name         = "tf-testing-${var.testing_id}-accdis-${random_id.rng.hex}"
#   enabled      = "Disabled"
#   access_level = "Full Access"
# }

# // Site Access account
# resource "jamfpro_account" "account_site_access" {
#   name         = "tf-testing-${var.testing_id}-saccess-${random_id.rng.hex}"
#   enabled      = "Enabled"
#   access_level = "Site Access"
#   site_id      = jsondecode(data.local_file.site_and_computer_ids.content).site
# }

// Full account with all optional fields
resource "jamfpro_account" "account_max" {
  name                  = "tf-testing-${var.testing_id}-max-${random_id.rng.hex}"
  directory_user        = false
  full_name            = "Test User"
  email                = "test.user@example.com"
  enabled              = "Enabled"
  identity_server_id   = 1
  force_password_change = true
  access_level         = "Full Access"
  password             = "SecurePassword123!"
  privilege_set        = "Administrator"
  # site_id               = jsondecode(data.local_file.site_and_computer_ids.content).site

  jss_objects_privileges = [
    "Create Categories",
    "Read Categories",
    "Update Categories",
    "Delete Categories",
    "Read Directory Bindings",
    "Read Dock Items",
    "Read Packages",
    "Read Printers",
    "Read Scripts"
  ]

  jss_settings_privileges = [
    "Read JSS Settings",
    "Update JSS Settings",
    "Read Activation Code"
  ]

  jss_actions_privileges = [
    "Read JSS Actions",
    "Execute JSS Actions"
  ]

  # casper_admin_privileges = [
  #   "Use Casper Admin",
  #   "Save With Casper Admin"
  # ]

  # casper_remote_privileges = [
  #   "Use Casper Remote",
  #   "Save With Casper Remote"
  # ]

  casper_imaging_privileges = [
    "Use Casper Imaging",
    "Save With Casper Imaging"
  ]

  recon_privileges = [
    "Read Recon",
    "Update Recon"
  ]
}


# // Group Access account with custom privileges
# resource "jamfpro_account" "account_group_access" {
#   name         = "tf-testing-${var.testing_id}-gaccess-${random_id.rng.hex}"
#   enabled      = "Enabled"
#   access_level = "Group Access"
#   privilege_set = "Custom"

#   jss_objects_privileges = [
#     "Read Categories",
#     "Read Directory Bindings",
#     "Read Dock Items"
#   ]

#   jss_settings_privileges = [
#     "Read JSS Settings"
#   ]
# }



# // ========================================================================== //
# // Multiple accounts

# // Multiple minimal accounts
# resource "jamfpro_account" "account_multiple_min" {
#   count        = 10
#   name         = "tf-testing-${var.testing_id}-min-${count.index}-${random_id.rng.hex}"
#   enabled      = "Enabled"
#   access_level = "Full Access"
# }

# // Multiple full accounts
# resource "jamfpro_account" "account_multiple_max" {
#   count                 = 10
#   name                  = "tf-testing--${var.testing_id}-max-${count.index}-${random_id.rng.hex}"
#   directory_user        = false
#   full_name            = "Test User ${count.index}"
#   email                = "test.user${count.index}@example.com"
#   enabled              = "Enabled"
#   identity_server_id   = 1
#   force_password_change = true
#   access_level         = "Full Access"
#   password             = "SecurePassword123!"
#   privilege_set        = "Administrator"
#   site_id              = 1

#   jss_objects_privileges = [
#     "Create Categories",
#     "Read Categories",
#     "Update Categories",
#     "Delete Categories"
#   ]

#   jss_settings_privileges = [
#     "Read JSS Settings",
#     "Update JSS Settings"
#   ]

#   jss_actions_privileges = [
#     "Read JSS Actions",
#     "Execute JSS Actions"
#   ]

#   casper_admin_privileges = [
#     "Use Casper Admin",
#     "Save With Casper Admin"
#   ]

#   casper_remote_privileges = [
#     "Use Casper Remote",
#     "Save With Casper Remote"
#   ]

#   casper_imaging_privileges = [
#     "Use Casper Imaging",
#     "Save With Casper Imaging"
#   ]

#   recon_privileges = [
#     "Read Recon",
#     "Update Recon"
#   ]
# }

# // Multiple site access accounts
# resource "jamfpro_account" "account_multiple_site_access" {
#   count        = 10
#   name         = "tf-testing-${var.testing_id}-saccess-${count.index}-${random_id.rng.hex}"
#   enabled      = "Enabled"
#   access_level = "Site Access"
#   site_id      = 1
# }

# // Multiple group access accounts
# resource "jamfpro_account" "account_multiple_group_access" {
#   count         = 10
#   name          = "tf-testing--${var.testing_id}-gaccess-${count.index}-${random_id.rng.hex}"
#   enabled       = "Enabled"
#   access_level  = "Group Access"
#   privilege_set = "Custom"

#   jss_objects_privileges = [
#     "Read Categories",
#     "Read Directory Bindings",
#     "Read Dock Items"
#   ]

#   jss_settings_privileges = [
#     "Read JSS Settings"
#   ]
# }
