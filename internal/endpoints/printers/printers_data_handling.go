// printers_rdata_handling.go
package printers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateJamfProResourcePrinterDataFields enforces specific rules for printer configuration fields.
func validateJamfProResourcePrinterDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Retrieve the value of 'use_generic'
	useGeneric, ok := diff.GetOk("use_generic")
	if !ok {
		// If 'use_generic' is not set, no further validation is needed
		return nil
	}

	// Scenario 1: When 'use_generic' is true, 'ppd_path' must be set to the generic path
	if useGeneric.(bool) {
		expectedPPDPath := "/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/PrintCore.framework/Resources/Generic.ppd"
		if ppdPath, ok := diff.GetOk("ppd_path"); !ok || ppdPath.(string) != expectedPPDPath {
			return fmt.Errorf("when 'use_generic' is true, 'ppd_path' must be set to '%s'", expectedPPDPath)
		}
	}

	// Scenario 2: When 'use_generic' is true, 'ppd' must be empty
	if useGeneric.(bool) {
		if ppd, ok := diff.GetOk("ppd"); ok && ppd.(string) != "" {
			return fmt.Errorf("when 'use_generic' is true, 'ppd' must be empty")
		}
	}

	// Scenario 3: When 'use_generic' is false, 'ppd_path' must not be empty
	if !useGeneric.(bool) {
		if ppdPath, ok := diff.GetOk("ppd_path"); !ok || ppdPath == "" {
			return fmt.Errorf("when 'use_generic' is false, 'ppd_path' must not be empty")
		}
	}

	// Scenario 4: When 'use_generic' is false, 'ppd' must not be empty
	if !useGeneric.(bool) {
		if ppd, ok := diff.GetOk("ppd"); !ok || ppd == "" {
			return fmt.Errorf("when 'use_generic' is false, 'ppd' must not be empty")
		}
	}

	return nil
}
