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

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/accountgroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/accounts"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/filesharedistributionpoints"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/macosconfigurationprofiles"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/advancedcomputersearches"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/advancedmobiledevicesearches"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/advancedusersearches"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/allowedfileextensions"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/apiintegrations"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/apiroles"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/buildings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/byoprofiles"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computercheckin"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computerextensionattributes"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computergroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computerinventory"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/computerprestages"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/departments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/diskencryptionconfigurations"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/dockitems"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/policies"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/printers"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/scripts"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/sites"
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

// Schema defines the configuration attributes for the http_client within the JamfPro provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JAMFPRO_INSTANCE_NAME", ""),
				Description: "The Jamf Pro instance name. For mycompany.jamfcloud.com, define mycompany in this field.",
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
			"log_level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "warning", // Set default log level as warning to align with http_client package
				ValidateFunc: validation.StringInSlice([]string{
					"debug", "info", "warning", "none",
				}, false),
				Description: "The logging level: debug, info, warning, or none",
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
				Default:     5,
				Description: "The buffer period for token refresh.",
			},
			"total_retry_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "The total retry duration.",
			},
			"custom_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "The custom timeout.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jamfpro_accounts":                      accounts.DataSourceJamfProAccounts(),
			"jamfpro_account_groups":                accountgroups.DataSourceJamfProAccountGroups(),
			"jamfpro_api_integrations":              apiintegrations.DataSourceJamfProApiIntegrations(),
			"jamfpro_api_roles":                     apiroles.DataSourceJamfProAPIRoles(),
			"jamfpro_buildings":                     buildings.DataSourceJamfProBuildings(),
			"jamfpro_computer_extension_attributes": computerextensionattributes.DataSourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_groups":               computergroups.DataSourceJamfProComputerGroups(),
			"jamfpro_computer_inventory":            computerinventory.DataSourceJamfProComputerInventory(),
			//"jamfpro_computer_prestages":            computerprestages.DataSourceJamfProComputerInventory(),
			"jamfpro_departments":                    departments.DataSourceJamfProDepartments(),
			"jamfpro_disk_encryption_configurations": diskencryptionconfigurations.DataSourceJamfProDiskEncryptionConfigurations(),
			"jamfpro_dock_items":                     dockitems.DataSourceJamfProDockItems(),
			"jamfpro_file_share_distribution_points": filesharedistributionpoints.DataSourceJamfProFileShareDistributionPoints(),
			"jamfpro_sites":                          sites.DataSourceJamfProSites(),
			"jamfpro_scripts":                        scripts.DataSourceJamfProScripts(),
			//"jamfpro_macos_configuration_profiles":  macosconfigurationprofiles.DataSourceJamfProMacOSConfigurationProfiles(),
			"jamfpro_policies": policies.DataSourceJamfProPolicies(),
			"jamfpro_printers": printers.DataSourceJamfProPrinters(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"jamfpro_accounts":                        accounts.ResourceJamfProAccounts(),
			"jamfpro_account_groups":                  accountgroups.ResourceJamfProAccountGroups(),
			"jamfpro_advanced_computer_searches":      advancedcomputersearches.ResourceJamfProAdvancedComputerSearches(),
			"jamfpro_advanced_mobile_device_searches": advancedmobiledevicesearches.ResourceJamfProAdvancedMobileDeviceSearches(),
			"jamfpro_advanced_user_searches":          advancedusersearches.ResourceJamfProAdvancedUserSearches(),
			"jamfpro_allowed_file_extension":          allowedfileextensions.ResourceJamfProAllowedFileExtensions(),
			"jamfpro_api_integrations":                apiintegrations.ResourceJamfProApiIntegrations(),
			"jamfpro_api_roles":                       apiroles.ResourceJamfProAPIRoles(),
			"jamfpro_buildings":                       buildings.ResourceJamfProBuildings(),
			"jamfpro_byoprofiles":                     byoprofiles.ResourceJamfProBYOProfiles(),
			"jamfpro_computer_checkin":                computercheckin.ResourceJamfProComputerCheckin(),
			"jamfpro_computer_extension_attributes":   computerextensionattributes.ResourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_groups":                 computergroups.ResourceJamfProComputerGroups(),
			"jamfpro_computer_prestages":              computerprestages.ResourceJamfProComputerPrestage(),
			"jamfpro_departments":                     departments.ResourceJamfProDepartments(),
			"jamfpro_disk_encryption_configurations":  diskencryptionconfigurations.ResourceJamfProDiskEncryptionConfigurations(),
			"jamfpro_dock_items":                      dockitems.ResourceJamfProDockItems(),
			"jamfpro_file_share_distribution_points":  filesharedistributionpoints.ResourceJamfProFileShareDistributionPoints(),
			"jamfpro_sites":                           sites.ResourceJamfProSites(),
			"jamfpro_scripts":                         scripts.ResourceJamfProScripts(),
			"jamfpro_macos_configuration_profile":     macosconfigurationprofiles.ResourceJamfProMacOSConfigurationProfiles(),
			"jamfpro_policies":                        policies.ResourceJamfProPolicies(),
			"jamfpro_printers":                        printers.ResourceJamfProPrinters(),
		},
	}

	provider.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		instanceName, err := GetInstanceName(d)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
				Detail:   err.Error(),
			})
			return nil, diags
		}

		clientID, err := GetClientID(d)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
				Detail:   err.Error(),
			})
			return nil, diags
		}

		clientSecret, err := GetClientSecret(d)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
				Detail:   err.Error(),
			})
			return nil, diags
		}

		MaxRetryAttempts := d.Get("max_retry_attempts").(int)
		EnableDynamicRateLimiting := d.Get("enable_dynamic_rate_limiting").(bool)
		MaxConcurrentRequests := d.Get("max_concurrent_requests").(int)

		TokenRefreshBufferPeriod := d.Get("token_refresh_buffer_period").(int)
		TokenRefreshBufferPeriodTyped := time.Duration(TokenRefreshBufferPeriod) * time.Minute

		TotalRetryDuration := d.Get("total_retry_duration").(int)
		TotalRetryDurationTyped := time.Duration(TotalRetryDuration) * time.Second

		CustomTimeout := d.Get("custom_timeout").(int)
		CustomTimeoutTyped := time.Duration(CustomTimeout) * time.Second

		// Convert the log level from string to the LogLevel type.
		// (Assuming there's a function in your client package that does this)
		logLevel := d.Get("log_level").(string)
		parsedLogLevel, err := client.ConvertToLogLevel(logLevel)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid log level",
				Detail:   err.Error(),
			})
			return nil, diags
		}

		httpClientConfig := http_client.Config{
			InstanceName:              instanceName,
			Auth:                      http_client.AuthConfig{ClientID: clientID, ClientSecret: clientSecret},
			LogLevel:                  parsedLogLevel,
			MaxRetryAttempts:          MaxRetryAttempts,
			EnableDynamicRateLimiting: EnableDynamicRateLimiting,
			MaxConcurrentRequests:     MaxConcurrentRequests,
			TokenRefreshBufferPeriod:  TokenRefreshBufferPeriodTyped,
			TotalRetryDuration:        TotalRetryDurationTyped,
			CustomTimeout:             CustomTimeoutTyped,
		}

		httpclient, err := jamfpro.NewClient(httpClientConfig)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		var wrapper client.APIClient
		wrapper.Conn = httpclient
		return &wrapper, nil
	}
	return provider
}
