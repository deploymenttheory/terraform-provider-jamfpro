// staticcomputergroup_state.go
package staticcomputergroups

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the provided ResourceComputerGroup object.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceComputerGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", resource.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle Site
	site := []interface{}{}
	if resource.Site != nil && resource.Site.ID != -1 {
		site = append(site, map[string]interface{}{
			"id":   resource.Site.ID,
			"name": resource.Site.Name,
		})
	}

	if err := d.Set("site", site); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle Computers
	if resource.Computers != nil && resource.Computers.Computers != nil {
		computers := setComputerGroupSubsetComputersContainer(resource.Computers)
		if err := d.Set("computers", computers); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if err := d.Set("computers", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// setComputerGroupSubsetComputersContainer flattens a ComputerGroupSubsetComputersContainer object into a format suitable for Terraform state.
func setComputerGroupSubsetComputersContainer(computers *jamfpro.ComputerGroupSubsetComputersContainer) []interface{} {
	if computers == nil || computers.Computers == nil {
		return []interface{}{}
	}

	var computersList []interface{}
	for _, computer := range *computers.Computers {
		computerMap := map[string]interface{}{
			"id":              computer.ID,
			"name":            computer.Name,
			"serial_number":   computer.SerialNumber,
			"mac_address":     computer.MacAddress,
			"alt_mac_address": computer.AltMacAddress,
		}
		computersList = append(computersList, computerMap)
	}

	return computersList
}
