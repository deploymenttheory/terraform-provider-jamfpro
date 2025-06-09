resource "jamfpro_mobile_device_application" "self_service" {
  name         = "Jamf Self Service"
  display_name = "Jamf Self Service"
  bundle_id    = "com.jamfsoftware.selfservice"
  version      = "11.3.3"
  internal_app = false

  category_id = 9
  site_id     = -1

  itunes_store_url      = "https://apps.apple.com/gb/app/jamf-self-service/id718509958"
  external_url          = "https://apps.apple.com/gb/app/jamf-self-service/id718509958"
  itunes_country_region = "US"
  itunes_sync_time      = 0

  deploy_automatically                   = true
  deploy_as_managed_app                  = true
  remove_app_when_mdm_profile_is_removed = true
  prevent_backup_of_app_data             = false
  allow_user_to_delete                   = false
  require_network_tethered               = false
  keep_description_and_icon_up_to_date   = false
  keep_app_updated_on_devices            = false
  free                                   = true
  take_over_management                   = true
  host_externally                        = true
  make_available_after_install           = false

  vpp {
    assign_vpp_device_based_licenses = true
    vpp_admin_account_id             = 2
  }

  self_service {
    self_service_description = <<-EOT
Jamf Self Service empowers you to be more productive, successful, and self-sufficient with your iOS or iPadOS device.

Using the intuitive interface, you can browse and install trusted apps and books from your organization, update configurations, and receive real-time notifications for available services without having to contact IT. 

Jamf Self Service is a one stop shop to get everything you need on your iOS or iPadOS device to be successful in your organization.

Key Features:
⁃ Install apps
⁃ Install books
⁃ Update device configurations
⁃ Receive notifications

Minimum Requirements:
⁃ Mobile device with iOS 11 or later, or iPadOS 13 or later
⁃ Jamf Pro 9.4 or later

* Jamf Self Service requires that your device is managed by your organization's IT department using Jamf Pro.
    EOT
    feature_on_main_page     = false
    notification             = false
    self_service_icon {
      id = 157
    }
  }

  app_configuration {
    preferences = <<-EOT
<dict>
<key>INVITATION_STRING</key>
<string>$MOBILEDEVICEAPPINVITE</string>
<key>JSS_ID</key>
<string>$JSSID</string>
<key>SERIAL_NUMBER</key>
<string>$SERIALNUMBER</string>
<key>DEVICE_NAME</key>
<string>$DEVICENAME</string>
<key>MAC_ADDRESS</key>
<string>$MACADDRESS</string>
<key>MANAGEMENT_ID</key>
<string>$MANAGEMENTID</string>
<key>JSS_URL</key>      
<string>$JPS_URL</string>
</dict>
    EOT
  }

  scope {
    all_mobile_devices = false
    all_jss_users      = false

    mobile_device_group_ids = [9]

    exclusions {
      mobile_device_group_ids = [14]
    }
  }
}
