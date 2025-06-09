package app_installer

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs an AppCatalogDeployment object from the provided schema data.
func construct(d *schema.ResourceData, client *jamfpro.Client) (*jamfpro.ResourceJamfAppCatalogDeployment, error) {
	appTitleName := d.Get("app_title_name").(string)

	// Fetch app titles dynamically from the API and resolve the app_title_id
	titlesResp, err := client.GetJamfAppCatalogAppInstallerTitles(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch app titles from Jamf Pro: %v", err)
	}

	var appTitleID string
	var validTitles []string
	for _, result := range titlesResp.Results {
		validTitles = append(validTitles, result.TitleName)
		if result.TitleName == appTitleName {
			appTitleID = result.ID
		}
	}
	if appTitleID == "" {
		return nil, fmt.Errorf("in 'jamfpro_app_installer.%s': 'app_title_name' must be one of the following values: %s", appTitleName, strings.Join(validTitles, ", "))
	}

	resource := &jamfpro.ResourceJamfAppCatalogDeployment{
		Name:                            d.Get("name").(string),
		Enabled:                         jamfpro.BoolPtr(d.Get("enabled").(bool)),
		AppTitleId:                      appTitleID,
		DeploymentType:                  d.Get("deployment_type").(string),
		UpdateBehavior:                  d.Get("update_behavior").(string),
		CategoryId:                      d.Get("category_id").(string),
		SiteId:                          d.Get("site_id").(string),
		SmartGroupId:                    d.Get("smart_group_id").(string),
		InstallPredefinedConfigProfiles: jamfpro.BoolPtr(d.Get("install_predefined_config_profiles").(bool)),
		TriggerAdminNotifications:       jamfpro.BoolPtr(d.Get("trigger_admin_notifications").(bool)),
		SelectedVersion:                 d.Get("selected_version").(string),
	}

	// Construct notification settings
	if v, ok := d.GetOk("notification_settings"); ok && len(v.([]interface{})) > 0 {
		ns := v.([]interface{})[0].(map[string]interface{})
		resource.NotificationSettings = jamfpro.JamfAppCatalogDeploymentSubsetNotificationSettings{
			NotificationMessage:  ns["notification_message"].(string),
			NotificationInterval: ns["notification_interval"].(int),
			DeadlineMessage:      ns["deadline_message"].(string),
			Deadline:             ns["deadline"].(int),
			QuitDelay:            ns["quit_delay"].(int),
			CompleteMessage:      ns["complete_message"].(string),
			Relaunch:             jamfpro.BoolPtr(ns["relaunch"].(bool)),
			Suppress:             jamfpro.BoolPtr(ns["suppress"].(bool)),
		}
	}

	// Construct self-service settings
	if v, ok := d.GetOk("self_service_settings"); ok && len(v.([]interface{})) > 0 {
		ss := v.([]interface{})[0].(map[string]interface{})
		resource.SelfServiceSettings = jamfpro.JamfAppCatalogDeploymentSubsetSelfServiceSettings{
			IncludeInFeaturedCategory:   jamfpro.BoolPtr(ss["include_in_featured_category"].(bool)),
			IncludeInComplianceCategory: jamfpro.BoolPtr(ss["include_in_compliance_category"].(bool)),
			ForceViewDescription:        jamfpro.BoolPtr(ss["force_view_description"].(bool)),
			Description:                 ss["description"].(string),
		}

		// Construct categories
		if categories, ok := ss["categories"].(*schema.Set); ok {
			for _, cat := range categories.List() {
				category := cat.(map[string]interface{})
				resource.SelfServiceSettings.Categories = append(resource.SelfServiceSettings.Categories, jamfpro.JamfAppCatalogDeploymentSubsetCategory{
					ID:       category["id"].(string),
					Featured: jamfpro.BoolPtr(category["featured"].(bool)),
				})
			}
		}
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro App Installer Deployment '%s' to JSON: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro App Installer Deployment JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
