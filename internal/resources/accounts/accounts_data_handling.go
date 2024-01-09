// accounts_data_handling.go
package accounts

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customDiffAccounts is a custom diff function for the Jamf Pro Account resource.
// This function is used during the Terraform plan phase to apply custom validation rules
// that are not covered by the basic schema validation.
func customDiffAccounts(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	accessLevel, ok := d.GetOk("access_level")
	if !ok || accessLevel == nil {
		// If access_level is not set, no further checks required
		return nil
	}

	// Enforce that the 'site' attribute must be set if access_level is 'Site Access'
	if accessLevel.(string) == "Site Access" {
		if _, ok := d.GetOk("site"); !ok {
			return fmt.Errorf("'site' must be set when 'access_level' is 'Site Access'")
		}
	}

	return nil
}
