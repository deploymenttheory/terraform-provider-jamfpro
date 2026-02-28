// packages_custom_diff.go
package packages

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations for packages.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
	if err := customValidateFilePath(ctx, diff, meta); err != nil {
		return err
	}

	return computeFileHash(ctx, diff, meta)
}

// customValidateFilePath is a custom validation function for the package_file_source field.
// It ensures that the package_file_source field ends with .dmg if fill_user_template or fill_existing_users are set to true.
func customValidateFilePath(_ context.Context, d *schema.ResourceDiff, _ any) error {
	filePath, ok := d.Get("package_file_source").(string)
	if !ok {
		return fmt.Errorf("invalid type for package_file_source")
	}

	if strings.HasSuffix(filePath, ".dmg") {
		return nil
	}

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

// computeFileHash calculates the SHA3-512 hash of the local package file during the plan phase
// and compares it against the hash_value in state (JCDS-computed SHA3-512). If they differ,
// hash_value is updated which triggers a file re-upload during update.
//
// Guards:
//   - Skips on new resources (no ID yet) — create always uploads.
//   - Skips HTTP/HTTPS sources — URL string changes already trigger updates.
//   - Skips when hash_value in state is empty — JCDS hasn't computed the hash yet,
//     or this is a provider upgrade where the field hasn't been populated yet.
func computeFileHash(_ context.Context, d *schema.ResourceDiff, _ any) error {
	if d.Id() == "" {
		return nil
	}

	filePath, ok := d.Get("package_file_source").(string)
	if !ok || filePath == "" {
		return nil
	}

	if strings.HasPrefix(filePath, "http") {
		return nil
	}

	currentHash, ok := d.Get("hash_value").(string)
	if !ok || currentHash == "" {
		return nil
	}

	localHash, err := jamfpro.CalculateSHA3_512(filePath)
	if err != nil {
		return fmt.Errorf("failed to calculate SHA3-512 hash of package file: %v", err)
	}

	if localHash != currentHash {
		if err := d.SetNew("hash_value", localHash); err != nil {
			return fmt.Errorf("failed to set hash_value: %v", err)
		}
	}

	return nil
}
