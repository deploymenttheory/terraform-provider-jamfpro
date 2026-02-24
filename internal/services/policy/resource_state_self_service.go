package policy

import (
	"maps"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// stateSelfService Reads response and states self-service items and states only if non-default
func stateSelfService(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	if !resp.SelfService.UseForSelfService {
		return
	}

	out_self_service := []map[string]any{{}}
	out_self_service_slice := out_self_service[0]

	maps.Copy(out_self_service_slice, map[string]any{
		"use_for_self_service":            resp.SelfService.UseForSelfService,
		"self_service_display_name":       resp.SelfService.SelfServiceDisplayName,
		"install_button_text":             resp.SelfService.InstallButtonText,
		"reinstall_button_text":           resp.SelfService.ReinstallButtonText,
		"self_service_description":        utils.NormalizeWhitespace(resp.SelfService.SelfServiceDescription),
		"force_users_to_view_description": resp.SelfService.ForceUsersToViewDescription,
		"self_service_icon_id":            resp.SelfService.SelfServiceIcon.ID,
		"feature_on_main_page":            resp.SelfService.FeatureOnMainPage,
		"notification":                    resp.SelfService.Notification,
		"notification_type":               resp.SelfService.NotificationType,
		"notification_subject":            resp.SelfService.NotificationSubject,
		"notification_message":            resp.SelfService.NotificationMessage,
	})

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
