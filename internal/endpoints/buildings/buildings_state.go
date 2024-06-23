// buildings_state.go
package buildings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Building information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceBuilding) diag.Diagnostics {

	var diags diag.Diagnostics

	// Map the configuration fields from the API response to a structured map
	buildingData := map[string]interface{}{
		"name":            resp.Name,
		"street_address1": resp.StreetAddress1,
		"street_address2": resp.StreetAddress2,
		"city":            resp.City,
		"state_province":  resp.StateProvince,
		"zip_postal_code": resp.ZipPostalCode,
		"country":         resp.Country,
	}

	// Set the structured map in the Terraform state
	for key, val := range buildingData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags

}
