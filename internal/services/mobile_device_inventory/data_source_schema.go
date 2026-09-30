package mobile_device_inventory

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceJamfProMobileDeviceInventory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"warn_if_not_found": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enabling this setting will cause the provider to only WARN if a mobile device is not found. By default the provider will ERROR.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the mobile device",
			},
			"serial_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The serial number of the mobile device",
			},
			"mobile_device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mobile device ID",
			},
			"udid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique Device Identifier",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name of the device",
			},
			"asset_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Asset tag assigned to the device",
			},
			"device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Device ID",
			},
			"model": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Device model",
			},
			"model_identifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Model identifier",
			},
			"model_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Model number",
			},
			"os_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating system version",
			},
			"os_build": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS build number",
			},
			"os_supplemental_build_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS supplemental build version",
			},
			"os_rapid_security_response": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS rapid security response version",
			},
			"software_update_device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Software update device ID",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address",
			},
			"wifi_mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Wi-Fi MAC address",
			},
			"bluetooth_mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bluetooth MAC address",
			},
			"managed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the device is managed",
			},
			"supervised": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the device is supervised",
			},
			"device_ownership_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Device ownership type",
			},
			"enrollment_method_prestage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enrollment method prestage",
			},
			"enrollment_session_token_valid": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether enrollment session token is valid",
			},
			"last_inventory_update_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last inventory update date",
			},
			"last_enrolled_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last enrollment date",
			},
			"mdm_profile_expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MDM profile expiration date",
			},
			"device_phone_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Device phone number",
			},
			"carrier_settings_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Carrier settings version",
			},
			"current_carrier_network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current carrier network",
			},
			"home_carrier_network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Home carrier network",
			},
			"iccid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Integrated Circuit Card Identifier",
			},
			"imei": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "International Mobile Equipment Identity",
			},
			"imei2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Second IMEI for dual SIM devices",
			},
			"meid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mobile Equipment Identifier",
			},
			"eid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "eSIM identifier",
			},
			"current_mobile_country_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current mobile country code",
			},
			"current_mobile_network_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current mobile network code",
			},
			"home_mobile_country_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Home mobile country code",
			},
			"home_mobile_network_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Home mobile network code",
			},
			"cellular_technology": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cellular technology type",
			},
			"data_roaming_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether data roaming is enabled",
			},
			"roaming": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the device is roaming",
			},
			"voice_roaming_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Voice roaming enabled status",
			},
			"personal_hotspot_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether personal hotspot is enabled",
			},
			"modem_firmware_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modem firmware version",
			},
			"capacity_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total storage capacity in megabytes",
			},
			"available_space_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Available storage space in megabytes",
			},
			"used_space_percentage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Percentage of storage space used",
			},
			"battery_level": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current battery level percentage",
			},
			"battery_health": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Battery health status",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username",
			},
			"full_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User full name",
			},
			"email_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User email address",
			},
			"phone_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User phone number",
			},
			"position": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User position",
			},
			"department": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Department",
			},
			"building": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Building location",
			},
			"room": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Room location",
			},
			"apple_care_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AppleCare ID",
			},
			"po_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Purchase order number",
			},
			"po_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Purchase order date",
			},
			"purchase_price": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Purchase price",
			},
			"purchased_or_leased": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the device is purchased or leased",
			},
			"purchasing_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Purchasing account",
			},
			"purchasing_contact": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Purchasing contact",
			},
			"lease_expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Lease expiration date",
			},
			"life_expectancy_years": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Expected life expectancy in years",
			},
			"vendor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vendor",
			},
			"warranty_expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Warranty expiration date",
			},
			"activation_lock_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether activation lock is enabled",
			},
			"block_encryption_capable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the device supports block encryption",
			},
			"data_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether data protection is enabled",
			},
			"file_encryption_capable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the device supports file encryption",
			},
			"hardware_encryption_supported": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether hardware encryption is supported",
			},
			"jailbreak_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Jailbreak status",
			},
			"passcode_compliant": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether passcode is compliant",
			},
			"passcode_compliant_with_profile": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether passcode is compliant with profile",
			},
			"passcode_present": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether a passcode is present",
			},
			"passcode_lock_grace_period_enforced_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Passcode lock grace period in seconds",
			},
			"personal_device_profile_current": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether personal device profile is current",
			},
			"lost_mode_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether lost mode is enabled",
			},
			"lost_mode_enabled_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date lost mode was enabled",
			},
			"device_locator_service_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether device locator service is enabled",
			},
			"do_not_disturb_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether Do Not Disturb is enabled",
			},
			"cloud_backup_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether cloud backup is enabled",
			},
			"last_cloud_backup_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last cloud backup date",
			},
			"last_backup_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last backup date",
			},
			"itunes_store_account_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether iTunes Store account is active",
			},
			"bluetooth_low_energy_capable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the device is Bluetooth Low Energy capable",
			},
			"exchange_device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exchange device ID",
			},
			"shared_ipad": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is a shared iPad",
			},
			"quota_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Quota size for shared iPad",
			},
			"resident_users": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of resident users for shared iPad",
			},
			"declarative_device_management_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether declarative device management is enabled",
			},
			"management_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Management ID",
			},
			"time_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Device time zone",
			},
			"languages": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Device languages",
			},
			"locales": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Device locales",
			},
			"tethered": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the device is tethered",
			},
			"extension_attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Custom extension attributes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Extension attribute display name",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Extension attribute values",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}
