package computer_extension_attribute

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	inputTypeScript  = "SCRIPT"
	inputTypePopup   = "POPUP"
	inputTypeLdapMap = "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING"
)

var (
	errScriptContentsOnlyWithScriptInput   = errors.New("'script_contents' can only be set when 'input_type' is 'SCRIPT'")
	errScriptContentsRequiredForScriptDiff = errors.New("'script_contents' must be set when 'input_type' is 'SCRIPT'")
	errPopupChoicesOnlyWithPopupInput      = errors.New("'popup_menu_choices' can only be set when 'input_type' is 'POPUP'")
	errPopupChoicesRequiredForPopupDiff    = errors.New("'popup_menu_choices' must be set when 'input_type' is 'POPUP'")
	errLDAPMappingOnlyWithDirectoryInput   = errors.New("'ldap_attribute_mapping' can only be set when 'input_type' is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
	errLDAPMappingRequiredForDirectoryDiff = errors.New("'ldap_attribute_mapping' must be set when 'input_type' is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
	errLDAPAllowedOnlyWithDirectoryInput   = errors.New("'ldap_extension_attribute_allowed' can only be true when 'input_type' is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'")
	errLDAPAllowedRequiresMapping          = errors.New("'ldap_extension_attribute_allowed' requires 'ldap_attribute_mapping' to be set")
)

// mainCustomDiffFunc orchestrates all custom diff validations for computer extension attributes.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
	if err := validateInputTypeSpecificAttributes(ctx, diff, meta); err != nil {
		return err
	}

	if err := validateLDAPExtensionAttributeAllowed(ctx, diff, meta); err != nil {
		return err
	}

	return nil
}

func validateInputTypeSpecificAttributes(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	inputType := diff.Get("input_type").(string)

	scriptContents := ""
	if script, ok := diff.GetOk("script_contents"); ok {
		scriptContents = script.(string)
	}
	if scriptContents != "" && inputType != inputTypeScript {
		return fmt.Errorf("in 'jamfpro_computer_extension_attribute.%s': %w", resourceName, errScriptContentsOnlyWithScriptInput)
	}
	if inputType == inputTypeScript && scriptContents == "" {
		return fmt.Errorf("in 'jamfpro_computer_extension_attribute.%s': %w", resourceName, errScriptContentsRequiredForScriptDiff)
	}

	choicesLen := 0
	if choices, ok := diff.GetOk("popup_menu_choices"); ok {
		choicesLen = len(choices.([]any))
	}
	if choicesLen > 0 && inputType != inputTypePopup {
		return fmt.Errorf("in 'jamfpro_computer_extension_attribute.%s': %w", resourceName, errPopupChoicesOnlyWithPopupInput)
	}
	if inputType == inputTypePopup && choicesLen == 0 {
		return fmt.Errorf("in 'jamfpro_computer_extension_attribute.%s': %w", resourceName, errPopupChoicesRequiredForPopupDiff)
	}

	ldapMapping := ""
	if mapping, ok := diff.GetOk("ldap_attribute_mapping"); ok {
		ldapMapping = mapping.(string)
	}
	if ldapMapping != "" && inputType != inputTypeLdapMap {
		return fmt.Errorf("in 'jamfpro_computer_extension_attribute.%s': %w", resourceName, errLDAPMappingOnlyWithDirectoryInput)
	}

	if inputType == inputTypeLdapMap && ldapMapping == "" {
		return fmt.Errorf("in 'jamfpro_computer_extension_attribute.%s': %w", resourceName, errLDAPMappingRequiredForDirectoryDiff)
	}

	return nil
}

func validateLDAPExtensionAttributeAllowed(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	inputType := diff.Get("input_type").(string)
	ldapAllowed := diff.Get("ldap_extension_attribute_allowed").(bool)

	if !ldapAllowed {
		return nil
	}

	if inputType != inputTypeLdapMap {
		return fmt.Errorf("in 'jamfpro_computer_extension_attribute.%s': %w", resourceName, errLDAPAllowedOnlyWithDirectoryInput)
	}

	if mapping, ok := diff.GetOk("ldap_attribute_mapping"); !ok || mapping.(string) == "" {
		return fmt.Errorf("in 'jamfpro_computer_extension_attribute.%s': %w", resourceName, errLDAPAllowedRequiresMapping)
	}

	return nil
}
