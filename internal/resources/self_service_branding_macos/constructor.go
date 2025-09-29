// self_service_branding_macos_constructor.go
package self_service_branding_macos

import (
	"errors"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrIconIDNotInteger                = errors.New("icon_id must be integer")
	ErrBrandingHeaderImageIDNotInteger = errors.New("branding_header_image_id must be integer")
)

// construct builds a ResourceSelfServiceBrandingDetail from the Terraform schema.
func construct(d *schema.ResourceData) (*jamfpro.ResourceSelfServiceBrandingDetail, error) {
	branding := &jamfpro.ResourceSelfServiceBrandingDetail{
		BrandingName:          d.Get("sidebar_heading").(string),
		BrandingNameSecondary: d.Get("sidebar_subheading").(string),
		ApplicationName:       d.Get("application_header").(string),
		HomeHeading:           d.Get("home_page_heading").(string),
		HomeSubheading:        d.Get("home_page_subheading").(string),
	}

	if v, ok := d.GetOk("icon_id"); ok {
		if id, ok := v.(int); ok {
			branding.IconId = &id
		} else {
			return nil, ErrIconIDNotInteger
		}
	}

	if v, ok := d.GetOk("home_page_banner_image_id"); ok {
		if id, ok := v.(int); ok {
			branding.BrandingHeaderImageId = &id
		} else {
			return nil, ErrBrandingHeaderImageIDNotInteger
		}
	}

	return branding, nil
}
