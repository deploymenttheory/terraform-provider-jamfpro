// ========================================================================== //
// User-intiated settings
// ========================================================================== //

// ========================================================================== //
// Configure user-initiated settings

resource "jamfpro_user_initiated_enrollment_settings" "settings" {

  // General
  restrict_reenrollment_to_authorized_users_only  = false
  skip_certificate_installation_during_enrollment = true
  // third_party_signing_certificate = {}

  // Messaging
  messaging {
    language_code                                   = "en"
    language_name                                   = "english"
    page_title                                      = "Welcome to Device Enrollment"
    username_text                                   = "Username"
    password_text                                   = "Password"
    login_button_text                               = "Log In"
    device_ownership_page_text                      = "Select your device type"
    personal_device_button_name                     = "Personal Device"
    institutional_ownership_button_name             = "Company Device"
    personal_device_management_description          = "Your personal device will be managed with minimal restrictions"
    institutional_device_management_description     = "This company device will be fully managed"
    enroll_device_button_name                       = "Enroll Device"
    eula_personal_devices                           = "By enrolling, you agree to allow management of your personal device"
    eula_institutional_devices                      = "This device is subject to management policies as per company guidelines"
    accept_button_text                              = "Accept"
    site_selection_text                             = "Select your site"
    ca_certificate_installation_text                = "Install CA Certificate"
    ca_certificate_name                             = "Company Root CA"
    ca_certificate_description                      = "This certificate allows secure communication with company servers"
    ca_certificate_install_button_name              = "Install CA"
    institutional_mdm_profile_installation_text     = "Install Management Profile"
    institutional_mdm_profile_name                  = "Company MDM Profile"
    institutional_mdm_profile_description           = "This profile allows management of your company device"
    institutional_mdm_profile_pending_text          = "Installing MDM profile..."
    institutional_mdm_profile_install_button_name   = "Install"
    personal_mdm_profile_installation_text          = "Install Personal Device Profile"
    personal_mdm_profile_name                       = "Personal Device Profile"
    personal_mdm_profile_description                = "Limited management profile for personal devices"
    personal_mdm_profile_install_button_name        = "Install Profile"
    user_enrollment_mdm_profile_installation_text   = "Install User Enrollment Profile"
    user_enrollment_mdm_profile_name                = "User Enrollment Profile"
    user_enrollment_mdm_profile_description         = "Profile for user-based enrollment"
    user_enrollment_mdm_profile_install_button_name = "Install"
    quickadd_package_installation_text              = "Install Management Software"
    quickadd_package_name                           = "Company MDM Agent"
    quickadd_package_progress_text                  = "Installing management software..."
    quickadd_package_install_button_name            = "Install Software"
    enrollment_complete_text                        = "Enrollment Complete! Your device is now managed."
    enrollment_failed_text                          = "Enrollment Failed. Please try again."
    try_again_button_name                           = "Try Again"
    view_enrollment_status_button_name              = "Check Status"
    view_enrollment_status_text                     = "Check your enrollment status"
    log_out_button_name                             = "Log Out"
  }

  // Computers
  user_initiated_enrollment_for_computers {
    enable_user_initiated_enrollment_for_computers = true
    managed_local_administrator_account {
      create_managed_local_administrator_account                    = true
      management_account_username                                   = "jamfadmin"
      hide_managed_local_administrator_account                      = true
      allow_ssh_access_for_managed_local_administrator_account_only = false
    }
    ensure_ssh_is_enabled         = false
    launch_self_service_when_done = true
    // quickadd_package = {}
    account_driven_device_enrollment = false
  }

  // Devices
  user_initiated_enrollment_for_devices {
    profile_driven_enrollment_via_url {
      enable_for_institutionally_owned_devices = true
      enable_for_personally_owned_devices      = false
      personal_device_enrollment_type          = "USERENROLLMENT"
    }
    account_driven_user_enrollment {
      enable_for_personally_owned_mobile_devices     = false
      enable_for_personally_owned_vision_pro_devices = false
    }
    account_driven_device_enrollment {
      enable_for_institutionally_owned_mobile_devices = false
      enable_for_personally_owned_vision_pro_devices  = false
    }
  }
}