// staticmobiledevicegroup_data_source.go
package staticmobiledevicegroups

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceJamfProStaticMobileDeviceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the static mobile device group.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the static mobile device group.",
			},
			"is_smart": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is a smart group.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The site ID for the group.",
			},
			"assigned_mobile_device_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of assigned mobile device IDs.",
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	resourceID := d.Get("id").(string)
	name := d.Get("name").(string)

	if resourceID == "" && name == "" {

		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceMobileDeviceGroup
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if name != "" {
			resource, apiErr = client.GetMobileDeviceGroupByName(name)
		} else {
			resource, apiErr = client.GetMobileDeviceGroupByID(resourceID)
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

		return diag.FromErr(fmt.Errorf("failed to read Static Mobile Device Group with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")

		return diag.FromErr(fmt.Errorf("the Jamf Pro Static Mobile Device Group was not found"))
	}

	d.SetId(fmt.Sprintf("%d", resource.ID))
	return updateState(d, resource)
}
