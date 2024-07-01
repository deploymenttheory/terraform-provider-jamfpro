package policies

// TODO remove log.prints, debug use only
// TODO maybe review error handling here too?

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// stateSelfService Reads response and states self-service items and states only if non-default
func stateSelfService(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	if resp.SelfService == nil {
		return
	}

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

	nonDefault := false
	for key, value := range current {
		if value != defaults[key] {
			nonDefault = true
			break
		}
	}

	if !nonDefault {
		log.Println("[DEBUG] Self-service block has only default values, skipping state")
		return
	}

	log.Println("[DEBUG] Initializing self-service block in state")
	out_ss := make([]map[string]interface{}, 0)
	out_ss = append(out_ss, make(map[string]interface{}, 1))

	out_ss[0]["use_for_self_service"] = resp.SelfService.UseForSelfService
	out_ss[0]["self_service_display_name"] = resp.SelfService.SelfServiceDisplayName
	out_ss[0]["install_button_text"] = resp.SelfService.InstallButtonText
	out_ss[0]["self_service_description"] = resp.SelfService.SelfServiceDescription
	out_ss[0]["force_users_to_view_description"] = resp.SelfService.ForceUsersToViewDescription
	out_ss[0]["feature_on_main_page"] = resp.SelfService.FeatureOnMainPage

	err := d.Set("self_service", out_ss)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}
