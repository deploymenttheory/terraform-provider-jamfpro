// providers.go
package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/deploymenttheory/terraform-provider-jamfpro/version"
)

// TerraformProviderProductUserAgent is included in the User-Agent header for
// any API requests made by the provider.
const TerraformProviderProductUserAgent = "terraform-provider-jamfpro"

// GetInstanceName retrieves the 'instance_name' value from the Terraform configuration.
// If it's not present in the configuration, it attempts to fetch it from the JAMFPRO_INSTANCE environment variable.
func GetInstanceName(d *schema.ResourceData) (string, error) {
	instanceName := d.Get("instance_name").(string)
	if instanceName == "" {
		instanceName = os.Getenv("JAMFPRO_INSTANCE")
		if instanceName == "" {
			return "", fmt.Errorf("instance_name must be provided either as an environment variable (JAMFPRO_INSTANCE) or in the Terraform configuration")
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
				DefaultFunc: schema.EnvDefaultFunc("JAMFPRO_INSTANCE", ""),
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
			"debug_mode": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable or disable debug mode for verbose logging.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jamfpro_departments": dataSourceJamfProDepartments(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"jamfpro_departments": resourceJamfProDepartments(),
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

		config := ProviderConfig{
			InstanceName: instanceName,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			DebugMode:    d.Get("debug_mode").(bool),
			UserAgent:    provider.UserAgent(TerraformProviderProductUserAgent, version.ProviderVersion),
		}
		return config.Client()
	}
	return provider
}
