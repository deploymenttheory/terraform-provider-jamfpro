// printers_state.go
package printers

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Printer information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourcePrinter) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	resourceData := map[string]interface{}{
		"id":           strconv.Itoa(resp.ID),
		"name":         resp.Name,
		"category":     resp.Category,
		"uri":          resp.URI,
		"cups_name":    resp.CUPSName,
		"location":     resp.Location,
		"model":        resp.Model,
		"info":         resp.Info,
		"notes":        resp.Notes,
		"make_default": resp.MakeDefault,
		"use_generic":  resp.UseGeneric,
		"ppd":          resp.PPD,
		"ppd_path":     resp.PPDPath,
		"ppd_contents": resp.PPDContents,
	}

	// Iterate over the map and set each key-value pair in the Terraform state
	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
