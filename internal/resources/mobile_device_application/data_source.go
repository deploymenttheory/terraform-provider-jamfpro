// mobile_device_applications_data_source.go
package mobile_device_application

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMobileDeviceApplications provides information about a specific Jamf Pro Mobile Device Application
func DataSourceJamfProMobileDeviceApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the mobile device application.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the mobile device application.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the mobile device application.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bundle ID of the mobile device application.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the mobile device application.",
			},
			"scope": {
				Type:        schema.TypeList,
				Description: "The scope of the mobile device application.",
				Computed:    true,
				Elem:        sharedschemas.GetSharedMobileDeviceSchemaScope(),
			},
		},
	}
}

// dataSourceRead fetches the details of a specific Jamf Pro mobile device application
// from Jamf Pro using its ID. Once the details are fetched, they are set in the data source's state.
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

	resourceID := d.Get("id").(string)
	name := d.Get("name").(string)

	if resourceID == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceMobileDeviceApplication
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if name != "" {
			resource, apiErr = client.GetMobileDeviceApplicationByName(name)
		} else {
			resource, apiErr = client.GetMobileDeviceApplicationByID(resourceID)
		}

		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if name != "" {
			lookupMethod = "name"
			lookupValue = name
		}
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Mobile Device Application with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the Jamf Pro Mobile Device Application was not found"))
	}

	d.SetId(fmt.Sprintf("%d", resource.General.ID))
	return updateState(d, resource)
}
