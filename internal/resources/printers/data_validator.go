// printers_data_validator.go
package printers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	// Validate printer configuration fields
	if err := validateJamfProResourcePrinterDataFields(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateJamfProResourcePrinterDataFields enforces specific rules for printer configuration fields.
func validateJamfProResourcePrinterDataFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)

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
			return fmt.Errorf("in 'jamfpro_printer.%s': when 'use_generic' is true, 'ppd_path' must be set to '%s'", resourceName, expectedPPDPath)
		}
	}

	// Scenario 2: When 'use_generic' is true, 'ppd' must be empty
	if useGeneric.(bool) {
		if ppd, ok := diff.GetOk("ppd"); ok && ppd.(string) != "" {

			return fmt.Errorf("in 'jamfpro_printer.%s': when 'use_generic' is true, 'ppd' must be empty", resourceName)
		}
	}

	// Scenario 3: When 'use_generic' is false, 'ppd_path' must not be empty
	if !useGeneric.(bool) {
		if ppdPath, ok := diff.GetOk("ppd_path"); !ok || ppdPath == "" {

			return fmt.Errorf("in 'jamfpro_printer.%s': when 'use_generic' is false, 'ppd_path' must not be empty", resourceName)
		}
	}

	// Scenario 4: When 'use_generic' is false, 'ppd' must not be empty
	if !useGeneric.(bool) {
		if ppd, ok := diff.GetOk("ppd"); !ok || ppd == "" {

			return fmt.Errorf("in 'jamfpro_printer.%s': when 'use_generic' is false, 'ppd' must not be empty", resourceName)
		}
	}

	return nil
}
