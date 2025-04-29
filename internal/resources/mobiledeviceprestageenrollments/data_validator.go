// mobiledeviceprestageenrollments_data_validator.go
package mobiledeviceprestageenrollments

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	if err := validateAuthenticationPrompt(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateAuthenticationPrompt checks that the 'authentication_prompt' is only set when 'require_authentication' is true.
func validateAuthenticationPrompt(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("display_name").(string)
	requireAuth, ok := diff.GetOk("require_authentication")

	if !ok {
		return nil
	}

	authPrompt, authPromptOk := diff.GetOk("authentication_prompt")

	if requireAuth.(bool) && !authPromptOk {
		return fmt.Errorf("in 'jamfpro_mobile_device_prestage_enrollment.%s': 'authentication_prompt' is required when 'require_authentication' is set to true", resourceName)
	}

	if !requireAuth.(bool) && authPromptOk && authPrompt.(string) != "" {
		return fmt.Errorf("in 'jamfpro_mobile_device_prestage_enrollment.%s': 'authentication_prompt' is not allowed when 'require_authentication' is set to false", resourceName)
	}

	return nil
}

// validateDateFormat checks that the date is in the format YYYY-MM-DD, but only if the value is not null or empty.
func validateDateFormat(v interface{}, k string) (ws []string, errors []error) {
	dateString, ok := v.(string)
	if !ok {
		return
	}

	if dateString == "" {
		return
	}

	datePattern := `^\d{4}-\d{2}-\d{2}$`
	match, _ := regexp.MatchString(datePattern, dateString)

	if !match {
		errors = append(errors, fmt.Errorf("%q must be in the format YYYY-MM-DD, got: %s", k, dateString))
	}

	return
}
