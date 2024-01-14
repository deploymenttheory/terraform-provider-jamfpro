// buildings_data_source.go
package buildings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProBuildings provides information about a specific building in Jamf Pro.
func DataSourceJamfProBuildings() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceBuildingRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the building.",
			},
			"name": {
				Type:        schema.TypeString,
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

// DataSourceBuildingRead fetches the details of a specific building from Jamf Pro using either its unique Name or its Id.
func DataSourceBuildingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var building *jamfpro.ResourceBuilding
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok {
		buildingName, ok := v.(string)
		if !ok {
			return diag.Errorf("expected 'name' to be a string")
		}
		if buildingName != "" {
			building, err = conn.GetBuildingByName(buildingName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to fetch building by name: %v", err))
			}
		}
	} else if v, ok := d.GetOk("id"); ok {
		buildingID, ok := v.(string)
		if !ok {
			return diag.Errorf("expected 'id' to be a string")
		}
		building, err = conn.GetBuildingByID(buildingID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch building by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	if building == nil {
		return diag.FromErr(fmt.Errorf("building not found"))
	}

	// Set the data source attributes using the fetched data
	if err := d.Set("name", building.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'name': %v", err))
	}
	if err := d.Set("streetAddress1", building.StreetAddress1); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'streetAddress1': %v", err))
	}
	if err := d.Set("streetAddress2", building.StreetAddress2); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'streetAddress2': %v", err))
	}
	if err := d.Set("city", building.City); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'city': %v", err))
	}
	if err := d.Set("stateProvince", building.StateProvince); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'stateProvince': %v", err))
	}
	if err := d.Set("zipPostalCode", building.ZipPostalCode); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'zipPostalCode': %v", err))
	}
	if err := d.Set("country", building.Country); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'country': %v", err))
	}

	d.SetId(building.ID)

	return nil
}
