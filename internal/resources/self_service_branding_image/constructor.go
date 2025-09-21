// self_service_branding_image_constructor.go
package self_service_branding_image

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrBothSources           = fmt.Errorf("cannot specify both branding_image_file_path and branding_image_file_web_source, choose one source only")
	ErrEitherSourceRequired  = fmt.Errorf("either branding_image_file_path or branding_image_file_web_source must be specified")
	ErrDownloadBrandingImage = fmt.Errorf("failed to download Self Service branding image")
)

// construct constructs a ResourceSelfServiceBrandingImage object from the provided schema data.
func construct(d *schema.ResourceData) (string, error) {
	filePath := d.Get("self_service_branding_image_file_path").(string)
	webSource := d.Get("self_service_branding_image_file_web_source").(string)

	if filePath != "" && webSource != "" {
		return "", fmt.Errorf("%w", ErrBothSources)
	}

	if filePath != "" {
		return filePath, nil
	}

	if webSource != "" {
		localPath, err := common.DownloadFile(webSource)
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrDownloadBrandingImage, webSource)
		}
		return localPath, nil
	}
	return "", fmt.Errorf("%w", ErrEitherSourceRequired)
}
