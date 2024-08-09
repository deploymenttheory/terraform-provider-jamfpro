package policies

// TODO remove log.prints, debug use only
// TODO maybe review error handling here too?

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// stateSelfService Reads response and states self-service items and states only if non-default
func stateSelfService(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	defaults := map[string]interface{}{
		"use_for_self_service":            false,
		"self_service_display_name":       "",
		"install_button_text":             "Install",
		"self_service_description":        "",
		"force_users_to_view_description": false,
		"feature_on_main_page":            false,
	}

	current := map[string]interface{}{
		"use_for_self_service":            resp.SelfService.UseForSelfService,
		"self_service_display_name":       resp.SelfService.SelfServiceDisplayName,
		"install_button_text":             resp.SelfService.InstallButtonText,
		"self_service_description":        resp.SelfService.SelfServiceDescription,
		"force_users_to_view_description": resp.SelfService.ForceUsersToViewDescription,
		"feature_on_main_page":            resp.SelfService.FeatureOnMainPage,
	}

	allDefault := false
	for key, value := range current {
		if value != defaults[key] {
			allDefault = true
			break
		}
	}

	if allDefault {
		return
	}

	out_ss := make([]map[string]interface{}, 0)
	out_ss = append(out_ss, make(map[string]interface{}, 1))

	out_ss[0]["use_for_self_service"] = resp.SelfService.UseForSelfService
	out_ss[0]["self_service_display_name"] = resp.SelfService.SelfServiceDisplayName
	out_ss[0]["install_button_text"] = resp.SelfService.InstallButtonText
	out_ss[0]["self_service_description"] = resp.SelfService.SelfServiceDescription
	out_ss[0]["force_users_to_view_description"] = resp.SelfService.ForceUsersToViewDescription
	out_ss[0]["feature_on_main_page"] = resp.SelfService.FeatureOnMainPage

	out_ss[0]["self_service_category"] = make([]map[string]interface{}, 0)
	if resp.SelfService.SelfServiceCategories != nil {
		for _, v := range resp.SelfService.SelfServiceCategories {
			var categoryBlock []map[string]interface{}
			categoryItem := map[string]interface{}{
				"id":         v.ID,
				"display_in": v.DisplayIn,
				"feature_in": v.FeatureIn,
			}
			categoryBlock = append(categoryBlock, categoryItem)
			out_ss[0]["self_service_category"] = categoryBlock
		}
	}

	err := d.Set("self_service", out_ss)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}
