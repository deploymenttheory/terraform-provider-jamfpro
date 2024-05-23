// mobiledeviceconfigurationprofiles_data_validator.go
package mobiledeviceconfigurationprofiles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	// Validate configuration profile level
	if err := validateConfigurationProfileLevel(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateConfigurationProfileLevel validates that the 'PayloadScope' key in the payload matches the 'level' attribute.
func validateConfigurationProfileLevel(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string) // Assuming 'name' is always set and is unique
	level := diff.Get("level").(string)
	payloads := diff.Get("payloads").(string)

	// Decode the plist payload
	plistData, err := configurationprofiles.DecodePlist([]byte(payloads))
	if err != nil {
		return fmt.Errorf("in 'jamfpro_mobile_device_configuration_profile.%s': error decoding plist data: %v", resourceName, err)
	}

	// Check the PayloadScope in the plist
	payloadScope, err := getPayloadScope(plistData)
	if err != nil {
		return fmt.Errorf("in 'jamfpro_mobile_device_configuration_profile.%s': error getting 'PayloadScope' from plist: %v", resourceName, err)
	}

	if payloadScope != level {
		return fmt.Errorf("in 'jamfpro_mobile_device_configuration_profile.%s': 'level' attribute (%s) does not match the 'PayloadScope' in the plist (%s)", resourceName, level, payloadScope)
	}

	return nil
}

// getPayloadScope retrieves the 'PayloadScope' key from the decoded plist data.
func getPayloadScope(plistData map[string]interface{}) (string, error) {
	if scope, ok := plistData["PayloadScope"].(string); ok {
		return scope, nil
	}
	return "", fmt.Errorf("'PayloadScope' key not found in plist")
}
