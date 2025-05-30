package selfservicesettings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructSelfServiceSettings constructs a ResourceSelfServiceSettings object from the provided schema data
func constructSelfServiceSettings(d *schema.ResourceData) (*jamfpro.ResourceSelfServiceSettings, error) {
	installSettings := jamfpro.InstallSettings{
		InstallAutomatically: d.Get("install_automatically").(bool),
		InstallLocation:      d.Get("install_location").(string),
	}

	loginSettings := jamfpro.LoginSettings{
		UserLoginLevel:  d.Get("user_login_level").(string),
		AllowRememberMe: d.Get("allow_remember_me").(bool),
		UseFido2:        d.Get("use_fido2").(bool),
		AuthType:        d.Get("auth_type").(string),
	}

	configSettings := jamfpro.ConfigurationSettings{
		NotificationsEnabled:  d.Get("notifications_enabled").(bool),
		AlertUserApprovedMdm:  d.Get("alert_user_approved_mdm").(bool),
		DefaultLandingPage:    d.Get("default_landing_page").(string),
		DefaultHomeCategoryId: d.Get("default_home_category_id").(int),
		BookmarksName:         d.Get("bookmarks_name").(string),
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
