package enrollment_customization

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructImageUpload returns the path to an image suitable for upload from Terraform configuration
// It validates the image source and calls resizeImage if needed
func constructImageUpload(d *schema.ResourceData) (string, error) {
	imagePath := d.Get("enrollment_customization_image_source").(string)
	if imagePath == "" {
		return "", fmt.Errorf("enrollment_customization_image_source cannot be empty when specified")
	}

	ext := strings.ToLower(filepath.Ext(imagePath))
	if ext != ".png" && ext != ".gif" {
		return "", fmt.Errorf("image file must be PNG or GIF format, got: %s", ext)
	}

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	// Decode the image to determine its format and size
	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Check if resizing is needed (if image is larger than 180x180)
	if origWidth > 180 || origHeight > 180 {
		return common.ResizeImage(img, format, imagePath, 180, 180)
	}

	return imagePath, nil
}

// constructBaseResource creates a ResourceEnrollmentCustomization struct from Terraform configuration
func constructBaseResource(d *schema.ResourceData) (*jamfpro.ResourceEnrollmentCustomization, error) {
	brandingSettingsList := d.Get("branding_settings").([]interface{})
	if len(brandingSettingsList) == 0 {
		return nil, fmt.Errorf("branding_settings is required")
	}
	brandingSettingsData := brandingSettingsList[0].(map[string]interface{})

	brandingSettings := jamfpro.EnrollmentCustomizationSubsetBrandingSettings{
		TextColor:       brandingSettingsData["text_color"].(string),
		ButtonColor:     brandingSettingsData["button_color"].(string),
		ButtonTextColor: brandingSettingsData["button_text_color"].(string),
		BackgroundColor: brandingSettingsData["background_color"].(string),
		IconUrl:         brandingSettingsData["icon_url"].(string),
	}

	resource := &jamfpro.ResourceEnrollmentCustomization{
		SiteID:           d.Get("site_id").(string),
		DisplayName:      d.Get("display_name").(string),
		Description:      d.Get("description").(string),
		BrandingSettings: brandingSettings,
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Enrollment Customization '%s' to JSON: %v", resource.DisplayName, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Enrollment Customization JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// constructTextPane creates a ResourceEnrollmentCustomizationTextPane struct from Terraform configuration
func constructTextPane(data map[string]interface{}) (*jamfpro.ResourceEnrollmentCustomizationTextPane, error) {
	textPane := &jamfpro.ResourceEnrollmentCustomizationTextPane{
		Type:               "text",
		DisplayName:        data["display_name"].(string),
		Rank:               data["rank"].(int),
		Title:              data["title"].(string),
		Body:               data["body"].(string),
		Subtext:            data["subtext"].(string),
		BackButtonText:     data["back_button_text"].(string),
		ContinueButtonText: data["continue_button_text"].(string),
	}

	return textPane, nil
}

// constructLDAPPane creates a ResourceEnrollmentCustomizationLDAPPane struct from Terraform configuration
func constructLDAPPane(data map[string]interface{}) (*jamfpro.ResourceEnrollmentCustomizationLDAPPane, error) {
	ldapPane := &jamfpro.ResourceEnrollmentCustomizationLDAPPane{
		Type:               "ldap",
		DisplayName:        data["display_name"].(string),
		Rank:               data["rank"].(int),
		Title:              data["title"].(string),
		UsernameLabel:      data["username_label"].(string),
		PasswordLabel:      data["password_label"].(string),
		BackButtonText:     data["back_button_text"].(string),
		ContinueButtonText: data["continue_button_text"].(string),
	}

	// Process LDAP group access settings if present
	if groupsData, ok := data["ldap_group_access"].([]interface{}); ok && len(groupsData) > 0 {
		for _, groupData := range groupsData {
			group := groupData.(map[string]interface{})
			ldapGroupAccess := jamfpro.EnrollmentCustomizationLDAPGroupAccess{
				GroupName:    group["group_name"].(string),
				LDAPServerID: group["ldap_server_id"].(int),
			}
			ldapPane.LDAPGroupAccess = append(ldapPane.LDAPGroupAccess, ldapGroupAccess)
		}
	}

	return ldapPane, nil
}

// constructSSOPane creates a ResourceEnrollmentCustomizationSSOPane struct from Terraform configuration
func constructSSOPane(data map[string]interface{}) (*jamfpro.ResourceEnrollmentCustomizationSSOPane, error) {
	ssoPane := &jamfpro.ResourceEnrollmentCustomizationSSOPane{
		Type:                           "sso",
		DisplayName:                    data["display_name"].(string),
		Rank:                           data["rank"].(int),
		IsGroupEnrollmentAccessEnabled: data["is_group_enrollment_access_enabled"].(bool),
		GroupEnrollmentAccessName:      data["group_enrollment_access_name"].(string),
		IsUseJamfConnect:               data["is_use_jamf_connect"].(bool),
		ShortNameAttribute:             data["short_name_attribute"].(string),
		LongNameAttribute:              data["long_name_attribute"].(string),
	}

	return ssoPane, nil
}
