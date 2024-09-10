// computerprestageenrollments_data_validator.go
package computerprestageenrollments

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

	if err := validateRotateRecoveryLockPassword(ctx, diff, i); err != nil {
		return err
	}

	if err := validateMinimumOSSpecificVersion(ctx, diff, i); err != nil {
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
		return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'authentication_prompt' is required when 'require_authentication' is set to true", resourceName)
	}

	if !requireAuth.(bool) && authPromptOk && authPrompt.(string) != "" {
		return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'authentication_prompt' is not allowed when 'require_authentication' is set to false", resourceName)
	}

	return nil
}

// validateRotateRecoveryLockPassword checks that 'rotate_recovery_lock_password' is only set when 'recovery_lock_password_type' is 'RANDOM'.
// Not part of the mainCustomDiffFunc as it is not comparing different schema values.
func validateRotateRecoveryLockPassword(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("display_name").(string)
	passwordType, passwordTypeOk := diff.GetOk("recovery_lock_password_type")
	rotate, rotateOk := diff.GetOk("rotate_recovery_lock_password")

	if !passwordTypeOk {
		return nil
	}

	if passwordType.(string) == "RANDOM" && !rotateOk {
		return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'rotate_recovery_lock_password' is required when 'recovery_lock_password_type' is set to 'RANDOM'", resourceName)
	}

	if passwordType.(string) != "RANDOM" && rotateOk && rotate.(bool) {
		return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'rotate_recovery_lock_password' is not allowed when 'recovery_lock_password_type' is not set to 'RANDOM'", resourceName)
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

func validateMinimumOSSpecificVersion(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("display_name").(string)
	versionType := diff.Get("prestate_minimum_os_target_version_type").(string)
	specificVersion := diff.Get("minimum_os_specific_version").(string)

	if versionType == "MINIMUM_OS_SPECIFIC_VERSION" {
		if specificVersion == "" {
			return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'minimum_os_specific_version' must be set when 'prestate_minimum_os_target_version_type' is MINIMUM_OS_SPECIFIC_VERSION", resourceName)
		}

		validVersions := map[string]bool{
			"14.5":   true,
			"14.6":   true,
			"14.6.1": true,
		}

		if !validVersions[specificVersion] {
			return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'minimum_os_specific_version' must be one of '14.5', '14.6', or '14.6.1', got: %s", resourceName, specificVersion)
		}
	}

	return nil
}
