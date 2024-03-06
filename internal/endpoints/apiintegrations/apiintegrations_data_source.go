// apiintegrations_data_source.go
package apiintegrations

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

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
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var resource *jamfpro.ResourceApiIntegration

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetApiIntegrationByID(resourceIDInt)
		if apiErr != nil {
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the resource, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro API Integration with ID '%s' after retries: %v", resourceID, err))
	}

	// Check if resource data exists and set the Terraform state
	if resource != nil {
		d.SetId(resourceID) // Confirm the ID in the Terraform state
		if err := d.Set("display_name", resource.DisplayName); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro API Integration with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("") // Data not found, unset the ID in the Terraform state
	}

	return diags
}
