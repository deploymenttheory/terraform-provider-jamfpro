// printers_state.go
package printer

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Printer information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourcePrinter) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"id":            strconv.Itoa(resp.ID),
		"name":          resp.Name,
		"category_name": resp.Category,
		"uri":           resp.URI,
		"cups_name":     resp.CUPSName,
		"location":      resp.Location,
		"model":         resp.Model,
		"info":          resp.Info,
		"notes":         resp.Notes,
		"make_default":  resp.MakeDefault,
		"use_generic":   resp.UseGeneric,
		"ppd":           resp.PPD,
		"ppd_path":      resp.PPDPath,
		"ppd_contents":  resp.PPDContents,
	}

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
