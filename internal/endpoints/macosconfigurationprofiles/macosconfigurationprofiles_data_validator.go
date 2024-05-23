// macosconfigurationprofiles_data_validator.go
package macosconfigurationprofiles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles/datavalidators"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	// Validate based on 'distribution_method' attribute.
	if err := validateDistributionMethod(ctx, diff, i); err != nil {
		return err
	}

	// Validate configuration profile level
	if err := validateConfigurationProfileLevel(ctx, diff, i); err != nil {
		return err
	}

	// Validate configuration profile indentation
	if err := validateConfigurationProfileFormatting(ctx, diff, i); err != nil {
		return err
	}

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

// validateConfigurationProfileLevel validates that the 'PayloadScope' key in the payload matches the 'level' attribute.
func validateConfigurationProfileLevel(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	level := diff.Get("level").(string)
	payloads := diff.Get("payloads").(string)

	// Decode the plist payload
	plistData, err := plist.DecodePlist([]byte(payloads))
	if err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': error decoding plist data: %v", resourceName, err)
	}

	// Check the PayloadScope in the plist
	payloadScope, err := datavalidators.GetPayloadScope(plistData)
	if err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': error getting 'PayloadScope' from plist: %v", resourceName, err)
	}

	if payloadScope != level {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': 'level' attribute (%s) does not match the 'PayloadScope' in the plist (%s)", resourceName, level, payloadScope)
	}

	return nil
}

// validateConfigurationProfileFormatting validates the indentation of the plist XML.
func validateConfigurationProfileFormatting(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	payloads := diff.Get("payloads").(string)

	// Check if the XML is well-formed and properly indented
	if err := datavalidators.CheckPlistIndentationAndWhiteSpace(payloads); err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': %v", resourceName, err)
	}

	return nil
}
