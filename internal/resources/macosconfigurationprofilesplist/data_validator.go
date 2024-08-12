// macosconfigurationprofilesplist_data_validator.go
package macosconfigurationprofilesplist

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/datavalidators"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	if diff.Get("payload_validate").(bool) {
		if err := normalizePayloadState(ctx, diff, i); err != nil {
			return err
		}

		if err := validateDistributionMethod(ctx, diff, i); err != nil {
			return err
		}

		if err := validateMacOSConfigurationProfileLevel(ctx, diff, i); err != nil {
			return err
		}

		if err := validateConfigurationProfileFormatting(ctx, diff, i); err != nil {
			return err
		}

	}

	return nil
}

func normalizePayloadState(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	diff.SetNew("payloads", plist.NormalizePayloadState(diff.Get("payloads").(string)))
	return nil
}

// validateDistributionMethod checks that the 'self_service' block is only used when 'distribution_method' is "Make Available in Self Service".
func validateDistributionMethod(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	distributionMethod, ok := diff.GetOk("distribution_method")

	if !ok {
		return nil
	}

	selfServiceBlockExists := len(diff.Get("self_service").([]interface{})) > 0

	if distributionMethod == "Make Available in Self Service" && !selfServiceBlockExists {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': 'self_service' block is required when 'distribution_method' is set to 'Make Available in Self Service'", resourceName)
	}

	if distributionMethod != "Make Available in Self Service" && selfServiceBlockExists {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': 'self_service' block is not allowed when 'distribution_method' is set to '%s'", resourceName, distributionMethod)
	}

	return nil
}

// validateMacOSConfigurationProfileLevel validates that the 'PayloadScope' key in the payload matches the 'level' attribute.
func validateMacOSConfigurationProfileLevel(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	level := diff.Get("level").(string)
	payloads := diff.Get("payloads").(string)

	plistData, err := plist.DecodePlist([]byte(payloads))
	if err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': error decoding plist data: %v", resourceName, err)
	}

	payloadScope, err := datavalidators.GetPayloadScope(plistData)
	if err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': error getting 'PayloadScope' from plist: %v", resourceName, err)
	}

	if payloadScope != level {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': hcl 'level' attribute (%s) does not match the 'PayloadScope' in the root dict of the plist (%s); the values must be identical", resourceName, level, payloadScope)
	}

	return nil
}

// validateConfigurationProfileFormatting validates the indentation of the plist XML.
func validateConfigurationProfileFormatting(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	payloads := diff.Get("payloads").(string)

	if err := datavalidators.CheckPlistIndentationAndWhiteSpace(payloads); err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': %v", resourceName, err)
	}

	return nil
}
