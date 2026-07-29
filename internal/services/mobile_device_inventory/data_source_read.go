package mobile_device_inventory

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	var ident string

	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	var mobileDevice *jamfpro.ResourceMobileDeviceInventory
	var err error

	warn_if_not_found := d.Get("warn_if_not_found").(bool)

	if val, ok := d.GetOk("name"); ok {
		ident = val.(string)
		mobileDevice, err = client.GetMobileDeviceInventoryByName(ident)

	} else if val, ok := d.GetOk("serial_number"); ok {
		ident = val.(string)
		mobileDevice, err = client.GetMobileDeviceInventoryBySerialNumber(ident)

	} else if val, ok := d.GetOk("id"); ok {
		ident = val.(string)
		mobileDevice, err = client.GetMobileDeviceInventoryByID(ident)

	} else {
		return diag.Errorf("Either 'name', 'serial_number', or 'id' must be provided")
	}

	if err != nil {
		if warn_if_not_found {
			return append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Mobile device at %s not found", ident),
				Detail:   fmt.Sprintf("Not erroring due to warn_if_not_found enabled\nerr: %v", err),
			})
		}

		return diag.FromErr(fmt.Errorf("failed to fetch mobile device inventory: %v", err))
	}

	d.SetId(mobileDevice.MobileDeviceId)

	return setMobileDeviceInventoryData(d, mobileDevice)
}

func setMobileDeviceInventoryData(d *schema.ResourceData, device *jamfpro.ResourceMobileDeviceInventory) diag.Diagnostics {
	var diags diag.Diagnostics

	attributes := map[string]interface{}{
		"id":                                      device.MobileDeviceId,
		"mobile_device_id":                        device.MobileDeviceId,
		"serial_number":                           device.Hardware.SerialNumber,
		"model":                                   device.Hardware.Model,
		"model_identifier":                        device.Hardware.ModelIdentifier,
		"model_number":                            device.Hardware.ModelNumber,
		"capacity_mb":                             device.Hardware.CapacityMb,
		"available_space_mb":                      device.Hardware.AvailableSpaceMb,
		"used_space_percentage":                   device.Hardware.UsedSpacePercentage,
		"battery_level":                           device.Hardware.BatteryLevel,
		"battery_health":                          device.Hardware.BatteryHealth,
		"wifi_mac_address":                        device.Hardware.WifiMacAddress,
		"bluetooth_mac_address":                   device.Hardware.BluetoothMacAddress,
		"bluetooth_low_energy_capable":            device.Hardware.BluetoothLowEnergyCapable,
		"udid":                                    device.General.Udid,
		"display_name":                            device.General.DisplayName,
		"asset_tag":                               device.General.AssetTag,
		"device_id":                               device.General.DeviceId,
		"management_id":                           device.General.ManagementId,
		"ip_address":                              device.General.IpAddress,
		"managed":                                 device.General.Managed,
		"supervised":                              device.General.Supervised,
		"device_ownership_type":                   device.General.DeviceOwnershipType,
		"enrollment_method_prestage":              device.General.EnrollmentMethod,
		"enrollment_session_token_valid":          device.General.EnrollmentSessionTokenValid,
		"last_enrolled_date":                      device.General.LastEnrolledDate,
		"mdm_profile_expiration_date":             device.General.MdmProfileExpirationDate,
		"time_zone":                               device.General.TimeZone,
		"declarative_device_management_enabled":   device.General.DeclarativeDeviceManagementEnabled,
		"os_version":                              device.General.OsVersion,
		"os_build":                                device.General.OsBuild,
		"os_supplemental_build_version":           device.General.OsSupplementalBuildVersion,
		"os_rapid_security_response":              device.General.OsRapidSecurityResponse,
		"last_inventory_update_date":              device.General.LastInventoryUpdateDate,
		"last_cloud_backup_date":                  device.General.LastCloudBackupDate,
		"last_backup_date":                        device.General.LastBackupDate,
		"cloud_backup_enabled":                    device.General.CloudBackupEnabled,
		"device_locator_service_enabled":          device.General.DeviceLocatorServiceEnabled,
		"do_not_disturb_enabled":                  device.General.DoNotDisturbEnabled,
		"lost_mode_enabled":                       device.General.LostModeEnabled,
		"lost_mode_enabled_date":                  device.General.LostModeEnabledDate,
		"itunes_store_account_active":             device.General.ItunesStoreAccountActive,
		"languages":                               device.General.Languages,
		"locales":                                 device.General.Locales,
		"shared_ipad":                             device.General.SharedIpad,
		"quota_size":                              device.General.QuotaSize,
		"resident_users":                          device.General.ResidentUsers,
		"exchange_device_id":                      device.General.ExchangeDeviceId,
		"tethered":                                device.General.Tethered,
		"username":                                device.UserAndLocation.Username,
		"full_name":                               device.UserAndLocation.FullName,
		"email_address":                           device.UserAndLocation.EmailAddress,
		"phone_number":                            device.UserAndLocation.PhoneNumber,
		"position":                                device.UserAndLocation.Position,
		"department":                              device.UserAndLocation.Department,
		"building":                                device.UserAndLocation.Building,
		"room":                                    device.UserAndLocation.Room,
		"purchased_or_leased":                     device.Purchasing.PurchasedOrLeased,
		"po_number":                               device.Purchasing.PoNumber,
		"po_date":                                 device.Purchasing.PoDate,
		"vendor":                                  device.Purchasing.Vendor,
		"purchase_price":                          device.Purchasing.PurchasePrice,
		"purchasing_account":                      device.Purchasing.PurchasingAccount,
		"purchasing_contact":                      device.Purchasing.PurchasingContact,
		"warranty_expiration_date":                device.Purchasing.WarrantyExpirationDate,
		"apple_care_id":                           device.Purchasing.AppleCareId,
		"lease_expiration_date":                   device.Purchasing.LeaseExpirationDate,
		"life_expectancy_years":                   device.Purchasing.LifeExpectancyYears,
		"activation_lock_enabled":                 device.Security.ActivationLockEnabled,
		"data_protection":                         device.Security.DataProtection,
		"block_encryption_capable":                device.Security.BlockEncryptionCapable,
		"file_encryption_capable":                 device.Security.FileEncryptionCapable,
		"hardware_encryption_supported":           device.Security.HardwareEncryptionSupported,
		"passcode_present":                        device.Security.PasscodePresent,
		"passcode_compliant":                      device.Security.PasscodeCompliant,
		"passcode_compliant_with_profile":         device.Security.PasscodeCompliantWithProfile,
		"passcode_lock_grace_period_enforced_seconds": device.Security.PasscodeLockGracePeriodEnforcedSeconds,
		"personal_device_profile_current":         device.Security.PersonalDeviceProfileCurrent,
		"jailbreak_status":                        device.Security.JailbreakStatus,
		"cellular_technology":                     device.Network.CellularTechnology,
		"voice_roaming_enabled":                   device.Network.VoiceRoamingEnabled,
		"imei":                                    device.Network.Imei,
		"imei2":                                   device.Network.Imei2,
		"iccid":                                   device.Network.Iccid,
		"eid":                                     device.Network.Eid,
		"meid":                                    device.Network.Meid,
		"current_carrier_network":                 device.Network.CurrentCarrierNetwork,
		"home_carrier_network":                    device.Network.HomeCarrierNetwork,
		"current_mobile_country_code":             device.Network.CurrentMobileCountryCode,
		"current_mobile_network_code":             device.Network.CurrentMobileNetworkCode,
		"home_mobile_country_code":                device.Network.HomeMobileCountryCode,
		"home_mobile_network_code":                device.Network.HomeMobileNetworkCode,
		"carrier_settings_version":                device.Network.CarrierSettingsVersion,
		"data_roaming_enabled":                    device.Network.DataRoamingEnabled,
		"roaming":                                 device.Network.Roaming,
		"personal_hotspot_enabled":                device.Network.PersonalHotspotEnabled,
		"device_phone_number":                     device.Network.DevicePhoneNumber,
		"modem_firmware_version":                  device.Network.ModemFirmwareVersion,
	}

	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("failed to set %s: %v", key, err))...)
		}
	}

	if len(device.ExtensionAttributes) > 0 {
		extAttrs := make([]map[string]interface{}, len(device.ExtensionAttributes))
		for i, attr := range device.ExtensionAttributes {
			extAttrs[i] = map[string]interface{}{
				"display_name": attr.DisplayName,
				"value":        attr.Value,
			}
		}
		if err := d.Set("extension_attributes", extAttrs); err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("failed to set extension_attributes: %v", err))...)
		}
	}

	return diags
}
