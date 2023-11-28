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

// DataSourceJamfProBuilding provides information about a specific building in Jamf Pro.
func DataSourceJamfProBuilding() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceBuildingRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the building.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the building.",
				Computed:    true,
			},
			"streetAddress1": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The first line of the street address of the building.",
			},
			"streetAddress2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The second line of the street address of the building.",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The city in which the building is located.",
			},
			"stateProvince": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state or province in which the building is located.",
			},
			"zipPostalCode": {
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
	conn := meta.(*client.APIClient).Conn

	var building *jamfpro.ResponseBuilding
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		buildingName := v.(string)
		building, err = conn.GetBuildingByNameByID(buildingName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch building by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		buildingID := v.(string)
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
	d.SetId(building.ID)
	d.Set("name", building.Name)
	d.Set("streetAddress1", building.StreetAddress1)
	d.Set("streetAddress2", building.StreetAddress2)
	d.Set("city", building.City)
	d.Set("stateProvince", building.StateProvince)
	d.Set("zipPostalCode", building.ZipPostalCode)
	d.Set("country", building.Country)

	return nil
}
