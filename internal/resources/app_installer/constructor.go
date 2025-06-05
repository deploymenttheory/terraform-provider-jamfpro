package app_installer

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//go:embed app_catalog_app_installer_titles.json
var appTitlesFS embed.FS

var appTitles struct {
	TotalCount int                                          `json:"totalCount"`
	Results    []jamfpro.ResourceJamfAppCatalogAppInstaller `json:"results"`
}

func init() {
	// Read the embedded JSON file
	data, err := appTitlesFS.ReadFile("app_catalog_app_installer_titles.json")
	if err != nil {
		log.Fatalf("Failed to read app_catalog_app_installer_titles.json: %v", err)
	}

	if err := json.Unmarshal(data, &appTitles); err != nil {
		log.Fatalf("Failed to unmarshal app titles data: %v", err)
	}
}

func getAppTitleID(name string) (string, error) {
	for _, result := range appTitles.Results {
		if result.TitleName == name {
			return result.ID, nil
		}
	}
	return "", fmt.Errorf("no matching app title found for name: %s", name)
}

// construct constructs an AppCatalogDeployment object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceJamfAppCatalogDeployment, error) {
	name := d.Get("name").(string)
	appTitleID, err := getAppTitleID(name)
	if err != nil {
		return nil, err
	}

	resource := &jamfpro.ResourceJamfAppCatalogDeployment{
		Name:                            name,
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

	// Serialize and pretty-print the AppCatalogDeployment object as JSON for logging
	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro App Installer Deployment '%s' to JSON: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro App Installer Deployment JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
