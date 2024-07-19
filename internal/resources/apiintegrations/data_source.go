// apiintegrations_data_source.go
package apiintegrations

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
				Computed:    true,
				Description: "The display name of the API integration.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "client id",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "client secret",
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
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	var resource *jamfpro.ResourceApiIntegration

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetApiIntegrationByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro API Integration with ID '%s' after retries: %v", resourceID, err))
	}

	if resource == nil {
		d.SetId("")
		return append(diags, diag.FromErr(fmt.Errorf("recieved empty resource"))...)
	}

	d.SetId(resourceID)

	resp, err := client.RefreshClientCredentialsByApiRoleID(resourceID)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if err = d.Set("display_name", resource.DisplayName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err = d.Set("client_id", resource.ClientID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err = d.Set("client_secret", resp.ClientSecret); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
