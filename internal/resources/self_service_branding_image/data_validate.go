// self_service_branding_image_data_validate.go
package self_service_branding_image

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	ErrExpectedPNGNoExt  = fmt.Errorf("expected .png file, got no extension")
	ErrExpectedPNGBadExt = fmt.Errorf("expected .png file, got invalid extension")
)

// validateImageFilePath ensures the provided file path has a .png extension
func validateImageFilePath() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.Any(
		validation.StringMatch(
			regexp.MustCompile(`^.*\.png$`),
			"Expected .png file, got .%s",
		),
		func(i interface{}, k string) ([]string, []error) {
			v := i.(string)
			ext := filepath.Ext(v)
			if ext == "" {
				return nil, []error{fmt.Errorf("%w", ErrExpectedPNGNoExt)}
			}
			return nil, []error{fmt.Errorf("%w: %s", ErrExpectedPNGBadExt, ext)}
		},
	))
}
