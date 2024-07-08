// computercheckin_data_validation.go
package computercheckin

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateComputerCheckinDependencies checks the interdependencies of fields in the computer checkin configuration.
func validateComputerCheckinDependencies(d *schema.ResourceDiff) error {
	var errorMessages []string

	// Check 'log_startup_event' dependency
	if d.Get("log_startup_event").(bool) && !d.Get("create_startup_script").(bool) {
		errorMessages = append(errorMessages, "log_startup_event requires create_startup_script to be true")
	}

	// Check 'ensure_ssh_is_enabled' dependency
	if d.Get("ensure_ssh_is_enabled").(bool) && !d.Get("create_startup_script").(bool) {
		errorMessages = append(errorMessages, "ensure_ssh_is_enabled requires create_startup_script to be true")
	}

	// Check 'check_for_policies_at_startup' dependency
	if d.Get("check_for_policies_at_startup").(bool) && !d.Get("create_startup_script").(bool) {
		errorMessages = append(errorMessages, "check_for_policies_at_startup requires create_startup_script to be true")
	}

	// Check 'log_username' dependency
	if d.Get("log_username").(bool) && !d.Get("create_login_logout_hooks").(bool) {
		errorMessages = append(errorMessages, "log_username requires create_login_logout_hooks to be true")
	}

	// Check 'check_for_policies_at_login_logout' dependency
	if d.Get("check_for_policies_at_login_logout").(bool) && !d.Get("create_login_logout_hooks").(bool) {
		errorMessages = append(errorMessages, "check_for_policies_at_login_logout requires create_login_logout_hooks to be true")
	}

	// If there are any error messages, join them into a single error message
	if len(errorMessages) > 0 {
		return fmt.Errorf(strings.Join(errorMessages, "; "))
	}

	return nil
}
