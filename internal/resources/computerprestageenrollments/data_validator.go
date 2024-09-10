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

	if err := validateRecoveryLockPasswordType(ctx, diff, i); err != nil {
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

func validateRecoveryLockPasswordType(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("display_name").(string)
	enableRecoveryLock, enableRecoveryLockOk := diff.GetOk("enable_recovery_lock")
	passwordType, passwordTypeOk := diff.GetOk("recovery_lock_password_type")
	password, passwordOk := diff.GetOk("recovery_lock_password")
	rotate, rotateOk := diff.GetOk("rotate_recovery_lock_password")

	// Scenario 1: When enable_recovery_lock is false
	if !enableRecoveryLockOk || !enableRecoveryLock.(bool) {
		if rotateOk && rotate.(bool) {
			return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'rotate_recovery_lock_password' must be false when 'enable_recovery_lock' is false", resourceName)
		}
		if passwordOk && password.(string) != "" {
			return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'recovery_lock_password' must be empty when 'enable_recovery_lock' is false", resourceName)
		}
		if passwordTypeOk && passwordType.(string) != "MANUAL" {
			return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'recovery_lock_password_type' must be 'MANUAL' when 'enable_recovery_lock' is false (this is the default value)", resourceName)
		}
		return nil
	}

	// For scenarios where enable_recovery_lock is true
	if !passwordTypeOk {
		return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'recovery_lock_password_type' must be set when 'enable_recovery_lock' is true", resourceName)
	}

	switch passwordType.(string) {
	case "RANDOM":
		// Scenario 2: Random password with rotation
		if passwordOk && password.(string) != "" {
			return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'recovery_lock_password' must be empty when 'recovery_lock_password_type' is 'RANDOM'", resourceName)
		}
		// Note: rotate_recovery_lock_password can be either true or false for RANDOM
	case "MANUAL":
		// Scenario 3: Manual password without rotation
		if !passwordOk || password.(string) == "" {
			return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'recovery_lock_password' must be set and non-empty when 'recovery_lock_password_type' is 'MANUAL'", resourceName)
		}
		if rotateOk && rotate.(bool) {
			return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'rotate_recovery_lock_password' must be false when 'recovery_lock_password_type' is 'MANUAL'", resourceName)
		}
	default:
		return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'recovery_lock_password_type' must be either 'MANUAL' or 'RANDOM'", resourceName)
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
	versionType := diff.Get("prestage_minimum_os_target_version_type").(string)
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
