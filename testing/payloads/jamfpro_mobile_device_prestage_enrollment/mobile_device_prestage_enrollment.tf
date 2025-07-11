resource "jamfpro_mobile_device_prestage_enrollment" "example_prestage" {
  display_name                            = "iOS Device Prestage"
  mandatory                               = true
  mdm_removable                           = false
  support_phone_number                    = "+44-1234-567890"
  support_email_address                   = "support@mycompany.org"
  department                              = "IT Department"
  default_prestage                        = false
  keep_existing_site_membership           = false
  keep_existing_location_information      = false
  require_authentication                  = false
  authentication_prompt                   = ""
  prevent_activation_lock                 = true
  enable_device_based_activation_lock     = false
  device_enrollment_program_instance_id   = "1"
  enrollment_customization_id             = "0"
  language                                = ""
  region                                  = ""
  auto_advance_setup                      = false
  allow_pairing                           = false
  multi_user                              = false
  supervised                              = true
  maximum_shared_accounts                 = 10
  configure_device_before_setup_assistant = true
  anchor_certificates                     = []

  skip_setup_items {
    location                = false
    privacy                 = true
    biometric               = false
    software_update         = false
    diagnostics             = false
    imessage_and_facetime   = true
    intelligence            = true
    tv_room                 = true
    passcode                = false
    sim_setup               = true
    screen_time             = true
    restore_completed       = true
    tv_provider_sign_in     = true
    siri                    = true
    restore                 = true
    screen_saver            = true
    home_button_sensitivity = true
    cloud_storage           = true
    action_button           = true
    transfer_data           = true
    enable_lockdown_mode    = true
    zoom                    = true
    preferred_language      = true
    voice_selection         = true
    tv_home_screen_sync     = true
    safety                  = true
    terms_of_address        = true
    express_language        = true
    camera_button           = true
    apple_id                = true
    display_tone            = true
    watch_migration         = true
    update_completed        = false
    appearance              = true
    android                 = true
    payment                 = true
    onboarding              = true
    tos                     = false
    welcome                 = true
    tap_to_setup            = true
    safety_and_handling     = true
  }

  location_information {
    username      = ""
    realname      = ""
    phone         = ""
    email         = ""
    room          = ""
    position      = ""
    department_id = "1"
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

  names {
    assign_names_using       = "Serial Numbers"
    device_name_prefix       = "ABC-"
    device_name_suffix       = ""
    manage_names             = true
    device_naming_configured = true
  }

  timezone                     = "UTC"
  storage_quota_size_megabytes = 1024

  prestage_minimum_os_target_version_type_ios  = "MINIMUM_OS_SPECIFIC_VERSION"
  minimum_os_specific_version_ios              = "18.4.1"
  prestage_minimum_os_target_version_type_ipad = "NO_ENFORCEMENT"
}
