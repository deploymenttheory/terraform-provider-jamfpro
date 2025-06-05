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

	// Check 'startup_log' dependency
	startupLog, ok := d.Get("startup_log").(bool)
	createStartupScript, ok2 := d.Get("create_startup_script").(bool)
	if ok && ok2 && startupLog && !createStartupScript {
		errorMessages = append(errorMessages, "startup_log requires create_startup_script to be true")
	}

	// Check 'startup_ssh' dependency
	startupSsh, ok := d.Get("startup_ssh").(bool)
	if ok && ok2 && startupSsh && !createStartupScript {
		errorMessages = append(errorMessages, "startup_ssh requires create_startup_script to be true")
	}

	// Check 'startup_policies' dependency
	startupPolicies, ok := d.Get("startup_policies").(bool)
	if ok && ok2 && startupPolicies && !createStartupScript {
		errorMessages = append(errorMessages, "startup_policies requires create_startup_script to be true")
	}

	// Check 'hook_log' dependency
	hookLog, ok := d.Get("hook_log").(bool)
	createHooks, ok2 := d.Get("create_hooks").(bool)
	if ok && ok2 && hookLog && !createHooks {
		errorMessages = append(errorMessages, "hook_log requires create_hooks to be true")
	}

	// Check 'hook_policies' dependency
	hookPolicies, ok := d.Get("hook_policies").(bool)
	if ok && ok2 && hookPolicies && !createHooks {
		errorMessages = append(errorMessages, "hook_policies requires create_hooks to be true")
	}

	// If there are any error messages, return them with a constant format string
	if len(errorMessages) > 0 {

		return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
	}

	return nil
}
