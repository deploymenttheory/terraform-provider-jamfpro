// printers_state.go
package printers

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Printer information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourcePrinter) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	resourceData := map[string]interface{}{
		"id":           strconv.Itoa(resource.ID),
		"name":         resource.Name,
		"category":     resource.Category,
		"uri":          resource.URI,
		"cups_name":    resource.CUPSName,
		"location":     resource.Location,
		"model":        resource.Model,
		"info":         resource.Info,
		"notes":        resource.Notes,
		"make_default": resource.MakeDefault,
		"use_generic":  resource.UseGeneric,
		"ppd":          resource.PPD,
		"ppd_path":     resource.PPDPath,
		"ppd_contents": resource.PPDContents,
	}

	// Iterate over the map and set each key-value pair in the Terraform state
	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
