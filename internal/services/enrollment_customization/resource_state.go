// enrollment_customization_state.go
package enrollment_customization

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Enrollment Customization information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceEnrollmentCustomization) diag.Diagnostics {
	var diags diag.Diagnostics

	// Set main attributes
	resourceData := map[string]any{
		"id":           resp.ID,
		"site_id":      resp.SiteID,
		"display_name": resp.DisplayName,
		"description":  resp.Description,
	}

	// Set branding settings as a list with one item
	brandingSettings := []any{
		map[string]any{
			"text_color":        resp.BrandingSettings.TextColor,
			"button_color":      resp.BrandingSettings.ButtonColor,
			"button_text_color": resp.BrandingSettings.ButtonTextColor,
			"background_color":  resp.BrandingSettings.BackgroundColor,
			"icon_url":          resp.BrandingSettings.IconUrl,
		},
	}

	// Set all values in state
	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set branding settings separately
	if err := d.Set("branding_settings", brandingSettings); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// stateTextPrestagePane converts API response to Terraform state format
func stateTextPrestagePane(textPane *jamfpro.ResourceEnrollmentCustomizationTextPane) map[string]any {
	return map[string]any{
		"id":                   textPane.ID,
		"display_name":         textPane.DisplayName,
		"rank":                 textPane.Rank,
		"title":                textPane.Title,
		"body":                 textPane.Body,
		"subtext":              textPane.Subtext,
		"back_button_text":     textPane.BackButtonText,
		"continue_button_text": textPane.ContinueButtonText,
	}
}

// stateLDAPPrestagePane converts API response to Terraform state format
func stateLDAPPrestagePane(ldapPane *jamfpro.ResourceEnrollmentCustomizationLDAPPane) map[string]any {
	ldapGroupAccess := make([]map[string]any, 0)
	for _, group := range ldapPane.LDAPGroupAccess {
		groupMap := map[string]any{
			"group_name":     group.GroupName,
			"ldap_server_id": group.LDAPServerID,
		}
		ldapGroupAccess = append(ldapGroupAccess, groupMap)
	}

	return map[string]any{
		"id":                   ldapPane.ID,
		"display_name":         ldapPane.DisplayName,
		"rank":                 ldapPane.Rank,
		"title":                ldapPane.Title,
		"username_label":       ldapPane.UsernameLabel,
		"password_label":       ldapPane.PasswordLabel,
		"back_button_text":     ldapPane.BackButtonText,
		"continue_button_text": ldapPane.ContinueButtonText,
		"ldap_group_access":    ldapGroupAccess,
	}
}

// stateSSOPrestagePane converts API response to Terraform state format
func stateSSOPrestagePane(ssoPane *jamfpro.ResourceEnrollmentCustomizationSSOPane) map[string]any {
	return map[string]any{
		"id":                                 ssoPane.ID,
		"display_name":                       ssoPane.DisplayName,
		"rank":                               ssoPane.Rank,
		"is_group_enrollment_access_enabled": ssoPane.IsGroupEnrollmentAccessEnabled,
		"group_enrollment_access_name":       ssoPane.GroupEnrollmentAccessName,
		"is_use_jamf_connect":                ssoPane.IsUseJamfConnect,
		"short_name_attribute":               ssoPane.ShortNameAttribute,
		"long_name_attribute":                ssoPane.LongNameAttribute,
	}
}
