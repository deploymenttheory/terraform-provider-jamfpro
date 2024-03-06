// apiroles_data_source.go
package apiroles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProAPIRoles provides information about a specific Jamf Pro API role by its ID or Name.
func DataSourceJamfProAPIRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProAPIRolesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the API role.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the Jamf Pro API role.",
			},
		},
	}
}

// DataSourceJamfProAPIRolesRead fetches the details of a specific API role from Jamf Pro using either its unique Name or its Id.
// The function prioritizes the 'name' attribute over the 'id' attribute for fetching details. If neither 'name' nor 'id' is provided,
// it returns an error. Once the details are fetched, they are set in the data source's state.
func DataSourceJamfProAPIRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	var resource *jamfpro.ResourceAPIRole

	// Read operation with retry
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetJamfApiRoleByID(resourceID)
		if apiErr != nil {
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the resource, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro API Role with ID '%s' after retries: %v", resourceID, err))
	}

	// Check if resource data exists and set the Terraform state
	if resource != nil {
		d.SetId(resourceID) // Confirm the ID in the Terraform state
		if err := d.Set("display_name", resource.DisplayName); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'display_name' for Jamf Pro API Role with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("") // Data not found, unset the ID in the Terraform state
	}

	return diags
}
