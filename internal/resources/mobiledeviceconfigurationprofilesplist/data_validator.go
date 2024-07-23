// mobiledeviceconfigurationprofilesplist_data_validator.go
package mobiledeviceconfigurationprofilesplist

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/datavalidators"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	if err := validateMobileDeviceConfigurationProfileLevel(ctx, diff, i); err != nil {
		return err
	}

	if err := validateConfigurationProfileFormatting(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateMobileDeviceConfigurationProfileLevel validates that the 'PayloadScope' key in the payload matches the 'level' attribute.
func validateMobileDeviceConfigurationProfileLevel(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	level := diff.Get("level").(string)
	payloads := diff.Get("payloads").(string)

	plistData, err := plist.DecodePlist([]byte(payloads))
	if err != nil {
		return fmt.Errorf("in 'jamfpro_mobile_device_configuration_profile.%s': error decoding plist data: %v", resourceName, err)
	}

	payloadScope, err := datavalidators.GetPayloadScope(plistData)
	if err != nil {
		return fmt.Errorf("in 'jamfpro_mobile_device_configuration_profile.%s': error getting 'PayloadScope' from plist: %v", resourceName, err)
	}

	expectedScope := ""
	switch level {
	case "Device Level":
		expectedScope = "System"
	case "User Level":
		expectedScope = "User"
	default:
		return fmt.Errorf("in 'jamfpro_mobile_device_configuration_profile.%s': invalid 'level' attribute (%s)", resourceName, level)
	}

	if payloadScope != expectedScope {
		return fmt.Errorf("in 'jamfpro_mobile_device_configuration_profile.%s': .hcl 'level' attribute (%s) does not match the 'PayloadScope' in the plist (%s). When .hcl 'level' attribute is 'Device Level', the payload value must be 'System'. When .hcl 'level' attribute is 'User Level', the payload value must be 'User'.", resourceName, level, payloadScope)
	}

	return nil
}

// validateConfigurationProfileFormatting validates the indentation of the plist XML.
func validateConfigurationProfileFormatting(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	payloads := diff.Get("payloads").(string)

	if err := datavalidators.CheckPlistIndentationAndWhiteSpace(payloads); err != nil {
		return fmt.Errorf("in 'jamfpro_mobile_device_configuration_profile.%s': %v", resourceName, err)
	}

	return nil
}
