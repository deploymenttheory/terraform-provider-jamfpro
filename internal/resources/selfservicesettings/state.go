package selfservicesettings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Self Service settings information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceSelfServiceSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	installSettings := map[string]interface{}{
		"install_automatically": resp.InstallSettings.InstallAutomatically,
		"install_location":      resp.InstallSettings.InstallLocation,
	}

	loginSettings := map[string]interface{}{
		"user_login_level":  resp.LoginSettings.UserLoginLevel,
		"allow_remember_me": resp.LoginSettings.AllowRememberMe,
		"use_fido2":         resp.LoginSettings.UseFido2,
		"auth_type":         resp.LoginSettings.AuthType,
	}

	configurationSettings := map[string]interface{}{
		"notifications_enabled":    resp.ConfigurationSettings.NotificationsEnabled,
		"alert_user_approved_mdm":  resp.ConfigurationSettings.AlertUserApprovedMdm,
		"default_landing_page":     resp.ConfigurationSettings.DefaultLandingPage,
		"default_home_category_id": resp.ConfigurationSettings.DefaultHomeCategoryId,
		"bookmarks_name":           resp.ConfigurationSettings.BookmarksName,
	}

	if err := d.Set("install_settings", []interface{}{installSettings}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("login_settings", []interface{}{loginSettings}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("configuration_settings", []interface{}{configurationSettings}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
