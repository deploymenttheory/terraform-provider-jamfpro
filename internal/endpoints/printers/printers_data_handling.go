// printers_rdata_handling.go
package printers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateJamfProResourcePrinterDataFields enforces specific rules for printer configuration fields based on the use_generic flag.
// Scenario 1:
// When 'use_generic' is set to true, this function enforces the following rules:
// - 'ppd_path' must be set to the default generic PPD path: "/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/PrintCore.framework/Resources/Generic.ppd"
// - 'ppd' must be empty.
// This ensures that when a generic printer driver is used, the printer configuration adheres to the expected standards for generic drivers in Jamf Pro.
//
// Scenario 2:
// When 'use_generic' is set to false, this function enforces that 'ppd' and 'ppd_path' must not be empty, ensuring that specific printer driver details are provided.
func validateJamfProResourcePrinterDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Scenario 1: Enforce rules when 'use_generic' is true
	if useGeneric, ok := diff.GetOk("use_generic"); ok && useGeneric.(bool) {
		expectedPPDPath := "/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/PrintCore.framework/Resources/Generic.ppd"
		ppdPath, ppdPathOk := diff.GetOk("ppd_path")
		ppd, ppdOk := diff.GetOk("ppd")

		if ppdPathOk && ppdPath.(string) != expectedPPDPath {
			return fmt.Errorf("when 'use_generic' is true, 'ppd_path' must be set to '%s'", expectedPPDPath)
		}

		if ppdOk && ppd.(string) != "" {
			return fmt.Errorf("when 'use_generic' is true, 'ppd' must be empty")
		}
	}

	// Scenario 2: Enforce rules when 'use_generic' is false
	if useGeneric, ok := diff.GetOk("use_generic"); ok && !useGeneric.(bool) {
		ppdPath, ppdPathOk := diff.GetOk("ppd_path")
		ppd, ppdOk := diff.GetOk("ppd")

		if !ppdPathOk || ppdPath.(string) == "" {
			return fmt.Errorf("when 'use_generic' is false, 'ppd_path' must not be empty")
		}

		if !ppdOk || ppd.(string) == "" {
			return fmt.Errorf("when 'use_generic' is false, 'ppd' must not be empty")
		}
	}

	return nil
}
