# Configure Jamf Pro User Initiated Enrollment Settings
resource "jamfpro_user_initiated_enrollment_settings" "jamfpro_uie_settings" {
  # General Settings
  restrict_reenrollment_to_authorized_users_only  = true
  skip_certificate_installation_during_enrollment = true

  # Flush Settings
  flush_location_information         = true
  flush_location_history_information = true
  flush_policy_history               = true
  flush_extension_attributes         = true
  flush_software_update_plans        = true
  flush_mdm_commands_on_reenroll     = "DELETE_EVERYTHING"

  # Third-party MDM signing certificate
  third_party_signing_certificate {
    enabled           = true
    filename          = "my_mdm_signing_cert.p12"
    identity_keystore = filebase64("${path.module}/cert/path/test_certificate.p12")
    keystore_password = "your-cert-password"
  }

  # Computer Enrollment Settings
  user_initiated_enrollment_for_computers {
    enable_user_initiated_enrollment_for_computers = true
    ensure_ssh_is_enabled                          = true
    launch_self_service_when_done                  = true
    account_driven_device_enrollment               = true

    # Managed Local Administrator Account
    managed_local_administrator_account {
      create_managed_local_administrator_account                    = true
      management_account_username                                   = "jamfadmin"
      hide_managed_local_administrator_account                      = true
      allow_ssh_access_for_managed_local_administrator_account_only = true
    }

    # QuickAdd Package Signing
    quickadd_package {
      sign_quickadd_package = true
      filename              = "quickadd_signing_cert.p12"
      identity_keystore     = filebase64("${path.module}/cert/path/test_certificate.p12")
      keystore_password     = "your-cert-password"
    }
  }

  # Mobile Device Enrollment Settings
  user_initiated_enrollment_for_devices {
    # Profile-Driven Enrollment
    profile_driven_enrollment_via_url {
      enable_for_institutionally_owned_devices = true
      enable_for_personally_owned_devices      = true
      personal_device_enrollment_type          = "USERENROLLMENT"
    }

    # Account-Driven User Enrollment
    account_driven_user_enrollment {
      enable_for_personally_owned_mobile_devices     = true
      enable_for_personally_owned_vision_pro_devices = true
    }

    # Account-Driven Device Enrollment
    account_driven_device_enrollment {
      enable_for_institutionally_owned_mobile_devices = true
      enable_for_personally_owned_vision_pro_devices  = false
    }
  }

  # Enrollment Messaging - English (this block is always required. It's built into the jamf gui)
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

  # Enrollment Messaging - French (All additional languages are optional)
  messaging {
    language_code                                   = "fr"
    language_name                                   = "french"
    page_title                                      = "Welcome to Device Enrollment"
    username_text                                   = "Username"
    password_text                                   = "Password"
    login_button_text                               = "Log In"
    device_ownership_page_text                      = "Select device type"
    personal_device_button_name                     = "Personal Device"
    institutional_ownership_button_name             = "Company Device"
    personal_device_management_description          = "Personal device with minimal management"
    institutional_device_management_description     = "Company device with full management"
    enroll_device_button_name                       = "Enroll Device"
    eula_personal_devices                           = "Agreement for personal device management"
    eula_institutional_devices                      = "Agreement for company device management"
    accept_button_text                              = "Accept"
    site_selection_text                             = "Select site"
    ca_certificate_installation_text                = "Install Certificate"
    ca_certificate_name                             = "Root Certificate"
    ca_certificate_description                      = "Security certificate"
    ca_certificate_install_button_name              = "Install"
    institutional_mdm_profile_installation_text     = "Install Profile"
    institutional_mdm_profile_name                  = "MDM Profile"
    institutional_mdm_profile_description           = "Management profile"
    institutional_mdm_profile_pending_text          = "Installing..."
    institutional_mdm_profile_install_button_name   = "Install"
    personal_mdm_profile_installation_text          = "Install Profile"
    personal_mdm_profile_name                       = "Personal Profile"
    personal_mdm_profile_description                = "Management profile"
    personal_mdm_profile_install_button_name        = "Install"
    user_enrollment_mdm_profile_installation_text   = "Install Profile"
    user_enrollment_mdm_profile_name                = "User Profile"
    user_enrollment_mdm_profile_description         = "User enrollment profile"
    user_enrollment_mdm_profile_install_button_name = "Install"
    quickadd_package_installation_text              = "Install Software"
    quickadd_package_name                           = "Management Agent"
    quickadd_package_progress_text                  = "Installing..."
    quickadd_package_install_button_name            = "Install"
    enrollment_complete_text                        = "Enrollment Complete"
    enrollment_failed_text                          = "Enrollment Failed"
    try_again_button_name                           = "Try Again"
    view_enrollment_status_button_name              = "Check Status"
    view_enrollment_status_text                     = "View enrollment status"
    log_out_button_name                             = "Log Out"
  }

  //Directory Service Group Enrollment Settings
  directory_service_group_enrollment_settings {
    allow_group_to_enroll_institutionally_owned_devices                      = false
    allow_group_to_enroll_personally_owned_devices                           = false
    allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade = false
    require_eula                                                             = true
    ldap_server_id                                                           = data.jamfpro_cloud_idp.by_name.id // LDAP or cloud idp
    directory_service_group_name                                             = "Test M365 account"
    directory_service_group_id                                               = "27230740-e063-4931-be75-f5e9b2e4ad53"
    site_id                                                                  = "-1"
  }

  directory_service_group_enrollment_settings {
    allow_group_to_enroll_institutionally_owned_devices                      = true
    allow_group_to_enroll_personally_owned_devices                           = false
    allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade = false
    require_eula                                                             = true
    ldap_server_id                                                           = data.jamfpro_cloud_idp.by_name.id // LDAP or cloud idp
    directory_service_group_name                                             = "Test Team"
    directory_service_group_id                                               = "a2327741-8784-40bf-aa3b-7fb979ea8658"
    site_id                                                                  = "-1"
  }
}
