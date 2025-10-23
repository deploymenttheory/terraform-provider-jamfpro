package managed_software_update

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	if err := validateGroupOrDevice(ctx, diff, i); err != nil {
		return err
	}

	if err := validateConfigFields(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateGroupOrDevice ensures that either 'group' or 'device' is specified, but not both.
func validateGroupOrDevice(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	_, hasGroup := diff.GetOk("group")
	_, hasDevice := diff.GetOk("device")

	if hasGroup && hasDevice {
		return fmt.Errorf("in 'jamfpro_managed_software_update.%s': only one of 'group' or 'device' can be specified, not both", resourceName)
	}

	if !hasGroup && !hasDevice {
		return fmt.Errorf("in 'jamfpro_managed_software_update.%s': either 'group' or 'device' must be specified", resourceName)
	}

	return nil
}

// validateConfigFields performs validation on the 'config' block fields.
func validateConfigFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	config := diff.Get("config").([]interface{})

	if len(config) == 0 {
		return fmt.Errorf("in 'jamfpro_managed_software_update.%s': 'config' block is required", resourceName)
	}

	configMap := config[0].(map[string]interface{})
	updateAction := configMap["update_action"].(string)
	versionType := configMap["version_type"].(string)
	specificVersion := configMap["specific_version"].(string)
	maxDeferrals := configMap["max_deferrals"].(int)

	if updateAction == "DOWNLOAD_INSTALL_ALLOW_DEFERRAL" && maxDeferrals == 0 {
		return fmt.Errorf("in 'jamfpro_managed_software_update.%s': 'max_deferrals' must be set when 'update_action' is 'DOWNLOAD_INSTALL_ALLOW_DEFERRAL'", resourceName)
	}

	if (versionType == "SPECIFIC_VERSION" || versionType == "CUSTOM_VERSION") && specificVersion == "" {
		return fmt.Errorf("in 'jamfpro_managed_software_update.%s': 'specific_version' must be set when 'version_type' is 'SPECIFIC_VERSION' or 'CUSTOM_VERSION'", resourceName)
	}

	return nil
}
