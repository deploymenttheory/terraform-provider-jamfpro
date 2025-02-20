package enrollmentcustomizations

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructImageUpload creates a file path for image upload from Terraform configuration
func constructImageUpload(d *schema.ResourceData) (string, error) {
	imagePath := d.Get("enrollment_customization_image_source").(string)
	if imagePath == "" {
		return "", fmt.Errorf("enrollment_customization_image_source cannot be empty when specified")
	}
	return imagePath, nil
}

// construct creates a ResourceEnrollmentCustomization struct from Terraform configuration
func construct(d *schema.ResourceData) (*jamfpro.ResourceEnrollmentCustomization, error) {
	// Get branding settings from the schema
	brandingSettingsList := d.Get("branding_settings").([]interface{})
	if len(brandingSettingsList) == 0 {
		return nil, fmt.Errorf("branding_settings is required")
	}
	brandingSettingsData := brandingSettingsList[0].(map[string]interface{})

	// Create branding settings
	brandingSettings := jamfpro.EnrollmentCustomizationSubsetBrandingSettings{
		TextColor:       brandingSettingsData["text_color"].(string),
		ButtonColor:     brandingSettingsData["button_color"].(string),
		ButtonTextColor: brandingSettingsData["button_text_color"].(string),
		BackgroundColor: brandingSettingsData["background_color"].(string),
		IconUrl:         brandingSettingsData["icon_url"].(string),
	}

	// Create main resource
	resource := &jamfpro.ResourceEnrollmentCustomization{
		SiteID:           d.Get("site_id").(string),
		DisplayName:      d.Get("display_name").(string),
		Description:      d.Get("description").(string),
		BrandingSettings: brandingSettings,
	}

	// Include ID if it exists (for updates)
	if id := d.Id(); id != "" {
		resource.ID = id
	}

	return resource, nil
}
