package macos_onboarding_settings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceMacOSOnboardingSettingsV0 returns the V0 schema where onboarding_items was a TypeList.
func resourceMacOSOnboardingSettingsV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"onboarding_items": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"entity_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_service_entity_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"OS_X_POLICY", "OS_X_MAC_APP", "OS_X_CONFIGURATION_PROFILE"}, false),
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
		},
	}
}

// upgradeMacOSOnboardingSettingsV0toV1 migrates state from V0 (TypeList) to V1 (TypeSet).
// The underlying data structure is identical; re-encoding is sufficient for Terraform to
// re-hash the items using the new TypeSet schema.
func upgradeMacOSOnboardingSettingsV0toV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return rawState, nil
}
