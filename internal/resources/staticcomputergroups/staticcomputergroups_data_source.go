// staticcomputergroup_data_source.go
package staticcomputergroups

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProStaticComputerGroups provides information about a specific computer group in Jamf Pro.
func DataSourceJamfProStaticComputerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProStaticComputerGroupsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the Jamf Pro static computer group.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the Jamf Pro static computer group.",
			},
		},
	}
}

// DataSourceJamfProStaticComputerGroupsRead fetches the details of a specific computer group
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
func DataSourceJamfProStaticComputerGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	var resource *jamfpro.ResourceComputerGroup
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetComputerGroupByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Jamf Pro Static Computer Group with ID '%s' after retries: %v", resourceID, err))
	}

	if resource != nil {
		d.SetId(resourceID)
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Static Computer Group with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}
