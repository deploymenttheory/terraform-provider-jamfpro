package computer_inventory

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead fetches the details of a specific macOS Configuration Profile
// from Jamf Pro using either its unique Name or its ID. The function prioritizes the 'name' attribute over the 'id'
// attribute for fetching details. If neither 'name' nor 'id' is provided, it returns an error.
// Once the details are fetched, they are set in the data source's state.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Asserts 'meta' as '*client.client'
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	var profile *jamfpro.ResourceComputerInventory
	var err error

	// Fetch profile by 'name' or 'id'
	if v, ok := d.GetOk("name"); ok {
		profileName, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'name' as string")
		}
		profile, err = client.GetComputerInventoryByName(profileName)
	} else if v, ok := d.GetOk("id"); ok {
		profileID, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'id' as string")
		}
		profile, err = client.GetComputerInventoryByID(profileID)
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch macOS Configuration Profile: %v", err))
	}

	// Set top-level attributes
	d.SetId(profile.ID)
	d.Set("id", profile.ID)
	d.Set("udid", profile.UDID)

	// Set 'general' section
	if err := setGeneralSection(d, profile.General); err != nil {
		return diag.FromErr(err)
	}

	// Set 'diskEncryption' section
	if err := setDiskEncryptionSection(d, profile.DiskEncryption); err != nil {
		return diag.FromErr(err)
	}

	// Set 'purchasing' section
	if err := setPurchasingSection(d, profile.Purchasing); err != nil {
		return diag.FromErr(err)
	}

	// Set 'applications' section
	if err := setApplicationsSection(d, profile.Applications); err != nil {
		return diag.FromErr(err)
	}

	// Set 'storage' section
	if err := setStorageSection(d, profile.Storage); err != nil {
		return diag.FromErr(err)
	}

	// Set 'userAndLocation' section
	if err := setUserAndLocationSection(d, profile.UserAndLocation); err != nil {
		return diag.FromErr(err)
	}

	// Set 'hardware' section
	if err := setHardwareSection(d, profile.Hardware); err != nil {
		return diag.FromErr(err)
	}

	// Set 'localUserAccounts' section
	if err := setLocalUserAccountsSection(d, profile.LocalUserAccounts); err != nil {
		return diag.FromErr(err)
	}

	// Set 'certificates' section
	if err := setCertificatesSection(d, profile.Certificates); err != nil {
		return diag.FromErr(err)
	}

	// Set 'attachments' section
	if err := setAttachmentsSection(d, profile.Attachments); err != nil {
		return diag.FromErr(err)
	}

	// Set 'plugins' section
	if err := setPluginsSection(d, profile.Plugins); err != nil {
		return diag.FromErr(err)
	}

	// Set 'packageReceipts' section
	if err := setPackageReceiptsSection(d, profile.PackageReceipts); err != nil {
		return diag.FromErr(err)
	}

	// Set 'fonts' section
	if err := setFontsSection(d, profile.Fonts); err != nil {
		return diag.FromErr(err)
	}

	// Set 'security' section
	if err := setSecuritySection(d, profile.Security); err != nil {
		return diag.FromErr(err)
	}

	// Set 'operatingSystem' section
	if err := setOperatingSystemSection(d, profile.OperatingSystem); err != nil {
		return diag.FromErr(err)
	}

	// Set 'licensedSoftware' section
	if err := setLicensedSoftwareSection(d, profile.LicensedSoftware); err != nil {
		return diag.FromErr(err)
	}

	// Set 'ibeacons' section
	if err := setIBeaconsSection(d, profile.Ibeacons); err != nil {
		return diag.FromErr(err)
	}

	// Set 'softwareUpdates' section
	if err := setSoftwareUpdatesSection(d, profile.SoftwareUpdates); err != nil {
		return diag.FromErr(err)
	}

	// Set 'extensionAttributes' section
	if err := setExtensionAttributesSection(d, profile.ExtensionAttributes); err != nil {
		return diag.FromErr(err)
	}

	// Set 'groupMemberships' section
	if err := setGroupMembershipsSection(d, profile.GroupMemberships); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// setGeneralSection maps the 'general' section of the computer inventory response to the Terraform resource data and updates the state.
func setGeneralSection(d *schema.ResourceData, general jamfpro.ComputerInventorySubsetGeneral) error {
	// Initialize a map to hold the 'general' section attributes.
	gen := make(map[string]interface{})

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
	gen["mdm_capable"] = general.MdmCapable.Capable
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

	// Handle nested object 'remoteManagement'.
	remoteManagement := make(map[string]interface{})
	remoteManagement["managed"] = general.RemoteManagement.Managed
	remoteManagement["management_username"] = general.RemoteManagement.ManagementUsername
	gen["remote_management"] = []interface{}{remoteManagement}

	// Handle nested object 'site'.
	if general.Site.ID != "" || general.Site.Name != "" {
		site := make(map[string]interface{})
		site["id"] = general.Site.ID
		site["name"] = general.Site.Name
		gen["site_id"] = []interface{}{site}
	}

	// Handle nested object 'enrollmentMethod'.
	if general.EnrollmentMethod.ID != "" || general.EnrollmentMethod.ObjectName != "" || general.EnrollmentMethod.ObjectType != "" {
		enrollmentMethod := make(map[string]interface{})
		enrollmentMethod["id"] = general.EnrollmentMethod.ID
		enrollmentMethod["object_name"] = general.EnrollmentMethod.ObjectName
		enrollmentMethod["object_type"] = general.EnrollmentMethod.ObjectType
		gen["enrollment_method"] = []interface{}{enrollmentMethod}
	}

	// Set the 'general' section in the Terraform resource data.
	return d.Set("general", []interface{}{gen})
}

// setDiskEncryptionSection maps the 'diskEncryption' section of the computer inventory response to the Terraform resource data and updates the state.
func setDiskEncryptionSection(d *schema.ResourceData, diskEncryption jamfpro.ComputerInventorySubsetDiskEncryption) error {
	// Initialize a map to hold the 'diskEncryption' section attributes.
	diskEnc := make(map[string]interface{})

	// Map each attribute of the 'diskEncryption' section from the API response to the corresponding Terraform schema attribute.
	diskEnc["individual_recovery_key_validity_status"] = diskEncryption.IndividualRecoveryKeyValidityStatus
	diskEnc["institutional_recovery_key_present"] = diskEncryption.InstitutionalRecoveryKeyPresent
	diskEnc["disk_encryption_configuration_name"] = diskEncryption.DiskEncryptionConfigurationName
	diskEnc["file_vault2_eligibility_message"] = diskEncryption.FileVault2EligibilityMessage

	// Handle nested object 'bootPartitionEncryptionDetails'.
	bootPartitionDetails := make(map[string]interface{})
	bootPartitionDetails["partition_name"] = diskEncryption.BootPartitionEncryptionDetails.PartitionName
	bootPartitionDetails["partition_file_vault2_state"] = diskEncryption.BootPartitionEncryptionDetails.PartitionFileVault2State
	bootPartitionDetails["partition_file_vault2_percent"] = diskEncryption.BootPartitionEncryptionDetails.PartitionFileVault2Percent
	diskEnc["boot_partition_encryption_details"] = []interface{}{bootPartitionDetails}

	// Map 'fileVault2EnabledUserNames' as a list of strings.
	fileVaultUserNames := make([]string, len(diskEncryption.FileVault2EnabledUserNames))
	copy(fileVaultUserNames, diskEncryption.FileVault2EnabledUserNames)

	// Set 'fileVault2EnabledUserNames' in the 'diskEnc' map.
	diskEnc["file_vault2_enabled_user_names"] = fileVaultUserNames

	// Set the 'diskEncryption' section in the Terraform resource data.
	return d.Set("disk_encryption", []interface{}{diskEnc})
}

// setPurchasingSection maps the 'purchasing' section of the computer inventory response to the Terraform resource data and updates the state.
func setPurchasingSection(d *schema.ResourceData, purchasing jamfpro.ComputerInventorySubsetPurchasing) error {
	// Initialize a map to hold the 'purchasing' section attributes.
	purchasingMap := make(map[string]interface{})

	// Map each attribute of the 'purchasing' section from the API response to the corresponding Terraform schema attribute.
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

	// Map 'extensionAttributes' as a list of maps.
	extAttrs := make([]map[string]interface{}, len(purchasing.ExtensionAttributes))
	for i, attr := range purchasing.ExtensionAttributes {
		attrMap := make(map[string]interface{})
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

	// Set the 'purchasing' section in the Terraform resource data.
	return d.Set("purchasing", []interface{}{purchasingMap})
}

// setApplicationsSection maps the 'applications' section of the computer inventory response to the Terraform resource data and updates the state.
func setApplicationsSection(d *schema.ResourceData, applications []jamfpro.ComputerInventorySubsetApplication) error {
	// Create a slice to hold the application maps.
	apps := make([]interface{}, len(applications))

	for i, app := range applications {
		// Initialize a map for each application.
		appMap := make(map[string]interface{})

		// Map each attribute of the application from the API response to the corresponding Terraform schema attribute.
		appMap["name"] = app.Name
		appMap["path"] = app.Path
		appMap["version"] = app.Version
		appMap["mac_app_store"] = app.MacAppStore
		appMap["size_megabytes"] = app.SizeMegabytes
		appMap["bundle_id"] = app.BundleId
		appMap["update_available"] = app.UpdateAvailable
		appMap["external_version_id"] = app.ExternalVersionId

		// Add the application map to the slice.
		apps[i] = appMap
	}

	// Set the 'applications' section in the Terraform resource data.
	return d.Set("applications", apps)
}

// setStorageSection maps the 'storage' section of the computer inventory response to the Terraform resource data and updates the state.
func setStorageSection(d *schema.ResourceData, storage jamfpro.ComputerInventorySubsetStorage) error {
	storageMap := make(map[string]interface{})

	storageMap["boot_drive_available_space_megabytes"] = storage.BootDriveAvailableSpaceMegabytes

	// Mapping 'disks' array
	disks := make([]interface{}, len(storage.Disks))
	for i, disk := range storage.Disks {
		diskMap := make(map[string]interface{})
		diskMap["id"] = disk.ID
		diskMap["device"] = disk.Device
		diskMap["model"] = disk.Model
		diskMap["revision"] = disk.Revision
		diskMap["serial_number"] = disk.SerialNumber
		diskMap["size_megabytes"] = disk.SizeMegabytes
		diskMap["smart_status"] = disk.SmartStatus
		diskMap["type"] = disk.Type

		// Map 'partitions' if present
		partitions := make([]interface{}, len(disk.Partitions))
		for j, partition := range disk.Partitions {
			partitionMap := make(map[string]interface{})
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

	// Set the 'storage' section in the Terraform resource data.
	return d.Set("storage", []interface{}{storageMap})
}

// setUserAndLocationSection maps the 'userAndLocation' section of the computer inventory response to the Terraform resource data and updates the state.
func setUserAndLocationSection(d *schema.ResourceData, userAndLocation jamfpro.ComputerInventorySubsetUserAndLocation) error {
	userLocationMap := make(map[string]interface{})

	// Map each attribute from the 'userAndLocation' object to the corresponding schema attribute
	userLocationMap["username"] = userAndLocation.Username
	userLocationMap["realname"] = userAndLocation.Realname
	userLocationMap["email"] = userAndLocation.Email
	userLocationMap["position"] = userAndLocation.Position
	userLocationMap["phone"] = userAndLocation.Phone
	userLocationMap["department_id"] = userAndLocation.DepartmentId
	userLocationMap["building_id"] = userAndLocation.BuildingId
	userLocationMap["room"] = userAndLocation.Room

	// Map extension attributes if present
	if len(userAndLocation.ExtensionAttributes) > 0 {
		extAttrs := make([]map[string]interface{}, len(userAndLocation.ExtensionAttributes))
		for i, attr := range userAndLocation.ExtensionAttributes {
			attrMap := make(map[string]interface{})
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

	// Set the 'userAndLocation' section in the Terraform resource data
	return d.Set("user_and_location", []interface{}{userLocationMap})
}

// setHardwareSection maps the 'hardware' section of the computer inventory response to the Terraform resource data and updates the state.
func setHardwareSection(d *schema.ResourceData, hardware jamfpro.ComputerInventorySubsetHardware) error {
	hardwareMap := make(map[string]interface{})

	// Map each attribute from the 'hardware' object to the corresponding schema attribute
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

	// Map extension attributes if present
	if len(hardware.ExtensionAttributes) > 0 {
		extAttrs := make([]map[string]interface{}, len(hardware.ExtensionAttributes))
		for i, attr := range hardware.ExtensionAttributes {
			attrMap := make(map[string]interface{})
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

	// Set the 'hardware' section in the Terraform resource data
	return d.Set("hardware", []interface{}{hardwareMap})
}

// setLocalUserAccountsSection maps the 'localUserAccounts' section of the computer inventory response to the Terraform resource data and updates the state.
func setLocalUserAccountsSection(d *schema.ResourceData, localUserAccounts []jamfpro.ComputerInventorySubsetLocalUserAccount) error {
	accounts := make([]interface{}, len(localUserAccounts))
	for i, account := range localUserAccounts {
		acc := make(map[string]interface{})
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
	return d.Set("localUserAccounts", accounts)
}

// setCertificatesSection maps the 'certificate' section of the computer inventory response to the Terraform resource data and updates the state.
func setCertificatesSection(d *schema.ResourceData, certificates []jamfpro.ComputerInventorySubsetCertificate) error {
	certs := make([]interface{}, len(certificates))
	for i, cert := range certificates {
		certMap := make(map[string]interface{})
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
	atts := make([]interface{}, len(attachments))
	for i, att := range attachments {
		attMap := make(map[string]interface{})
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
	pluginList := make([]interface{}, len(plugins))
	for i, plugin := range plugins {
		pluginMap := make(map[string]interface{})
		pluginMap["name"] = plugin.Name
		pluginMap["version"] = plugin.Version
		pluginMap["path"] = plugin.Path
		pluginList[i] = pluginMap
	}
	return d.Set("plugins", pluginList)
}

// setPackageReceiptsSection maps the 'package receipts' section of the computer inventory response to the Terraform resource data and updates the state.
func setPackageReceiptsSection(d *schema.ResourceData, packageReceipts jamfpro.ComputerInventorySubsetPackageReceipts) error {
	packageReceiptMap := make(map[string]interface{})
	packageReceiptMap["installed_by_jamf_pro"] = packageReceipts.InstalledByJamfPro
	packageReceiptMap["installed_by_installer_swu"] = packageReceipts.InstalledByInstallerSwu
	packageReceiptMap["cached"] = packageReceipts.Cached
	return d.Set("package_receipts", []interface{}{packageReceiptMap})
}

// setFontsSection maps the 'fonts' section of the computer inventory response to the Terraform resource data and updates the state.
func setFontsSection(d *schema.ResourceData, fonts []jamfpro.ComputerInventorySubsetFont) error {
	fontsList := make([]interface{}, len(fonts))
	for i, font := range fonts {
		fontMap := make(map[string]interface{})
		fontMap["name"] = font.Name
		fontMap["version"] = font.Version
		fontMap["path"] = font.Path
		fontsList[i] = fontMap
	}
	return d.Set("fonts", fontsList)
}

// setSecuritySection maps the 'security' section of the computer inventory response to the Terraform resource data and updates the state.
func setSecuritySection(d *schema.ResourceData, security jamfpro.ComputerInventorySubsetSecurity) error {
	securityMap := make(map[string]interface{})
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
	return d.Set("security", []interface{}{securityMap})
}

// setOperatingSystemSection maps the 'Operating System' section of the computer inventory response to the Terraform resource data and updates the state.
func setOperatingSystemSection(d *schema.ResourceData, operatingSystem jamfpro.ComputerInventorySubsetOperatingSystem) error {
	osMap := make(map[string]interface{})
	osMap["name"] = operatingSystem.Name
	osMap["version"] = operatingSystem.Version
	osMap["build"] = operatingSystem.Build
	osMap["supplemental_build_version"] = operatingSystem.SupplementalBuildVersion
	osMap["rapid_security_response"] = operatingSystem.RapidSecurityResponse
	osMap["active_directory_status"] = operatingSystem.ActiveDirectoryStatus
	osMap["filevault2_status"] = operatingSystem.FileVault2Status
	osMap["softwareUpdate_device_id"] = operatingSystem.SoftwareUpdateDeviceId
	// Map extension attributes if present
	extAttrs := make([]map[string]interface{}, len(operatingSystem.ExtensionAttributes))
	for i, attr := range operatingSystem.ExtensionAttributes {
		attrMap := make(map[string]interface{})
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
	return d.Set("operating_system", []interface{}{osMap})
}

// setLicensedSoftwareSection maps the 'Licensed Software' section of the computer inventory response to the Terraform resource data and updates the state.
func setLicensedSoftwareSection(d *schema.ResourceData, licensedSoftware []jamfpro.ComputerInventorySubsetLicensedSoftware) error {
	softwareList := make([]interface{}, len(licensedSoftware))
	for i, software := range licensedSoftware {
		softwareMap := make(map[string]interface{})
		softwareMap["id"] = software.ID
		softwareMap["name"] = software.Name
		softwareList[i] = softwareMap
	}
	return d.Set("licensed_software", softwareList)
}

// setIBeaconsSection maps the 'IBeacons' section of the computer inventory response to the Terraform resource data and updates the state.
func setIBeaconsSection(d *schema.ResourceData, ibeacons []jamfpro.ComputerInventorySubsetIBeacon) error {
	ibeaconList := make([]interface{}, len(ibeacons))
	for i, ibeacon := range ibeacons {
		ibeaconMap := make(map[string]interface{})
		ibeaconMap["name"] = ibeacon.Name
		ibeaconList[i] = ibeaconMap
	}
	return d.Set("ibeacons", ibeaconList)
}

// setSoftwareUpdatesSection maps the 'Software Updates' section of the computer inventory response to the Terraform resource data and updates the state.
func setSoftwareUpdatesSection(d *schema.ResourceData, softwareUpdates []jamfpro.ComputerInventorySubsetSoftwareUpdate) error {
	updateList := make([]interface{}, len(softwareUpdates))
	for i, update := range softwareUpdates {
		updateMap := make(map[string]interface{})
		updateMap["name"] = update.Name
		updateMap["version"] = update.Version
		updateMap["package_name"] = update.PackageName
		updateList[i] = updateMap
	}
	return d.Set("software_updates", updateList)
}

// setExtensionAttributesSection maps the 'Extension Attributes' section of the computer inventory response to the Terraform resource data and updates the state.
func setExtensionAttributesSection(d *schema.ResourceData, extensionAttributes []jamfpro.ComputerInventorySubsetExtensionAttribute) error {
	attrList := make([]interface{}, len(extensionAttributes))
	for i, attr := range extensionAttributes {
		attrMap := make(map[string]interface{})
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
	memberships := make([]interface{}, len(groupMemberships))
	for i, group := range groupMemberships {
		groupMap := make(map[string]interface{})
		groupMap["group_id"] = group.GroupId
		groupMap["group_name"] = group.GroupName
		groupMap["smart_group"] = group.SmartGroup

		memberships[i] = groupMap
	}
	return d.Set("group_memberships", memberships)
}
