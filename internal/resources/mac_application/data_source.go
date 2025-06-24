package mac_application

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

// DataSourceJamfProMacApplications provides information about a specific Jamf Pro Mac Application
func DataSourceJamfProMacApplications() *schema.Resource {
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
				Description: "The unique identifier of the mac application.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the mac application.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the mac application.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bundle identifier of the mac application.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the mac application.",
			},
			"is_free": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the application is free.",
			},
			"deployment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deployment type of the mac application.",
			},
			"scope": {
				Type:        schema.TypeList,
				Description: "The scope of the mac application.",
				Computed:    true,
				Elem:        sharedschemas.GetSharedmacOSComputerSchemaScope(),
			},
		},
	}
}

// dataSourceRead fetches the details of a specific Jamf Pro mac application
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
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided")) //nolint:err113
	}

	var resource *jamfpro.ResourceMacApplications
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if name != "" {
			resource, apiErr = client.GetMacApplicationByName(name)
		} else {
			resource, apiErr = client.GetMacApplicationByID(resourceID)
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
		return diag.FromErr(fmt.Errorf("failed to read Mac Application with %s '%s' after retries: %v", lookupMethod, lookupValue, err)) //nolint:err113
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the Jamf Pro Mac Application was not found")) //nolint:err113
	}

	d.SetId(fmt.Sprintf("%d", resource.General.ID))
	return updateState(d, resource)
}
