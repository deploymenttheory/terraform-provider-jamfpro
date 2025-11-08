resource "jamfpro_mobile_device_prestage_enrollment" "test_prestage_enrollment" {
  display_name                            = "iOS Device Prestage ${var.testing_id}"
  mandatory                               = true
  mdm_removable                           = false
  support_phone_number                    = "+44-1234-567890"
  support_email_address                   = "support@mycompany.org"
  department                              = "IT Services"
  default_prestage                        = false
  keep_existing_site_membership           = false
  keep_existing_location_information      = false
  require_authentication                  = false
  prevent_activation_lock                 = true
  enable_device_based_activation_lock     = false
  device_enrollment_program_instance_id   = "1"
  auto_advance_setup                      = false
  allow_pairing                           = false
  multi_user                              = false
  supervised                              = true
  maximum_shared_accounts                 = 10
  configure_device_before_setup_assistant = true
  send_timezone                           = false
  use_storage_quota_size                  = false

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

  names {
    assign_names_using = "Default Names"
  }
}
