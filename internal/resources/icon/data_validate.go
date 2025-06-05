package icons

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// validateIconFilePath ensures the provided file path has a .png extension
func validateIconFilePath() schema.SchemaValidateFunc {
	return validation.Any(
		validation.StringMatch(
			regexp.MustCompile(`^.*\.png$`),
			"Expected .png file, got .%s",
		),
		func(i interface{}, k string) ([]string, []error) {
			v := i.(string)
			ext := filepath.Ext(v)
			if ext == "" {
				return nil, []error{fmt.Errorf("expected .png file, got no extension")}
			}
			return nil, []error{fmt.Errorf("expected .png file, got %s", ext)}
		},
	)
}
