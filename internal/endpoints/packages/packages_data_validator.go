// packages_data_handling.go
package packages

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customValidateFilePath is a custom validation function for the package_file_source field.
// It ensures that the package_file_source field ends with .dmg if fill_user_template or fill_existing_users are set to true.
func customValidateFilePath(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	filePath, ok := d.Get("package_file_source").(string)
	if !ok {
		return fmt.Errorf("invalid type for package_file_sourceh")
	}

	// Check if the file path ends with .dmg
	if strings.HasSuffix(filePath, ".dmg") {
		// File path ends with .dmg, allow fill_user_template and fill_existing_users to be true
		return nil
	}

	// If file path does not end with .dmg, ensure fill_user_template and fill_existing_users are not set to true
	fillUserTemplate, fillUserTemplateOk := d.GetOk("fill_user_template")
	fillExistingUsers, fillExistingUsersOk := d.GetOk("fill_existing_users")

	if fillUserTemplateOk && fillUserTemplate.(bool) {
		return fmt.Errorf("fill_user_template can only be set to true if the package defined in package_file_source ends with .dmg")
	}

	if fillExistingUsersOk && fillExistingUsers.(bool) {
		return fmt.Errorf("fill_existing_users can only be set to true if the package defined in package_file_source ends with .dmg")
	}

	return nil
}
