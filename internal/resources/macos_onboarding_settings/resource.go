// macos_onboarding_settings_resource.go
package macos_onboarding_settings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProMacOSOnboardingSettings defines the schema and CRUD operations for the macOS onboarding settings resource in Jamf Pro.
func ResourceJamfProMacOSOnboardingSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable or disable macOS onboarding",
			},
			"onboarding_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of onboarding items to display during device setup",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier for the onboarding item (computed)",
						},
						"entity_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the Self Service entity (policy, app, or configuration profile)",
						},
						"entity_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the entity (computed)",
						},
						"scope_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the scope (computed)",
						},
						"site_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the site (computed)",
						},
						"self_service_entity_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of Self Service entity. Supported values are 'OS_X_POLICY', 'OS_X_MAC_APP', 'OS_X_CONFIGURATION_PROFILE'",
							ValidateFunc: validation.StringInSlice([]string{
								"OS_X_POLICY",
								"OS_X_MAC_APP",
								"OS_X_CONFIGURATION_PROFILE",
							}, false),
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Priority order for the onboarding item (lower numbers appear first)",
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
		},
	}
}
