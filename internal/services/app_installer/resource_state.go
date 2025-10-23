// state.go
package app_installer

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest App Catalog Deployment information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceJamfAppCatalogDeployment) diag.Diagnostics {
	var diags diag.Diagnostics

	deploymentData := map[string]interface{}{
		"name":                               resp.Name,
		"enabled":                            resp.Enabled,
		"app_title_id":                       resp.AppTitleId,
		"app_title_name":                     d.Get("app_title_name"),
		"deployment_type":                    resp.DeploymentType,
		"update_behavior":                    resp.UpdateBehavior,
		"category_id":                        resp.CategoryId,
		"site_id":                            resp.SiteId,
		"smart_group_id":                     resp.SmartGroupId,
		"install_predefined_config_profiles": resp.InstallPredefinedConfigProfiles,
		"title_available_in_ais":             resp.TitleAvailableInAis,
		"trigger_admin_notifications":        resp.TriggerAdminNotifications,
		"selected_version":                   resp.SelectedVersion,
		"latest_available_version":           resp.LatestAvailableVersion,
		"version_removed":                    resp.VersionRemoved,
	}

	for key, val := range deploymentData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Update notification settings
	notificationSettings := map[string]interface{}{
		"notification_message":  resp.NotificationSettings.NotificationMessage,
		"notification_interval": resp.NotificationSettings.NotificationInterval,
		"deadline_message":      resp.NotificationSettings.DeadlineMessage,
		"deadline":              resp.NotificationSettings.Deadline,
		"quit_delay":            resp.NotificationSettings.QuitDelay,
		"complete_message":      resp.NotificationSettings.CompleteMessage,
		"relaunch":              resp.NotificationSettings.Relaunch,
		"suppress":              resp.NotificationSettings.Suppress,
	}
	if err := d.Set("notification_settings", []interface{}{notificationSettings}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update self-service settings
	selfServiceSettings := map[string]interface{}{
		"include_in_featured_category":   resp.SelfServiceSettings.IncludeInFeaturedCategory,
		"include_in_compliance_category": resp.SelfServiceSettings.IncludeInComplianceCategory,
		"force_view_description":         resp.SelfServiceSettings.ForceViewDescription,
		"description":                    resp.SelfServiceSettings.Description,
	}

	// Update categories
	var categories []interface{}
	for _, cat := range resp.SelfServiceSettings.Categories {
		category := map[string]interface{}{
			"id": cat.ID,
		}
		if cat.Featured != nil {
			category["featured"] = *cat.Featured
		}
		categories = append(categories, category)
	}
	selfServiceSettings["categories"] = schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"id":       {Type: schema.TypeString},
			"featured": {Type: schema.TypeBool},
		},
	}), categories)

	if err := d.Set("self_service_settings", []interface{}{selfServiceSettings}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
