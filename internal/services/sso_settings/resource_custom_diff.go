package sso_settings

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	errMetadataFileNameSetForURL             = errors.New("in 'jamfpro_sso_settings': metadata_file_name should not be set when metadata_source is URL")
	errFederationMetadataFileSetForURL       = errors.New("in 'jamfpro_sso_settings': federation_metadata_file should not be set when metadata_source is URL")
	errIdpURLRequiredForURL                  = errors.New("in 'jamfpro_sso_settings': idp_url is required when metadata_source is URL")
	errMetadataFileNameRequiredForFile       = errors.New("in 'jamfpro_sso_settings': metadata_file_name is required when metadata_source is FILE")
	errFederationMetadataFileRequiredForFile = errors.New("in 'jamfpro_sso_settings': federation_metadata_file is required when metadata_source is FILE")
	errIdpURLSetForFile                      = errors.New("in 'jamfpro_sso_settings': idp_url should not be set when metadata_source is FILE")
	errGroupEnrollmentNameRequired           = errors.New("in 'jamfpro_sso_settings': group_enrollment_access_name is required when group_enrollment_access_enabled is true")
)

// mainCustomDiffFunc orchestrates all custom diff validations for SSO settings
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i any) error {
	if err := validateSamlMetadataSettings(ctx, diff, i); err != nil {
		return err
	}

	if err := validateGroupEnrollmentSettings(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateSamlMetadataSettings ensures SAML metadata settings are properly configured
func validateSamlMetadataSettings(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	samlSettings := diff.Get("saml_settings").([]any)
	if len(samlSettings) == 0 {
		return nil
	}

	samlConfig := samlSettings[0].(map[string]any)
	metadataSource := samlConfig["metadata_source"].(string)
	resourceName, ok := diff.Get("name").(string)
	if !ok {
		resourceName = "unknown"
	}

	switch metadataSource {
	case "URL":
		if samlConfig["metadata_file_name"].(string) != "" {
			return fmt.Errorf("%w: %s", errMetadataFileNameSetForURL, resourceName)
		}
		if samlConfig["federation_metadata_file"].(string) != "" {
			return fmt.Errorf("%w: %s", errFederationMetadataFileSetForURL, resourceName)
		}
		if samlConfig["idp_url"].(string) == "" {
			return fmt.Errorf("%w: %s", errIdpURLRequiredForURL, resourceName)
		}
	case "FILE":
		if samlConfig["metadata_file_name"].(string) == "" {
			return fmt.Errorf("%w: %s", errMetadataFileNameRequiredForFile, resourceName)
		}
		if samlConfig["federation_metadata_file"].(string) == "" {
			return fmt.Errorf("%w: %s", errFederationMetadataFileRequiredForFile, resourceName)
		}
		if samlConfig["idp_url"].(string) != "" {
			return fmt.Errorf("%w: %s", errIdpURLSetForFile, resourceName)
		}
	}

	return nil
}

// validateGroupEnrollmentSettings ensures group enrollment settings are properly configured
func validateGroupEnrollmentSettings(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	groupEnabled := diff.Get("group_enrollment_access_enabled").(bool)
	groupName := diff.Get("group_enrollment_access_name").(string)
	resourceName, ok := diff.Get("name").(string)
	if !ok {
		resourceName = "unknown"
	}

	if groupEnabled && groupName == "" {
		return fmt.Errorf("%w: %s", errGroupEnrollmentNameRequired, resourceName)
	}

	return nil
}

// getConfigurationTypes returns a list of supported configuration types
func getConfigurationTypes() []string {
	return []string{
		"SAML",
		"OIDC",
		"OIDC_WITH_SAML",
	}
}

// getIdpProviderTypes returns a list of supported Identity Provider types
func getIdpProviderTypes() []string {
	return []string{
		"ADFS",
		"OKTA",
		"GOOGLE",
		"SHIBBOLETH",
		"ONELOGIN",
		"PING",
		"CENTRIFY",
		"AZURE",
		"OTHER",
	}
}

// getMetadataSourceTypes returns a list of supported metadata source types
func getMetadataSourceTypes() []string {
	return []string{
		"URL",
		"FILE",
		"UNKNOWN",
	}
}

// getUserMappingTypes returns a list of supported user mapping types
func getUserMappingTypes() []string {
	return []string{
		"USERNAME",
		"EMAIL",
	}
}
