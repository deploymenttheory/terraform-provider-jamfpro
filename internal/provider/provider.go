// providers.go
package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-http-client-integrations/jamf/jamfprointegration"
	"github.com/deploymenttheory/go-api-http-client/httpclient"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/account"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/account_driven_user_enrollment_settings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/account_group"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/activation_code"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/advanced_computer_search"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/advanced_mobile_device_search"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/advanced_user_search"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/allowed_file_extension"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/api_integration"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/api_role"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/app_installer"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/app_installer_global_settings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/building"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/category"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/client_checkin"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/cloud_ldap"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/cloudidp"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computer_extension_attribute"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computer_inventory_collection_settings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computer_prestage_enrollment"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/department"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/device_communication_settings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/device_enrollments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/deviceenrollmentspublickey"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/disk_encryption_configuration"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/dock_item"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/engage_settings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/enrollment_customization"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/file_share_distribution_point"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/icon"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/jamf_connect"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/jamf_protect"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/ldap_server"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/local_admin_password_settings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/macos_configuration_profile_plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/macos_configuration_profile_plist_generator"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/managed_software_update"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/mobile_device_configuration_profile_plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/mobile_device_extension_attribute"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/mobile_device_prestage_enrollment"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/network_segment"
	packages "github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/package"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/policy"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/printer"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/restricted_software"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/script"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/self_service_settings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/site"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/smart_computer_group"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/smart_mobile_device_group"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/smtp_server"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/sso_certificate"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/sso_failover"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/sso_settings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/static_computer_group"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/static_mobile_device_group"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/user_group"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/userinitiatedenrollment"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/volume_purchasing_locations"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/webhook"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"go.uber.org/zap"
)

// TerraformProviderProductUserAgent is included in the User-Agent header for
// any API requests made by the provider.
const (
	terraformProviderProductUserAgent = "terraform-provider-jamfpro"
	envVarOAuthClientId               = "JAMFPRO_CLIENT_ID"
	envVarOAuthClientSecret           = "JAMFPRO_CLIENT_SECRET"
	envVarBasicAuthUsername           = "JAMFPRO_BASIC_USERNAME"
	envVarBasicAuthPassword           = "JAMFPRO_BASIC_PASSWORD"
	envVarJamfProFQDN                 = "JAMFPRO_INSTANCE_FQDN"
	envVarJamfProAuthMethod           = "JAMFPRO_AUTH_METHOD"
	jamfLoadBalancerCookieName        = "jpro-ingress"
)

/*
GetJamfFqdn retrieves the instance domain name from the provided schema resource data.

If the instance domain is not found, it appends an error diagnostic to the diagnostics slice.

Parameters:

	d      - A pointer to the schema.ResourceData object which contains the resource data.
	diags  - A pointer to a slice of diag.Diagnostics where error messages will be appended.

Returns:

	A string representing the instance domain name. If the instance domain name is not provided,
	an error diagnostic is appended to diags and an empty string is returned.
*/
func GetJamfFqdn(d *schema.ResourceData, diags *diag.Diagnostics) string {
	jamf_fqdn, ok := d.GetOk("jamfpro_instance_fqdn")
	if jamf_fqdn.(string) == "" || !ok {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error getting instance name",
			Detail:   "instance_name must be provided either as an environment variable (JAMFPRO_INSTANCE_FQDN) or in the Terraform configuration",
		})
		return ""
	}
	return jamf_fqdn.(string)
}

/*
GetClientID retrieves the client ID from the provided schema resource data.
If the client ID is not found, it appends an error diagnostic to the diagnostics slice.

Parameters:

	d      - A pointer to the schema.ResourceData object which contains the resource data.
	diags  - A pointer to a slice of diag.Diagnostics where error messages will be appended.

Returns:

	A string representing the client ID. If the client ID is not provided,
	an error diagnostic is appended to diags and an empty string is returned.
*/
func GetClientID(d *schema.ResourceData, diags *diag.Diagnostics) string {
	clientID := d.Get("client_id").(string)
	if clientID == "" {

		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error getting client id",
			Detail:   "client_id must be provided either as an environment variable (JAMFPRO_CLIENT_ID) or in the Terraform configuration",
		})

		return ""

	}
	return clientID
}

/*
GetClientSecret retrieves the client secret from the provided schema resource data.
If the client ID is not found, it appends an error diagnostic to the diagnostics slice.

Parameters:

	d      - A pointer to the schema.ResourceData object which contains the resource data.
	diags  - A pointer to a slice of diag.Diagnostics where error messages will be appended.

Returns:

	A string representing the client ID. If the client ID is not provided,
	an error diagnostic is appended to diags and an empty string is returned.
*/
func GetClientSecret(d *schema.ResourceData, diags *diag.Diagnostics) string {
	clientSecret := d.Get("client_secret").(string)
	if clientSecret == "" {

		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error getting client secret",
			Detail:   "client_secret must be provided either as an environment variable (JAMFPRO_CLIENT_SECRET) or in the Terraform configuration",
		})

		return ""

	}
	return clientSecret
}

/*
GetBasicAuthUsername retrieves the basic auth username from the provided schema resource data.
If the client ID is not found, it appends an error diagnostic to the diagnostics slice.

Parameters:

	d      - A pointer to the schema.ResourceData object which contains the resource data.
	diags  - A pointer to a slice of diag.Diagnostics where error messages will be appended.

Returns:

	A string representing the client ID. If the client ID is not provided,
	an error diagnostic is appended to diags and an empty string is returned.
*/
func GetBasicAuthUsername(d *schema.ResourceData, diags *diag.Diagnostics) string {
	username, ok := d.GetOk("basic_auth_username")
	if !ok || username.(string) == "" {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error getting basic auth username",
			Detail:   "basic_auth_username must be provided either as an environment variable (JAMFPRO_BASIC_USERNAME) or in the Terraform configuration",
		})
		return ""
	}
	return username.(string)
}

/*
GetBasicAuthPassword retrieves the basic auth password from the provided schema resource data.
If the client ID is not found, it appends an error diagnostic to the diagnostics slice.

Parameters:

	d      - A pointer to the schema.ResourceData object which contains the resource data.
	diags  - A pointer to a slice of diag.Diagnostics where error messages will be appended.

Returns:

	A string representing the client ID. If the client ID is not provided,
	an error diagnostic is appended to diags and an empty string is returned.
*/
func GetBasicAuthPassword(d *schema.ResourceData, diags *diag.Diagnostics) string {
	password, ok := d.GetOk("basic_auth_password")
	if !ok || password.(string) == "" {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error getting basic auth password",
			Detail:   "basic_auth_password must be provided either as an environment variable (JAMFPRO_BASIC_PASSWORD) or in the Terraform configuration",
		})
		return ""
	}
	return password.(string)
}

/*
GetAuthMethod retrieves the auth method from the provided schema resource data.
If the auth method is not found, it appends an error diagnostic to the diagnostics slice.

Parameters:

	d      - A pointer to the schema.ResourceData object which contains the resource data.
	diags  - A pointer to a slice of diag.Diagnostics where error messages will be appended.

Returns:

	A string representing the auth method. If the auth method is not provided,
	an error diagnostic is appended to diags and an empty string is returned.
*/
func GetAuthMethod(d *schema.ResourceData, diags *diag.Diagnostics) string {
	authMethod, ok := d.GetOk("auth_method")
	if !ok || authMethod.(string) == "" {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error getting auth method",
			Detail:   "auth_method must be provided either as an environment variable (JAMFPRO_AUTH_METHOD) or in the Terraform configuration",
		})
		return ""
	}
	return authMethod.(string)
}

// Schema defines the configuration attributes for the  within the JamfPro provider.
func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"jamfpro_instance_fqdn": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(envVarJamfProFQDN, ""),
				Description: "The Jamf Pro FQDN (fully qualified domain name). example: https://mycompany.jamfcloud.com",
			},
			"auth_method": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(envVarJamfProAuthMethod, ""),
				Description: "The auth method chosen for interacting with Jamf Pro. Options are 'basic' for username/password or 'oauth2' for client id/secret.",
				ValidateFunc: validation.StringInSlice([]string{
					"basic", "oauth2",
				}, true),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(envVarOAuthClientId, ""),
				Description: "The Jamf Pro Client ID for authentication when auth_method is 'oauth2'.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(envVarOAuthClientSecret, ""),
				Description: "The Jamf Pro Client secret for authentication when auth_method is 'oauth2'.",
			},
			"basic_auth_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(envVarBasicAuthUsername, ""),
				Description: "The Jamf Pro username used for authentication when auth_method is 'basic'.",
			},
			"basic_auth_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(envVarBasicAuthPassword, ""),
				Description: "The Jamf Pro password used for authentication when auth_method is 'basic'.",
			},
			"enable_client_sdk_logs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Debug option to propogate logs from the SDK and HttpClient",
			},
			"client_sdk_log_export_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Specify the path to export http client logs to.",
			},
			"hide_sensitive_data": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Define whether sensitive fields should be hidden in logs. Default to hiding sensitive data in logs",
			},
			"custom_cookies": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				Description: "Persistent custom cookies used by HTTP Client in all requests.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "cookie key",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "cookie value",
						},
					},
				},
			},
			"jamfpro_load_balancer_lock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Programatically determines all available web app members in the load balancer and locks all instances of httpclient to the app for faster executions. \nTEMP SOLUTION UNTIL JAMF PROVIDES SOLUTION",
			},
			"token_refresh_buffer_period_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "The buffer period in seconds for token refresh.",
			},

			"mandatory_request_delay_milliseconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "A mandatory delay after each request before returning to reduce high volume of requests in a short time",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{

			"jamfpro_account":                                   account.DataSourceJamfProAccounts(),
			"jamfpro_account_group":                             account_group.DataSourceJamfProAccountGroups(),
			"jamfpro_advanced_computer_search":                  advanced_computer_search.DataSourceJamfProAdvancedComputerSearches(),
			"jamfpro_advanced_mobile_device_search":             advanced_mobile_device_search.DataSourceJamfProAdvancedMobileDeviceSearches(),
			"jamfpro_advanced_user_search":                      advanced_user_search.DataSourceJamfProAdvancedUserSearches(),
			"jamfpro_api_integration":                           api_integration.DataSourceJamfProApiIntegrations(),
			"jamfpro_api_role":                                  api_role.DataSourceJamfProAPIRoles(),
			"jamfpro_app_installer":                             app_installer.DataSourceJamfProAppInstallers(),
			"jamfpro_building":                                  building.DataSourceJamfProBuildings(),
			"jamfpro_category":                                  category.DataSourceJamfProCategories(),
			"jamfpro_cloud_idp":                                 cloudidp.DataSourceJamfProCloudIdp(),
			"jamfpro_computer_extension_attribute":              computer_extension_attribute.DataSourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_prestage_enrollment":              computer_prestage_enrollment.DataSourceJamfProComputerPrestageEnrollment(),
			"jamfpro_department":                                department.DataSourceJamfProDepartments(),
			"jamfpro_device_enrollments":                        device_enrollments.DataSourceJamfProDeviceEnrollments(),
			"jamfpro_device_enrollments_public_key":             deviceenrollmentspublickey.DataSourceJamfProDeviceEnrollmentsPublicKey(),
			"jamfpro_disk_encryption_configuration":             disk_encryption_configuration.DataSourceJamfProDiskEncryptionConfigurations(),
			"jamfpro_dock_item":                                 dock_item.DataSourceJamfProDockItems(),
			"jamfpro_file_share_distribution_point":             file_share_distribution_point.DataSourceJamfProFileShareDistributionPoints(),
			"jamfpro_ldap_server":                               ldap_server.DataSourceJamfProLDAPServers(),
			"jamfpro_network_segment":                           network_segment.DataSourceJamfProNetworkSegments(),
			"jamfpro_macos_configuration_profile_plist":         macos_configuration_profile_plist.DataSourceJamfProMacOSConfigurationProfilesPlist(),
			"jamfpro_mobile_device_configuration_profile_plist": mobile_device_configuration_profile_plist.DataSourceJamfProMobileDeviceConfigurationProfilesPlist(),
			"jamfpro_mobile_device_prestage_enrollment":         mobile_device_prestage_enrollment.DataSourceJamfProMobileDevicePrestageEnrollment(),
			"jamfpro_package":                                   packages.DataSourceJamfProPackages(),
			"jamfpro_policy":                                    policy.DataSourceJamfProPolicies(),
			"jamfpro_printer":                                   printer.DataSourceJamfProPrinters(),
			"jamfpro_script":                                    script.DataSourceJamfProScripts(),
			"jamfpro_site":                                      site.DataSourceJamfProSites(),
			"jamfpro_smart_computer_group":                      smart_computer_group.DataSourceJamfProSmartComputerGroups(),
			"jamfpro_smart_mobile_device_group":                 smart_mobile_device_group.DataSourceJamfProSmartMobileGroups(),
			"jamfpro_sso_certificate":                           sso_certificate.DataSourceJamfProSSOCertificate(),
			"jamfpro_sso_failover":                              sso_failover.DataSourceJamfProSSOFailover(),
			"jamfpro_static_computer_group":                     static_computer_group.DataSourceJamfProStaticComputerGroups(),
			"jamfpro_static_mobile_device_group":                static_mobile_device_group.DataSourceJamfProStaticMobileDeviceGroups(),
			"jamfpro_restricted_software":                       restricted_software.DataSourceJamfProRestrictedSoftwares(),
			"jamfpro_user_group":                                user_group.DataSourceJamfProUserGroups(),
			"jamfpro_volume_purchasing_locations":               volume_purchasing_locations.DataSourceJamfProVolumePurchasingLocations(),
			"jamfpro_webhook":                                   webhook.DataSourceJamfProWebhooks(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"jamfpro_account": account.ResourceJamfProAccounts(),
			"jamfpro_account_driven_user_enrollment_settings":     account_driven_user_enrollment_settings.ResourceJamfProAccountDrivenUserEnrollmentSettings(),
			"jamfpro_account_group":                               account_group.ResourceJamfProAccountGroups(),
			"jamfpro_activation_code":                             activation_code.ResourceJamfProActivationCode(),
			"jamfpro_advanced_computer_search":                    advanced_computer_search.ResourceJamfProAdvancedComputerSearches(),
			"jamfpro_advanced_mobile_device_search":               advanced_mobile_device_search.ResourceJamfProAdvancedMobileDeviceSearches(),
			"jamfpro_advanced_user_search":                        advanced_user_search.ResourceJamfProAdvancedUserSearches(),
			"jamfpro_allowed_file_extension":                      allowed_file_extension.ResourceJamfProAllowedFileExtensions(),
			"jamfpro_api_integration":                             api_integration.ResourceJamfProApiIntegrations(),
			"jamfpro_api_role":                                    api_role.ResourceJamfProAPIRoles(),
			"jamfpro_app_installer":                               app_installer.ResourceJamfProAppInstallers(),
			"jamfpro_app_installer_global_settings":               app_installer_global_settings.ResourceJamfProAppInstallerGlobalSettings(),
			"jamfpro_building":                                    building.ResourceJamfProBuildings(),
			"jamfpro_category":                                    category.ResourceJamfProCategories(),
			"jamfpro_client_checkin":                              client_checkin.ResourceJamfProClientCheckin(),
			"jamfpro_cloud_ldap":                                  cloud_ldap.ResourceJamfProCloudLdap(),
			"jamfpro_computer_extension_attribute":                computer_extension_attribute.ResourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_inventory_collection_settings":      computer_inventory_collection_settings.ResourceJamfProComputerInventoryCollectionSettings(),
			"jamfpro_computer_prestage_enrollment":                computer_prestage_enrollment.ResourceJamfProComputerPrestageEnrollment(),
			"jamfpro_department":                                  department.ResourceJamfProDepartments(),
			"jamfpro_device_communication_settings":               device_communication_settings.ResourceJamfProDeviceCommunicationSettings(),
			"jamfpro_device_enrollments":                          device_enrollments.ResourceJamfProDeviceEnrollments(),
			"jamfpro_disk_encryption_configuration":               disk_encryption_configuration.ResourceJamfProDiskEncryptionConfigurations(),
			"jamfpro_dock_item":                                   dock_item.ResourceJamfProDockItems(),
			"jamfpro_engage_settings":                             engage_settings.ResourceEngageSettings(),
			"jamfpro_enrollment_customization":                    enrollment_customization.ResourceJamfProEnrollmentCustomization(),
			"jamfpro_file_share_distribution_point":               file_share_distribution_point.ResourceJamfProFileShareDistributionPoints(),
			"jamfpro_icon":                                        icon.ResourceJamfProIcons(),
			"jamfpro_jamf_connect":                                jamf_connect.ResourceJamfConnectConfigProfile(),
			"jamfpro_jamf_protect":                                jamf_protect.ResourceJamfProtect(),
			"jamfpro_ldap_server":                                 ldap_server.ResourceJamfProLDAPServers(),
			"jamfpro_local_admin_password_settings":               local_admin_password_settings.ResourceLocalAdminPasswordSettings(),
			"jamfpro_network_segment":                             network_segment.ResourceJamfProNetworkSegments(),
			"jamfpro_macos_configuration_profile_plist":           macos_configuration_profile_plist.ResourceJamfProMacOSConfigurationProfilesPlist(),
			"jamfpro_macos_configuration_profile_plist_generator": macos_configuration_profile_plist_generator.ResourceJamfProMacOSConfigurationProfilesPlistGenerator(),
			"jamfpro_managed_software_update":                     managed_software_update.ResourceJamfProManagedSoftwareUpdate(),
			"jamfpro_mobile_device_configuration_profile_plist":   mobile_device_configuration_profile_plist.ResourceJamfProMobileDeviceConfigurationProfilesPlist(),
			"jamfpro_mobile_device_extension_attribute":           mobile_device_extension_attribute.ResourceJamfProMobileDeviceExtensionAttributes(),
			"jamfpro_mobile_device_prestage_enrollment":           mobile_device_prestage_enrollment.ResourceJamfProMobileDevicePrestageEnrollment(),
			"jamfpro_package":                                     packages.ResourceJamfProPackages(),
			"jamfpro_policy":                                      policy.ResourceJamfProPolicies(),
			"jamfpro_printer":                                     printer.ResourceJamfProPrinters(),
			"jamfpro_script":                                      script.ResourceJamfProScripts(),
			"jamfpro_self_service_settings":                       self_service_settings.ResourceJamfProSelfServiceSettings(),
			"jamfpro_smtp_server":                                 smtp_server.ResourceJamfProSMTPServer(),
			"jamfpro_site":                                        site.ResourceJamfProSites(),
			"jamfpro_smart_computer_group":                        smart_computer_group.ResourceJamfProSmartComputerGroups(),
			"jamfpro_smart_mobile_device_group":                   smart_mobile_device_group.ResourceJamfProSmartMobileGroups(),
			"jamfpro_sso_certificate":                             sso_certificate.ResourceJamfProSSOCertificate(),
			"jamfpro_sso_failover":                                sso_failover.ResourceJamfProSSOFailover(),
			"jamfpro_sso_settings":                                sso_settings.ResourceJamfProSsoSettings(),
			"jamfpro_static_computer_group":                       static_computer_group.ResourceJamfProStaticComputerGroups(),
			"jamfpro_static_mobile_device_group":                  static_mobile_device_group.ResourceJamfProStaticMobileDeviceGroups(),
			"jamfpro_restricted_software":                         restricted_software.ResourceJamfProRestrictedSoftwares(),
			"jamfpro_user_initiated_enrollment_settings":          userinitiatedenrollment.ResourceJamfProUserInitatedEnrollmentSettings(),
			"jamfpro_user_group":                                  user_group.ResourceJamfProUserGroups(),
			"jamfpro_volume_purchasing_locations":                 volume_purchasing_locations.ResourceJamfProVolumePurchasingLocations(),
			"jamfpro_webhook":                                     webhook.ResourceJamfProWebhooks(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var err error
		var diags diag.Diagnostics
		var jamfIntegration *jamfprointegration.Integration
		var jamfFQDN,
			clientId,
			clientSecret,
			basicAuthUsername,
			basicAuthPassword string

		// Logger
		// Probably should move this into it's own function.
		enableClientLogs := d.Get("enable_client_sdk_logs").(bool)
		logFilePath := d.Get("client_sdk_log_export_path").(string)

		if !enableClientLogs && logFilePath != "" {
			return nil, append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Bad configuration",
				Detail:   "Cannot have enable_client_sdk_logs disabled with client_sdk_log_export_path set",
			})
		}

		defaultLoggerConfig := zap.NewProductionConfig()
		var logLevel zap.AtomicLevel
		if enableClientLogs {
			logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
		} else {
			logLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
		}

		if logFilePath != "" {
			if _, err := os.Stat(logFilePath); err != nil {
				return nil, append(diags, diag.FromErr(err)...)
			}
			defaultLoggerConfig.OutputPaths = append(defaultLoggerConfig.OutputPaths, logFilePath)
		}

		defaultLoggerConfig.Level = logLevel
		defaultLogger, err := defaultLoggerConfig.Build()

		if err != nil {
			return nil, append(diags, diag.FromErr(err)...)
		}

		sugaredLogger := defaultLogger.Sugar()

		// Auth
		jamfFQDN = GetJamfFqdn(d, &diags)
		authMethod := GetAuthMethod(d, &diags)
		tokenRefrshBufferPeriod := time.Duration(d.Get("token_refresh_buffer_period_seconds").(int)) * time.Second

		hide_sensitive_data := d.Get("hide_sensitive_data").(bool)
		bootstrapClient := http.Client{}
		switch authMethod {
		case "oauth2":
			clientId = GetClientID(d, &diags)
			clientSecret = GetClientSecret(d, &diags)
			jamfIntegration, err = jamfprointegration.BuildWithOAuth(
				jamfFQDN,
				sugaredLogger,
				tokenRefrshBufferPeriod,
				clientId,
				clientSecret,
				hide_sensitive_data,
				bootstrapClient,
			)

		case "basic":
			basicAuthUsername = GetBasicAuthUsername(d, &diags)
			basicAuthPassword = GetBasicAuthPassword(d, &diags)
			jamfIntegration, err = jamfprointegration.BuildWithBasicAuth(
				jamfFQDN,
				sugaredLogger,
				tokenRefrshBufferPeriod,
				basicAuthUsername,
				basicAuthPassword,
				hide_sensitive_data,
				bootstrapClient,
			)

		default:
			return nil, append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "invalid auth method supplied",
				Detail:   "You should not be able to find this error. If you have, please raise an issue with the schema.",
			})

		}

		if err != nil {
			return nil, append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error building jamf integration",
				Detail:   fmt.Sprintf("error: %v", err),
			})
		}

		// Cookies
		var cookiesList []*http.Cookie
		load_balancer_lock_enabled := d.Get("jamfpro_load_balancer_lock").(bool)
		customCookies := d.Get("custom_cookies")

		if load_balancer_lock_enabled {
			cookies, err := jamfIntegration.GetSessionCookies()
			if err != nil {
				return nil, append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error getting session cookies",
					Detail:   fmt.Sprintf("error: %v", err),
				})
			}

			cookiesList = append(cookiesList, cookies...)

		}

		if customCookies != nil && len(customCookies.([]interface{})) > 0 {
			for _, v := range customCookies.([]interface{}) {
				name := v.(map[string]interface{})["name"]
				value := v.(map[string]interface{})["value"]

				if name == jamfLoadBalancerCookieName && load_balancer_lock_enabled {
					return nil, append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cannot have load balancer lock and custom cookie of same name. (jpro-ingress)",
					})
				}

				httpCookie := &http.Cookie{
					Name:  name.(string),
					Value: value.(string),
				}

				cookiesList = append(cookiesList, httpCookie)
			}
		}

		// Timeout overrides for resourccs which will take longer than the default ones.
		// Adjusts for if the load balancer lock is on.
		timeoutOverrides := TimeoutOverrides(load_balancer_lock_enabled)

		for k, r := range provider.ResourcesMap {

			r.Timeouts = &schema.ResourceTimeout{}
			r.Timeouts.Create = new(time.Duration)
			r.Timeouts.Read = new(time.Duration)
			r.Timeouts.Update = new(time.Duration)
			r.Timeouts.Delete = new(time.Duration)

			if override, ok := timeoutOverrides[k]; ok {

				if override.Create != 0 {
					*r.Timeouts.Create = override.Create
				}

				if override.Read != 0 {
					*r.Timeouts.Read = override.Read
				}

				if override.Update != 0 {
					*r.Timeouts.Update = override.Update
				}

				if override.Delete != 0 {
					*r.Timeouts.Delete = override.Delete
				}

				continue
			}

			*r.Timeouts.Create = Timeout(load_balancer_lock_enabled)
			*r.Timeouts.Read = Timeout(load_balancer_lock_enabled)
			*r.Timeouts.Update = Timeout(load_balancer_lock_enabled)
			*r.Timeouts.Delete = Timeout(load_balancer_lock_enabled)
		}

		// Packaging
		config := httpclient.ClientConfig{
			Integration:              jamfIntegration,
			Sugar:                    sugaredLogger,
			HideSensitiveData:        d.Get("hide_sensitive_data").(bool),
			TokenRefreshBufferPeriod: tokenRefrshBufferPeriod,
			CustomCookies:            cookiesList,
			MandatoryRequestDelay:    time.Duration(d.Get("mandatory_request_delay_milliseconds").(int)) * time.Millisecond,
			RetryEligiableRequests:   false, // Forced because terraform handles concurrency
			HTTP:                     http.Client{},
		}

		httpClient, err := config.Build()
		if err != nil {
			return nil, append(diags, diag.FromErr(err)...)
		}

		jamfProSdk := jamfpro.Client{
			HTTP: httpClient,
		}

		return &jamfProSdk, diags
	}

	return provider
}
