resource "jamfpro_user_initiated_enrollment_settings" "user_initiated_enrollment_settings" {
  restrict_reenrollment_to_authorized_users_only  = false
  skip_certificate_installation_during_enrollment = true

  user_initiated_enrollment_for_computers {
    enable_user_initiated_enrollment_for_computers = true
    ensure_ssh_is_enabled                          = false
    launch_self_service_when_done                  = true
    account_driven_device_enrollment               = false

    managed_local_administrator_account {
      create_managed_local_administrator_account                    = true
      management_account_username                                   = "lapsadmin"
      hide_managed_local_administrator_account                      = true
      allow_ssh_access_for_managed_local_administrator_account_only = false
    }
  }

  user_initiated_enrollment_for_devices {
    profile_driven_enrollment_via_url {
      enable_for_institutionally_owned_devices = true
      enable_for_personally_owned_devices      = false
    }
    account_driven_user_enrollment {
      enable_for_personally_owned_mobile_devices     = true
      enable_for_personally_owned_vision_pro_devices = false
      enable_maid_username_merge                     = true
    }
  }
  messaging {
    language_code                                   = "en"
    language_name                                   = "English"
    page_title                                      = "Enroll Your Device"
    username_text                                   = "Username"
    password_text                                   = "Password"
    login_button_text                               = "Log In"
    personal_device_button_name                     = "Personally Owned"
    institutional_ownership_button_name             = "Institutionally Owned"
    enroll_device_button_name                       = "Continue"
    accept_button_text                              = "Accept"
    ca_certificate_name                             = "CA Certificate"
    ca_certificate_install_button_name              = "Continue"
    institutional_mdm_profile_name                  = "MDM Profile"
    institutional_mdm_profile_install_button_name   = "Install"
    personal_mdm_profile_name                       = "MDM Profile"
    personal_mdm_profile_install_button_name        = "Enroll"
    user_enrollment_mdm_profile_name                = "MDM Profile"
    user_enrollment_mdm_profile_install_button_name = "Continue"
    quickadd_package_name                           = "QuickAdd.pkg"
    quickadd_package_install_button_name            = "Download"
    log_out_button_name                             = "Log Out"
  }
}
