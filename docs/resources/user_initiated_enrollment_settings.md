---
page_title: "jamfpro_user_initiated_enrollment_settings"
description: |-
  
---

# jamfpro_user_initiated_enrollment_settings (Resource)


## Example Usage
```terraform
# Configure Jamf Pro User Initiated Enrollment Settings
resource "jamfpro_user_initiated_enrollment_settings" "jamfpro_uie_settings" {
  # General Settings
  restrict_reenrollment_to_authorized_users_only  = true
  skip_certificate_installation_during_enrollment = true

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

  # Directory Service Group Enrollment Settings

  # All Directory Service Users (Required when managing other directory_service_group_enrollment_settings blocks)
  # This is a special group that always exists as ID 1 and cannot be created or destroyed
  directory_service_group_enrollment_settings {
    id                                                                       = "1" // Must be provided 
    allow_group_to_enroll_institutionally_owned_devices                      = true
    allow_group_to_enroll_personally_owned_devices                           = false
    allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade = true
    require_eula                                                             = true
    ldap_server_id                                                           = "-1"
    directory_service_group_name                                             = "All Directory Service Users"
    directory_service_group_id                                               = "-1"
    site_id                                                                  = "-1"
  }

  # Other groups
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `directory_service_group_enrollment_settings` (Block Set) Directory Service groups to configure enrollment access for (see [below for nested schema](#nestedblock--directory_service_group_enrollment_settings))
- `messaging` (Block Set) Localized text configuration for enrollment pages and messages (see [below for nested schema](#nestedblock--messaging))
- `restrict_reenrollment_to_authorized_users_only` (Boolean) Maps to payload field 'restrictReenrollment'.Restrict re-enrollment to authorized users only by only allowing re-enrollment of mobile devices and computers if the user has the applicable privilege (“Mobile Devices” or “Computers”) or their username matches the Username field in User and Location information.
- `skip_certificate_installation_during_enrollment` (Boolean) Maps to payload field 'installSingleProfile'. Certificate installation step is skipped during enrollment if your environment has an SSL certificate that was obtained from an internal CA or a trusted third-party vendor.
- `third_party_signing_certificate` (Block List, Max: 1) Third-party signing certificate configuration to ensure that the certificate signs configuration profiles sent to computers and mobile devices, and appears as verified to users during user-initiated enrollment. Maps to 'mdmSigningCertificate'. (see [below for nested schema](#nestedblock--third_party_signing_certificate))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `user_initiated_enrollment_for_computers` (Block Set, Max: 1) Configuration for user initiated enrollment for computers (see [below for nested schema](#nestedblock--user_initiated_enrollment_for_computers))
- `user_initiated_enrollment_for_devices` (Block Set, Max: 1) Allow users to enroll iOS/iPadOS devicesby going to https://JAMF_PRO_URL.jamfcloud.com/enroll (hosted in Jamf Cloud) or https://JAMF_PRO_URL.com:8443/enroll (hosted on-premise) (see [below for nested schema](#nestedblock--user_initiated_enrollment_for_devices))

### Read-Only

- `id` (String) The ID of this resource.
- `mdm_signing_certificate_details` (Set of Object) Details of the MDM signing certificate (see [below for nested schema](#nestedatt--mdm_signing_certificate_details))

<a id="nestedblock--directory_service_group_enrollment_settings"></a>
### Nested Schema for `directory_service_group_enrollment_settings`

Optional:

- `allow_group_to_enroll_institutionally_owned_devices` (Boolean) Whether directory group can perform user initiated device enrollment for iOS/iPadOS institutionally owned devices. Maps to request field 'enterpriseEnrollmentEnabled'
- `allow_group_to_enroll_personal_and_institutionally_owned_devices_via_ade` (Boolean) Whether directory group can perform user initiated device enrollment for iOS/iPadOS by signing in to their device using a Managed Apple ID with ade. Maps to request field 'accountDrivenUserEnrollmentEnabled'
- `allow_group_to_enroll_personally_owned_devices` (Boolean) Whether directory group can perform user initiated device enrollment for iOS/iPadOS personally owned devices. Maps to request field 'personalEnrollmentEnabled'
- `directory_service_group_id` (String) id of the Directory Service group to configure enrollment access for. Maps to request field 'groupId'
- `directory_service_group_name` (String) Name of the Directory Service group to configure enrollment access for. Maps to request field 'name'
- `id` (String) The unique identifier of the Directory Service group enrollment setting id. Only ID '1' is valid as it represents the built-in 'All Directory Service Users' group.
- `ldap_server_id` (String) The unique identifier of the Directory Service group id.
- `require_eula` (Boolean) Upon enrollment is the eula required to be accepted
- `site_id` (String) ID of the Site to allow this LDAP user group to select during enrollment. Set to '-1' for no site.


<a id="nestedblock--messaging"></a>
### Nested Schema for `messaging`

Required:

- `accept_button_text` (String) Name for the button that users tap/click to accept the End User License Agreement. Maps to request field 'eulaButton'
- `ca_certificate_install_button_name` (String) Name for the button that users tap to install the CA certificate. Maps to request field 'certificateButton'
- `ca_certificate_name` (String) Name to display for the CA certificate during enrollment. Maps to request field 'certificateProfileName'
- `enroll_device_button_name` (String) Name for the button that users tap to start enrollment. Maps to request field 'deviceClassButton'
- `institutional_mdm_profile_install_button_name` (String) Name for the button that users tap to install the MDM profile. Maps to request field 'enterpriseButton'
- `institutional_mdm_profile_name` (String) Name to display for the MDM profile during enrollment of a institutionally owned device. Maps to request field 'enterpriseProfileName'
- `institutional_ownership_button_name` (String) Name for the button that users tap to enroll an institutionally owned device. Maps to request field 'deviceClassEnterprise'
- `language_code` (String) Two letter ISO 639-1 Language 'set 1' code for this messaging configuration language (e.g., 'en') for the 'English' language. Ref. https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes. Maps to request field 'languageCode'
- `language_name` (String) The ISO 639 language name for this message configuration. (e.g., 'English'). Maps to request field 'name'
- `log_out_button_name` (String) Name for the button that users tap/click to log out. Maps to request field 'logoutButton'
- `login_button_text` (String) Name for the button that users tap/click to log in. Maps to request field 'loginButton'
- `page_title` (String) Title to display on all enrollment pages. Maps to request field 'title'
- `password_text` (String) Text to display for the password field on the login page during enrollment. Maps to request field 'password'
- `personal_device_button_name` (String) Name for the button that users tap to enroll a personally owned device. Maps to request field 'deviceClassPersonal'
- `personal_mdm_profile_install_button_name` (String) Name for the button that users tap to install the MDM profile. Maps to request field 'personalButton'
- `personal_mdm_profile_name` (String) Name to display for the MDM profile during enrollment of a personally owned device. Maps to request field 'personalProfileName'
- `quickadd_package_install_button_name` (String) Name for the button that users tap to install the QuickAdd package. Maps to request field 'quickAddButton'
- `quickadd_package_name` (String) Name to display for the QuickAdd package during enrollment. Maps to request field 'quickAddName'
- `user_enrollment_mdm_profile_install_button_name` (String) Name for the button that users tap to install the MDM profile. Maps to request field 'userEnrollmentButton'
- `user_enrollment_mdm_profile_name` (String) Name to display for the MDM profile. Maps to request field 'userEnrollmentProfileName'
- `username_text` (String) Text to display for the username field on the login page during enrollment. Maps to request field 'username'

Optional:

- `ca_certificate_description` (String) Description to display for the CA certificate during enrollment. Maps to request field 'certificateProfileDescription'
- `ca_certificate_installation_text` (String) Text to display when installing the CA certificate during enrollment. Maps to request field 'certificateText'
- `device_ownership_page_text` (String) Text to display during enrollment that prompts the user to specify the device ownership type. Maps to request field 'deviceClassDescription'
- `enrollment_complete_text` (String) Text to display when enrollment is complete. Maps to request field 'completeMessage'
- `enrollment_failed_text` (String) Text to display when enrollment fails. Maps to request field 'failedMessage'
- `eula_institutional_devices` (String) End User License Agreement to display during enrollment of institutionally owned devices and computers. Maps to request field 'personalEula'
- `eula_personal_devices` (String) End User License Agreement to display during enrollment of personally owned devices. Maps to request field 'enterpriseEula'
- `institutional_device_management_description` (String) Description to display for institutional device management when users enroll an institutionally owned device. Maps to request field 'deviceClassEnterpriseDescription'
- `institutional_mdm_profile_description` (String) Description to display for the MDM profile during enrollment of a institutionally owned device. Maps to request field 'enterpriseProfileDescription'
- `institutional_mdm_profile_installation_text` (String) Text to display when installing the MDM profile during enrollment of a institutionally owned device. Maps to request field 'enterpriseText'
- `institutional_mdm_profile_pending_text` (String) Text to display when the user is installing the MDM profile on their computer. Maps to request field 'enterprisePending'
- `login_page_text` (String) Text to display below the title on the login page during enrollment. Maps to request field 'loginDescription'
- `personal_device_management_description` (String) Description to display for personal device management when users enroll a personally owned device. Maps to request field 'deviceClassPersonalDescription'
- `personal_mdm_profile_description` (String) Description to display for the MDM profile during enrollment of a personally owned device. Maps to request field 'personalProfileDescription'
- `personal_mdm_profile_installation_text` (String) Text to display when installing the MDM profile during enrollment of a personally owned device. Maps to request field 'personalText'
- `quickadd_package_installation_text` (String) Text to display when installing the QuickAdd package during enrollment. Maps to request field 'quickAddText'
- `quickadd_package_progress_text` (String) Text to display when the QuickAdd package is downloading. Maps to request field 'quickAddPending'
- `site_selection_text` (String) Text to display that prompts the user to select a site if the user has more than one site to choose from during enrollment. Maps to request field 'siteDescription'
- `try_again_button_name` (String) Name for the button that users tap/click to try enrolling again. Maps to request field 'tryAgainButton'
- `user_enrollment_mdm_profile_description` (String) Description to display for the MDM profile. Maps to request field 'userEnrollmentProfileDescription'
- `user_enrollment_mdm_profile_installation_text` (String) Text to display when prompting to install the MDM profile. Maps to request field 'userEnrollmentText'
- `view_enrollment_status_button_name` (String) Name for the button that users tap to view the enrollment status for the device. Maps to request field 'checkNowButton'
- `view_enrollment_status_text` (String) Text to display during enrollment that prompts the user to view the enrollment status for the device. Maps to request field 'checkEnrollmentMessage'


<a id="nestedblock--third_party_signing_certificate"></a>
### Nested Schema for `third_party_signing_certificate`

Required:

- `enabled` (Boolean) Maps to payload field 'signingMdmProfileEnabled'. Whether to use a third-party signing certificate.
- `filename` (String) Name of the certificate file in .p12 format. e.g 'my_test_certificate.p12'.
- `identity_keystore` (String, Sensitive) Base64-encoded certificate in .p12 format.
- `keystore_password` (String, Sensitive) Password for the certificate keystore


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)


<a id="nestedblock--user_initiated_enrollment_for_computers"></a>
### Nested Schema for `user_initiated_enrollment_for_computers`

Required:

- `account_driven_device_enrollment` (Boolean) Account-Driven Device Enrollment enablement for institutionally owned computers. Maps to request field 'accountDrivenDeviceMacosEnrollmentEnabled'
- `enable_user_initiated_enrollment_for_computers` (Boolean) Allow users to enroll computers by going to https://JAMF_PRO_URL.jamfcloud.com/enroll (hosted in Jamf Cloud) or https://JAMF_PRO_URL.com:8443/enroll (hosted on-premise). Maps to request field 'macOsEnterpriseEnrollmentEnabled'
- `ensure_ssh_is_enabled` (Boolean) Enable SSH (Remote Login) on computers that have it disabled. Maps to request field 'ensureSshRunning'
- `launch_self_service_when_done` (Boolean) Ensure that computers launch Self Service immediately after they are enrolled. Maps to request field 'launchSelfService'

Optional:

- `managed_local_administrator_account` (Block Set, Max: 1) Configuration for the managed local administrator account (see [below for nested schema](#nestedblock--user_initiated_enrollment_for_computers--managed_local_administrator_account))
- `quickadd_package` (Block List, Max: 1) Third-party signing certificate configuration to ensure that the certificate signs configuration profiles sent to computers and mobile devices, and appears as verified to users during user-initiated enrollment.Maps to 'developerCertificateIdentity'. (see [below for nested schema](#nestedblock--user_initiated_enrollment_for_computers--quickadd_package))

<a id="nestedblock--user_initiated_enrollment_for_computers--managed_local_administrator_account"></a>
### Nested Schema for `user_initiated_enrollment_for_computers.managed_local_administrator_account`

Required:

- `allow_ssh_access_for_managed_local_administrator_account_only` (Boolean) Make the managed local administrator account the only account that has SSH (Remote Login) access to computers. Maps to request field 'allowSshOnlyManagementAccount'
- `hide_managed_local_administrator_account` (Boolean) Hide the managed local administrator account from users. Maps to request field 'hideManagementAccount'

Optional:

- `create_managed_local_administrator_account` (Boolean) Enables managed local administrator account to be used for computers enrolled via PreStage enrollment or user-initiated enrollment. Passwords are set to a random 29 characters that automatically rotate every 1 hour. You can change the rotation settings on the https://your-instance.jamfcloud.com/computerSecurity page. Maps to request field 'createManagementAccount'
- `management_account_username` (String) Account username to be used for computers enrolled via PreStage enrollment or user-initiated enrollment. This account is created by the Jamf management framework. Maps to request field 'managementUsername'


<a id="nestedblock--user_initiated_enrollment_for_computers--quickadd_package"></a>
### Nested Schema for `user_initiated_enrollment_for_computers.quickadd_package`

Required:

- `filename` (String) Name of the certificate file in .p12 format. e.g 'my_test_certificate.p12'.
- `identity_keystore` (String, Sensitive) Base64-encoded certificate in .p12 format.
- `keystore_password` (String, Sensitive) Password for the certificate keystore
- `sign_quickadd_package` (Boolean) Ensure that the QuickAdd package is signed and appears as verified to users when enrolling via user-initiated enrollment with a QuickAdd package. Maps to request field 'signQuickAdd'



<a id="nestedblock--user_initiated_enrollment_for_devices"></a>
### Nested Schema for `user_initiated_enrollment_for_devices`

Optional:

- `account_driven_device_enrollment` (Block Set, Max: 1) Configuration for the Account-Driven Device Enrollment (see [below for nested schema](#nestedblock--user_initiated_enrollment_for_devices--account_driven_device_enrollment))
- `account_driven_user_enrollment` (Block Set, Max: 1) Configuration for the Account-Driven User Enrollment (see [below for nested schema](#nestedblock--user_initiated_enrollment_for_devices--account_driven_user_enrollment))
- `profile_driven_enrollment_via_url` (Block Set, Max: 1) Configuration for the Profile-Driven Enrollment via URL (see [below for nested schema](#nestedblock--user_initiated_enrollment_for_devices--profile_driven_enrollment_via_url))

<a id="nestedblock--user_initiated_enrollment_for_devices--account_driven_device_enrollment"></a>
### Nested Schema for `user_initiated_enrollment_for_devices.account_driven_device_enrollment`

Optional:

- `enable_for_institutionally_owned_mobile_devices` (Boolean) Whether user initiated device enrollment for iOS/iPadOS institutionally owned devices is enabled. Maps to request field 'accountDrivenDeviceIosEnrollmentEnabled'
- `enable_for_personally_owned_vision_pro_devices` (Boolean) Whether user initiated device enrollment for visionOS iOS/iPadOS institutionally owned devices is enabled. Maps to request field 'accountDrivenDeviceVisionosEnrollmentEnabled'


<a id="nestedblock--user_initiated_enrollment_for_devices--account_driven_user_enrollment"></a>
### Nested Schema for `user_initiated_enrollment_for_devices.account_driven_user_enrollment`

Optional:

- `enable_for_personally_owned_mobile_devices` (Boolean) Whether user initiated device enrollment for iOS/iPadOS personally owned devices is enabled. Maps to request field 'accountDrivenUserEnrollmentEnabled'
- `enable_for_personally_owned_vision_pro_devices` (Boolean) Whether user initiated device enrollment for visionOS personally owned devices is enabled. Maps to request field 'accountDrivenUserVisionosEnrollmentEnabled'
- `enable_maid_username_merge` (Boolean) Maps to API-only field 'maidUsernameMergeEnabled'. Whether to enable merging of the Managed Apple account with the Jamf Pro user account during enrollment.


<a id="nestedblock--user_initiated_enrollment_for_devices--profile_driven_enrollment_via_url"></a>
### Nested Schema for `user_initiated_enrollment_for_devices.profile_driven_enrollment_via_url`

Optional:

- `enable_for_institutionally_owned_devices` (Boolean) Whether user initiated device enrollment for iOS/iPadOS institutionally owned devices is enabled. Maps to request field 'iosEnterpriseEnrollmentEnabled'
- `enable_for_personally_owned_devices` (Boolean) Whether user initiated device enrollment for iOS/iPadOS personally owned devices is enabled. Maps to request field 'iosPersonalEnrollmentEnabled'
- `personal_device_enrollment_type` (String) Enrollment type for user initiated device enrollment for iOS/iPadOS personally owned devices. Note: Personal data on devices enrolled via User Enrollment is retained when the MDM profile is removed. (iOS 13.1 or later, iPadOS 13.1 or later). Maps to request field 'personalDeviceEnrollmentType'



<a id="nestedatt--mdm_signing_certificate_details"></a>
### Nested Schema for `mdm_signing_certificate_details`

Read-Only:

- `serial_number` (String)
- `subject` (String)