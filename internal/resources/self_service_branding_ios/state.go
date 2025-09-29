// self_service_branding_ios_state.go
package self_service_branding_ios

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func updateState(d *schema.ResourceData, resp *jamfpro.ResourceSelfServiceBrandingIOSDetail) diag.Diagnostics {
	var diags diag.Diagnostics

	if resp == nil {
		d.SetId("")
		return diags
	}

	if err := d.Set("main_header", resp.BrandingName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if resp.ID != "" {
		d.SetId(resp.ID)
	}

	if resp.IconId != nil {
		if err := d.Set("icon_id", *resp.IconId); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if resp.IconId == nil {
		if err := d.Set("icon_id", nil); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if resp.HeaderBackgroundColorCode != "" {
		if err := d.Set("header_background_color_code", resp.HeaderBackgroundColorCode); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if resp.MenuIconColorCode != "" {
		if err := d.Set("menu_icon_color_code", resp.MenuIconColorCode); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if resp.BrandingNameColorCode != "" {
		if err := d.Set("branding_name_color_code", resp.BrandingNameColorCode); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if resp.StatusBarTextColor != "" {
		if err := d.Set("status_bar_text_color", resp.StatusBarTextColor); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
