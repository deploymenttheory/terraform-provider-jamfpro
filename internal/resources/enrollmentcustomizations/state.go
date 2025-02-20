// enrollment_customization_state.go
package enrollmentcustomizations

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Enrollment Customization information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceEnrollmentCustomization) diag.Diagnostics {
	var diags diag.Diagnostics

	// Set main attributes
	resourceData := map[string]interface{}{
		"id":           resp.ID,
		"site_id":      resp.SiteID,
		"display_name": resp.DisplayName,
		"description":  resp.Description,
	}

	// Set branding settings as a list with one item
	brandingSettings := []interface{}{
		map[string]interface{}{
			"text_color":        resp.BrandingSettings.TextColor,
			"button_color":      resp.BrandingSettings.ButtonColor,
			"button_text_color": resp.BrandingSettings.ButtonTextColor,
			"background_color":  resp.BrandingSettings.BackgroundColor,
			"icon_url":          resp.BrandingSettings.IconUrl,
		},
	}

	// Set all values in state
	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set branding settings separately
	if err := d.Set("branding_settings", brandingSettings); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
