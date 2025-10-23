package file_share_distribution_point

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead reads the current state of a
// Jamf Pro File Share Distribution Point Resource from the remote system.
// The function:
// 1. Fetches the file share distribution point's current state using its ID.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the file share distribution point being deleted outside of Terraform, to keep the Terraform state synchronized.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	var resource *jamfpro.ResourceFileShareDistributionPoint
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetDistributionPointByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro File Share Distribution Point with ID '%s' after retries: %v", resourceID, err))
	}

	if resource != nil {
		d.SetId(resourceID)
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro File Share Distribution Point with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}
