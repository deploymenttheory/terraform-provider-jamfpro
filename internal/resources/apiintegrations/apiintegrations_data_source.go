// apiintegrations_data_source.go
package apiintegrations

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProApiIntegrations provides information about a specific API integration by its ID or Name.
func DataSourceJamfProApiIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProApiIntegrationsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the API integration.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the API integration.",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Indicates if the API integration is enabled.",
				Computed:    true,
			},
			"access_token_lifetime_seconds": {
				Type:        schema.TypeInt,
				Description: "The access token lifetime in seconds for the API integration.",
				Computed:    true,
			},
			"app_type": {
				Type:        schema.TypeString,
				Description: "The app type of the API integration.",
				Computed:    true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "The client ID of the API integration.",
				Computed:    true,
			},
			"authorization_scopes": {
				Type:        schema.TypeSet,
				Description: "The list of authorization scopes for the API integration.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// dataSourceJamfProApiIntegrationsRead fetches the details of a specific API integration
// from Jamf Pro using either its unique Name or its Id. The function prioritizes the 'display_name' attribute over the 'id'
// attribute for fetching details. If neither 'display_name' nor 'id' is provided, it returns an error.
// Once the details are fetched, they are set in the data source's state.
//
// Parameters:
// - ctx: The context within which the function is called. It's used for timeouts and cancellation.
// - d: The current state of the data source.
// - meta: The meta object that can be used to retrieve the API client connection.
//
// Returns:
// - diag.Diagnostics: Returns any diagnostics (errors or warnings) encountered during the function's execution.
func dataSourceJamfProApiIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var integration *jamfpro.ApiIntegration
	var err error

	// Check if DisplayName is provided in the data source configuration
	if v, ok := d.GetOk("display_name"); ok && v.(string) != "" {
		integrationName := v.(string)
		integration, err = conn.GetApiIntegrationNameByID(integrationName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch API integration by display name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		integrationID := v.(int) // Correctly cast to int
		integration, err = conn.GetApiIntegrationByID(integrationID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch API integration by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'display_name' or 'id' must be provided")
	}

	// Set the data source attributes using the fetched data
	d.SetId(fmt.Sprintf("%d", integration.ID))

	if err := d.Set("display_name", integration.DisplayName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'display_name': %v", err))
	}

	if err := d.Set("enabled", integration.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'enabled': %v", err))
	}

	if err := d.Set("access_token_lifetime_seconds", integration.AccessTokenLifetimeSeconds); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'access_token_lifetime_seconds': %v", err))
	}

	if err := d.Set("app_type", integration.AppType); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'app_type': %v", err))
	}

	if err := d.Set("client_id", integration.ClientID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'client_id': %v", err))
	}

	// Convert the authorization scopes to a schema.Set before setting it
	authorizationScopesSet := schema.NewSet(schema.HashString, []interface{}{})
	for _, scope := range integration.AuthorizationScopes {
		authorizationScopesSet.Add(scope)
	}

	if err := d.Set("authorization_scopes", authorizationScopesSet); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'authorization_scopes': %v", err))
	}

	return nil
}
