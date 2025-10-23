// state.go
package building

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Building information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceBuilding) diag.Diagnostics {
	var diags diag.Diagnostics

	buildingData := map[string]any{
		"name":            resp.Name,
		"street_address1": resp.StreetAddress1,
		"street_address2": resp.StreetAddress2,
		"city":            resp.City,
		"state_province":  resp.StateProvince,
		"zip_postal_code": resp.ZipPostalCode,
		"country":         resp.Country,
	}

	for key, val := range buildingData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags

}
