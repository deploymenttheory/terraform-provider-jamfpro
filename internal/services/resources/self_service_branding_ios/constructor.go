// self_service_branding_ios_constructor.go
package self_service_branding_ios

import (
	"errors"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrIconIDNotInteger = errors.New("icon_id must be integer")
)

// construct builds a ResourceSelfServiceBrandingIOSDetail from the Terraform schema.
func construct(d *schema.ResourceData) (*jamfpro.ResourceSelfServiceBrandingIOSDetail, error) {
	branding := &jamfpro.ResourceSelfServiceBrandingIOSDetail{
		BrandingName:              d.Get("main_header").(string),
		HeaderBackgroundColorCode: d.Get("header_background_color_code").(string),
		MenuIconColorCode:         d.Get("menu_icon_color_code").(string),
		BrandingNameColorCode:     d.Get("branding_name_color_code").(string),
		StatusBarTextColor:        d.Get("status_bar_text_color").(string),
	}

	if v, ok := d.GetOk("icon_id"); ok {
		if id, ok := v.(int); ok {
			branding.IconId = &id
		} else {
			return nil, ErrIconIDNotInteger
		}
	}

	return branding, nil
}
