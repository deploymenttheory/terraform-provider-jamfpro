package building

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProBuildings provides information about a specific building in Jamf Pro.
func DataSourceJamfProBuildings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the building.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the building.",
			},
			"street_address1": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The first line of the street address of the building.",
			},
			"street_address2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The second line of the street address of the building.",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The city in which the building is located.",
			},
			"state_province": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state or province in which the building is located.",
			},
			"zip_postal_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ZIP or postal code of the building.",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The country in which the building is located.",
			},
		},
	}
}

// dataSourceRead fetches the details of a specific building from Jamf Pro using either its unique Name or its Id.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	resourceID := d.Get("id").(string)
	name := d.Get("name").(string)

	if resourceID == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceBuilding
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if name != "" {
			resource, apiErr = client.GetBuildingByName(name)
		} else {
			resource, apiErr = client.GetBuildingByID(resourceID)
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
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Building with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the Jamf Pro Building was not found"))
	}

	d.SetId(resource.ID)
	return updateState(d, resource)
}
