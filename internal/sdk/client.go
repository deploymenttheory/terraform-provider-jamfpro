// Package sdk is the provider's boundary with the Jamf Pro SDK. Provider code
// depends on the Client interface here and does not import the SDK module.
package sdk

import (
	"context"
	"net/url"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

var _ Client = (*jamfpro.Client)(nil)

// Client is the set of Jamf Pro API calls the provider makes. Methods are
// grouped by the SDK source file that implements them.
//
// To use an SDK call that is not listed, add its method here with the SDK's
// signature.
type Client interface {
	// classicapi_accounts.go
	GetAccountByID(id string) (*jamfpro.ResourceAccount, error)
	GetAccountByName(name string) (*jamfpro.ResourceAccount, error)
	CreateAccount(account *jamfpro.ResourceAccount) (*jamfpro.ResponseAccountCreatedAndUpdated, error)
	UpdateAccountByID(id string, account *jamfpro.ResourceAccount) (*jamfpro.ResponseAccountCreatedAndUpdated, error)
	DeleteAccountByID(id string) error

	// classicapi_accounts_groups.go
	GetAccountGroupByID(id string) (*jamfpro.ResourceAccountGroup, error)
	GetAccountGroupByName(name string) (*jamfpro.ResourceAccountGroup, error)
	CreateAccountGroup(accountGroup *jamfpro.ResourceAccountGroup) (*jamfpro.ResponseAccountGroupCreated, error)
	UpdateAccountGroupByID(id string, accountGroup *jamfpro.ResourceAccountGroup) (*jamfpro.ResourceAccountGroup, error)
	DeleteAccountGroupByID(id string) error

	// classicapi_activation_code.go
	GetActivationCode() (*jamfpro.ResourceActivationCode, error)
	UpdateActivationCode(activationCode *jamfpro.ResourceActivationCode) error

	// classicapi_advanced_computer_searches.go
	GetAdvancedComputerSearchByID(id string) (*jamfpro.ResourceAdvancedComputerSearch, error)
	CreateAdvancedComputerSearch(search *jamfpro.ResourceAdvancedComputerSearch) (*jamfpro.ResponseAdvancedComputerSearchCreatedAndUpdated, error)
	UpdateAdvancedComputerSearchByID(id string, search *jamfpro.ResourceAdvancedComputerSearch) (*jamfpro.ResponseAdvancedComputerSearchCreatedAndUpdated, error)
	DeleteAdvancedComputerSearchByID(id string) error

	// classicapi_advanced_user_searches.go
	GetAdvancedUserSearchByID(id string) (*jamfpro.ResourceAdvancedUserSearch, error)
	CreateAdvancedUserSearch(search *jamfpro.ResourceAdvancedUserSearch) (*jamfpro.ResponseAdvancedUserSearchCreatedAndUpdated, error)
	UpdateAdvancedUserSearchByID(id string, search *jamfpro.ResourceAdvancedUserSearch) (*jamfpro.ResponseAdvancedUserSearchCreatedAndUpdated, error)
	DeleteAdvancedUserSearchByID(id string) error

	// classicapi_allowed_file_extensions.go
	GetAllowedFileExtensionByID(id string) (*jamfpro.ResourceAllowedFileExtension, error)
	CreateAllowedFileExtension(extension *jamfpro.ResourceAllowedFileExtension) (*jamfpro.ResourceAllowedFileExtension, error)
	DeleteAllowedFileExtensionByID(id string) error

	// classicapi_computer_groups.go
	GetComputerGroupByID(id string) (*jamfpro.ResourceComputerGroup, error)
	GetComputerGroupByName(name string) (*jamfpro.ResourceComputerGroup, error)
	CreateComputerGroup(group *jamfpro.ResourceComputerGroup) (*jamfpro.ResponseComputerGroupreatedAndUpdated, error)
	UpdateComputerGroupByID(id string, group *jamfpro.ResourceComputerGroup) (*jamfpro.ResponseComputerGroupreatedAndUpdated, error)
	DeleteComputerGroupByID(id string) error

	// classicapi_disk_encryption_configurations.go
	GetDiskEncryptionConfigurationByID(id string) (*jamfpro.ResourceDiskEncryptionConfiguration, error)
	GetDiskEncryptionConfigurationByName(name string) (*jamfpro.ResourceDiskEncryptionConfiguration, error)
	CreateDiskEncryptionConfiguration(config *jamfpro.ResourceDiskEncryptionConfiguration) (*jamfpro.ResponseDiskEncryptionConfigurationCreatedAndUpdated, error)
	UpdateDiskEncryptionConfigurationByID(id string, config *jamfpro.ResourceDiskEncryptionConfiguration) (*jamfpro.ResponseDiskEncryptionConfigurationCreatedAndUpdated, error)
	DeleteDiskEncryptionConfigurationByID(id string) error

	// classicapi_dock_items.go
	GetDockItemByID(id string) (*jamfpro.ResourceDockItem, error)
	GetDockItemByName(name string) (*jamfpro.ResourceDockItem, error)
	CreateDockItem(dockItem *jamfpro.ResourceDockItem) (*jamfpro.ResourceDockItem, error)
	UpdateDockItemByID(id string, dockItem *jamfpro.ResourceDockItem) (*jamfpro.ResourceDockItem, error)
	DeleteDockItemByID(id string) error

	// classicapi_file_share_distribution_points.go
	GetDistributionPointByID(id string) (*jamfpro.ResourceFileShareDistributionPoint, error)
	CreateDistributionPoint(dp *jamfpro.ResourceFileShareDistributionPoint) (*jamfpro.ResponseFileShareDistributionPointCreatedAndUpdated, error)
	UpdateDistributionPointByID(id string, dp *jamfpro.ResourceFileShareDistributionPoint) (*jamfpro.ResponseFileShareDistributionPointCreatedAndUpdated, error)
	DeleteDistributionPointByID(id string) error

	// classicapi_ldap_servers.go
	GetLDAPServerByID(id string) (*jamfpro.ResourceLDAPServers, error)
	CreateLDAPServer(ldapServer *jamfpro.ResourceLDAPServers) (*jamfpro.LDAPServersListItem, error)
	UpdateLDAPServerByID(id string, ldapServer *jamfpro.ResourceLDAPServers) (*jamfpro.ResourceLDAPServers, error)
	DeleteLDAPServerByID(id string) error

	// classicapi_mac_applications.go
	GetMacApplicationByID(id string) (*jamfpro.ResourceMacApplications, error)
	GetMacApplicationByName(name string) (*jamfpro.ResourceMacApplications, error)
	CreateMacApplication(macApp *jamfpro.ResourceMacApplications) (*jamfpro.MacApplicationsListItem, error)
	UpdateMacApplicationByID(id string, macApp *jamfpro.ResourceMacApplications) (*jamfpro.ResourceMacApplications, error)
	DeleteMacApplicationByID(id string) error

	// classicapi_macos_configuration_profiles.go
	GetMacOSConfigurationProfileByID(id string) (*jamfpro.ResourceMacOSConfigurationProfile, error)
	GetMacOSConfigurationProfileByName(name string) (*jamfpro.ResourceMacOSConfigurationProfile, error)
	CreateMacOSConfigurationProfile(profile *jamfpro.ResourceMacOSConfigurationProfile) (*jamfpro.ResponseMacOSConfigurationProfileCreationUpdate, error)
	UpdateMacOSConfigurationProfileByID(id string, profile *jamfpro.ResourceMacOSConfigurationProfile) (int, error)
	DeleteMacOSConfigurationProfileByID(id string) error
	DeleteMacOSConfigurationProfileByName(name string) error

	// classicapi_mobile_device_applications.go
	GetMobileDeviceApplicationByID(id string) (*jamfpro.ResourceMobileDeviceApplication, error)
	GetMobileDeviceApplicationByName(name string) (*jamfpro.ResourceMobileDeviceApplication, error)
	CreateMobileDeviceApplication(app *jamfpro.ResourceMobileDeviceApplication) (*jamfpro.MobileDeviceApplicationsListItem, error)
	UpdateMobileDeviceApplicationByID(id string, app *jamfpro.ResourceMobileDeviceApplication) (*jamfpro.ResourceMobileDeviceApplication, error)
	DeleteMobileDeviceApplicationpByID(id string) error

	// classicapi_mobile_device_configuration_profiles.go
	GetMobileDeviceConfigurationProfileByID(id string) (*jamfpro.ResourceMobileDeviceConfigurationProfile, error)
	CreateMobileDeviceConfigurationProfile(profile *jamfpro.ResourceMobileDeviceConfigurationProfile) (*jamfpro.ResponseMobileDeviceConfigurationProfileCreateAndUpdate, error)
	UpdateMobileDeviceConfigurationProfileByID(id string, profile *jamfpro.ResourceMobileDeviceConfigurationProfile) (*jamfpro.ResponseMobileDeviceConfigurationProfileCreateAndUpdate, error)
	DeleteMobileDeviceConfigurationProfileByID(id string) error

	// classicapi_mobile_device_groups.go
	GetMobileDeviceGroupByID(id string) (*jamfpro.ResourceMobileDeviceGroup, error)
	GetMobileDeviceGroupByName(name string) (*jamfpro.ResourceMobileDeviceGroup, error)
	CreateMobileDeviceGroup(group *jamfpro.ResourceMobileDeviceGroup) (*jamfpro.ResourceMobileDeviceGroup, error)
	UpdateMobileDeviceGroupByID(id string, group *jamfpro.ResourceMobileDeviceGroup) (*jamfpro.ResourceMobileDeviceGroup, error)
	DeleteMobileDeviceGroupByID(id string) error

	// classicapi_network_segments.go
	GetNetworkSegmentByID(id string) (*jamfpro.ResourceNetworkSegment, error)
	CreateNetworkSegment(segment *jamfpro.ResourceNetworkSegment) (*jamfpro.ResponseNetworkSegmentCreatedAndUpdated, error)
	UpdateNetworkSegmentByID(id string, segment *jamfpro.ResourceNetworkSegment) (*jamfpro.ResponseNetworkSegmentCreatedAndUpdated, error)
	DeleteNetworkSegmentByID(id string) error

	// classicapi_policies.go
	GetPolicyByID(id string) (*jamfpro.ResourcePolicy, error)
	CreatePolicy(policy *jamfpro.ResourcePolicy) (*jamfpro.ResponsePolicyCreateAndUpdate, error)
	UpdatePolicyByID(id string, policy *jamfpro.ResourcePolicy) (*jamfpro.ResponsePolicyCreateAndUpdate, error)
	DeletePolicyByID(id string) error

	// classicapi_printers.go
	GetPrinterByID(id string) (*jamfpro.ResourcePrinter, error)
	CreatePrinter(printer *jamfpro.ResourcePrinter) (*jamfpro.ResponsePrinterCreateAndUpdate, error)
	UpdatePrinterByID(id string, printer *jamfpro.ResourcePrinter) (*jamfpro.ResponsePrinterCreateAndUpdate, error)
	DeletePrinterByID(id string) error

	// classicapi_restricted_software.go
	GetRestrictedSoftwareByID(id string) (*jamfpro.ResourceRestrictedSoftware, error)
	CreateRestrictedSoftware(restrictedSoftware *jamfpro.ResourceRestrictedSoftware) (*jamfpro.ResponseRestrictedSoftwareCreateAndUpdate, error)
	UpdateRestrictedSoftwareByID(id string, restrictedSoftware *jamfpro.ResourceRestrictedSoftware) (*jamfpro.ResponseRestrictedSoftwareCreateAndUpdate, error)
	DeleteRestrictedSoftwareByID(id string) error

	// classicapi_sites.go
	GetSiteByID(id string) (*jamfpro.SharedResourceSite, error)
	GetSiteByName(name string) (*jamfpro.SharedResourceSite, error)
	CreateSite(site *jamfpro.SharedResourceSite) (*jamfpro.SharedResourceSite, error)
	UpdateSiteByID(id string, site *jamfpro.SharedResourceSite) (*jamfpro.SharedResourceSite, error)
	DeleteSiteByID(id string) error

	// classicapi_usergroups.go
	GetUserGroupByID(id string) (*jamfpro.ResourceUserGroup, error)
	CreateUserGroup(userGroup *jamfpro.ResourceUserGroup) (*jamfpro.ResponseUserGroupCreateAndUpdate, error)
	UpdateUserGroupByID(id string, userGroup *jamfpro.ResourceUserGroup) (*jamfpro.ResponseUserGroupCreateAndUpdate, error)
	DeleteUserGroupByID(id string) error

	// classicapi_users.go
	GetUsers() (*jamfpro.ResponseUsersList, error)
	GetUserByID(id string) (*jamfpro.ResourceUser, error)
	GetUserByName(name string) (*jamfpro.ResourceUser, error)
	GetUserByEmail(email string) (*jamfpro.ResourceUser, error)

	// classicapi_webhooks.go
	GetWebhookByID(id string) (*jamfpro.ResourceWebhook, error)
	CreateWebhook(webhook *jamfpro.ResourceWebhook) (*jamfpro.ResourceWebhook, error)
	UpdateWebhookByID(id string, webhook *jamfpro.ResourceWebhook) (*jamfpro.ResourceWebhook, error)
	DeleteWebhookByID(id string) error

	// jamfproapi_access_management_settings.go
	GetAccessManagementSettings() (*jamfpro.ResourceAccessManagementSettings, error)
	CreateAccessManagementSettings(settingsCreate jamfpro.ResourceAccessManagementSettings) (*jamfpro.ResourceAccessManagementSettings, error)

	// jamfproapi_adcs_settings.go
	CreateAdcsSettingsV1(settings *jamfpro.ResourceAdcsSettingsV1) (*jamfpro.ResponseAdcsSettingsCreatedV1, error)
	GetAdcsSettingsByIDV1(id string) (*jamfpro.ResponseAdcsSettingsV1, error)
	UpdateAdcsSettingsByIDV1(id string, settings *jamfpro.ResourceAdcsSettingsV1) error
	DeleteAdcsSettingsByIDV1(id string) error

	// jamfproapi_advanced_mobile_device_searches.go
	GetAdvancedMobileDeviceSearchByID(id string) (*jamfpro.ResourceAdvancedMobileDeviceSearch, error)
	CreateAdvancedMobileDeviceSearch(search jamfpro.ResourceAdvancedMobileDeviceSearch) (*jamfpro.ResponseAdvancedMobileDeviceSearchCreate, error)
	UpdateAdvancedMobileDeviceSearchByID(id string, search jamfpro.ResourceAdvancedMobileDeviceSearch) (*jamfpro.ResourceAdvancedMobileDeviceSearch, error)
	DeleteAdvancedMobileDeviceSearchByID(id string) error

	// jamfproapi_api_integrations.go
	GetApiIntegrationByID(id string) (*jamfpro.ResourceApiIntegration, error)
	CreateApiIntegration(integration *jamfpro.ResourceApiIntegration) (*jamfpro.ResourceApiIntegration, error)
	UpdateApiIntegrationByID(id string, integrationUpdate *jamfpro.ResourceApiIntegration) (*jamfpro.ResourceApiIntegration, error)
	DeleteApiIntegrationByID(id string) error
	RefreshClientCredentialsByApiRoleID(id string) (*jamfpro.ResourceClientCredentials, error)

	// jamfproapi_api_role_privileges.go
	GetJamfAPIPrivileges() (*jamfpro.ResourceApiRolePrivilegesList, error)

	// jamfproapi_api_roles.go
	GetJamfApiRoleByID(id string) (*jamfpro.ResourceAPIRole, error)
	GetJamfApiRoleByName(name string) (*jamfpro.ResourceAPIRole, error)
	CreateJamfApiRole(role *jamfpro.ResourceAPIRole) (*jamfpro.ResourceAPIRole, error)
	UpdateJamfApiRoleByID(id string, roleUpdate *jamfpro.ResourceAPIRole) (*jamfpro.ResourceAPIRole, error)
	DeleteJamfApiRoleByID(id string) error

	// jamfproapi_app_installers.go
	GetJamfAppCatalogAppInstallerTermsAndConditionsStatus() (*jamfpro.ResponseJamfAppCatalogDeploymentTermsAndConditionsStatus, error)
	AcceptJamfAppCatalogAppInstallerTermsAndConditions() (*jamfpro.ResponseJamfAppCatalogDeploymentTermsAndConditionsStatus, error)
	GetJamfAppCatalogAppInstallerTitles(params url.Values) (*jamfpro.ResponseJamfAppCatalogTitleList, error)
	GetJamfAppCatalogAppInstallerByID(id string) (*jamfpro.ResourceJamfAppCatalogDeployment, error)
	GetJamfAppCatalogAppInstallerByName(name string) (*jamfpro.ResourceJamfAppCatalogDeployment, error)
	GetJamfAppCatalogAppInstallerGlobalSettings() (*jamfpro.JamfAppCatalogDeploymentSubsetNotificationSettings, error)
	UpdateJamfAppCatalogAppInstallerGlobalSettings(payload *jamfpro.ResponseJamfAppCatalogGlobalSettings) (*jamfpro.ResponseJamfAppCatalogGlobalSettings, error)
	CreateJamfAppCatalogAppInstallerDeployment(payload *jamfpro.ResourceJamfAppCatalogDeployment) (*jamfpro.ResponseJamfAppCatalogDeploymentCreateAndUpdate, error)
	UpdateJamfAppCatalogAppInstallerDeploymentByID(id string, payload *jamfpro.ResourceJamfAppCatalogDeployment) (*jamfpro.ResourceJamfAppCatalogDeployment, error)
	DeleteJamfAppCatalogAppInstallerDeploymentByID(id string) error

	// jamfproapi_buildings.go
	GetBuildingByID(id string) (*jamfpro.ResourceBuilding, error)
	GetBuildingByName(name string) (*jamfpro.ResourceBuilding, error)
	CreateBuilding(building *jamfpro.ResourceBuilding) (*jamfpro.ResponseBuildingCreate, error)
	UpdateBuildingByID(id string, buildingUpdate *jamfpro.ResourceBuilding) (*jamfpro.ResourceBuilding, error)
	DeleteBuildingByID(id string) error

	// jamfproapi_categories.go
	GetCategoryByID(id string) (*jamfpro.ResourceCategory, error)
	GetCategoryByName(name string) (*jamfpro.ResourceCategory, error)
	CreateCategory(category *jamfpro.ResourceCategory) (*jamfpro.ResponseCategoryCreateAndUpdate, error)
	UpdateCategoryByID(id string, categoryUpdate *jamfpro.ResourceCategory) (*jamfpro.ResponseCategoryCreateAndUpdate, error)
	DeleteCategoryByID(id string) error

	// jamfproapi_client_checkin.go
	GetClientCheckinSettings() (*jamfpro.ResourceClientCheckinSettings, error)
	UpdateClientCheckinSettings(settingsUpdate jamfpro.ResourceClientCheckinSettings) (*jamfpro.ResourceClientCheckinSettings, error)

	// jamfproapi_cloud_distribution_point.go
	GetCloudDistributionPointV1() (*jamfpro.ResponseCloudDistributionPointV1, error)
	CreateCloudDistributionPointV1(payload *jamfpro.ResourceCloudDistributionPointV1) (*jamfpro.ResponseCloudDistributionPointV1, error)
	UpdateCloudDistributionPointV1(payload *jamfpro.ResourceCloudDistributionPointV1) (*jamfpro.ResponseCloudDistributionPointV1, error)
	DeleteCloudDistributionPointV1() error
	GetCloudDistributionPointUploadCapabilityV1() (*jamfpro.ResourceCloudDistributionPointUploadCapabilityV1, error)
	RefreshCloudDistributionPointInventoryV1(fileName string) error

	// jamfproapi_cloud_idp.go
	GetCloudIdentityProviderConfigurationByID(id string) (*jamfpro.ResourceCloudIdentityProviderDetails, error)
	GetCloudIdentityProviderConfigurationByName(name string) (*jamfpro.ResourceCloudIdentityProvider, error)

	// jamfproapi_cloud_ldap.go
	CreateCloudIdentityProviderLdap(config *jamfpro.ResourceCloudLdap) (*jamfpro.ResponseCloudIdentityProviderLdapCreated, error)
	GetCloudIdentityProviderLdapByID(id string) (*jamfpro.ResourceCloudLdap, error)
	UpdateCloudIdentityProviderLdap(id string, config *jamfpro.ResourceCloudLdap) (*jamfpro.ResourceCloudLdap, error)
	DeleteCloudIdentityProviderLdapByID(id string) error

	// jamfproapi_computer_extension_attributes.go
	GetComputerExtensionAttributeByID(id string) (*jamfpro.ResourceComputerExtensionAttribute, error)
	GetComputerExtensionAttributeByName(name string) (*jamfpro.ResourceComputerExtensionAttribute, error)
	CreateComputerExtensionAttribute(attribute *jamfpro.ResourceComputerExtensionAttribute) (*jamfpro.ResponseComputerExtensionAttributeCreated, error)
	UpdateComputerExtensionAttributeByID(id string, attribute *jamfpro.ResourceComputerExtensionAttribute) (*jamfpro.ResourceComputerExtensionAttribute, error)
	DeleteComputerExtensionAttributeByID(id string) error

	// jamfproapi_computer_groups_v2_context.go
	GetSmartComputerGroupByIDV2WithContext(ctx context.Context, id string) (*jamfpro.ResourceSmartComputerGroupV2, error)
	CreateSmartComputerGroupV2WithContext(ctx context.Context, input jamfpro.ResourceSmartComputerGroupV2) (*jamfpro.ResponseSmartComputerGroupCreateV2, error)
	UpdateSmartComputerGroupByIDV2WithContext(ctx context.Context, id string, input jamfpro.ResourceSmartComputerGroupV2) (*jamfpro.ResourceSmartComputerGroupV2, error)
	DeleteSmartComputerGroupByIDV2WithContext(ctx context.Context, id string) error

	// jamfproapi_computer_groups_v3_context.go
	GetSmartComputerGroupByIDV3WithContext(ctx context.Context, id string) (*jamfpro.ResourceSmartComputerGroupV3, error)
	CreateSmartComputerGroupV3WithContext(ctx context.Context, input jamfpro.ResourceSmartComputerGroupV3) (*jamfpro.ResponseSmartComputerGroupCreateV3, error)
	UpdateSmartComputerGroupByIDV3WithContext(ctx context.Context, id string, input jamfpro.ResourceSmartComputerGroupV3) (*jamfpro.ResourceSmartComputerGroupV3, error)
	DeleteSmartComputerGroupByIDV3WithContext(ctx context.Context, id string) error
	GetSmartComputerGroupByNameV3WithContext(ctx context.Context, name string) (*jamfpro.ResponseSmartComputerGroupListItemV3, error)
	GetComputerSmartGroupPlatformIDByJamfProIDWithContext(ctx context.Context, id string) (string, error)
	UpdateGroupDescriptionByIDV2WithContext(ctx context.Context, platformID string, description string) error
	CreateComputerGroupWithContext(ctx context.Context, group *jamfpro.ResourceComputerGroup) (*jamfpro.ResponseComputerGroupreatedAndUpdated, error)
	UpdateComputerGroupByIDWithContext(ctx context.Context, id string, group *jamfpro.ResourceComputerGroup) (*jamfpro.ResponseComputerGroupreatedAndUpdated, error)

	// jamfproapi_computer_inventory.go
	GetComputerInventoryByID(id string) (*jamfpro.ResourceComputerInventory, error)
	GetComputerInventoryByName(name string) (*jamfpro.ResourceComputerInventory, error)
	GetComputerInventoryBySerialNumber(serialNumber string) (*jamfpro.ResourceComputerInventory, error)

	// jamfproapi_computer_inventory_collection_settings.go
	CreateComputerInventoryCollectionSettingsCustomPath(customPath *jamfpro.ResourceComputerInventoryCollectionSettingsCustomPath) (*jamfpro.ComputerInventoryCollectionSettingsSubsetPathItem, error)
	DeleteComputerInventoryCollectionSettingsCustomPathByID(id string) error
	GetComputerInventoryCollectionSettings() (*jamfpro.ResourceComputerInventoryCollectionSettings, error)
	UpdateComputerInventoryCollectionSettings(settingsUpdate *jamfpro.ResourceComputerInventoryCollectionSettings) (*jamfpro.ResourceComputerInventoryCollectionSettings, error)

	// jamfproapi_computer_prestages.go
	GetComputerPrestageByID(id string) (*jamfpro.ResourceComputerPrestage, error)
	CreateComputerPrestage(prestage *jamfpro.ResourceComputerPrestage) (*jamfpro.ResponseComputerPrestageCreate, error)
	UpdateComputerPrestageByID(id string, prestageUpdate *jamfpro.ResourceComputerPrestage) (*jamfpro.ResourceComputerPrestage, error)
	DeleteComputerPrestageByID(id string) error

	// jamfproapi_departments.go
	GetDepartmentByID(id string) (*jamfpro.ResourceDepartment, error)
	GetDepartmentByName(name string) (*jamfpro.ResourceDepartment, error)
	CreateDepartment(department *jamfpro.ResourceDepartment) (*jamfpro.ResponseDepartmentCreate, error)
	UpdateDepartmentByID(id string, departmentUpdate *jamfpro.ResourceDepartment) (*jamfpro.ResourceDepartment, error)
	DeleteDepartmentByID(id string) error

	// jamfproapi_device_communication_settings.go
	GetDeviceCommunicationSettings() (*jamfpro.ResourceDeviceCommunicationSettings, error)
	UpdateDeviceCommunicationSettings(updatedSettings jamfpro.ResourceDeviceCommunicationSettings) (*jamfpro.ResourceDeviceCommunicationSettings, error)

	// jamfproapi_device_enrollments.go
	GetDeviceEnrollmentByID(id string) (*jamfpro.ResourceDeviceEnrollment, error)
	GetDeviceEnrollmentByName(name string) (*jamfpro.ResourceDeviceEnrollment, error)
	GetDeviceEnrollmentsPublicKey() (string, error)
	CreateDeviceEnrollmentWithMDMServerToken(tokenUpload *jamfpro.ResourceDeviceEnrollmentTokenUpload) (*jamfpro.ResponseDeviceEnrollmentTokenUpload, error)
	UpdateDeviceEnrollmentMetadataByID(id string, deviceEnrollment *jamfpro.ResourceDeviceEnrollmentUpdate) (*jamfpro.ResourceDeviceEnrollment, error)
	UpdateDeviceEnrollmentMDMServerToken(id string, tokenUpload *jamfpro.ResourceDeviceEnrollmentTokenUpload) (*jamfpro.ResponseDeviceEnrollmentTokenUpload, error)
	DeleteDeviceEnrollmentByID(id string) error

	// jamfproapi_engage.go
	GetEngageSettings() (*jamfpro.ResourceEngageSettings, error)
	UpdateEngageSettings(settingsUpdate jamfpro.ResourceEngageSettings) (*jamfpro.ResourceEngageSettings, error)

	// jamfproapi_enrollment.go
	GetAccountDrivenUserEnrollmentAccessGroups(params url.Values) (*jamfpro.ResponseAccountDrivenUserEnrollmentAccessGroupsList, error)
	CreateAccountDrivenUserEnrollmentAccessGroup(script *jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup) (*jamfpro.ResponseAccountDrivenUserEnrollmentAccessGroupCreateAndUpdate, error)
	UpdateAccountDrivenUserEnrollmentAccessGroupByID(id string, groupUpdate *jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup) (*jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup, error)
	DeleteAccountDrivenUserEnrollmentAccessGroupByID(id string) error
	GetEnrollment() (*jamfpro.ResourceEnrollment, error)
	UpdateEnrollment(enrollmentUpdate *jamfpro.ResourceEnrollment) (*jamfpro.ResourceEnrollment, error)
	GetEnrollmentMessages() ([]jamfpro.ResourceEnrollmentLanguage, error)
	UpdateEnrollmentMessageByLanguageID(languageId string, enrollmentMessage *jamfpro.ResourceEnrollmentLanguage) (*jamfpro.ResourceEnrollmentLanguage, error)
	DeleteEnrollmentMessageByLanguageID(languageId string) error
	GetEnrollmentLanguageCodes() ([]jamfpro.ResourceLanguageCode, error)

	// jamfproapi_enrollment_customizations.go
	GetEnrollmentCustomizationByID(id string) (*jamfpro.ResourceEnrollmentCustomization, error)
	CreateEnrollmentCustomization(enrollmentCustomization jamfpro.ResourceEnrollmentCustomization) (*jamfpro.ResponseEnrollmentCustomizationCreate, error)
	UpdateEnrollmentCustomizationByID(id string, updatedCustomization jamfpro.ResourceEnrollmentCustomization) (*jamfpro.ResourceEnrollmentCustomization, error)
	DeleteEnrollmentCustomizationByID(id string) error
	UploadEnrollmentCustomizationsImage(filePath string) (*jamfpro.ResponseEnrollmentCustomizationUpload, error)
	GetPrestagePanes(customizationResourceID string) (*jamfpro.ResponsePrestagePanesList, error)
	CreateTextPrestagePane(customizationResourceID string, textPane jamfpro.ResourceEnrollmentCustomizationTextPane) (*jamfpro.ResponseEnrollmentCustomizationTextPane, error)
	CreateLDAPPrestagePane(customizationID string, ldapPane jamfpro.ResourceEnrollmentCustomizationLDAPPane) (*jamfpro.ResponseEnrollmentCustomizationLDAPPane, error)
	CreateSSOPrestagePane(customizationID string, ssoPane jamfpro.ResourceEnrollmentCustomizationSSOPane) (*jamfpro.ResponseEnrollmentCustomizationSSOPane, error)
	UpdateTextPrestagePaneByID(customizationResourceID string, panelID string, textPane jamfpro.ResourceEnrollmentCustomizationTextPane) (*jamfpro.ResourceEnrollmentCustomizationTextPane, error)
	UpdateLDAPPrestagePaneByID(customizationResourceID string, panelID string, ldapPane jamfpro.ResourceEnrollmentCustomizationLDAPPane) (*jamfpro.ResourceEnrollmentCustomizationLDAPPane, error)
	UpdateSSOPrestagePaneByID(customizationResourceID string, panelID string, ssoPane jamfpro.ResourceEnrollmentCustomizationSSOPane) (*jamfpro.ResourceEnrollmentCustomizationSSOPane, error)
	GetTextPrestagePaneByID(customizationResourceID string, paneID string) (*jamfpro.ResourceEnrollmentCustomizationTextPane, error)
	GetLDAPPrestagePaneByID(customizationResourceID string, paneID string) (*jamfpro.ResourceEnrollmentCustomizationLDAPPane, error)
	GetSSOPrestagePaneByID(customizationResourceID string, paneID string) (*jamfpro.ResourceEnrollmentCustomizationSSOPane, error)
	DeleteSSOPrestagePane(customizationResourceID string, paneID string) error
	DeletePrestagePane(customizationResourceID string, paneID string) error

	// jamfproapi_enrollment_settings.go
	GetADUESessionTokenSettings() (*jamfpro.ResourceADUETokenSettings, error)
	UpdateADUESessionTokenSettings(updatedSettings jamfpro.ResourceADUETokenSettings) (*jamfpro.ResourceADUETokenSettings, error)

	// jamfproapi_groups.go
	GetGroupByID(id string) (*jamfpro.ResourceGroup, error)
	GetComputerGroupByJamfProName(name string) (*jamfpro.ResourceGroup, error)
	GetMobileGroupByJamfProName(name string) (*jamfpro.ResourceGroup, error)
	GetComputerGroupByJamfProID(id string) (*jamfpro.ResourceGroup, error)
	GetMobileGroupByJamfProID(id string) (*jamfpro.ResourceGroup, error)

	// jamfproapi_icon.go
	GetIconByID(id int) (*jamfpro.ResponseIconUpload, error)
	UploadIcon(filePath string) (*jamfpro.ResponseIconUpload, error)

	// jamfproapi_impact_alert_notification_settings.go
	GetImpactAlertNotificationSettings() (*jamfpro.ResourceImpactAlertNotificationSettings, error)
	UpdateImpactAlertNotificationSettings(settings jamfpro.ResourceImpactAlertNotificationSettings) error

	// jamfproapi_jamf_connect.go
	GetJamfConnectConfigProfileByID(profileID int) (*jamfpro.ResourceJamfConnectConfigProfile, error)
	GetJamfConnectConfigProfileByName(name string) (*jamfpro.ResourceJamfConnectConfigProfile, error)
	GetJamfConnectConfigProfiles(params url.Values) (*jamfpro.ResponseJamfConnectConfigProfilesList, error)
	UpdateJamfConnectConfigProfileByConfigProfileUUID(configProfileUUID string, profileUpdate *jamfpro.ResourceJamfConnectConfigProfileUpdate) (*jamfpro.ResourceJamfConnectConfigProfile, error)

	// jamfproapi_jamf_content_distribution_server.go
	GetJCDS2Packages() ([]jamfpro.ResponseJCDS2List, error)

	// jamfproapi_jamf_pro_version.go
	GetJamfProVersion() (*jamfpro.ResponseJamfProVersion, error)

	// jamfproapi_jamf_protect.go
	GetJamfProtectSettings() (*jamfpro.ResponseJamfProtectSettings, error)
	UpdateJamfProtectSettings(settings jamfpro.ResourceJamfProtectSettings) (*jamfpro.ResponseJamfProtectSettings, error)
	CreateJamfProtectIntegration(registration jamfpro.ResourceJamfProtectRegistration, autoInstall bool) (*jamfpro.ResponseJamfProtectSettings, error)
	GetJamfProtectPlans(params url.Values) (*jamfpro.ResponseJamfProtectPlansList, error)
	DeleteJamfProtectIntegration() error

	// jamfproapi_jamf_sso_failover.go
	GetSSOFailoverSettings() (*jamfpro.ResponseSSOFailover, error)
	UpdateFailoverUrl() (*jamfpro.ResponseSSOFailover, error)

	// jamfproapi_ldap.go
	GetLdapGroupsV1(search string) (*jamfpro.ResponseLdapGroupSearchResultsV1, error)

	// jamfproapi_local_admin_password.go
	GetLocalAdminPasswordSettings() (*jamfpro.ResourceLocalAdminPasswordSettings, error)
	UpdateLocalAdminPasswordSettings(settings *jamfpro.ResourceLocalAdminPasswordSettings) error

	// jamfproapi_managed_software_updates.go
	CreateManagedSoftwareUpdatePlanByDeviceID(plan *jamfpro.ResourceManagedSoftwareUpdatePlan) (*jamfpro.ResponseManagedSoftwareUpdatePlanCreate, error)
	GetManagedSoftwareUpdateFeatureToggle() (*jamfpro.ResourceManagedSoftwareUpdateFeatureToggle, error)
	UpdateManagedSoftwareUpdateFeatureToggle(payload *jamfpro.ResourceManagedSoftwareUpdateFeatureToggle) (*jamfpro.ResponseManagedSoftwareUpdateFeatureToggle, error)
	CreateManagedSoftwareUpdatePlanByGroupID(plan *jamfpro.ResourceManagedSoftwareUpdatePlan) (*jamfpro.ResponseManagedSoftwareUpdatePlanCreate, error)
	GetManagedSoftwareUpdatePlanByUUID(UUID string) (*jamfpro.ResponseManagedSoftwareUpdatePlan, error)

	// jamfproapi_mobile_device_extension_attributes.go
	GetMobileDeviceExtensionAttributeByID(id string) (*jamfpro.ResourceMobileDeviceExtensionAttribute, error)
	CreateMobileDeviceExtensionAttribute(attribute *jamfpro.ResourceMobileDeviceExtensionAttribute) (*jamfpro.ResponseMobileDeviceExtensionAttributeCreated, error)
	UpdateMobileDeviceExtensionAttributeByID(id string, attribute *jamfpro.ResourceMobileDeviceExtensionAttribute) (*jamfpro.ResourceMobileDeviceExtensionAttribute, error)
	DeleteMobileDeviceExtensionAttributeByID(id string) error

	// jamfproapi_mobile_device_groups_context.go
	GetSmartMobileDeviceGroupByIDV1WithContext(ctx context.Context, id string) (*jamfpro.ResourceSmartMobileDeviceGroupV1, error)
	DeleteSmartMobileDeviceGroupByIDV1WithContext(ctx context.Context, id string) error

	// jamfproapi_mobile_device_prestages.go
	GetMobileDevicePrestageByID(id string) (*jamfpro.ResourceMobileDevicePrestage, error)
	CreateMobileDevicePrestage(newPrestage jamfpro.ResourceMobileDevicePrestage) (*jamfpro.ResponseMobileDevicePrestageCreate, error)
	UpdateMobileDevicePrestageByID(id string, prestageUpdate *jamfpro.ResourceMobileDevicePrestage) (*jamfpro.ResourceMobileDevicePrestage, error)
	DeleteMobileDevicePrestageByID(id string) error

	// jamfproapi_mobile_devices_inventory.go
	GetMobileDeviceInventoryByID(id string) (*jamfpro.ResourceMobileDeviceInventory, error)
	GetMobileDeviceInventoryBySerialNumber(serialNumber string) (*jamfpro.ResourceMobileDeviceInventory, error)
	GetMobileDeviceInventoryByName(name string) (*jamfpro.ResourceMobileDeviceInventory, error)

	// jamfproapi_onboarding.go
	GetOnboardingSettings() (*jamfpro.ResponseOnboardingSettings, error)
	UpdateOnboardingSettings(request jamfpro.ResourceUpdateOnboardingSettings) (*jamfpro.ResponseOnboardingSettings, error)

	// jamfproapi_packages.go
	GetPackageByID(id string) (*jamfpro.ResourcePackage, error)
	GetPackageByName(name string) (*jamfpro.ResourcePackage, error)
	CreatePackage(packageMetadata jamfpro.ResourcePackage) (*jamfpro.ResponsePackageCreatedAndUpdated, error)
	UploadPackage(id string, filePaths []string) (*jamfpro.ResponsePackageCreatedAndUpdated, error)
	UpdatePackageByID(id string, packageMetadata jamfpro.ResourcePackage) (*jamfpro.ResourcePackage, error)
	DeletePackageByID(id string) error

	// jamfproapi_policy_properties.go
	GetPolicyProperties() (*jamfpro.ResourcePolicyProperties, error)
	UpdatePolicyProperties(settingsUpdate jamfpro.ResourcePolicyProperties) (*jamfpro.ResourcePolicyProperties, error)

	// jamfproapi_reenrollment.go
	GetReenrollmentSettings() (*jamfpro.ResourceReenrollmentSettings, error)
	UpdateReenrollmentSettings(settingsUpdate jamfpro.ResourceReenrollmentSettings) (*jamfpro.ResourceReenrollmentSettings, error)

	// jamfproapi_scripts.go
	GetScriptByID(id string) (*jamfpro.ResourceScript, error)
	GetScriptByName(name string) (*jamfpro.ResourceScript, error)
	CreateScript(script *jamfpro.ResourceScript) (*jamfpro.ResponseScriptCreate, error)
	UpdateScriptByID(id string, scriptUpdate *jamfpro.ResourceScript) (*jamfpro.ResourceScript, error)
	DeleteScriptByID(id string) error

	// jamfproapi_self_service_branding_macos.go
	GetSelfServiceBrandingMacOSByID(id string) (*jamfpro.ResourceSelfServiceBrandingDetail, error)
	CreateSelfServiceBrandingMacOS(branding *jamfpro.ResourceSelfServiceBrandingDetail) (*jamfpro.ResourceSelfServiceBrandingDetail, error)
	UpdateSelfServiceBrandingMacOSByID(id string, brandingUpdate *jamfpro.ResourceSelfServiceBrandingDetail) (*jamfpro.ResourceSelfServiceBrandingDetail, error)
	DeleteSelfServiceBrandingMacOSByID(id string) error

	// jamfproapi_self_service_branding_mobile.go
	GetSelfServiceBrandingIOSByID(id string) (*jamfpro.ResourceSelfServiceBrandingIOSDetail, error)
	CreateSelfServiceBrandingIOS(branding *jamfpro.ResourceSelfServiceBrandingIOSDetail) (*jamfpro.CreateResponseSelfServiceBrandingIOS, error)
	UpdateSelfServiceBrandingIOSByID(id string, brandingUpdate *jamfpro.ResourceSelfServiceBrandingIOSDetail) (*jamfpro.ResourceSelfServiceBrandingIOSDetail, error)
	DeleteSelfServiceBrandingIOSByID(id string) error

	// jamfproapi_self_service_branding_upload.go
	UploadSelfServiceBrandingImage(filePath string) (*jamfpro.ResponseSelfServiceBrandingImage, error)

	// jamfproapi_self_service_plus_settings.go
	GetSelfServicePlusSettings() (*jamfpro.ResourceSelfServicePlusSettings, error)
	UpdateSelfServicePlusSettings(settings jamfpro.ResourceSelfServicePlusSettings) error

	// jamfproapi_self_service_settings.go
	GetSelfServiceSettings() (*jamfpro.ResourceSelfServiceSettings, error)
	UpdateSelfServiceSettings(settingsUpdate *jamfpro.ResourceSelfServiceSettings) (*jamfpro.ResourceSelfServiceSettings, error)

	// jamfproapi_service_discovery_enrollment.go
	GetServiceDiscoveryEnrollmentWellKnownSettingsV1() (*jamfpro.ResponseServiceDiscoveryEnrollmentWellKnownSettingsV1, error)
	UpdateServiceDiscoveryEnrollmentWellKnownSettingsV1(request jamfpro.ResponseServiceDiscoveryEnrollmentWellKnownSettingsV1) error

	// jamfproapi_smart_computer_groups.go
	GetSmartComputerGroupByIDV2(id string) (*jamfpro.ResourceSmartComputerGroupV2, error)
	GetSmartComputerGroupByNameV2(name string) (*jamfpro.ResponseSmartComputerGroupListItemV2, error)

	// jamfproapi_smart_mobile_device_groups.go
	GetSmartMobileDeviceGroupByIDV1(id string) (*jamfpro.ResourceSmartMobileDeviceGroupV1, error)
	GetSmartMobileDeviceGroupByNameV1(groupName string) (*jamfpro.ResponseSmartMobileDeviceGroupListItemV1, error)
	CreateSmartMobileDeviceGroupV1(request jamfpro.ResourceSmartMobileDeviceGroupV1) (*jamfpro.ResponseSmartMobileDeviceGroupCreateV1, error)
	UpdateSmartMobileDeviceGroupByIDV1(id string, request jamfpro.ResourceSmartMobileDeviceGroupV1) (*jamfpro.ResourceSmartMobileDeviceGroupV1, error)

	// jamfproapi_smtp_server.go
	GetSMTPServerInformation() (*jamfpro.ResourceSMTPServer, error)
	UpdateSMTPServerInformation(settings *jamfpro.ResourceSMTPServer) (*jamfpro.ResourceSMTPServer, error)

	// jamfproapi_sso_certificate.go
	GetSSOCertificate() (*jamfpro.ResourceSSOKeystoreResponse, error)
	CreateSSOCertificate() (*jamfpro.ResourceSSOKeystoreResponse, error)
	DeleteSSOCertificate() error

	// jamfproapi_sso_settings.go
	GetSsoSettings() (*jamfpro.ResourceSsoSettings, error)
	UpdateSsoSettings(updatedSettings jamfpro.ResourceSsoSettings) (*jamfpro.ResourceSsoSettings, error)
	DisableSso() error

	// jamfproapi_volume_purchasing_locations.go
	GetVolumePurchasingLocationByID(id string) (*jamfpro.ResourceVolumePurchasingLocation, error)
	CreateVolumePurchasingLocation(request *jamfpro.VolumePurchasingLocationCreateUpdateRequest) (*jamfpro.ResponseVolumePurchasingLocationCreate, error)
	UpdateVolumePurchasingLocationByID(id string, request *jamfpro.VolumePurchasingLocationCreateUpdateRequest) (*jamfpro.ResourceVolumePurchasingLocation, error)
	DeleteVolumePurchasingLocationByID(id string) error
}
