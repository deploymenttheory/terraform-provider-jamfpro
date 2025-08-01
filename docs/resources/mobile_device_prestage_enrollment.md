---
page_title: "jamfpro_mobile_device_prestage_enrollment"
description: |-
  
---

# jamfpro_mobile_device_prestage_enrollment (Resource)


## Example Usage
```terraform
resource "jamfpro_mobile_device_prestage_enrollment" "example_prestage" {
  display_name                            = "iOS Device Prestage"
  mandatory                               = true
  mdm_removable                           = false
  support_phone_number                    = "+44-1234-567890"
  support_email_address                   = "support@mycompany.org"
  department                              = "IT Department"
  default_prestage                        = false
  enrollment_site_id                      = "-1"
  keep_existing_site_membership           = false
  keep_existing_location_information      = false
  require_authentication                  = true
  authentication_prompt                   = "Sign in to continue setup"
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
  site_id                                 = "-1"
  anchor_certificates                     = []

  skip_setup_items {
    location                = false
    privacy                 = true
    biometric               = false
    software_update         = false
    diagnostics             = true
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
    tos                     = true
    welcome                 = true
    tap_to_setup            = true
  }

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
    lease_date         = "1970-01-01"
    po_date            = "1970-01-01"
    warranty_date      = "1970-01-01"
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `allow_pairing` (Boolean) Allow device pairing.
- `authentication_prompt` (String) Message displayed when authentication is required.
- `auto_advance_setup` (Boolean) Indicates if setup should auto-advance.
- `configure_device_before_setup_assistant` (Boolean) Configure device before Setup Assistant.
- `default_prestage` (Boolean) Whether this is the default prestage enrollment configuration.
- `department` (String) Department associated with the prestage.
- `device_enrollment_program_instance_id` (String) The Automated Device Enrollment instance ID.
- `display_name` (String) The display name of the mobile device prestage enrollment.
- `enable_device_based_activation_lock` (Boolean) Enable device-based Activation Lock.
- `enrollment_customization_id` (String) The enrollment customization ID. Set to 0 if unused.
- `keep_existing_location_information` (Boolean) Maintain existing location information during enrollment.
- `keep_existing_site_membership` (Boolean) Maintain existing site membership during enrollment.
- `language` (String) The language setting defined for the mobile device prestage. Leverages ISO 639-1 (two-letter language codes): https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes . Ensure you define a code supported by jamf pro. Can be left blank.
- `location_information` (Block List, Min: 1) Location information associated with the Jamf Pro mobile device prestage. (see [below for nested schema](#nestedblock--location_information))
- `mandatory` (Boolean) Make MDM Profile Mandatory.
- `maximum_shared_accounts` (Number) Maximum number of shared accounts.
- `mdm_removable` (Boolean) Allow MDM Profile Removal.
- `multi_user` (Boolean) Enable multi-user mode.
- `names` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--names))
- `prestage_minimum_os_target_version_type_ios` (String) The type of minimum OS version enforcement for iOS devices.
- `prestage_minimum_os_target_version_type_ipad` (String) The type of minimum OS version enforcement for iPadOS devices.
- `prevent_activation_lock` (Boolean) Prevent Activation Lock on the device.
- `purchasing_information` (Block List, Min: 1, Max: 1) Purchasing information associated with the mobile device prestage. (see [below for nested schema](#nestedblock--purchasing_information))
- `region` (String) The region setting defined for the mobile device prestage. Leverages ISO 3166-1 alpha-2 (two-letter country codes): https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2 . Ensure you define a code supported by jamf pro. Can be left blank.
- `require_authentication` (Boolean) Require authentication during enrollment.
- `skip_setup_items` (Block List, Min: 1, Max: 1) Selected items are not displayed in the Setup Assistant during mobile device setup within Apple Device Enrollment (ADE). (see [below for nested schema](#nestedblock--skip_setup_items))
- `supervised` (Boolean) Device is supervised.
- `support_email_address` (String) Support email address for the organization.
- `support_phone_number` (String) Support phone number for the organization.
- `timezone` (String) The timezone to be set on the device. Default is UTC

### Optional

- `anchor_certificates` (List of String) List of Base64 encoded PEM Certificates.
- `enforce_temporary_session_timeout` (Boolean) Indicates if temporary session timeout should be enforced.
- `enforce_user_session_timeout` (Boolean) Indicates if user session timeout should be enforced.
- `enrollment_site_id` (String) Site ID for device enrollment.
- `minimum_os_specific_version_ios` (String) The specific minimum OS version required for iOS devices when using MINIMUM_OS_SPECIFIC_VERSION type.
- `minimum_os_specific_version_ipad` (String) The specific minimum OS version required for iPadOS devices when using MINIMUM_OS_SPECIFIC_VERSION type.
- `rts_config_profile_id` (String) The ID of the RTS configuration profile.
- `rts_enabled` (Boolean) Enable RTS.
- `send_timezone` (Boolean) Indicates if timezone should be sent to the device.
- `site_id` (String) The jamf pro site ID. Set to -1 if not used.
- `storage_quota_size_megabytes` (Number) The storage quota size in megabytes.
- `temporary_session_only` (Boolean) Indicates if the session should be temporary only.
- `temporary_session_timeout_seconds` (Number) The timeout duration for temporary sessions in minutes.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `use_storage_quota_size` (Boolean) Indicates if storage quota size should be enforced.
- `user_session_timeout` (Number) The timeout duration for user sessions in minutes.

### Read-Only

- `id` (String) The unique identifier of the mobile device prestage.
- `profile_uuid` (String) The profile UUID of the Automated Device Enrollment instance to associate with the PreStage enrollment. Devices associated with the selected Automated Device Enrollment instance can be assigned the PreStage enrollment
- `version_lock` (Number) The version lock value of the purchasing_information. Optimistic lockingis a mechanism that prevents concurrent operations from taking place on a givenresource. Jamf Pro does this to safeguard resources and workflows that aresensitive to frequent updates, ensuring that one update has completed beforeany additional requests can be processed. Valid request handling is managed bythe construct function.

<a id="nestedblock--location_information"></a>
### Nested Schema for `location_information`

Required:

- `email` (String) The email address associated with this location. Can be left blank.
- `phone` (String) The phone number associated with this location. Can be left blank.
- `position` (String) The position associated with this location. Can be left blank.
- `realname` (String) The real name associated with this location. Can be left blank.
- `room` (String) The room associated with this location. Can be left blank.
- `username` (String) The username for the location information. Can be left blank.

Optional:

- `building_id` (String) The building ID associated with this computer prestage. Set to -1 if not used.
- `department_id` (String) The jamf pro department ID associated with this computer prestage. Set to -1 if not used.

Read-Only:

- `id` (String) The ID of the location information.
- `version_lock` (Number) The version lock of the location information. Optimistic lockingis a mechanism that prevents concurrent operations from taking place on a givenresource. Jamf Pro does this to safeguard resources and workflows that aresensitive to frequent updates, ensuring that one update has completed beforeany additional requests can be processed. Valid request handling is managed bythe construct function.


<a id="nestedblock--names"></a>
### Nested Schema for `names`

Required:

- `assign_names_using` (String) Method to use for assigning device names. Valid values are: 'Default Names', 'List of Names', 'Serial Numbers' or 'Single Name'.
- `device_naming_configured` (Boolean) Indicates if device naming has been configured for this prestage.
- `manage_names` (Boolean) Indicates if device names should be managed by this prestage.

Optional:

- `device_name_prefix` (String) The prefix to use when naming devices with 'Serial Numbers' method.
- `device_name_suffix` (String) The suffix to use when naming devices with 'Serial Numbers' method.
- `prestage_device_names` (Block List) List of predefined device names when using 'List of Names' assignment method. (see [below for nested schema](#nestedblock--names--prestage_device_names))
- `single_device_name` (String) The name to use when using 'Single Name' assignment method.

<a id="nestedblock--names--prestage_device_names"></a>
### Nested Schema for `names.prestage_device_names`

Required:

- `device_name` (String) The name to be assigned to the device.

Read-Only:

- `id` (String) The unique identifier of the device name entry.
- `used` (Boolean) Indicates if this device name has been used.



<a id="nestedblock--purchasing_information"></a>
### Nested Schema for `purchasing_information`

Required:

- `apple_care_id` (String) The AppleCare ID. Can be left blank.
- `leased` (Boolean) Indicates if the item is leased. Default value to false if unused.
- `life_expectancy` (Number) The life expectancy in years. Set to 0 if unused.
- `po_number` (String) The purchase order number. Can be left blank.
- `purchase_price` (String) The purchase price. Can be left blank.
- `purchased` (Boolean) Indicates if the item is purchased. Default value to true if unused.
- `purchasing_account` (String) The purchasing account. Can be left blank.
- `purchasing_contact` (String) The purchasing contact. Can be left blank.
- `vendor` (String) The vendor name. Can be left blank.

Optional:

- `lease_date` (String) The lease date in YYYY-MM-DD format. Use '1970-01-01' if unused.
- `po_date` (String) The purchase order date in YYYY-MM-DD format. Use '1970-01-01' if unused
- `warranty_date` (String) The warranty date in YYYY-MM-DD format. Use '1970-01-01' if unused

Read-Only:

- `id` (String) The ID of the purchasing information.
- `version_lock` (Number) The version lock value of the purchasing_information. Optimistic lockingis a mechanism that prevents concurrent operations from taking place on a givenresource. Jamf Pro does this to safeguard resources and workflows that aresensitive to frequent updates, ensuring that one update has completed beforeany additional requests can be processed. Valid request handling is managed bythe construct function.


<a id="nestedblock--skip_setup_items"></a>
### Nested Schema for `skip_setup_items`

Required:

- `action_button` (Boolean) Skip Action Button setup during device enrollment.
- `android` (Boolean) Skip Android Migration setup during device enrollment.
- `appearance` (Boolean) Skip Appearance setup during device enrollment.
- `apple_id` (Boolean) Skip Apple ID setup during device enrollment.
- `biometric` (Boolean) Skip Biometric setup during device enrollment.
- `camera_button` (Boolean) Skip Camera Button setup during device enrollment.
- `cloud_storage` (Boolean) Skip Cloud Storage setup during device enrollment.
- `diagnostics` (Boolean) Skip Diagnostics setup during device enrollment.
- `display_tone` (Boolean) Skip Display Tone setup during device enrollment.
- `enable_lockdown_mode` (Boolean) Skip Enable Lockdown Mode setup during device enrollment.
- `express_language` (Boolean) Skip Express Language setup during device enrollment.
- `home_button_sensitivity` (Boolean) Skip Home Button Sensitivity setup during device enrollment.
- `imessage_and_facetime` (Boolean) Skip iMessage and FaceTime setup during device enrollment.
- `intelligence` (Boolean) Skip Intelligence setup during device enrollment.
- `location` (Boolean) Skip Location setup during device enrollment.
- `onboarding` (Boolean) Skip Onboarding setup during device enrollment.
- `passcode` (Boolean) Skip Passcode setup during device enrollment.
- `payment` (Boolean) Skip Payment setup during device enrollment.
- `preferred_language` (Boolean) Skip Preferred Language setup during device enrollment.
- `privacy` (Boolean) Skip Privacy setup during device enrollment.
- `restore` (Boolean) Skip Restore setup during device enrollment.
- `restore_completed` (Boolean) Skip Restore Completed setup during device enrollment.
- `safety` (Boolean) Skip Safety setup during device enrollment.
- `safety_and_handling` (Boolean) Skip Safety and Handling setup during device enrollment.
- `screen_saver` (Boolean) Skip Screen Saver setup during device enrollment.
- `screen_time` (Boolean) Skip Screen Time setup during device enrollment.
- `sim_setup` (Boolean) Skip SIM setup during device enrollment.
- `siri` (Boolean) Skip Siri setup during device enrollment.
- `software_update` (Boolean) Skip Software Update setup during device enrollment.
- `tap_to_setup` (Boolean) Skip Tap to Setup during device enrollment.
- `terms_of_address` (Boolean) Skip Terms of Address setup during device enrollment.
- `tos` (Boolean) Skip Terms of Service setup during device enrollment.
- `transfer_data` (Boolean) Skip Transfer Data setup during device enrollment.
- `tv_home_screen_sync` (Boolean) Skip TV Home Screen Sync during device enrollment.
- `tv_provider_sign_in` (Boolean) Skip TV Provider Sign In during device enrollment.
- `tv_room` (Boolean) Skip TV Room setup during device enrollment.
- `update_completed` (Boolean) Skip Update Completed setup during device enrollment.
- `voice_selection` (Boolean) Skip Voice Selection setup during device enrollment.
- `watch_migration` (Boolean) Skip Watch Migration setup during device enrollment.
- `welcome` (Boolean) Skip Welcome setup during device enrollment.
- `zoom` (Boolean) Skip Zoom setup during device enrollment.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)