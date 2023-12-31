// providers.go
package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/apiintegrations"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/apiroles"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/buildings"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computercheckin"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computerextensionattributes"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computergroups"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/computerinventory"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/departments"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/dockitems"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/policies"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/printers"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/scripts"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/sites"

	"github.com/deploymenttheory/terraform-provider-jamfpro/version"
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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jamfpro_api_integrations":              apiintegrations.DataSourceJamfProApiIntegrations(),
			"jamfpro_api_roles":                     apiroles.DataSourceJamfProAPIRoles(),
			"jamfpro_buildings":                     buildings.DataSourceJamfProBuilding(),
			"jamfpro_computer_extension_attributes": computerextensionattributes.DataSourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_groups":               computergroups.DataSourceJamfProComputerGroups(),
			"jamfpro_computer_inventory":            computerinventory.DataSourceJamfProComputerInventory(),
			"jamfpro_departments":                   departments.DataSourceJamfProDepartments(),
			"jamfpro_dock_items":                    dockitems.DataSourceJamfProDockItems(),
			"jamfpro_sites":                         sites.DataSourceJamfProSites(),
			"jamfpro_scripts":                       scripts.DataSourceJamfProScripts(),
			//"jamfpro_macos_configuration_profiles":  macosconfigurationprofiles.DataSourceJamfProMacOSConfigurationProfiles(),
			"jamfpro_policies": policies.DataSourceJamfProPolicies(),
			"jamfpro_printers": printers.DataSourceJamfProPrinters(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"jamfpro_api_integrations":              apiintegrations.ResourceJamfProApiIntegrations(),
			"jamfpro_api_roles":                     apiroles.ResourceJamfProAPIRoles(),
			"jamfpro_buildings":                     buildings.ResourceJamfProBuilding(),
			"jamfpro_computer_checkin":              computercheckin.ResourceJamfProComputerCheckin(),
			"jamfpro_computer_extension_attributes": computerextensionattributes.ResourceJamfProComputerExtensionAttributes(),
			"jamfpro_computer_groups":               computergroups.ResourceJamfProComputerGroups(),
			"jamfpro_departments":                   departments.ResourceJamfProDepartments(),
			"jamfpro_dock_items":                    dockitems.ResourceJamfProDockItems(),
			"jamfpro_sites":                         sites.ResourceJamfProSites(),
			"jamfpro_scripts":                       scripts.ResourceJamfProScripts(),
			//"jamfpro_macos_configuration_profiles":  macosconfigurationprofiles.ResourceJamfProMacOSConfigurationProfiles(),
			"jamfpro_policies": policies.ResourceJamfProPolicies(),
			"jamfpro_printers": printers.ResourceJamfProPrinters(),
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

		// Retrieve the log level from the configuration.
		logLevel := d.Get("log_level").(string)

		// Convert the log level from string to the LogLevel type.
		// (Assuming there's a function in your client package that does this)
		parsedLogLevel, err := client.ConvertToLogLevel(logLevel)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid log level",
				Detail:   err.Error(),
			})
			return nil, diags
		}

		config := client.ProviderConfig{
			InstanceName: instanceName,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			LogLevel:     parsedLogLevel,
			UserAgent:    provider.UserAgent(TerraformProviderProductUserAgent, version.ProviderVersion),
		}
		return config.Client()
	}
	return provider
}
