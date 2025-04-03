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
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/accountgroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/accounts"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/activationcode"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/advancedcomputersearches"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/advancedmobiledevicesearches"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/advancedusersearches"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/allowedfileextensions"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/apiintegrations"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/apiroles"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/appinstallerglobalsettings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/appinstallers"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/buildings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/categories"
	computercheckin "github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/clientcheckin"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computerextensionattributes"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computerinventory"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computerinventorycollection"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computerinventorycollectionsettings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computerprestageenrollments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/departments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/devicecommunicationsettings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/deviceenrollments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/diskencryptionconfigurations"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/dockitems"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/enrollmentcustomizations"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/filesharedistributionpoints"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/icons"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/jamfconnect"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/localadminpasswordsettings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/macosconfigurationprofilesplist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/macosconfigurationprofilesplistgenerator"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/managedsoftwareupdates"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/mobiledeviceconfigurationprofilesplist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/mobiledeviceextensionattributes"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/networksegments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/packages"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/policies"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/printers"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/restrictedsoftware"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/scripts"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/sites"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/smartcomputergroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/smartmobiledevicegroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/smtpserver"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/staticcomputergroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/usergroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/webhooks"
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
				Description: "Auth method chosen for Jamf.",
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
				Description: "Programatically determines all available web app members in the load balance and locks all instances of httpclient to the app for faster executions. \nTEMP SOLUTION UNTIL JAMF PROVIDES SOLUTION",
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

			"jamfpro_account":                                   accounts.DataSourceJamfProAccounts(),
			"jamfpro_account_group":                             accountgroups.DataSourceJamfProAccountGroups(),
			"jamfpro_advanced_computer_search":                  advancedcomputersearches.DataSourceJamfProAdvancedComputerSearches(),
			"jamfpro_advanced_mobile_device_search":             advancedmobiledevicesearches.DataSourceJamfProAdvancedMobileDeviceSearches(),
			"jamfpro_advanced_user_search":                      advancedusersearches.DataSourceJamfProAdvancedUserSearches(),
			"jamfpro_api_integration":                           apiintegrations.DataSourceJamfProApiIntegrations(),
			"jamfpro_api_role":                                  apiroles.DataSourceJamfProAPIRoles(),
			"jamfpro_app_installer":                             appinstallers.DataSourceJamfProAppInstallers(),
			"jamfpro_building":                                  buildings.DataSourceJamfProBuildings(),
			"jamfpro_category":                                  categories.DataSourceJamfProCategories(),
			"jamfpro_computer_extension_attribute":              computerextensionattributes.DataSourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_inventory":                        computerinventory.DataSourceJamfProComputerInventory(),
			"jamfpro_computer_prestage_enrollment":              computerprestageenrollments.DataSourceJamfProComputerPrestageEnrollmentEnrollment(),
			"jamfpro_department":                                departments.DataSourceJamfProDepartments(),
			"jamfpro_device_enrollments":                        deviceenrollments.DataSourceJamfProDeviceEnrollments(),
			"jamfpro_disk_encryption_configuration":             diskencryptionconfigurations.DataSourceJamfProDiskEncryptionConfigurations(),
			"jamfpro_dock_item":                                 dockitems.DataSourceJamfProDockItems(),
			"jamfpro_file_share_distribution_point":             filesharedistributionpoints.DataSourceJamfProFileShareDistributionPoints(),
			"jamfpro_network_segment":                           networksegments.DataSourceJamfProNetworkSegments(),
			"jamfpro_macos_configuration_profile_plist":         macosconfigurationprofilesplist.DataSourceJamfProMacOSConfigurationProfilesPlist(),
			"jamfpro_mobile_device_configuration_profile_plist": mobiledeviceconfigurationprofilesplist.DataSourceJamfProMobileDeviceConfigurationProfilesPlist(),

			/* "jamfpro_mobile_device_extension_attribute":         mobiledeviceextensionattribute.DataSourceJamfProMobileDeviceExtensionAttributes(), */
			"jamfpro_package":                   packages.DataSourceJamfProPackages(),
			"jamfpro_policy":                    policies.DataSourceJamfProPolicies(),
			"jamfpro_printer":                   printers.DataSourceJamfProPrinters(),
			"jamfpro_script":                    scripts.DataSourceJamfProScripts(),
			"jamfpro_site":                      sites.DataSourceJamfProSites(),
			"jamfpro_smart_computer_group":      smartcomputergroups.DataSourceJamfProSmartComputerGroups(),
			"jamfpro_smart_mobile_device_group": smartmobiledevicegroups.DataSourceJamfProSmartMobileGroups(),
			"jamfpro_static_computer_group":     staticcomputergroups.DataSourceJamfProStaticComputerGroups(),
			"jamfpro_restricted_software":       restrictedsoftware.DataSourceJamfProRestrictedSoftwares(),
			"jamfpro_user_group":                usergroups.DataSourceJamfProUserGroups(),
			"jamfpro_webhook":                   webhooks.DataSourceJamfProWebhooks(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"jamfpro_account":                                     accounts.ResourceJamfProAccounts(),
			"jamfpro_account_group":                               accountgroups.ResourceJamfProAccountGroups(),
			"jamfpro_activation_code":                             activationcode.ResourceJamfProActivationCode(),
			"jamfpro_advanced_computer_search":                    advancedcomputersearches.ResourceJamfProAdvancedComputerSearches(),
			"jamfpro_advanced_mobile_device_search":               advancedmobiledevicesearches.ResourceJamfProAdvancedMobileDeviceSearches(),
			"jamfpro_advanced_user_search":                        advancedusersearches.ResourceJamfProAdvancedUserSearches(),
			"jamfpro_allowed_file_extension":                      allowedfileextensions.ResourceJamfProAllowedFileExtensions(),
			"jamfpro_api_integration":                             apiintegrations.ResourceJamfProApiIntegrations(),
			"jamfpro_api_role":                                    apiroles.ResourceJamfProAPIRoles(),
			"jamfpro_app_installer":                               appinstallers.ResourceJamfProAppInstallers(),
			"jamfpro_app_installer_global_settings":               appinstallerglobalsettings.ResourceJamfProAppInstallerGlobalSettings(),
			"jamfpro_building":                                    buildings.ResourceJamfProBuildings(),
			"jamfpro_category":                                    categories.ResourceJamfProCategories(),
			"jamfpro_client_checkin":                              computercheckin.ResourceJamfProClientCheckin(),
			"jamfpro_computer_extension_attribute":                computerextensionattributes.ResourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_inventory_collection":               computerinventorycollection.ResourceJamfProComputerInventoryCollection(),
			"jamfpro_computer_inventory_collection_settings":      computerinventorycollectionsettings.ResourceJamfProComputerInventoryCollectionSettings(),
			"jamfpro_computer_prestage_enrollment":                computerprestageenrollments.ResourceJamfProComputerPrestageEnrollmentEnrollment(),
			"jamfpro_department":                                  departments.ResourceJamfProDepartments(),
			"jamfpro_device_communication_settings":               devicecommunicationsettings.ResourceJamfProDeviceCommunicationSettings(),
			"jamfpro_disk_encryption_configuration":               diskencryptionconfigurations.ResourceJamfProDiskEncryptionConfigurations(),
			"jamfpro_dock_item":                                   dockitems.ResourceJamfProDockItems(),
			"jamfpro_enrollment_customization":                    enrollmentcustomizations.ResourceJamfProEnrollmentCustomization(),
			"jamfpro_file_share_distribution_point":               filesharedistributionpoints.ResourceJamfProFileShareDistributionPoints(),
			"jamfpro_icon":                                        icons.ResourceJamfProIcons(),
			"jamfpro_jamf_connect":                                jamfconnect.ResourceJamfConnectConfigProfile(),
			"jamfpro_network_segment":                             networksegments.ResourceJamfProNetworkSegments(),
			"jamfpro_macos_configuration_profile_plist":           macosconfigurationprofilesplist.ResourceJamfProMacOSConfigurationProfilesPlist(),
			"jamfpro_local_admin_password_settings":               localadminpasswordsettings.ResourceLocalAdminPasswordSettings(),
			"jamfpro_macos_configuration_profile_plist_generator": macosconfigurationprofilesplistgenerator.ResourceJamfProMacOSConfigurationProfilesPlistGenerator(),
			"jamfpro_managed_software_update":                     managedsoftwareupdates.ResourceJamfProManagedSoftwareUpdate(),
			"jamfpro_mobile_device_configuration_profile_plist":   mobiledeviceconfigurationprofilesplist.ResourceJamfProMobileDeviceConfigurationProfilesPlist(),
			"jamfpro_mobile_device_extension_attribute":           mobiledeviceextensionattributes.ResourceJamfProMobileDeviceExtensionAttributes(),
			"jamfpro_package":                                     packages.ResourceJamfProPackages(),
			"jamfpro_policy":                                      policies.ResourceJamfProPolicies(),
			"jamfpro_printer":                                     printers.ResourceJamfProPrinters(),
			"jamfpro_script":                                      scripts.ResourceJamfProScripts(),
			"jamfpro_smtp_server":                                 smtpserver.ResourceJamfProSMTPServer(),
			"jamfpro_site":                                        sites.ResourceJamfProSites(),
			"jamfpro_smart_computer_group":                        smartcomputergroups.ResourceJamfProSmartComputerGroups(),
			"jamfpro_smart_mobile_device_group":                   smartmobiledevicegroups.ResourceJamfProSmartMobileGroups(),
			"jamfpro_static_computer_group":                       staticcomputergroups.ResourceJamfProStaticComputerGroups(),
			"jamfpro_restricted_software":                         restrictedsoftware.ResourceJamfProRestrictedSoftwares(),
			"jamfpro_user_group":                                  usergroups.ResourceJamfProUserGroups(),
			"jamfpro_webhook":                                     webhooks.ResourceJamfProWebhooks(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
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

		if customCookies != nil && len(customCookies.([]any)) > 0 {
			for _, v := range customCookies.([]any) {
				name := v.(map[string]any)["name"]
				value := v.(map[string]any)["value"]

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
