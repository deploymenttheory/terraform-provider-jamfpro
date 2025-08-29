---
page_title: "jamfpro_computer_prestage_enrollment"
description: |-
  
---

# jamfpro_computer_prestage_enrollment (Resource)


## Example Usage
```terraform
// Minimum Configuration

resource "jamfpro_computer_prestage_enrollment" "minimum_example" {
  display_name                          = "jamfpro-sdk-example-computerPrestageMinimum-config"
  mandatory                             = true
  mdm_removable                         = true
  support_phone_number                  = "111-222-3333"
  support_email_address                 = "email@company.com"
  department                            = "department name"
  default_prestage                      = false
  enrollment_site_id                    = "-1"
  keep_existing_site_membership         = false
  keep_existing_location_information    = false
  require_authentication                = false
  authentication_prompt                 = "hello welcome to your enterprise managed macOS device"
  prevent_activation_lock               = false
  enable_device_based_activation_lock   = false
  device_enrollment_program_instance_id = "1"
  skip_setup_items {
    biometric            = false
    terms_of_address     = false
    file_vault           = false
    icloud_diagnostics   = false
    diagnostics          = false
    accessibility        = false
    apple_id             = false
    screen_time          = false
    siri                 = false
    display_tone         = false
    restore              = false
    appearance           = false
    privacy              = false
    payment              = false
    registration         = false
    tos                  = false
    icloud_storage       = false
    location             = false
    intelligence         = false
    enable_lockdown_mode = false
    welcome              = false
    wallpaper            = false
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
  anchor_certificates                     = []
  enrollment_customization_id             = "0"
  language                                = ""
  region                                  = ""
  auto_advance_setup                      = false
  install_profiles_during_setup           = true
  prestage_installed_profile_ids          = []
  custom_package_ids                      = []
  custom_package_distribution_point_id    = "-1"
  enable_recovery_lock                    = false
  recovery_lock_password_type             = "MANUAL" // "MANUAL" / "RANDOM"
  recovery_lock_password                  = ""
  rotate_recovery_lock_password           = false
  prestage_minimum_os_target_version_type = "NO_ENFORCEMENT"
  minimum_os_specific_version             = ""
  site_id                                 = "-1"
  account_settings {
    payload_configured                           = true
    local_admin_account_enabled                  = false
    admin_username                               = ""
    admin_password                               = ""
    hidden_admin_account                         = false
    local_user_managed                           = false
    user_account_type                            = "ADMINISTRATOR"
    prefill_primary_account_info_feature_enabled = false
    prefill_type                                 = "UNKNOWN"
    prefill_account_full_name                    = ""
    prefill_account_user_name                    = ""
    prevent_prefill_info_from_modification       = false
  }
}

// Configured the Jamf Pro Computer Prestage Enrollment
resource "jamfpro_computer_prestage_enrollment" "configured_example_1" {
  display_name                          = "testing_values"
  mandatory                             = true
  mdm_removable                         = true
  support_phone_number                  = "111-222-3333"
  support_email_address                 = "email@company.com"
  department                            = "department name"
  default_prestage                      = false
  enrollment_site_id                    = "-1"
  keep_existing_site_membership         = false
  keep_existing_location_information    = false
  require_authentication                = false
  authentication_prompt                 = "hello welcome to your enterprise managed macOS device"
  prevent_activation_lock               = false
  enable_device_based_activation_lock   = false
  device_enrollment_program_instance_id = "1"
  skip_setup_items {
    biometric            = false
    terms_of_address     = false
    file_vault           = false
    icloud_diagnostics   = false
    diagnostics          = false
    accessibility        = false
    apple_id             = false
    screen_time          = false
    siri                 = false
    display_tone         = false
    restore              = false
    appearance           = false
    privacy              = true
    payment              = false
    registration         = false
    tos                  = false
    icloud_storage       = false
    location             = false
    intelligence         = true
    enable_lockdown_mode = false
    welcome              = false
    wallpaper            = false
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
  anchor_certificates                     = []
  enrollment_customization_id             = "0"
  language                                = ""
  region                                  = ""
  auto_advance_setup                      = false
  install_profiles_during_setup           = true
  prestage_installed_profile_ids          = [3114, jamfpro_macos_configuration_profile_plist.jamfpro_macos_configuration_profile_001.id]
  custom_package_ids                      = [1, 2]
  custom_package_distribution_point_id    = "-2" // "-1" - not used / "-2" - Cloud Distribution Point (Jamf Cloud) / "any other number" - Distribution Point ID
  enable_recovery_lock                    = true
  recovery_lock_password_type             = "MANUAL" // "MANUAL" / "RANDOM"
  recovery_lock_password                  = "thing"
  rotate_recovery_lock_password           = false
  prestage_minimum_os_target_version_type = "MINIMUM_OS_LATEST_VERSION" // "NO_ENFORCEMENT" / "MINIMUM_OS_LATEST_VERSION" / "MINIMUM_OS_LATEST_MAJOR_VERSION" / "MINIMUM_OS_SPECIFIC_VERSION"
  minimum_os_specific_version             = "14.6.1"

  site_id = "-1"
  account_settings {
    payload_configured                           = true
    local_admin_account_enabled                  = true
    admin_username                               = "thing"
    admin_password                               = "thing"
    hidden_admin_account                         = true
    local_user_managed                           = true
    user_account_type                            = "ADMINISTRATOR" // "STANDARD" / "ADMINISTRATOR" / "SKIP"
    prefill_primary_account_info_feature_enabled = true
    prefill_type                                 = "CUSTOM" // "UNKNOWN" / "CUSTOM" / "DEVICE_OWNER"
    prefill_account_full_name                    = "firstname.lastname"
    prefill_account_user_name                    = "firstname.lastname"
    prevent_prefill_info_from_modification       = false
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_settings` (Block List, Min: 1) (see [below for nested schema](#nestedblock--account_settings))
- `authentication_prompt` (String) Authentication Message to display to the user. Used when Require Authentication is enabled. Can be left blank.
- `auto_advance_setup` (Boolean) Indicates if setup should auto-advance.
- `custom_package_distribution_point_id` (String) Set the Enrollment Packages distribution point by it's ID.Valid values are: None using '-1', Cloud Distribution Point (Jamf Cloud)by using '-2', else all other valid valid values correspond to theID of the distribution point.
- `custom_package_ids` (Set of String) Define the Enrollment Packages by their package ID to add an enrollment package to the PreStage enrollment. Compatible packages must be built as flat, distribution style .pkg files and be signed by a certificate that is trusted by managed computers. Can be left blank.
- `default_prestage` (Boolean) Indicates if this is the default computer prestage enrollment configuration. If yes then new devices will be automatically assigned to this PreStage enrollment
- `department` (String) The department the computer prestage is assigned to. Can be left blank.
- `device_enrollment_program_instance_id` (String) The Automated Device Enrollment instance ID to associate with the PreStage enrollment. Devices associated with the selected Automated Device Enrollment instance can be assigned the PreStage enrollment
- `display_name` (String) The display name of the computer prestage enrollment.
- `enable_device_based_activation_lock` (Boolean) Indicates if device-based activation lock should be enabled.
- `enable_recovery_lock` (Boolean) Configure how the Recovery Lock password is set on computers with macOS 11.5 or later.
- `enrollment_customization_id` (String) The enrollment customization ID. Set to 0 if unused.
- `enrollment_site_id` (String) The jamf pro Site ID that computers will be added to during enrollment. Should be set to -1, if not used.
- `install_profiles_during_setup` (Boolean) Indicates if profiles should be installed during setup.
- `keep_existing_location_information` (Boolean) Indicates if enrolled should use existing location information, if applicable
- `keep_existing_site_membership` (Boolean) Indicates if enrolled should use existing site membership, if applicable
- `language` (String) The language setting defined for the computer prestage. Leverages ISO 639-1 (two-letter language codes): https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes . Ensure you define a code supported by jamf pro. Can be left blank.
- `location_information` (Block List, Min: 1) Location information associated with the Jamf Pro computer prestage. (see [below for nested schema](#nestedblock--location_information))
- `mandatory` (Boolean) Make MDM Profile Mandatory and require the user to apply the MDM profile. Computers with macOS 10.15 or later automatically require the user to apply the MDM profile
- `mdm_removable` (Boolean) Allow MDM Profile Removal and allow the user to remove the MDM profile.
- `prestage_installed_profile_ids` (Set of String) IDs of the macOS configuration profiles installed during PreStage enrollment. requires decending order of profile IDs so uses a set rather than a list. can be left blank.
- `prestage_minimum_os_target_version_type` (String) Enforce a minimum macOS target version type for the prestage enrollment. Required.
- `prevent_activation_lock` (Boolean) Prevent user from enabling Activation Lock.
- `purchasing_information` (Block List, Min: 1) Purchasing information associated with the computer prestage. (see [below for nested schema](#nestedblock--purchasing_information))
- `recovery_lock_password_type` (String) Method to use to set Recovery Lock password.'MANUAL' results in user having to enter a password. (Applies to all users) 'RANDOM' results inautomatic generation of a random password being set for the device. 'MANUAL' is the default.
- `region` (String) The region setting defined for the computer prestage. Leverages ISO 3166-1 alpha-2 (two-letter country codes): https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2 . Ensure you define a code supported by jamf pro. Can be left blank.
- `require_authentication` (Boolean) Indicates if the user is required to provide username and password on computers with macOS 10.10 or later.
- `rotate_recovery_lock_password` (Boolean) Indicates if the recovery lock password should be rotated.
- `site_id` (String) The jamf pro site ID. Set to -1 if not used.
- `skip_setup_items` (Block List, Min: 1, Max: 1) Selected items are not displayed in the Setup Assistant during macOS device setup within Apple Device Enrollment (ADE). (see [below for nested schema](#nestedblock--skip_setup_items))
- `support_email_address` (String) The Support email address for the organization. Can be left blank.
- `support_phone_number` (String) The Support phone number for the organization. Can be left blank.

### Optional

- `anchor_certificates` (List of String) List of Base64 encoded PEM Certificates.
- `minimum_os_specific_version` (String) The minimum macOS version to enforce for the prestage enrollment. Only used if prestate_minimum_os_target_version_type is set to MINIMUM_OS_SPECIFIC_VERSION.
- `recovery_lock_password` (String) Generate new Recovery Lock password 60 minutes after the password is viewed in Jamf Pro. Can be left blank.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The unique identifier of the computer prestage.
- `profile_uuid` (String) The profile UUID of the Automated Device Enrollment instance to associate with the PreStage enrollment. Devices associated with the selected Automated Device Enrollment instance can be assigned the PreStage enrollment
- `version_lock` (Number) The version lock value of the purchasing_information. Optimistic lockingis a mechanism that prevents concurrent operations from taking place on a givenresource. Jamf Pro does this to safeguard resources and workflows that aresensitive to frequent updates, ensuring that one update has completed beforeany additional requests can be processed. Valid request handling is managed bythe construct function.

<a id="nestedblock--account_settings"></a>
### Nested Schema for `account_settings`

Required:

- `admin_password` (String, Sensitive) The admin password. Can be left blank if not used.
- `admin_username` (String, Sensitive) The admin username. Can be left blank if not used.
- `hidden_admin_account` (Boolean) Indicates if the admin account is hidden.
- `local_admin_account_enabled` (Boolean) Indicates if the local admin account is enabled.
- `local_user_managed` (Boolean) Indicates if the local user is managed.
- `payload_configured` (Boolean) Indicates if the payload is configured.
- `prefill_account_full_name` (String) Type of information to use to pre-fill the primary account full name with. Can be left blank.
- `prefill_account_user_name` (String) Type of information to use to pre-fill the primary account user name with. Can be left blank.
- `prefill_primary_account_info_feature_enabled` (Boolean) Indicates if prefilling primary account info feature is enabled.
- `prefill_type` (String) Pre-fill primary account information type (CUSTOM, DEVICE_OWNER, or UNKNOWN). Set as UNKNOWN if you wish to leave it unconfigured.
- `prevent_prefill_info_from_modification` (Boolean) Lock prefill primary account information from modification.
- `user_account_type` (String) Type of user account (ADMINISTRATOR, STANDARD, SKIP).

Read-Only:

- `id` (String) ID of Account Settings.
- `version_lock` (Number) The version lock value of the account settings block. Optimistic lockingis a mechanism that prevents concurrent operations from taking place on a givenresource. Jamf Pro does this to safeguard resources and workflows that aresensitive to frequent updates, ensuring that one update has completed beforeany additional requests can be processed. Valid request handling is managed bythe construct function.


<a id="nestedblock--location_information"></a>
### Nested Schema for `location_information`

Required:

- `department_id` (String) The jamf pro department ID associated with this computer prestage. Set to -1 if not used.
- `email` (String) The email address associated with this location. Can be left blank.
- `phone` (String) The phone number associated with this location. Can be left blank.
- `position` (String) The position associated with this location. Can be left blank.
- `realname` (String) The real name associated with this location. Can be left blank.
- `room` (String) The room associated with this location. Can be left blank.
- `username` (String) The username for the location information. Can be left blank.

Optional:

- `building_id` (String) The building ID associated with this computer prestage. Set to -1 if not used.

Read-Only:

- `id` (String) The ID of the location information.
- `version_lock` (Number) The version lock of the location information. Optimistic lockingis a mechanism that prevents concurrent operations from taking place on a givenresource. Jamf Pro does this to safeguard resources and workflows that aresensitive to frequent updates, ensuring that one update has completed beforeany additional requests can be processed. Valid request handling is managed bythe construct function.


<a id="nestedblock--purchasing_information"></a>
### Nested Schema for `purchasing_information`

Required:

- `apple_care_id` (String) The AppleCare ID. Can be left blank.
- `lease_date` (String) The lease date in YYYY-MM-DD format. Use '1970-01-01' if unused.
- `leased` (Boolean) Indicates if the item is leased. Default value to false if unused.
- `life_expectancy` (Number) The life expectancy in years. Set to 0 if unused.
- `po_date` (String) The purchase order date in YYYY-MM-DD format. Use '1970-01-01' if unused
- `po_number` (String) The purchase order number. Can be left blank.
- `purchase_price` (String) The purchase price. Can be left blank.
- `purchased` (Boolean) Indicates if the item is purchased. Default value to true if unused.
- `purchasing_account` (String) The purchasing account. Can be left blank.
- `purchasing_contact` (String) The purchasing contact. Can be left blank.
- `vendor` (String) The vendor name. Can be left blank.
- `warranty_date` (String) The warranty date in YYYY-MM-DD format. Use '1970-01-01' if unused

Read-Only:

- `id` (String) The ID of the purchasing information.
- `version_lock` (Number) The version lock value of the purchasing_information. Optimistic lockingis a mechanism that prevents concurrent operations from taking place on a givenresource. Jamf Pro does this to safeguard resources and workflows that aresensitive to frequent updates, ensuring that one update has completed beforeany additional requests can be processed. Valid request handling is managed bythe construct function.


<a id="nestedblock--skip_setup_items"></a>
### Nested Schema for `skip_setup_items`

Required:

- `accessibility` (Boolean) Skip accessibility setup.
- `additional_privacy_settings` (Boolean) Skip additional privacy settings setup.
- `appearance` (Boolean) Skip Appearance setup.
- `apple_id` (Boolean) Skip Apple ID setup.
- `biometric` (Boolean) Skip biometric setup.
- `diagnostics` (Boolean) Skip diagnostics setup.
- `display_tone` (Boolean) Skip Display Tone setup. (Deprecated)
- `enable_lockdown_mode` (Boolean) Skip lockdown mode setup.
- `file_vault` (Boolean) Skip FileVault setup.
- `icloud_diagnostics` (Boolean) Skip iCloud diagnostics setup.
- `icloud_storage` (Boolean) Skip iCloud Storage setup.
- `intelligence` (Boolean) Skip Apple Intelligence setup.
- `location` (Boolean) Skip Location setup.
- `payment` (Boolean) Skip Payment setup.
- `privacy` (Boolean) Skip Privacy setup.
- `registration` (Boolean) Skip Registration setup.
- `restore` (Boolean) Skip Restore setup.
- `screen_time` (Boolean) Skip Screen Time setup.
- `siri` (Boolean) Skip Siri setup.
- `software_update` (Boolean) Skip software update setup.
- `terms_of_address` (Boolean) Skip terms of address setup.
- `tos` (Boolean) Skip Terms of Service setup.
- `wallpaper` (Boolean) Skip wallpaper setup.
- `welcome` (Boolean) Skip welcome setup.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)