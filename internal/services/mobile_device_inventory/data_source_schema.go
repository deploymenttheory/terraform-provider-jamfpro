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
				Type:     schema.TypeString,
				Optional: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobile_device_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"udid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asset_tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"model_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"model_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_build": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_supplemental_build_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_rapid_security_response": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"software_update_device_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wifi_mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bluetooth_mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"supervised": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"device_ownership_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enrollment_method_prestage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enrollment_session_token_valid": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_inventory_update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_enrolled_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mdm_profile_expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_phone_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"carrier_settings_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_carrier_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"home_carrier_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iccid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"imei": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"imei2": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_mobile_country_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_mobile_network_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"home_mobile_country_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"home_mobile_network_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cellular_technology": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_roaming_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"roaming": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"voice_roaming_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"personal_hotspot_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"modem_firmware_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_space_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_space_percentage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"battery_level": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"battery_health": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phone_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"position": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"department": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"building": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"room": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"apple_care_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"po_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"po_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"purchase_price": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"purchased_or_leased": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"purchasing_account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"purchasing_contact": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lease_expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"life_expectancy_years": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"warranty_expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"activation_lock_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"block_encryption_capable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"data_protection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"file_encryption_capable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"hardware_encryption_supported": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"jailbreak_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"passcode_compliant": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"passcode_compliant_with_profile": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"passcode_present": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"passcode_lock_grace_period_enforced_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"personal_device_profile_current": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lost_mode_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lost_mode_enabled_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_locator_service_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"do_not_disturb_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cloud_backup_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_cloud_backup_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_backup_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"itunes_store_account_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"bluetooth_low_energy_capable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"exchange_device_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared_ipad": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"quota_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resident_users": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"declarative_device_management_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"management_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"languages": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locales": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tethered": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"extension_attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
