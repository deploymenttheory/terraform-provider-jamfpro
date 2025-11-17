resource "jamfpro_computer_prestage_enrollment" "minimum_example" {
  display_name                          = "jamfpro-sdk-example-computerPrestageMinimum-config"
  mandatory                             = true
  mdm_removable                         = true
  support_phone_number                  = "111-222-3333"
  support_email_address                 = "email@company.com"
  department                            = "department name"
  default_prestage                      = false
  keep_existing_site_membership         = false
  keep_existing_location_information    = false
  require_authentication                = false
  authentication_prompt                 = "hello welcome to your enterprise managed macOS device"
  prevent_activation_lock               = false
  enable_device_based_activation_lock   = false
  device_enrollment_program_instance_id = "1"
  anchor_certificates                   = []
  auto_advance_setup                    = false
  install_profiles_during_setup         = false
  prestage_installed_profile_ids        = []
  custom_package_ids                    = []
  custom_package_distribution_point_id  = "-1"
  enable_recovery_lock                  = false
  recovery_lock_password_type           = "MANUAL"
  recovery_lock_password                = ""
  rotate_recovery_lock_password         = false

  location_information {
    username      = ""
    realname      = ""
    phone         = ""
    email         = ""
    room          = ""
    position      = ""
    department_id = "-1"
    building_id   = "-1"
  }

  purchasing_information {
    leased             = false
    purchased          = true
    apple_care_id      = ""
    po_number          = ""
    vendor             = ""
    purchase_price     = ""
    life_expectancy    = 0
    purchasing_account = ""
    purchasing_contact = ""
  }
}
