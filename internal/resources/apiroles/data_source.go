// apiroles_data_source.go
package apiroles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProAPIRoles provides information about a specific Jamf Pro API role by its ID or Name.
func DataSourceJamfProAPIRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the Jamf API Role.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The display name of the Jamf API Role.",
			},
			"privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "List of privileges associated with the Jamf API Role.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// dataSourceRead fetches the details of a specific API role from Jamf Pro using either its unique Name or its Id.
// The function prioritizes the 'name' attribute over the 'id' attribute for fetching details. If neither 'name' nor 'id' is provided,
// it returns an error. Once the details are fetched, they are set in the data source's state.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	// Get the id and display_name from schema
	resourceID := d.Get("id").(string)
	displayName := d.Get("display_name").(string)

	// Validate that at least one identifier is provided
	if resourceID == "" && displayName == "" {

		return diag.FromErr(fmt.Errorf("either 'id' or 'display_name' must be provided"))
	}

	var resource *jamfpro.ResourceAPIRole
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		// Try to get by display_name first if provided
		if displayName != "" {
			resource, apiErr = client.GetJamfApiRoleByName(displayName)
		} else {
			// Fall back to ID lookup
			resource, apiErr = client.GetJamfApiRoleByID(resourceID)
		}

		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if displayName != "" {
			lookupMethod = "display name"
			lookupValue = displayName
		}

		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro API Role with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")

		return diag.FromErr(fmt.Errorf("jamf Pro API Role not found"))
	}

	d.SetId(resource.ID)
	return updateState(d, resource)
}
