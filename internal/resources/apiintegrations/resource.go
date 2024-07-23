// apiintegrations_resource.go
package apiintegrations

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProApiIntegrations defines the schema and CRUD operations for managing Jamf Pro API Integrations in Terraform.
func ResourceJamfProApiIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: validateResourceAPIIntegrationsDataFields,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the API integration.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the API integration.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if the API integration is enabled.",
			},
			"access_token_lifetime_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The access token lifetime in seconds for the API integration.",
			},
			"app_type": {
				Type:     schema.TypeString,
				Computed: true,
				//Required:     true,
				Description: "The app type of the API integration.",
				//ValidateFunc: validation.StringInSlice([]string{"CLIENT_CREDENTIALS", "NATIVE_APP_OAUTH", "NONE"}, false),
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The client ID of the API integration.",
			},
			"authorization_scopes": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of authorization roles scoped to the API integration.",
			},
		},
	}
}
