// filesharedistributionpoints_data_source.go
package filesharedistributionpoints

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProFileShareDistributionPoints defines the schema and CRUD operations for managing Jamf Pro Distribution Point in Terraform.
func DataSourceJamfProFileShareDistributionPoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the distribution point.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the distribution point.",
			},
		},
	}
}

// dataSourceReading the current state of a
// Jamf Pro File Share Distribution Point Resource from the remote system.
// The function:
// 1. Fetches the dock item's current state using its ID. If it fails then obtain dock item's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the dock item being deleted outside of Terraform, to keep the Terraform state synchronized.
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
