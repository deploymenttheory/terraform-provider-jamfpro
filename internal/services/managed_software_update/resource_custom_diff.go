package managed_software_update

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i any) error {
	if err := validateGroupOrDevice(ctx, diff, i); err != nil {
		return err
	}

	if err := validateUpdateActionFields(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateGroupOrDevice ensures that either 'group' or 'device' is specified, but not both.
func validateGroupOrDevice(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	_, hasGroup := diff.GetOk("group")
	_, hasDevice := diff.GetOk("device")

	if hasGroup && hasDevice {
		return fmt.Errorf("in 'jamfpro_managed_software_update': only one of 'group' or 'device' can be specified, not both")
	}

	if !hasGroup && !hasDevice {
		return fmt.Errorf("in 'jamfpro_managed_software_update': either 'group' or 'device' must be specified")
	}

	return nil
}

// validateUpdateActionFields validates the root-level update action fields.
func validateUpdateActionFields(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	updateAction := diff.Get("update_action").(string)
	versionType := diff.Get("version_type").(string)
	specificVersion := diff.Get("specific_version").(string)
	maxDeferrals := diff.Get("max_deferrals").(int)

	if updateAction == "DOWNLOAD_INSTALL_ALLOW_DEFERRAL" && maxDeferrals == 0 {
		return fmt.Errorf("in 'jamfpro_managed_software_update': 'max_deferrals' must be set when 'update_action' is 'DOWNLOAD_INSTALL_ALLOW_DEFERRAL'")
	}

	if (versionType == "SPECIFIC_VERSION" || versionType == "CUSTOM_VERSION") && specificVersion == "" {
		return fmt.Errorf("in 'jamfpro_managed_software_update': 'specific_version' must be set when 'version_type' is 'SPECIFIC_VERSION' or 'CUSTOM_VERSION'")
	}

	return nil
}
