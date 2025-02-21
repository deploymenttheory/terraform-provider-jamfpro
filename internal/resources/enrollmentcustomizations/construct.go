package enrollmentcustomizations

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/gif" // register GIF format
	"image/png"
	_ "image/png" // register PNG format
	"os"
	"path/filepath"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nfnt/resize"
)

// constructImageUpload loads and processes an image file for upload from Terraform configuration
func constructImageUpload(d *schema.ResourceData) (string, error) {
	imagePath := d.Get("enrollment_customization_image_source").(string)
	if imagePath == "" {
		return "", fmt.Errorf("enrollment_customization_image_source cannot be empty when specified")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(imagePath))
	if ext != ".png" && ext != ".gif" {
		return "", fmt.Errorf("image file must be PNG or GIF format, got: %s", ext)
	}

	// Open and read the image file
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
		// Calculate scaling to maintain aspect ratio while fitting within 180x180
		ratio := float64(origWidth) / float64(origHeight)
		var newWidth, newHeight uint
		if ratio > 1 { // wider than tall
			newWidth = 180
			newHeight = uint(float64(180) / ratio)
		} else { // taller than wide
			newHeight = 180
			newWidth = uint(float64(180) * ratio)
		}

		// Resize the image while preserving aspect ratio
		resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

		tempDir := os.TempDir()
		tempFilename := filepath.Join(tempDir, "resized_"+filepath.Base(imagePath))

		tempFile, err := os.Create(tempFilename)
		if err != nil {
			return "", fmt.Errorf("failed to create temporary file for resized image: %v", err)
		}
		defer tempFile.Close()

		// Encode the image back to the temporary file
		if format == "png" || ext == ".png" {
			err = png.Encode(tempFile, resizedImg)
		} else if format == "gif" || ext == ".gif" {
			err = gif.Encode(tempFile, resizedImg, nil)
		} else {
			err = png.Encode(tempFile, resizedImg)
		}

		if err != nil {
			return "", fmt.Errorf("failed to encode processed image: %v", err)
		}

		return tempFilename, nil
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
