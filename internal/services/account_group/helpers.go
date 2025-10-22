// accountgroups_data_handling.go
package account_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customDiffAccountGroups is a custom diff function for the Jamf Pro Account resource.
// This function is used during the Terraform plan phase to apply custom validation rules
// that are not covered by the basic schema validation.
func customDiffAccountGroups(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	accessLevel := d.Get("access_level").(string)

	if accessLevel == "Site Access" {
		if _, ok := d.GetOk("site_id"); !ok {
			return fmt.Errorf("'site' must be set when 'access_level' is 'Site Access'")
		}
	}

	return nil
}
