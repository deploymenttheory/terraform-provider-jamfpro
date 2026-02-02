package macos_configuration_profile_plist

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/plist"
	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations for macOS config profiles.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i any) error {
	if diff.Get("payload_validate").(bool) {
		if err := validatePayloadIdentifers(ctx, diff, i); err != nil {
			return err
		}

		if err := normalizePayloadState(ctx, diff, i); err != nil {
			return err
		}

		if err := validatePlistPayloadScope(ctx, diff, i); err != nil {
			return err
		}
	}

	if err := validateDistributionMethod(ctx, diff, i); err != nil {
		return err
	}

	if err := validateSelfServiceCategories(ctx, diff, i); err != nil {
		return err
	}

	if err := validateAllComputersScope(ctx, diff, i); err != nil {
		return err
	}

	if err := validateAllUsersScope(ctx, diff, i); err != nil {
		return err
	}

	if err := sharedschemas.ValidateScopeDirectoryServiceUserGroupNames(ctx, diff, i); err != nil {
		return fmt.Errorf("validating scope directory service user/group names: %w", err)
	}

	return nil
}

func normalizePayloadState(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	diff.SetNew("payloads", plist.NormalizePayloadState(diff.Get("payloads").(string)))
	return nil
}

// validatePayloadIdentifers performs the payload validation that was previously in the ValidateFunc.
func validatePayloadIdentifers(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	payload := diff.Get("payloads").(string)

	profile, err := plist.UnmarshalPayload(payload)
	if err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile_plist.%s': error unmarshalling payload: %v", resourceName, err)
	}

	if profile.PayloadIdentifier != profile.PayloadUUID {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile_plist.%s': root-level PayloadIdentifier and PayloadUUID within the plist do not match. Expected PayloadIdentifier to be '%s', but got '%s'", resourceName, profile.PayloadUUID, profile.PayloadIdentifier)
	}

	errs := plist.ValidatePayloadFields(profile)
	if len(errs) > 0 {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile_plist.%s': %v", resourceName, errs)
	}

	return nil
}

// validateDistributionMethod checks that the 'self_service' block is only used when 'distribution_method' is "Make Available in Self Service".
func validateDistributionMethod(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	distributionMethod, ok := diff.GetOk("distribution_method")

	if !ok {
		return nil
	}

	selfServiceBlockExists := len(diff.Get("self_service").([]any)) > 0

	if distributionMethod == "Make Available in Self Service" && !selfServiceBlockExists {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': 'self_service' block is required when 'distribution_method' is set to 'Make Available in Self Service'", resourceName)
	}

	if distributionMethod != "Make Available in Self Service" && selfServiceBlockExists {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': 'self_service' block is not allowed when 'distribution_method' is set to '%s'", resourceName, distributionMethod)
	}

	return nil
}

// validatePlistPayloadScope validates that the 'PayloadScope' key in the payload matches the 'level' attribute in the HCL.
func validatePlistPayloadScope(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	level := diff.Get("level").(string)
	payloads := diff.Get("payloads").(string)

	plistData, err := plist.DecodePlist([]byte(payloads))
	if err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': error decoding plist data: %v", resourceName, err)
	}

	payloadScope, err := plist.GetPayloadScope(plistData)
	if err != nil {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': error getting 'PayloadScope' from plist: %v", resourceName, err)
	}

	if payloadScope != level {
		return fmt.Errorf("in 'jamfpro_macos_configuration_profile.%s': the hcl 'level' attribute (%s) does not match the 'PayloadScope' in the root dict of the plist (%s); the values must be identical", resourceName, level, payloadScope)
	}

	return nil
}

// validateSelfServiceCategories validates the 'self_service_category' block.
func validateSelfServiceCategories(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	selfServiceRaw, ok := diff.GetOk("self_service")
	if !ok {
		return nil
	}

	selfService := selfServiceRaw.([]any)[0].(map[string]any)
	categories, ok := selfService["self_service_category"].([]any)
	if !ok {
		return nil
	}

	for i, catRaw := range categories {
		cat := catRaw.(map[string]any)
		displayIn, displayOk := cat["display_in"].(bool)
		featureIn, featureOk := cat["feature_in"].(bool)

		if displayOk && featureOk && featureIn && !displayIn {
			return fmt.Errorf("in 'jamfpro_macos_configuration_profile_plist.%s': self_service_category[%d]: feature_in can only be true if display_in is also true", resourceName, i)
		}
	}

	return nil
}

func validateAllComputersScope(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	scopeRaw, ok := diff.GetOk("scope")
	if !ok {
		return nil
	}

	scope := scopeRaw.([]any)[0].(map[string]any)
	allComputers := scope["all_computers"].(bool)

	if allComputers {
		fieldsToCheck := []string{"computer_ids", "computer_group_ids"}
		for _, field := range fieldsToCheck {
			if value, exists := scope[field]; exists {
				if setVal, ok := value.(*schema.Set); ok && setVal.Len() > 0 {
					return fmt.Errorf("in 'jamfpro_macos_configuration_profile_plist.%s': when 'all_computers' scope is set to true, '%s' should not be set", resourceName, field)
				}
			}
		}
	}

	return nil
}

func validateAllUsersScope(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	scopeRaw, ok := diff.GetOk("scope")
	if !ok {
		return nil
	}

	scope := scopeRaw.([]any)[0].(map[string]any)
	allJssUsers := scope["all_jss_users"].(bool)

	if allJssUsers {
		fieldsToCheck := []string{"jss_user_ids", "jss_user_group_ids", "building_ids", "department_ids"}
		for _, field := range fieldsToCheck {
			if value, exists := scope[field]; exists {
				if setVal, ok := value.(*schema.Set); ok && setVal.Len() > 0 {
					return fmt.Errorf("in 'jamfpro_macos_configuration_profile_plist.%s': when 'all_jss_users' scope is set to true, '%s' should not be set", resourceName, field)
				}
			}
		}
	}

	return nil
}
