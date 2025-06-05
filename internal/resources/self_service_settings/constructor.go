package self_service_settings

// selfservicesettings_resource.go
// Package selfservicesettings provides the schema and CRUD operations for managing Jamf Pro Self Service Settings in Terraform.
// This package includes functions to create, read, update, and delete the Self Service Settings configuration.
// It also includes a function to construct the ResourceSelfServiceSettings object from the schema data.

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructSelfServiceSettings constructs a ResourceSelfServiceSettings object from the provided schema data
func constructSelfServiceSettings(d *schema.ResourceData) (*jamfpro.ResourceSelfServiceSettings, error) {
	installSettingsList := d.Get("install_settings").([]interface{})
	var installSettings jamfpro.InstallSettings
	if len(installSettingsList) > 0 {
		installSettingsMap := installSettingsList[0].(map[string]interface{})
		installSettings = jamfpro.InstallSettings{
			InstallAutomatically: installSettingsMap["install_automatically"].(bool),
			InstallLocation:      installSettingsMap["install_location"].(string),
		}
	}

	loginSettingsList := d.Get("login_settings").([]interface{})
	var loginSettings jamfpro.LoginSettings
	if len(loginSettingsList) > 0 {
		loginSettingsMap := loginSettingsList[0].(map[string]interface{})
		loginSettings = jamfpro.LoginSettings{
			UserLoginLevel:  loginSettingsMap["user_login_level"].(string),
			AllowRememberMe: loginSettingsMap["allow_remember_me"].(bool),
			UseFido2:        loginSettingsMap["use_fido2"].(bool),
			AuthType:        loginSettingsMap["auth_type"].(string),
		}
	}

	configSettingsList := d.Get("configuration_settings").([]interface{})
	var configSettings jamfpro.ConfigurationSettings
	if len(configSettingsList) > 0 {
		configSettingsMap := configSettingsList[0].(map[string]interface{})
		configSettings = jamfpro.ConfigurationSettings{
			NotificationsEnabled:  configSettingsMap["notifications_enabled"].(bool),
			AlertUserApprovedMdm:  configSettingsMap["alert_user_approved_mdm"].(bool),
			DefaultLandingPage:    configSettingsMap["default_landing_page"].(string),
			DefaultHomeCategoryId: configSettingsMap["default_home_category_id"].(int),
			BookmarksName:         configSettingsMap["bookmarks_name"].(string),
		}
	}

	resource := &jamfpro.ResourceSelfServiceSettings{
		InstallSettings:       installSettings,
		LoginSettings:         loginSettings,
		ConfigurationSettings: configSettings,
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Self Service Settings to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Self Service Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
