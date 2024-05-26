// providers.go
package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/deploymenttheory/go-api-http-client/httpclient"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/accountgroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/accounts"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/advancedcomputersearches"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/advancedmobiledevicesearches"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/advancedusersearches"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/allowedfileextensions"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/apiintegrations"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/apiroles"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/buildings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/categories"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computercheckin"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computerextensionattributes"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computergroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computerinventory"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computerinventorycollection"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computerprestageenrollments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/departments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/diskencryptionconfigurations"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/dockitems"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/filesharedistributionpoints"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/macosconfigurationprofiles"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/macosconfigurationprofilesw0de"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/mobiledeviceconfigurationprofiles"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/networksegments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/packages"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/policies"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/printers"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/restrictedsoftware"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/scripts"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/sites"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/usergroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/webhooks"
)

// TerraformProviderProductUserAgent is included in the User-Agent header for
// any API requests made by the provider.
const TerraformProviderProductUserAgent = "terraform-provider-jamfpro"

// GetInstanceName retrieves the 'instance_name' value from the Terraform configuration.
// If it's not present in the configuration, it attempts to fetch it from the JAMFPRO_INSTANCE_NAME environment variable.
func GetInstanceName(d *schema.ResourceData) (string, error) {
	instanceName := d.Get("instance_name").(string)
	if instanceName == "" {
		instanceName = os.Getenv("JAMFPRO_INSTANCE_NAME")
		if instanceName == "" {
			return "", fmt.Errorf("instance_name must be provided either as an environment variable (JAMFPRO_INSTANCE_NAME) or in the Terraform configuration")
		}
	}
	return instanceName, nil
}

// GetClientID retrieves the 'client_id' value from the Terraform configuration.
// If it's not present in the configuration, it attempts to fetch it from the JAMFPRO_CLIENT_ID environment variable.
func GetClientID(d *schema.ResourceData) (string, error) {
	clientID := d.Get("client_id").(string)
	if clientID == "" {
		clientID = os.Getenv("JAMFPRO_CLIENT_ID")
		if clientID == "" {
			return "", fmt.Errorf("client_id must be provided either as an environment variable (JAMFPRO_CLIENT_ID) or in the Terraform configuration")
		}
	}
	return clientID, nil
}

// GetClientSecret retrieves the 'client_secret' value from the Terraform configuration.
// If it's not present in the configuration, it attempts to fetch it from the JAMFPRO_CLIENT_SECRET environment variable.
func GetClientSecret(d *schema.ResourceData) (string, error) {
	clientSecret := d.Get("client_secret").(string)
	if clientSecret == "" {
		clientSecret = os.Getenv("JAMFPRO_CLIENT_SECRET")
		if clientSecret == "" {
			return "", fmt.Errorf("client_secret must be provided either as an environment variable (JAMFPRO_CLIENT_SECRET) or in the Terraform configuration")
		}
	}
	return clientSecret, nil
}

// GetClientUsername retrieves the 'username' value from the Terraform configuration.
// If it's not present in the configuration, it attempts to fetch it from the JAMFPRO_USERNAME environment variable.
func GetClientUsername(d *schema.ResourceData) (string, error) {
	username := d.Get("username").(string)
	if username == "" {
		username = os.Getenv("JAMFPRO_USERNAME")
		if username == "" {
			return "", fmt.Errorf("username must be provided either as an environment variable (JAMFPRO_USERNAME) or in the Terraform configuration")
		}
	}
	return username, nil
}

// GetClientPassword retrieves the 'password' value from the Terraform configuration.
// If it's not present in the configuration, it attempts to fetch it from the JAMFPRO_PASSWORD environment variable.
func GetClientPassword(d *schema.ResourceData) (string, error) {
	password := d.Get("password").(string)
	if password == "" {
		password = os.Getenv("JAMFPRO_PASSWORD")
		if password == "" {
			return "", fmt.Errorf("password must be provided either as an environment variable (JAMFPRO_PASSWORD) or in the Terraform configuration")
		}
	}
	return password, nil
}

// Schema defines the configuration attributes for the  within the JamfPro provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JAMFPRO_INSTANCE_NAME", ""),
				Description: "The Jamf Pro instance name. For https://mycompany.jamfcloud.com, define 'mycompany' in this field.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JAMFPRO_CLIENT_ID", ""),
				Description: "The Jamf Pro Client ID for authentication.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JAMFPRO_CLIENT_SECRET", ""),
				Description: "The Jamf Pro Client secret for authentication.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JAMFPRO_USERNAME", ""),
				Description: "The Jamf Pro username used for authentication.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JAMFPRO_PASSWORD", ""),
				Description: "The Jamf Pro password used for authentication.",
			},
			"log_level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "warning", // Align with the default log level in the  package
				ValidateFunc: validation.StringInSlice([]string{
					"debug", "info", "warning", "none",
				}, false),
				Description: "The logging level: debug, info, warning, or none",
			},
			"log_output_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "console", // Default to console for human-readable format
				Description: "The output format of the logs. Use 'JSON' for JSON format, 'console' for human-readable format. Defaults to console if no value is supplied.",
			},
			"log_console_separator": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     " ", // Set a default value for the separator
				Description: "The separator character used in console log output.",
			},
			"log_export_path": {
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
			"enable_cookie_jar": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable or disable the cookie jar for the HTTP client.",
			},
			"custom_cookies": {
				Type:     schema.TypeMap,
				Optional: true,
				Default:  nil,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"max_retry_attempts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "The maximum number of retry request attempts for retryable HTTP methods.",
			},
			"enable_dynamic_rate_limiting": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable dynamic rate limiting.",
			},
			"max_concurrent_requests": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "The maximum number of concurrent requests allowed in the semaphore.",
			},
			"token_refresh_buffer_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5, // Convert minutes to time.Duration in code
				Description: "The buffer period in minutes for token refresh.",
			},
			"total_retry_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60, // Convert seconds to time.Duration in code
				Description: "The total retry duration in seconds.",
			},
			"custom_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60, // Convert seconds to time.Duration in code
				Description: "The custom timeout in seconds for the HTTP client.",
			},
			"override_base_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Base domain override used when the default in the API handler isn't suitable.",
			},
			"api_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the API integration handler to use for the http client.",
				Default:     "jamfpro",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jamfpro_account":                             accounts.DataSourceJamfProAccounts(),
			"jamfpro_account_group":                       accountgroups.DataSourceJamfProAccountGroups(),
			"jamfpro_advanced_computer_search":            advancedcomputersearches.DataSourceJamfProAdvancedComputerSearches(),
			"jamfpro_advanced_mobile_device_search":       advancedmobiledevicesearches.DataSourceJamfProAdvancedMobileDeviceSearches(),
			"jamfpro_advanced_user_search":                advancedusersearches.DataSourceJamfProAdvancedUserSearches(),
			"jamfpro_api_integration":                     apiintegrations.DataSourceJamfProApiIntegrations(),
			"jamfpro_api_role":                            apiroles.DataSourceJamfProAPIRoles(),
			"jamfpro_building":                            buildings.DataSourceJamfProBuildings(),
			"jamfpro_category":                            categories.DataSourceJamfProCategories(),
			"jamfpro_computer_extension_attribute":        computerextensionattributes.DataSourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_group":                      computergroups.DataSourceJamfProComputerGroups(),
			"jamfpro_computer_inventory":                  computerinventory.DataSourceJamfProComputerInventory(),
			"jamfpro_computer_prestage_enrollment":        computerprestageenrollments.DataSourceJamfProComputerPrestageEnrollmentEnrollment(),
			"jamfpro_department":                          departments.DataSourceJamfProDepartments(),
			"jamfpro_disk_encryption_configuration":       diskencryptionconfigurations.DataSourceJamfProDiskEncryptionConfigurations(),
			"jamfpro_dock_item":                           dockitems.DataSourceJamfProDockItems(),
			"jamfpro_file_share_distribution_point":       filesharedistributionpoints.DataSourceJamfProFileShareDistributionPoints(),
			"jamfpro_network_segment":                     networksegments.DataSourceJamfProNetworkSegments(),
			"jamfpro_mobile_device_configuration_profile": mobiledeviceconfigurationprofiles.DataSourceJamfProMobileDeviceConfigurationProfiles(),
			"jamfpro_package":                             packages.DataSourceJamfProPackages(),
			// "jamfpro_policy":                        policies.DataSourceJamfProPolicies(),
			"jamfpro_printer":             printers.DataSourceJamfProPrinters(),
			"jamfpro_script":              scripts.DataSourceJamfProScripts(),
			"jamfpro_site":                sites.DataSourceJamfProSites(),
			"jamfpro_restricted_software": restrictedsoftware.DataSourceJamfProRestrictedSoftwares(),
			"jamfpro_user_group":          usergroups.DataSourceJamfProUserGroups(),
			"jamfpro_webhook":             webhooks.DataSourceJamfProWebhooks(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"jamfpro_account":                             accounts.ResourceJamfProAccounts(),
			"jamfpro_account_group":                       accountgroups.ResourceJamfProAccountGroups(),
			"jamfpro_advanced_computer_search":            advancedcomputersearches.ResourceJamfProAdvancedComputerSearches(),
			"jamfpro_advanced_mobile_device_search":       advancedmobiledevicesearches.ResourceJamfProAdvancedMobileDeviceSearches(),
			"jamfpro_advanced_user_search":                advancedusersearches.ResourceJamfProAdvancedUserSearches(),
			"jamfpro_allowed_file_extension":              allowedfileextensions.ResourceJamfProAllowedFileExtensions(),
			"jamfpro_api_integration":                     apiintegrations.ResourceJamfProApiIntegrations(),
			"jamfpro_api_role":                            apiroles.ResourceJamfProAPIRoles(),
			"jamfpro_building":                            buildings.ResourceJamfProBuildings(),
			"jamfpro_category":                            categories.ResourceJamfProCategories(),
			"jamfpro_computer_checkin":                    computercheckin.ResourceJamfProComputerCheckin(),
			"jamfpro_computer_extension_attribute":        computerextensionattributes.ResourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_group":                      computergroups.ResourceJamfProComputerGroups(),
			"jamfpro_computer_inventory_collection":       computerinventorycollection.ResourceJamfProComputerInventoryCollection(),
			"jamfpro_computer_prestage_enrollment":        computerprestageenrollments.ResourceJamfProComputerPrestageEnrollmentEnrollment(),
			"jamfpro_department":                          departments.ResourceJamfProDepartments(),
			"jamfpro_disk_encryption_configuration":       diskencryptionconfigurations.ResourceJamfProDiskEncryptionConfigurations(),
			"jamfpro_dock_item":                           dockitems.ResourceJamfProDockItems(),
			"jamfpro_file_share_distribution_point":       filesharedistributionpoints.ResourceJamfProFileShareDistributionPoints(),
			"jamfpro_network_segment":                     networksegments.ResourceJamfProNetworkSegments(),
			"jamfpro_macos_configuration_profile":         macosconfigurationprofiles.ResourceJamfProMacOSConfigurationProfiles(),
			"jamfpro_macos_configuration_profile_w0de":    macosconfigurationprofilesw0de.ResourceJamfProMacOSConfigurationProfiles(),
			"jamfpro_mobile_device_configuration_profile": mobiledeviceconfigurationprofiles.ResourceJamfProMobileDeviceConfigurationProfiles(),
			"jamfpro_package":                             packages.ResourceJamfProPackages(),
			"jamfpro_policy":                              policies.ResourceJamfProPolicies(),
			"jamfpro_printer":                             printers.ResourceJamfProPrinters(),
			"jamfpro_script":                              scripts.ResourceJamfProScripts(),
			"jamfpro_site":                                sites.ResourceJamfProSites(),
			"jamfpro_restricted_software":                 restrictedsoftware.ResourceJamfProRestrictedSoftwares(),
			"jamfpro_user_group":                          usergroups.ResourceJamfProUserGroups(),
			"jamfpro_webhook":                             webhooks.ResourceJamfProWebhooks(),
		},
	}

	provider.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		instanceName, err := GetInstanceName(d)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error getting instance name",
				Detail:   err.Error(),
			})
			return nil, diags
		}

		// Attempt to get client credentials (Client ID and Secret) or user credentials (Username and Password)
		clientID, errClientID := GetClientID(d)
		clientSecret, errClientSecret := GetClientSecret(d)
		username, errUsername := GetClientUsername(d)
		password, errPassword := GetClientPassword(d)

		// extract value for httpclient build and for determining resource propagation time
		enableCookieJar := d.Get("enable_cookie_jar").(bool)

		// Check if either pair of credentials is provided, prioritizing Client ID/Secret
		if errClientID == nil && errClientSecret == nil && clientID != "" && clientSecret != "" {
			// Client ID and Client Secret are provided
			// Initialize client with OAuth credentials
		} else if errUsername == nil && errPassword == nil && username != "" && password != "" {
			// Username and Password are provided
			// Initialize client with Username/Password credentials
		} else {
			// Neither set of credentials provided or incomplete set provided
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid Authentication Configuration",
				Detail:   "You must provide either a valid 'client_id' and 'client_secret' pair or a 'username' and 'password' pair for authentication.",
			})
			return nil, diags
		}

		// Translate the log level from the Terraform configuration
		logLevelStr := d.Get("log_level").(string)
		logLevel := logging.TranslateLogLevel(logLevelStr)

		// Build the HTTP client configuration
		httpClientConfig := httpclient.ClientConfig{
			Environment: httpclient.EnvironmentConfig{
				InstanceName:       instanceName,
				OverrideBaseDomain: d.Get("override_base_domain").(string),
				APIType:            "jamfpro",
			},
			Auth: httpclient.AuthConfig{
				Username:     username,
				Password:     password,
				ClientID:     clientID,
				ClientSecret: clientSecret,
			},
			ClientOptions: httpclient.ClientOptions{
				Logging: httpclient.LoggingConfig{
					LogLevel:            logLevel,
					LogOutputFormat:     d.Get("log_output_format").(string),
					LogConsoleSeparator: d.Get("log_console_separator").(string),
					LogExportPath:       d.Get("log_export_path").(string),
					HideSensitiveData:   d.Get("hide_sensitive_data").(bool),
				},
				Cookies: httpclient.CookieConfig{
					EnableCookieJar: enableCookieJar,
					CustomCookies:   make(map[string]string),
				},
				Retry: httpclient.RetryConfig{
					MaxRetryAttempts:          d.Get("max_retry_attempts").(int),
					EnableDynamicRateLimiting: d.Get("enable_dynamic_rate_limiting").(bool),
				},
				Concurrency: httpclient.ConcurrencyConfig{
					MaxConcurrentRequests: d.Get("max_concurrent_requests").(int),
				},
				Timeout: httpclient.TimeoutConfig{
					TokenRefreshBufferPeriod: time.Duration(d.Get("token_refresh_buffer_period").(int)) * time.Minute,
					TotalRetryDuration:       time.Duration(d.Get("total_retry_duration").(int)) * time.Second,
					CustomTimeout:            time.Duration(d.Get("custom_timeout").(int)) * time.Second,
				},
				Redirect: httpclient.RedirectConfig{},
			},
		}

		if d.Get("custom_cookies") != nil {
			httpClientConfig.ClientOptions.Cookies.CustomCookies[d.Get("custom_cookies.key").(string)] = d.Get("custom_cookies.value").(string)
		}

		// Conditionally print debug information.
		if !d.Get("hide_sensitive_data").(bool) {
			fmt.Printf("Debug: Building HTTP client with config: %+v\n", httpClientConfig)
		}

		httpclient, err := jamfpro.BuildClient(httpClientConfig)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		// Initialize the provider's APIClient struct with the Jamf Pro HTTP client and cookie jar setting
		jamfProAPIClient := client.APIClient{
			Conn:            httpclient,
			EnableCookieJar: enableCookieJar, // Allows use the cookie jar value within provider outside of the client
		}

		return &jamfProAPIClient, diags
	}
	return provider
}
