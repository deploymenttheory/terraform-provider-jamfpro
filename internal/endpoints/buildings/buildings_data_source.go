// buildings_data_source.go
package buildings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProBuildings provides information about a specific building in Jamf Pro.
func DataSourceJamfProBuildings() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceBuildingRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the building.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the building.",
			},
		},
	}
}

// DataSourceBuildingRead fetches the details of a specific building from Jamf Pro using either its unique Name or its Id.
func DataSourceBuildingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	var resource *jamfpro.ResourceBuilding

	// Read operation with retry
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetBuildingByID(resourceID)
		if apiErr != nil {
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the resource, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Building with ID '%s' after retries: %v", resourceID, err))
	}

	// Check if resource data exists and set the Terraform state
	if resource != nil {
		d.SetId(resourceID) // Confirm the ID in the Terraform state
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Building with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("") // Data not found, unset the ID in the Terraform state
	}

	return diags
}
