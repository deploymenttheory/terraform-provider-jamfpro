resource "jamfpro_computer_prestage_enrollment" "example" {
  display_name                            = "Example Prestage Enrollment"
  mandatory                               = true
  mdm_removable                           = false
  support_phone_number                    = null
  support_email_address                   = null
  department                              = null
  default_prestage                        = false
  enrollment_site_id                      = null
  keep_existing_site_membership           = false
  keep_existing_location_information      = false
  require_authentication                  = false
  authentication_prompt                   = null
  prevent_activation_lock                 = false
  enable_device_based_activation_lock     = false
  device_enrollment_program_instance_id   = "00000000000000000000000000000000"
  skip_setup_items {
    biometric           = false
    terms_of_address    = false
    file_vault          = false
    icloud_diagnostics  = false
    diagnostics         = false
    accessibility       = false
    apple_id            = false
    screen_time         = false
    siri                = false
    display_tone        = false
    restore             = false
    appearance          = false
    privacy             = false
    payment             = false
    registration        = false
    tos                 = false
    icloud_storage      = false
    location            = false
  }
  location_information {
    id             = null
    version_lock   = null
    username       = null
    realname       = null
    phone          = null
    email          = null
    room           = null
    position       = null
    department_id  = null
    building_id    = null
  }
  purchasing_information {
    id                  = "1"
    version_lock        = 1
    leased              = false
    purchased           = true
    apple_care_id       = null
    po_number           = "PO12345"
    vendor              = null
    purchase_price      = null
    life_expectancy     = null
    purchasing_account  = null
    purchasing_contact  = "John Doe"
    lease_date          = null
    po_date             = null
    warranty_date       = null
  }
  anchor_certificates                     = null
  enrollment_customization_id             = "0"
  language                                = null
  region                                  = null
  auto_advance_setup                      = false
  install_profiles_during_setup           = false
  prestage_installed_profile_ids          = null
  custom_package_ids                      = null
  custom_package_distribution_point_id    = null
  enable_recovery_lock                    = null
  recovery_lock_password_type             = null
  recovery_lock_password                  = null
  rotate_recovery_lock_password           = null
  site_id                                 = null
  version_lock                            = 1
  account_settings {
    version_lock                                  = 1
    payload_configured                            = null
    local_admin_account_enabled                   = null
    admin_username                                = null
    admin_password                                = null
    hidden_admin_account                          = null
    local_user_managed                            = null
    user_account_type                             = null
    prefill_primary_account_info_feature_enabled  = null
    prefill_type                                  = "UNKNOWN"
    prefill_account_full_name                     = null
    prefill_account_user_name                     = null
    prevent_prefill_info_from_modification        = false
  }
}