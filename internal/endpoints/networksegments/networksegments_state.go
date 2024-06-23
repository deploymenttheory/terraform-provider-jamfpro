// networksegments_state.go
package networksegments

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Network Segment information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceNetworkSegment) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	resourceData := map[string]interface{}{
		"id":                   strconv.Itoa(resp.ID),
		"name":                 resp.Name,
		"starting_address":     resp.StartingAddress,
		"ending_address":       resp.EndingAddress,
		"distribution_server":  resp.DistributionServer,
		"distribution_point":   resp.DistributionPoint,
		"url":                  resp.URL,
		"swu_server":           resp.SWUServer,
		"building":             resp.Building,
		"department":           resp.Department,
		"override_buildings":   resp.OverrideBuildings,
		"override_departments": resp.OverrideDepartments,
	}

	// Iterate over the map and set each key-value pair in the Terraform state
	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
