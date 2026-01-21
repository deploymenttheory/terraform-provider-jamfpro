package mobile_device_extension_attribute

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	inputTypeText      = "TEXT"
	inputTypePopup     = "POPUP"
	inputTypeDirectory = "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING"
)

var (
	errPopupChoicesOnlyWithPopupInput  = errors.New("'popup_menu_choices' can only be set when 'input_type' is 'POPUP'")
	errPopupChoicesRequiredForPopup    = errors.New("'popup_menu_choices' must be set when 'input_type' is 'POPUP'")
	errLDAPMappingOnlyWithDirectory    = errors.New("'ldap_attribute_mapping' can only be set when 'input_type' is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
	errLDAPMappingRequiredForDirectory = errors.New("'ldap_attribute_mapping' must be set when 'input_type' is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
	errLDAPAllowedOnlyWithDirectory    = errors.New("'ldap_extension_attribute_allowed' can only be true when 'input_type' is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
)

// mainCustomDiffFunc orchestrates all custom diff validations for mobile device extension attributes.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
	if err := validateInputTypeSpecificAttributes(ctx, diff, meta); err != nil {
		return err
	}

	if err := validateLDAPExtensionAttributeAllowed(ctx, diff, meta); err != nil {
		return err
	}

	return nil
}

// validateInputTypeSpecificAttributes ensures that attributes specific to certain input types are only set when appropriate.
func validateInputTypeSpecificAttributes(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	inputType := diff.Get("input_type").(string)

	choicesLen := 0
	if choices, ok := diff.GetOk("popup_menu_choices"); ok {
		choicesLen = choices.(*schema.Set).Len()
	}

	if choicesLen > 0 && inputType != inputTypePopup {
		return fmt.Errorf("in 'jamfpro_mobile_device_extension_attribute.%s': %w", resourceName, errPopupChoicesOnlyWithPopupInput)
	}

	if inputType == inputTypePopup && choicesLen == 0 {
		return fmt.Errorf("in 'jamfpro_mobile_device_extension_attribute.%s': %w", resourceName, errPopupChoicesRequiredForPopup)
	}

	ldapMapping := ""
	if mapping, ok := diff.GetOk("ldap_attribute_mapping"); ok {
		ldapMapping = mapping.(string)
	}

	if ldapMapping != "" && inputType != inputTypeDirectory {
		return fmt.Errorf("in 'jamfpro_mobile_device_extension_attribute.%s': %w", resourceName, errLDAPMappingOnlyWithDirectory)
	}

	if inputType == inputTypeDirectory && ldapMapping == "" {
		return fmt.Errorf("in 'jamfpro_mobile_device_extension_attribute.%s': %w", resourceName, errLDAPMappingRequiredForDirectory)
	}

	return nil
}

// validateLDAPExtensionAttributeAllowed ensures that ldap_extension_attribute_allowed is only true when input_type is DIRECTORY_SERVICE_ATTRIBUTE_MAPPING.
func validateLDAPExtensionAttributeAllowed(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	inputType := diff.Get("input_type").(string)
	ldapAllowed := diff.Get("ldap_extension_attribute_allowed").(bool)

	if !ldapAllowed {
		return nil
	}

	if inputType != inputTypeDirectory {
		return fmt.Errorf("in 'jamfpro_mobile_device_extension_attribute.%s': %w", resourceName, errLDAPAllowedOnlyWithDirectory)
	}

	return nil
}
