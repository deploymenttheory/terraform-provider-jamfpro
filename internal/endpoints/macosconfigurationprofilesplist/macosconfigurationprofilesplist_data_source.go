// macosconfigurationprofilesplist_data_source.go
package macosconfigurationprofilesplist

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMacOSConfigurationProfilesPlist provides information about a specific department in Jamf Pro.
func DataSourceJamfProMacOSConfigurationProfilesPlist() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProMacOSConfigurationProfilePlistRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier for the macOS configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the macOS configuration profile.",
			},
		},
	}
}

// DataSourceJamfProMacOSConfigurationProfilePlistRead fetches the details of a macOS configuration profile.
func DataSourceJamfProMacOSConfigurationProfilePlistRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var resource *jamfpro.ResourceMacOSConfigurationProfile

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetMacOSConfigurationProfileByID(resourceIDInt)
		if apiErr != nil {
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the computer prestage, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro macOS Configuration Profile with ID '%s' after retries: %v", resourceID, err))
	}

	// Check if resource data exists and set the Terraform state
	if resource != nil {
		d.SetId(resourceID) // Confirm the ID in the Terraform state
		if err := d.Set("name", resource.General.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro macOS Configuration Profile with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("") // Data not found, unset the ID in the Terraform state
	}

	return diags
}
