package static_mobile_device_group

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Mobile Device Group information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceMobileDeviceGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", resp.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	d.Set("site_id", resp.Site.ID)

	var assignments []any
	if resp.MobileDevices != nil {
		for _, comp := range resp.MobileDevices {
			assignments = append(assignments, comp.ID)
		}

		if err := d.Set("assigned_mobile_device_ids", assignments); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
