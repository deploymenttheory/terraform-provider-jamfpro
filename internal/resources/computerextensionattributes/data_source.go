// computerextensionattributes_data_source.go
package computerextensionattributes

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProComputerExtensionAttributes provides information about a specific computer extension attribute by its ID or Name.
func DataSourceJamfProComputerExtensionAttributes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The unique identifier of the Jamf Pro computer extension attribute.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The unique name of the Jamf Pro computer extension attribute.",
			},
		},
	}
}

// GetComputerExtensionAttributeByName
// dataSourceRead fetches the details of a specific computer extension attribute
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
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)
	resourceName := d.Get("name").(string)

	var resource *jamfpro.ResourceComputerExtensionAttribute
	retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if resourceID != "" {
			resource, apiErr = client.GetComputerExtensionAttributeByID(resourceID)
			if apiErr != nil {
				return retry.NonRetryableError(fmt.Errorf("failed to read Jamf Pro Computer Extension Attribute with ID '%s': %v", resourceID, apiErr))
			}
		} else if resourceName != "" {
			resource, apiErr = client.GetComputerExtensionAttributeByName(resourceName)
			if apiErr != nil {
				return retry.NonRetryableError(fmt.Errorf("failed to read Jamf Pro Computer Extension Attribute with Name '%s': %v", resourceName, apiErr))
			}
		} else {
			return retry.NonRetryableError(fmt.Errorf("either 'id' or 'name' must be provided"))
		}

		return nil
	})

	if resource != nil {
		d.SetId(resource.ID)
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Computer Extension Attribute with ID '%s': %v", resource.ID, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}
