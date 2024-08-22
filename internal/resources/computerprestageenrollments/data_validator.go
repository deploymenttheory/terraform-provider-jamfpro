// computerprestageenrollments_data_validator.go
package computerprestageenrollments

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	if err := validateAuthenticationPrompt(ctx, diff, i); err != nil {
		return err
	}

	if err := validateRecoveryLockPassword(ctx, diff, i); err != nil {
		return err
	}

	if err := validateRotateRecoveryLockPassword(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateAuthenticationPrompt checks that the 'authentication_prompt' is only set when 'require_authentication' is true.
func validateAuthenticationPrompt(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
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

// validateRecoveryLockPassword checks that 'recovery_lock_password' is only set when 'recovery_lock_password_type' is 'MANUAL'.
func validateRecoveryLockPassword(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	passwordType, passwordTypeOk := diff.GetOk("recovery_lock_password_type")
	password, passwordOk := diff.GetOk("recovery_lock_password")

	if !passwordTypeOk {
		return nil
	}

	if passwordType.(string) == "MANUAL" && !passwordOk {
		return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'recovery_lock_password' is required when 'recovery_lock_password_type' is set to 'MANUAL'", resourceName)
	}

	if passwordType.(string) != "MANUAL" && passwordOk && password.(string) != "" {
		return fmt.Errorf("in 'jamfpro_computer_prestage_enrollment.%s': 'recovery_lock_password' is not allowed when 'recovery_lock_password_type' is not set to 'MANUAL'", resourceName)
	}

	return nil
}

// validateRotateRecoveryLockPassword checks that 'rotate_recovery_lock_password' is only set when 'recovery_lock_password_type' is 'RANDOM'.
func validateRotateRecoveryLockPassword(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
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
