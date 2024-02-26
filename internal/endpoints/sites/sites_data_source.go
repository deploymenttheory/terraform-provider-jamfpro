// sites_data_source.go
package sites

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProSites provides information about a specific Jamf Pro site by its ID or Name.
func DataSourceJamfProSites() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProSitesRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the Jamf Pro site.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the Jamf Pro site.",
			},
		},
	}
}

// DataSourceJamfProSitesRead fetches the details of a specific Jamf Pro site
// from Jamf Pro using either its unique Name or its Id. The function prioritizes the 'name' attribute over the 'id'
// attribute for fetching details. If neither 'name' nor 'id' is provided, it returns an error.
// Once the details are fetched, they are set in the data source's state.
//
// Parameters:
// - ctx: The context within which the function is called. It's used for timeouts and cancellation.
// - d: The current state of the data source.
// - meta: The meta object that can be used to retrieve the API client connection.
//
// Returns:
// - diag.Diagnostics: Returns any diagnostics (errors or warnings) encountered during the function's execution.
func DataSourceJamfProSitesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Get the site ID from the data source's arguments
	resourceID, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("'id' must be provided")
	}
	resourceIDInt, err := strconv.Atoi(resourceID.(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting 'id' to int: %v", err))
	}

	var site *jamfpro.SharedResourceSite

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		site, apiErr = conn.GetSiteByID(resourceIDInt)
		if apiErr != nil {
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the data, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Site with ID '%d' after retries: %v", resourceIDInt, err))
	}

	// Check if resource data exists and set the Terraform state
	if site != nil {
		d.SetId(fmt.Sprintf("%d", resourceIDInt)) // Set the id in the Terraform state
		if err := d.Set("name", site.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Site with ID '%d': %v", resourceIDInt, err))...)
		}
	} else {
		d.SetId("") // Data not found, unset the id in the Terraform state
	}

	return diags
}
