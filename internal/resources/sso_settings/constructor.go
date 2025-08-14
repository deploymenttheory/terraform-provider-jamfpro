// sso_settings_constructor.go
package sso_settings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs the SSO settings resource from the Terraform schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceSsoSettings, error) {
	var oidcSettings jamfpro.OidcSettings
	if v, ok := d.GetOk("oidc_settings"); ok && len(v.([]interface{})) > 0 {
		oidcMap := v.([]interface{})[0].(map[string]interface{})
		oidcSettings = jamfpro.OidcSettings{
			UserMapping:                   oidcMap["user_mapping"].(string),
			JamfIdAuthenticationEnabled:   jamfpro.BoolPtr(oidcMap["jamf_id_authentication_enabled"].(bool)),
			UsernameAttributeClaimMapping: oidcMap["username_attribute_claim_mapping"].(string),
		}
	}

	var samlSettings jamfpro.SamlSettings
	if v, ok := d.GetOk("saml_settings"); ok && len(v.([]interface{})) > 0 {
		samlMap := v.([]interface{})[0].(map[string]interface{})
		samlSettings = jamfpro.SamlSettings{
			IdpUrl:                  samlMap["idp_url"].(string),
			EntityId:                samlMap["entity_id"].(string),
			MetadataSource:          samlMap["metadata_source"].(string),
			UserMapping:             samlMap["user_mapping"].(string),
			IdpProviderType:         samlMap["idp_provider_type"].(string),
			GroupRdnKey:             samlMap["group_rdn_key"].(string),
			UserAttributeName:       samlMap["user_attribute_name"].(string),
			GroupAttributeName:      samlMap["group_attribute_name"].(string),
			UserAttributeEnabled:    samlMap["user_attribute_enabled"].(bool),
			MetadataFileName:        samlMap["metadata_file_name"].(string),
			OtherProviderTypeName:   samlMap["other_provider_type_name"].(string),
			FederationMetadataFile:  samlMap["federation_metadata_file"].(string),
			TokenExpirationDisabled: samlMap["token_expiration_disabled"].(bool),
			SessionTimeout:          samlMap["session_timeout"].(int),
		}
	}

	var enrollmentSsoConfig jamfpro.EnrollmentSsoConfig
	if v, ok := d.GetOk("enrollment_sso_config"); ok && len(v.([]interface{})) > 0 {
		if enrollmentMap, ok := v.([]interface{})[0].(map[string]interface{}); ok && enrollmentMap != nil {
			hosts := make([]string, 0)
			if hostsRaw, ok := enrollmentMap["hosts"]; ok && hostsRaw != nil {
				if hostsInterface, ok := hostsRaw.([]interface{}); ok {
					hosts = expandStringList(hostsInterface)
				}
			}

			enrollmentSsoConfig = jamfpro.EnrollmentSsoConfig{
				Hosts:          hosts,
				ManagementHint: enrollmentMap["management_hint"].(string),
			}
		}
	}

	resource := &jamfpro.ResourceSsoSettings{
		SsoEnabled:                    d.Get("sso_enabled").(bool),
		ConfigurationType:             d.Get("configuration_type").(string),
		OidcSettings:                  &oidcSettings,
		SamlSettings:                  &samlSettings,
		SsoBypassAllowed:              d.Get("sso_bypass_allowed").(bool),
		SsoForEnrollmentEnabled:       d.Get("sso_for_enrollment_enabled").(bool),
		SsoForMacOsSelfServiceEnabled: d.Get("sso_for_macos_self_service_enabled").(bool),
		EnrollmentSsoForAccountDrivenEnrollmentEnabled: d.Get("enrollment_sso_for_account_driven_enrollment_enabled").(bool),
		GroupEnrollmentAccessEnabled:                   d.Get("group_enrollment_access_enabled").(bool),
		GroupEnrollmentAccessName:                      d.Get("group_enrollment_access_name").(string),
		EnrollmentSsoConfig:                            &enrollmentSsoConfig,
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro SSO Settings to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro SSO Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// expandStringList converts a list of interfaces to a list of strings.
func expandStringList(list []interface{}) []string {
	result := make([]string, 0, len(list))
	for _, item := range list {
		result = append(result, item.(string))
	}
	return result
}
