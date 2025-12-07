package policy

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// stateSelfService Reads response and states self-service items and states only if non-default
func stateSelfService(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	policyName := resp.General.Name

	defaults := map[string]any{
		"use_for_self_service":            false,
		"self_service_display_name":       "",
		"install_button_text":             "Install",
		"reinstall_button_text":           "Reinstall",
		"self_service_description":        "",
		"force_users_to_view_description": false,
		"self_service_icon_id":            0,
		"feature_on_main_page":            false,
		"notification":                    false,
		"notification_type":               "Self Service",
		"notification_subject":            "",
		"notification_message":            "",
	}

	// Get icon ID - only track it when self service is enabled
	// Icons cannot be removed via the API once set, so we don't track orphaned icons on disabled policies
	iconID := 0
	if resp.SelfService.UseForSelfService && resp.SelfService.SelfServiceIcon != nil {
		iconID = resp.SelfService.SelfServiceIcon.ID
	}

	current := map[string]any{
		"use_for_self_service":            resp.SelfService.UseForSelfService,
		"self_service_display_name":       resp.SelfService.SelfServiceDisplayName,
		"install_button_text":             resp.SelfService.InstallButtonText,
		"reinstall_button_text":           resp.SelfService.ReinstallButtonText,
		"self_service_description":        resp.SelfService.SelfServiceDescription,
		"force_users_to_view_description": resp.SelfService.ForceUsersToViewDescription,
		"self_service_icon_id":            iconID,
		"feature_on_main_page":            resp.SelfService.FeatureOnMainPage,
		"notification":                    resp.SelfService.Notification,
		"notification_type":               resp.SelfService.NotificationType,
		"notification_subject":            resp.SelfService.NotificationSubject,
		"notification_message":            resp.SelfService.NotificationMessage,
	}

	allDefault := true
	for key, value := range current {
		// Special case: if self_service_display_name equals the policy name, Jamf Pro auto-populated it
		// Treat it as a default value (empty string) since the user didn't explicitly set it
		if key == "self_service_display_name" && value == policyName {
			continue
		}

		// Special case: if install_button_text is empty, Jamf Pro returns empty string instead of default
		// Treat it as the default "Install"
		if key == "install_button_text" && value == "" {
			continue
		}

		// Special case: if notification_type is empty, Jamf Pro returns empty string instead of default
		// Treat it as the default "Self Service"
		if key == "notification_type" && value == "" {
			continue
		}

		if value != defaults[key] {
			allDefault = false
			break
		}
	}

	if allDefault && len(resp.SelfService.SelfServiceCategories) == 0 {
		return
	}

	out_ss := make([]map[string]any, 0)
	out_ss = append(out_ss, make(map[string]any, 1))

	out_ss[0]["use_for_self_service"] = resp.SelfService.UseForSelfService
	out_ss[0]["self_service_display_name"] = resp.SelfService.SelfServiceDisplayName
	out_ss[0]["install_button_text"] = resp.SelfService.InstallButtonText
	out_ss[0]["reinstall_button_text"] = resp.SelfService.ReinstallButtonText
	out_ss[0]["self_service_description"] = resp.SelfService.SelfServiceDescription
	out_ss[0]["force_users_to_view_description"] = resp.SelfService.ForceUsersToViewDescription
	out_ss[0]["self_service_icon_id"] = iconID
	out_ss[0]["feature_on_main_page"] = resp.SelfService.FeatureOnMainPage
	out_ss[0]["notification"] = resp.SelfService.Notification
	out_ss[0]["notification_type"] = resp.SelfService.NotificationType
	out_ss[0]["notification_subject"] = resp.SelfService.NotificationSubject
	out_ss[0]["notification_message"] = resp.SelfService.NotificationMessage

	categoryBlock := make([]map[string]any, 0)
	if resp.SelfService.SelfServiceCategories != nil {
		for _, v := range resp.SelfService.SelfServiceCategories {
			categoryItem := map[string]any{
				"id":         v.ID,
				"display_in": v.DisplayIn,
				"feature_in": v.FeatureIn,
			}
			categoryBlock = append(categoryBlock, categoryItem)
		}
	}
	out_ss[0]["self_service_category"] = categoryBlock

	err := d.Set("self_service", out_ss)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}
