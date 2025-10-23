// self_service_branding_macos_state.go
package self_service_branding_macos

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func updateState(d *schema.ResourceData, resp *jamfpro.ResourceSelfServiceBrandingDetail) diag.Diagnostics {
	var diags diag.Diagnostics

	if resp == nil {
		d.SetId("")
		return diags
	}

	if err := d.Set("sidebar_heading", resp.BrandingName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sidebar_subheading", resp.BrandingNameSecondary); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("application_header", resp.ApplicationName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("home_page_heading", resp.HomeHeading); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("home_page_subheading", resp.HomeSubheading); err != nil {
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
	if resp.BrandingHeaderImageId != nil {
		if err := d.Set("home_page_banner_image_id", *resp.BrandingHeaderImageId); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if err := d.Set("home_page_banner_image_id", nil); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if resp.IconId != nil {
		if err := d.Set("icon_id", *resp.IconId); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if err := d.Set("icon_id", nil); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
