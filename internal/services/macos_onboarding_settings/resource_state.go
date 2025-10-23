// macos_onboarding_settings_state.go
package macos_onboarding_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the state of the macOS onboarding settings resource in Terraform
func updateState(d *schema.ResourceData, resp *jamfpro.ResponseOnboardingSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("enabled", resp.Enabled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	onboardingItems := make([]any, 0, len(resp.OnboardingItems))
	for _, item := range resp.OnboardingItems {
		itemMap := map[string]any{
			"id":                       item.ID,
			"entity_id":                item.EntityID,
			"entity_name":              item.EntityName,
			"scope_description":        item.ScopeDescription,
			"site_description":         item.SiteDescription,
			"self_service_entity_type": item.SelfServiceEntityType,
			"priority":                 item.Priority,
		}
		onboardingItems = append(onboardingItems, itemMap)
	}

	if err := d.Set("onboarding_items", onboardingItems); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
