// networksegments_state.go
package networksegments

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Network Segment information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceNetworkSegment) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	resourceData := map[string]interface{}{
		"id":                   strconv.Itoa(resource.ID),
		"name":                 resource.Name,
		"starting_address":     resource.StartingAddress,
		"ending_address":       resource.EndingAddress,
		"distribution_server":  resource.DistributionServer,
		"distribution_point":   resource.DistributionPoint,
		"url":                  resource.URL,
		"swu_server":           resource.SWUServer,
		"building":             resource.Building,
		"department":           resource.Department,
		"override_buildings":   resource.OverrideBuildings,
		"override_departments": resource.OverrideDepartments,
	}

	// Iterate over the map and set each key-value pair in the Terraform state
	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
