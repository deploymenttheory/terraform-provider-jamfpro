// providers.go
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/deploymenttheory/terraform-provider-jamfpro/version"
)

// TerraformProviderProductUserAgent is included in the User-Agent header for
// any API requests made by the provider.
const TerraformProviderProductUserAgent = "terraform-provider-jamfpro"

// Schema defines the configuration attributes for the http_client within the JamfPro provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Jamf Pro instance name. For mycompany.jamfcloud.com, define mycompany in this field.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JAMFPRO_CLIENT_ID", nil),
				Description: "The Jamf Pro Client ID for authentication.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JAMFPRO_CLIENT_SECRET", nil),
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

	provider.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		config := ProviderConfig{
			InstanceName: d.Get("instance_name").(string),
			ClientID:     d.Get("client_id").(string),
			ClientSecret: d.Get("client_secret").(string),
			DebugMode:    d.Get("debug_mode").(bool),
			UserAgent:    provider.UserAgent(TerraformProviderProductUserAgent, version.ProviderVersion),
		}
		return config.Client()
	}

	return provider
}
