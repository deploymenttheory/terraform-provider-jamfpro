package icon

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	ErrNoPNGExtension   = errors.New("expected .png file, got no extension")
	ErrInvalidExtension = errors.New("expected .png file")
)

// validateIconFilePath ensures the provided file path has a .png extension
func validateIconFilePath() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.Any(
		validation.StringMatch(
			regexp.MustCompile(`^.*\.png$`),
			"Expected .png file",
		),
		func(i any, k string) ([]string, []error) {
			v := i.(string)
			ext := filepath.Ext(v)
			if ext == "" {
				return nil, []error{ErrNoPNGExtension}
			}
			return nil, []error{fmt.Errorf("%w, got %s", ErrInvalidExtension, ext)}
		},
	))
}
