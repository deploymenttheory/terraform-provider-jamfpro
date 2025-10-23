// networksegments_state.go
package network_segment

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Network Segment information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceNetworkSegment) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]any{
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

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
