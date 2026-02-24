package computer_inventory

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead fetches the details of a specific computer inventory
// from Jamf Pro using either its unique Name, Serial Number, or its ID. The function prioritizes the 'name' attribute,
// then 'serial_number', then 'id' for fetching details. If none are provided, it returns an error.
// Once the details are fetched, they are set in the data source's state.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	var computer *jamfpro.ResourceComputerInventory
	var err error

	allow_not_found := d.Get("allow_not_found").(bool)

	if v, ok := d.GetOk("name"); ok {
		attrName, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'name' as string")
		}

		computer, err = client.GetComputerInventoryByName(attrName)
	} else if v, ok := d.GetOk("serial_number"); ok {

		serialNumber, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'serial_number' as string")
		}
		computer, err = client.GetComputerInventoryBySerialNumber(serialNumber)

	} else if v, ok := d.GetOk("id"); ok {
		profileID, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'id' as string")
		}
		computer, err = client.GetComputerInventoryByID(profileID)

	} else {
		return diag.Errorf("Either 'name', 'serial_number', or 'id' must be provided")
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch computer inventory: %v", err))
	}

	d.SetId(computer.ID)
	d.Set("id", computer.ID)
	d.Set("udid", computer.UDID)

	if err := setGeneralSection(d, computer.General); err != nil {
		return diag.FromErr(err)
	}

	if err := setDiskEncryptionSection(d, computer.DiskEncryption); err != nil {
		return diag.FromErr(err)
	}

	if err := setPurchasingSection(d, computer.Purchasing); err != nil {
		return diag.FromErr(err)
	}

	if err := setApplicationsSection(d, computer.Applications); err != nil {
		return diag.FromErr(err)
	}

	if err := setStorageSection(d, computer.Storage); err != nil {
		return diag.FromErr(err)
	}

	if err := setUserAndLocationSection(d, computer.UserAndLocation); err != nil {
		return diag.FromErr(err)
	}

	if err := setHardwareSection(d, computer.Hardware); err != nil {
		return diag.FromErr(err)
	}

	if err := setLocalUserAccountsSection(d, computer.LocalUserAccounts); err != nil {
		return diag.FromErr(err)
	}

	if err := setCertificatesSection(d, computer.Certificates); err != nil {
		return diag.FromErr(err)
	}

	if err := setAttachmentsSection(d, computer.Attachments); err != nil {
		return diag.FromErr(err)
	}

	if err := setPluginsSection(d, computer.Plugins); err != nil {
		return diag.FromErr(err)
	}

	if err := setPackageReceiptsSection(d, computer.PackageReceipts); err != nil {
		return diag.FromErr(err)
	}

	if err := setFontsSection(d, computer.Fonts); err != nil {
		return diag.FromErr(err)
	}

	if err := setSecuritySection(d, computer.Security); err != nil {
		return diag.FromErr(err)
	}

	if err := setOperatingSystemSection(d, computer.OperatingSystem); err != nil {
		return diag.FromErr(err)
	}

	if err := setLicensedSoftwareSection(d, computer.LicensedSoftware); err != nil {
		return diag.FromErr(err)
	}

	if err := setIBeaconsSection(d, computer.Ibeacons); err != nil {
		return diag.FromErr(err)
	}

	if err := setSoftwareUpdatesSection(d, computer.SoftwareUpdates); err != nil {
		return diag.FromErr(err)
	}

	if err := setExtensionAttributesSection(d, computer.ExtensionAttributes); err != nil {
		return diag.FromErr(err)
	}

	if err := setGroupMembershipsSection(d, computer.GroupMemberships); err != nil {
		return diag.FromErr(err)
	}

	if err := setConfigurationProfilesSection(d, computer.ConfigurationProfiles); err != nil {
		return diag.FromErr(err)
	}

	if err := setPrintersSection(d, computer.Printers); err != nil {
		return diag.FromErr(err)
	}

	if err := setServicesSection(d, computer.Services); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// setGeneralSection maps the 'general' section of the computer inventory response to the Terraform resource data and updates the state.
func setGeneralSection(d *schema.ResourceData, general jamfpro.ComputerInventorySubsetGeneral) error {
	// Initialize a map to hold the 'general' section attributes.
	gen := make(map[string]any)

	// Map each attribute of the 'general' section from the API response to the corresponding Terraform schema attribute.
	gen["name"] = general.Name
	gen["last_ip_address"] = general.LastIpAddress
	gen["last_reported_ip"] = general.LastReportedIp
	gen["jamf_binary_version"] = general.JamfBinaryVersion
	gen["platform"] = general.Platform
	gen["barcode1"] = general.Barcode1
	gen["barcode2"] = general.Barcode2
	gen["asset_tag"] = general.AssetTag
	gen["supervised"] = general.Supervised
	gen["mdm_capable"] = []interface{}{
		map[string]interface{}{
			"capable":       general.MdmCapable.Capable,
			"capable_users": general.MdmCapable.CapableUsers,
		},
	}
	gen["report_date"] = general.ReportDate
	gen["last_contact_time"] = general.LastContactTime
	gen["last_cloud_backup_date"] = general.LastCloudBackupDate
	gen["last_enrolled_date"] = general.LastEnrolledDate
	gen["mdm_profile_expiration"] = general.MdmProfileExpiration
	gen["initial_entry_date"] = general.InitialEntryDate
	gen["distribution_point"] = general.DistributionPoint
	gen["itunes_store_account_active"] = general.ItunesStoreAccountActive
	gen["enrolled_via_automated_device_enrollment"] = general.EnrolledViaAutomatedDeviceEnrollment
	gen["user_approved_mdm"] = general.UserApprovedMdm
	gen["declarative_device_management_enabled"] = general.DeclarativeDeviceManagementEnabled
	gen["management_id"] = general.ManagementId

	remoteManagement := make(map[string]any)
	remoteManagement["managed"] = general.RemoteManagement.Managed
	remoteManagement["management_username"] = general.RemoteManagement.ManagementUsername
	gen["remote_management"] = []any{remoteManagement}

	if general.Site.ID != "" || general.Site.Name != "" {
		site := make(map[string]any)
		site["id"] = general.Site.ID
		site["name"] = general.Site.Name
		gen["site_id"] = []any{site}
	}

	if general.EnrollmentMethod.ID != "" || general.EnrollmentMethod.ObjectName != "" || general.EnrollmentMethod.ObjectType != "" {
		enrollmentMethod := make(map[string]any)
		enrollmentMethod["id"] = general.EnrollmentMethod.ID
		enrollmentMethod["object_name"] = general.EnrollmentMethod.ObjectName
		enrollmentMethod["object_type"] = general.EnrollmentMethod.ObjectType
		gen["enrollment_method"] = []any{enrollmentMethod}
	}

	return d.Set("general", []any{gen})
}

// setDiskEncryptionSection maps the 'diskEncryption' section of the computer inventory response to the Terraform resource data and updates the state.
func setDiskEncryptionSection(d *schema.ResourceData, diskEncryption jamfpro.ComputerInventorySubsetDiskEncryption) error {
	diskEnc := make(map[string]any)

	diskEnc["individual_recovery_key_validity_status"] = diskEncryption.IndividualRecoveryKeyValidityStatus
	diskEnc["institutional_recovery_key_present"] = diskEncryption.InstitutionalRecoveryKeyPresent
	diskEnc["disk_encryption_configuration_name"] = diskEncryption.DiskEncryptionConfigurationName
	diskEnc["file_vault2_eligibility_message"] = diskEncryption.FileVault2EligibilityMessage

	bootPartitionDetails := make(map[string]any)
	bootPartitionDetails["partition_name"] = diskEncryption.BootPartitionEncryptionDetails.PartitionName
	bootPartitionDetails["partition_file_vault2_state"] = diskEncryption.BootPartitionEncryptionDetails.PartitionFileVault2State
	bootPartitionDetails["partition_file_vault2_percent"] = diskEncryption.BootPartitionEncryptionDetails.PartitionFileVault2Percent
	diskEnc["boot_partition_encryption_details"] = []any{bootPartitionDetails}

	fileVaultUserNames := make([]string, len(diskEncryption.FileVault2EnabledUserNames))
	copy(fileVaultUserNames, diskEncryption.FileVault2EnabledUserNames)

	diskEnc["file_vault2_enabled_user_names"] = fileVaultUserNames

	return d.Set("disk_encryption", []any{diskEnc})
}

// setPurchasingSection maps the 'purchasing' section of the computer inventory response to the Terraform resource data and updates the state.
func setPurchasingSection(d *schema.ResourceData, purchasing jamfpro.ComputerInventorySubsetPurchasing) error {
	purchasingMap := make(map[string]any)

	purchasingMap["leased"] = purchasing.Leased
	purchasingMap["purchased"] = purchasing.Purchased
	purchasingMap["po_number"] = purchasing.PoNumber
	purchasingMap["po_date"] = purchasing.PoDate
	purchasingMap["vendor"] = purchasing.Vendor
	purchasingMap["warranty_date"] = purchasing.WarrantyDate
	purchasingMap["apple_care_id"] = purchasing.AppleCareId
	purchasingMap["lease_date"] = purchasing.LeaseDate
	purchasingMap["purchase_price"] = purchasing.PurchasePrice
	purchasingMap["life_expectancy"] = purchasing.LifeExpectancy
	purchasingMap["purchasing_account"] = purchasing.PurchasingAccount
	purchasingMap["purchasing_contact"] = purchasing.PurchasingContact

	extAttrs := make([]map[string]any, len(purchasing.ExtensionAttributes))
	for i, attr := range purchasing.ExtensionAttributes {
		attrMap := make(map[string]any)
		attrMap["definition_id"] = attr.DefinitionId
		attrMap["name"] = attr.Name
		attrMap["description"] = attr.Description
		attrMap["enabled"] = attr.Enabled
		attrMap["multi_value"] = attr.MultiValue
		attrMap["values"] = attr.Values
		attrMap["data_type"] = attr.DataType
		attrMap["options"] = attr.Options
		attrMap["input_type"] = attr.InputType

		extAttrs[i] = attrMap
	}
	purchasingMap["extension_attributes"] = extAttrs

	return d.Set("purchasing", []any{purchasingMap})
}

// setApplicationsSection maps the 'applications' section of the computer inventory response to the Terraform resource data and updates the state.
func setApplicationsSection(d *schema.ResourceData, applications []jamfpro.ComputerInventorySubsetApplication) error {
	apps := make([]any, len(applications))

	for i, app := range applications {
		appMap := make(map[string]any)
		appMap["name"] = app.Name
		appMap["path"] = app.Path
		appMap["version"] = app.Version
		appMap["mac_app_store"] = app.MacAppStore
		appMap["size_megabytes"] = app.SizeMegabytes
		appMap["bundle_id"] = app.BundleId
		appMap["update_available"] = app.UpdateAvailable
		appMap["external_version_id"] = app.ExternalVersionId

		apps[i] = appMap
	}

	return d.Set("applications", apps)
}

// setStorageSection maps the 'storage' section of the computer inventory response to the Terraform resource data and updates the state.
func setStorageSection(d *schema.ResourceData, storage jamfpro.ComputerInventorySubsetStorage) error {
	storageMap := make(map[string]any)

	storageMap["boot_drive_available_space_megabytes"] = storage.BootDriveAvailableSpaceMegabytes

	disks := make([]any, len(storage.Disks))
	for i, disk := range storage.Disks {
		diskMap := make(map[string]any)
		diskMap["id"] = disk.ID
		diskMap["device"] = disk.Device
		diskMap["model"] = disk.Model
		diskMap["revision"] = disk.Revision
		diskMap["serial_number"] = disk.SerialNumber
		diskMap["size_megabytes"] = disk.SizeMegabytes
		diskMap["smart_status"] = disk.SmartStatus
		diskMap["type"] = disk.Type

		partitions := make([]any, len(disk.Partitions))
		for j, partition := range disk.Partitions {
			partitionMap := make(map[string]any)
			partitionMap["name"] = partition.Name
			partitionMap["size_megabytes"] = partition.SizeMegabytes
			partitionMap["available_megabytes"] = partition.AvailableMegabytes
			partitionMap["partition_type"] = partition.PartitionType
			partitionMap["percent_used"] = partition.PercentUsed
			partitionMap["file_vault2_state"] = partition.FileVault2State
			partitionMap["file_vault2_progress_percent"] = partition.FileVault2ProgressPercent
			partitionMap["lvm_managed"] = partition.LvmManaged
			partitions[j] = partitionMap
		}
		diskMap["partitions"] = partitions

		disks[i] = diskMap
	}
	storageMap["disks"] = disks

	return d.Set("storage", []any{storageMap})
}

// setUserAndLocationSection maps the 'userAndLocation' section of the computer inventory response to the Terraform resource data and updates the state.
func setUserAndLocationSection(d *schema.ResourceData, userAndLocation jamfpro.ComputerInventorySubsetUserAndLocation) error {
	userLocationMap := make(map[string]any)

	userLocationMap["username"] = userAndLocation.Username
	userLocationMap["realname"] = userAndLocation.Realname
	userLocationMap["email"] = userAndLocation.Email
	userLocationMap["position"] = userAndLocation.Position
	userLocationMap["phone"] = userAndLocation.Phone
	userLocationMap["department_id"] = userAndLocation.DepartmentId
	userLocationMap["building_id"] = userAndLocation.BuildingId
	userLocationMap["room"] = userAndLocation.Room

	if len(userAndLocation.ExtensionAttributes) > 0 {
		extAttrs := make([]map[string]any, len(userAndLocation.ExtensionAttributes))
		for i, attr := range userAndLocation.ExtensionAttributes {
			attrMap := make(map[string]any)
			attrMap["definition_id"] = attr.DefinitionId
			attrMap["name"] = attr.Name
			attrMap["description"] = attr.Description
			attrMap["enabled"] = attr.Enabled
			attrMap["multi_value"] = attr.MultiValue
			attrMap["values"] = attr.Values
			attrMap["data_type"] = attr.DataType
			attrMap["options"] = attr.Options
			attrMap["input_type"] = attr.InputType

			extAttrs[i] = attrMap
		}
		userLocationMap["extension_attributes"] = extAttrs
	}

	return d.Set("user_and_location", []any{userLocationMap})
}

// setHardwareSection maps the 'hardware' section of the computer inventory response to the Terraform resource data and updates the state.
func setHardwareSection(d *schema.ResourceData, hardware jamfpro.ComputerInventorySubsetHardware) error {
	hardwareMap := make(map[string]any)

	hardwareMap["make"] = hardware.Make
	hardwareMap["model"] = hardware.Model
	hardwareMap["model_identifier"] = hardware.ModelIdentifier
	hardwareMap["serial_number"] = hardware.SerialNumber
	hardwareMap["processor_speed_mhz"] = hardware.ProcessorSpeedMhz
	hardwareMap["processor_count"] = hardware.ProcessorCount
	hardwareMap["core_count"] = hardware.CoreCount
	hardwareMap["processor_type"] = hardware.ProcessorType
	hardwareMap["processor_architecture"] = hardware.ProcessorArchitecture
	hardwareMap["bus_speed_mhz"] = hardware.BusSpeedMhz
	hardwareMap["cache_size_kilobytes"] = hardware.CacheSizeKilobytes
	hardwareMap["network_adapter_type"] = hardware.NetworkAdapterType
	hardwareMap["mac_address"] = hardware.MacAddress
	hardwareMap["alt_network_adapter_type"] = hardware.AltNetworkAdapterType
	hardwareMap["alt_mac_address"] = hardware.AltMacAddress
	hardwareMap["total_ram_megabytes"] = hardware.TotalRamMegabytes
	hardwareMap["open_ram_slots"] = hardware.OpenRamSlots
	hardwareMap["battery_capacity_percent"] = hardware.BatteryCapacityPercent
	hardwareMap["smc_version"] = hardware.SmcVersion
	hardwareMap["nic_speed"] = hardware.NicSpeed
	hardwareMap["optical_drive"] = hardware.OpticalDrive
	hardwareMap["boot_rom"] = hardware.BootRom
	hardwareMap["ble_capable"] = hardware.BleCapable
	hardwareMap["supports_ios_app_installs"] = hardware.SupportsIosAppInstalls
	hardwareMap["apple_silicon"] = hardware.AppleSilicon

	if len(hardware.ExtensionAttributes) > 0 {
		extAttrs := make([]map[string]any, len(hardware.ExtensionAttributes))
		for i, attr := range hardware.ExtensionAttributes {
			attrMap := make(map[string]any)
			attrMap["definition_id"] = attr.DefinitionId
			attrMap["name"] = attr.Name
			attrMap["description"] = attr.Description
			attrMap["enabled"] = attr.Enabled
			attrMap["multi_value"] = attr.MultiValue
			attrMap["values"] = attr.Values
			attrMap["data_type"] = attr.DataType
			attrMap["options"] = attr.Options
			attrMap["input_type"] = attr.InputType

			extAttrs[i] = attrMap
		}
		hardwareMap["extension_attributes"] = extAttrs
	}

	return d.Set("hardware", []any{hardwareMap})
}

// setLocalUserAccountsSection maps the 'localUserAccounts' section of the computer inventory response to the Terraform resource data and updates the state.
func setLocalUserAccountsSection(d *schema.ResourceData, localUserAccounts []jamfpro.ComputerInventorySubsetLocalUserAccount) error {
	accounts := make([]any, len(localUserAccounts))
	for i, account := range localUserAccounts {
		acc := make(map[string]any)
		acc["uid"] = account.UID
		acc["user_guid"] = account.UserGuid
		acc["username"] = account.Username
		acc["full_name"] = account.FullName
		acc["admin"] = account.Admin
		acc["home_directory"] = account.HomeDirectory
		acc["home_directory_size_mb"] = account.HomeDirectorySizeMb
		acc["file_vault2_enabled"] = account.FileVault2Enabled
		acc["user_account_type"] = account.UserAccountType
		acc["password_min_length"] = account.PasswordMinLength
		acc["password_max_age"] = account.PasswordMaxAge
		acc["password_min_complex_characters"] = account.PasswordMinComplexCharacters
		acc["password_history_depth"] = account.PasswordHistoryDepth
		acc["password_require_alphanumeric"] = account.PasswordRequireAlphanumeric
		acc["computer_azure_active_directory_id"] = account.ComputerAzureActiveDirectoryId
		acc["user_azure_active_directory_id"] = account.UserAzureActiveDirectoryId
		acc["azure_active_directory_id"] = account.AzureActiveDirectoryId
		accounts[i] = acc
	}
	return d.Set("local_user_accounts", accounts)
}

// setCertificatesSection maps the 'certificate' section of the computer inventory response to the Terraform resource data and updates the state.
func setCertificatesSection(d *schema.ResourceData, certificates []jamfpro.ComputerInventorySubsetCertificate) error {
	certs := make([]any, len(certificates))
	for i, cert := range certificates {
		certMap := make(map[string]any)
		certMap["common_name"] = cert.CommonName
		certMap["identity"] = cert.Identity
		certMap["expiration_date"] = cert.ExpirationDate
		certMap["username"] = cert.Username
		certMap["lifecycle_status"] = cert.LifecycleStatus
		certMap["certificate_status"] = cert.CertificateStatus
		certMap["subject_name"] = cert.SubjectName
		certMap["serial_number"] = cert.SerialNumber
		certMap["sha1_fingerprint"] = cert.Sha1Fingerprint
		certMap["issued_date"] = cert.IssuedDate
		certs[i] = certMap
	}
	return d.Set("certificates", certs)
}

// setAttachmentsSection maps the 'attachments' section of the computer inventory response to the Terraform resource data and updates the state.
func setAttachmentsSection(d *schema.ResourceData, attachments []jamfpro.ComputerInventorySubsetAttachment) error {
	atts := make([]any, len(attachments))
	for i, att := range attachments {
		attMap := make(map[string]any)
		attMap["id"] = att.ID
		attMap["name"] = att.Name
		attMap["file_type"] = att.FileType
		attMap["size_bytes"] = att.SizeBytes
		atts[i] = attMap
	}
	return d.Set("attachments", atts)
}

// setPluginsSection maps the 'plugins' section of the computer inventory response to the Terraform resource data and updates the state.
func setPluginsSection(d *schema.ResourceData, plugins []jamfpro.ComputerInventorySubsetPlugin) error {
	pluginList := make([]any, len(plugins))
	for i, plugin := range plugins {
		pluginMap := make(map[string]any)
		pluginMap["name"] = plugin.Name
		pluginMap["version"] = plugin.Version
		pluginMap["path"] = plugin.Path
		pluginList[i] = pluginMap
	}
	return d.Set("plugins", pluginList)
}

// setPackageReceiptsSection maps the 'package receipts' section of the computer inventory response to the Terraform resource data and updates the state.
func setPackageReceiptsSection(d *schema.ResourceData, packageReceipts jamfpro.ComputerInventorySubsetPackageReceipts) error {
	packageReceiptMap := make(map[string]any)
	packageReceiptMap["installed_by_jamf_pro"] = packageReceipts.InstalledByJamfPro
	packageReceiptMap["installed_by_installer_swu"] = packageReceipts.InstalledByInstallerSwu
	packageReceiptMap["cached"] = packageReceipts.Cached
	return d.Set("package_receipts", []any{packageReceiptMap})
}

// setFontsSection maps the 'fonts' section of the computer inventory response to the Terraform resource data and updates the state.
func setFontsSection(d *schema.ResourceData, fonts []jamfpro.ComputerInventorySubsetFont) error {
	fontsList := make([]any, len(fonts))
	for i, font := range fonts {
		fontMap := make(map[string]any)
		fontMap["name"] = font.Name
		fontMap["version"] = font.Version
		fontMap["path"] = font.Path
		fontsList[i] = fontMap
	}
	return d.Set("fonts", fontsList)
}

// setSecuritySection maps the 'security' section of the computer inventory response to the Terraform resource data and updates the state.
func setSecuritySection(d *schema.ResourceData, security jamfpro.ComputerInventorySubsetSecurity) error {
	securityMap := make(map[string]any)
	securityMap["sip_status"] = security.SipStatus
	securityMap["gatekeeper_status"] = security.GatekeeperStatus
	securityMap["xprotect_version"] = security.XprotectVersion
	securityMap["auto_login_disabled"] = security.AutoLoginDisabled
	securityMap["remote_desktop_enabled"] = security.RemoteDesktopEnabled
	securityMap["activation_lock_enabled"] = security.ActivationLockEnabled
	securityMap["recovery_lock_enabled"] = security.RecoveryLockEnabled
	securityMap["firewall_enabled"] = security.FirewallEnabled
	securityMap["secure_boot_level"] = security.SecureBootLevel
	securityMap["external_boot_level"] = security.ExternalBootLevel
	securityMap["bootstrap_token_allowed"] = security.BootstrapTokenAllowed
	return d.Set("security", []any{securityMap})
}

// setOperatingSystemSection maps the 'Operating System' section of the computer inventory response to the Terraform resource data and updates the state.
func setOperatingSystemSection(d *schema.ResourceData, operatingSystem jamfpro.ComputerInventorySubsetOperatingSystem) error {
	osMap := make(map[string]any)
	osMap["name"] = operatingSystem.Name
	osMap["version"] = operatingSystem.Version
	osMap["build"] = operatingSystem.Build
	osMap["supplemental_build_version"] = operatingSystem.SupplementalBuildVersion
	osMap["rapid_security_response"] = operatingSystem.RapidSecurityResponse
	osMap["active_directory_status"] = operatingSystem.ActiveDirectoryStatus
	osMap["filevault2_status"] = operatingSystem.FileVault2Status
	osMap["software_update_device_id"] = operatingSystem.SoftwareUpdateDeviceId

	extAttrs := make([]map[string]any, len(operatingSystem.ExtensionAttributes))
	for i, attr := range operatingSystem.ExtensionAttributes {
		attrMap := make(map[string]any)
		attrMap["definition_id"] = attr.DefinitionId
		attrMap["name"] = attr.Name
		attrMap["description"] = attr.Description
		attrMap["enabled"] = attr.Enabled
		attrMap["multi_value"] = attr.MultiValue
		attrMap["values"] = attr.Values
		attrMap["data_type"] = attr.DataType
		attrMap["options"] = attr.Options
		attrMap["input_type"] = attr.InputType

		extAttrs[i] = attrMap
	}
	osMap["extension_attributes"] = extAttrs
	return d.Set("operating_system", []any{osMap})
}

// setLicensedSoftwareSection maps the 'Licensed Software' section of the computer inventory response to the Terraform resource data and updates the state.
func setLicensedSoftwareSection(d *schema.ResourceData, licensedSoftware []jamfpro.ComputerInventorySubsetLicensedSoftware) error {
	softwareList := make([]any, len(licensedSoftware))
	for i, software := range licensedSoftware {
		softwareMap := make(map[string]any)
		softwareMap["id"] = software.ID
		softwareMap["name"] = software.Name
		softwareList[i] = softwareMap
	}
	return d.Set("licensed_software", softwareList)
}

// setIBeaconsSection maps the 'IBeacons' section of the computer inventory response to the Terraform resource data and updates the state.
func setIBeaconsSection(d *schema.ResourceData, ibeacons []jamfpro.ComputerInventorySubsetIBeacon) error {
	ibeaconList := make([]any, len(ibeacons))
	for i, ibeacon := range ibeacons {
		ibeaconMap := make(map[string]any)
		ibeaconMap["name"] = ibeacon.Name
		ibeaconList[i] = ibeaconMap
	}
	return d.Set("ibeacons", ibeaconList)
}

// setSoftwareUpdatesSection maps the 'Software Updates' section of the computer inventory response to the Terraform resource data and updates the state.
func setSoftwareUpdatesSection(d *schema.ResourceData, softwareUpdates []jamfpro.ComputerInventorySubsetSoftwareUpdate) error {
	updateList := make([]any, len(softwareUpdates))
	for i, update := range softwareUpdates {
		updateMap := make(map[string]any)
		updateMap["name"] = update.Name
		updateMap["version"] = update.Version
		updateMap["package_name"] = update.PackageName
		updateList[i] = updateMap
	}
	return d.Set("software_updates", updateList)
}

// setExtensionAttributesSection maps the 'Extension Attributes' section of the computer inventory response to the Terraform resource data and updates the state.
func setExtensionAttributesSection(d *schema.ResourceData, extensionAttributes []jamfpro.ComputerInventorySubsetExtensionAttribute) error {
	attrList := make([]any, len(extensionAttributes))
	for i, attr := range extensionAttributes {
		attrMap := make(map[string]any)
		attrMap["definition_id"] = attr.DefinitionId
		attrMap["name"] = attr.Name
		attrMap["description"] = attr.Description
		attrMap["enabled"] = attr.Enabled
		attrMap["multi_value"] = attr.MultiValue
		attrMap["values"] = attr.Values
		attrMap["data_type"] = attr.DataType
		attrMap["options"] = attr.Options
		attrMap["input_type"] = attr.InputType
		attrList[i] = attrMap
	}
	return d.Set("extension_attributes", attrList)
}

// setGroupMembershipsSection maps the 'groupMemberships' section of the computer inventory response to the Terraform resource data and updates the state.
func setGroupMembershipsSection(d *schema.ResourceData, groupMemberships []jamfpro.ComputerInventorySubsetGroupMembership) error {
	memberships := make([]any, len(groupMemberships))
	for i, group := range groupMemberships {
		groupMap := make(map[string]any)
		groupMap["group_id"] = group.GroupId
		groupMap["group_name"] = group.GroupName
		groupMap["smart_group"] = group.SmartGroup

		memberships[i] = groupMap
	}
	return d.Set("group_memberships", memberships)
}

// setConfigurationProfilesSection maps the 'configurationProfiles' section of the computer inventory response to the Terraform resource data and updates the state.
func setConfigurationProfilesSection(d *schema.ResourceData, configurationProfiles []jamfpro.ComputerInventorySubsetConfigurationProfile) error {
	profiles := make([]any, len(configurationProfiles))
	for i, computer := range configurationProfiles {
		profileMap := map[string]interface{}{
			"id":                 computer.ID,
			"username":           computer.Username,
			"last_installed":     computer.LastInstalled,
			"removable":          computer.Removable,
			"display_name":       computer.DisplayName,
			"profile_identifier": computer.ProfileIdentifier,
		}
		profiles[i] = profileMap
	}
	return d.Set("configuration_profiles", profiles)
}

// setPrintersSection maps the 'printers' section of the computer inventory response to the Terraform resource data and updates the state.
func setPrintersSection(d *schema.ResourceData, printers []jamfpro.ComputerInventorySubsetPrinter) error {
	printerList := make([]any, len(printers))
	for i, printer := range printers {
		printerMap := map[string]interface{}{
			"name":     printer.Name,
			"type":     printer.Type,
			"uri":      printer.URI,
			"location": printer.Location,
		}
		printerList[i] = printerMap
	}
	return d.Set("printers", printerList)
}

// setServicesSection maps the 'services' section of the computer inventory response to the Terraform resource data and updates the state.
func setServicesSection(d *schema.ResourceData, services []jamfpro.ComputerInventorySubsetService) error {
	serviceList := make([]any, len(services))
	for i, service := range services {
		serviceMap := map[string]interface{}{
			"name": service.Name,
		}
		serviceList[i] = serviceMap
	}
	return d.Set("services", serviceList)
}
