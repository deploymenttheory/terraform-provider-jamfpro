// accountgroups_data_source.go
package account_group

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead fetches the details of specific account group from Jamf Pro using either their unique Name or Id.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var err error
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)
	resourceName := d.Get("name").(string)

	var resource *jamfpro.ResourceAccountGroup

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		if resourceID != "" {
				resource, apiErr = client.GetAccountGroupByID(resourceID)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
		}
		if resourceName != "" {
				resource, apiErr = client.GetAccountGroupByName(resourceName)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Account Group with ID '%s' after retries: %v", resourceID, err))
	}

	if resource != nil {
		d.SetId(strconv.Itoa(resource.ID))
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Account Group with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}
