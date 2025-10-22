// sso_settings_state.go
package sso_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the state of the SSO settings resource in Terraform
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceSsoSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := map[string]any{
		"sso_enabled":                                          resp.SsoEnabled,
		"configuration_type":                                   resp.ConfigurationType,
		"sso_bypass_allowed":                                   resp.SsoBypassAllowed,
		"sso_for_enrollment_enabled":                           resp.SsoForEnrollmentEnabled,
		"sso_for_macos_self_service_enabled":                   resp.SsoForMacOsSelfServiceEnabled,
		"enrollment_sso_for_account_driven_enrollment_enabled": resp.EnrollmentSsoForAccountDrivenEnrollmentEnabled,
		"group_enrollment_access_enabled":                      resp.GroupEnrollmentAccessEnabled,
		"group_enrollment_access_name":                         resp.GroupEnrollmentAccessName,
	}

	for key, val := range settings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	oidcSettings := []any{
		map[string]any{
			"user_mapping":                     resp.OidcSettings.UserMapping,
			"jamf_id_authentication_enabled":   resp.OidcSettings.JamfIdAuthenticationEnabled,
			"username_attribute_claim_mapping": resp.OidcSettings.UsernameAttributeClaimMapping,
		},
	}
	if err := d.Set("oidc_settings", oidcSettings); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	samlSettings := []any{
		map[string]any{
			"idp_url":                   resp.SamlSettings.IdpUrl,
			"entity_id":                 resp.SamlSettings.EntityId,
			"metadata_source":           resp.SamlSettings.MetadataSource,
			"user_mapping":              resp.SamlSettings.UserMapping,
			"idp_provider_type":         resp.SamlSettings.IdpProviderType,
			"group_rdn_key":             resp.SamlSettings.GroupRdnKey,
			"user_attribute_name":       resp.SamlSettings.UserAttributeName,
			"group_attribute_name":      resp.SamlSettings.GroupAttributeName,
			"user_attribute_enabled":    resp.SamlSettings.UserAttributeEnabled,
			"metadata_file_name":        resp.SamlSettings.MetadataFileName,
			"other_provider_type_name":  resp.SamlSettings.OtherProviderTypeName,
			"federation_metadata_file":  resp.SamlSettings.FederationMetadataFile,
			"token_expiration_disabled": resp.SamlSettings.TokenExpirationDisabled,
			"session_timeout":           resp.SamlSettings.SessionTimeout,
		},
	}
	if err := d.Set("saml_settings", samlSettings); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if resp.EnrollmentSsoConfig.Hosts != nil || resp.EnrollmentSsoConfig.ManagementHint != "" {
		enrollmentConfig := []any{
			map[string]any{
				"hosts":           resp.EnrollmentSsoConfig.Hosts,
				"management_hint": resp.EnrollmentSsoConfig.ManagementHint,
			},
		}
		if err := d.Set("enrollment_sso_config", enrollmentConfig); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
