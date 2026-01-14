package policy

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// stateSelfService Reads response and states self-service items and states only if non-default
func stateSelfService(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {

	if !resp.SelfService.UseForSelfService {
		d.Set("self_service", "")
		return
	}

	// This matches the UI behaviour as close as possible, I'm purposefully not obsecuring the logic.
	state_icon_val := d.Get("self_service.0.self_service_icon_id")
	server_icon_val := resp.SelfService.SelfServiceIcon.ID

	if state_icon_val == 0 && server_icon_val != 0 {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid configuration - API Limitation",
			Detail:   "Unable to unset icon ID once set. Please assign a different icon or replace the policy.",
		})
	}

	out_self_service := []map[string]any{{}}
	out_self_service_slice := out_self_service[0]

	out_self_service_slice = map[string]any{
		"use_for_self_service":            resp.SelfService.UseForSelfService,
		"self_service_display_name":       resp.SelfService.SelfServiceDisplayName,
		"install_button_text":             resp.SelfService.InstallButtonText,
		"reinstall_button_text":           resp.SelfService.ReinstallButtonText,
		"self_service_description":        resp.SelfService.SelfServiceDescription,
		"force_users_to_view_description": resp.SelfService.ForceUsersToViewDescription,
		"self_service_icon_id":            resp.SelfService.SelfServiceIcon.ID,
		"feature_on_main_page":            resp.SelfService.FeatureOnMainPage,
		"notification":                    resp.SelfService.Notification,
		"notification_type":               resp.SelfService.NotificationType,
		"notification_subject":            resp.SelfService.NotificationSubject,
		"notification_message":            resp.SelfService.NotificationMessage,
	}

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

	out_self_service_slice["self_service_category"] = categoryBlock

	err := d.Set("self_service", out_self_service)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}
