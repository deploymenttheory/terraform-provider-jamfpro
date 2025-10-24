// macos_onboarding_settings_constructor.go
package macos_onboarding_settings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs the macOS onboarding settings resource from the Terraform schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceUpdateOnboardingSettings, error) {
	var onboardingItems []jamfpro.SubsetOnboardingItemRequest

	if v, ok := d.GetOk("onboarding_items"); ok {
		itemsList := v.([]any)
		onboardingItems = make([]jamfpro.SubsetOnboardingItemRequest, 0, len(itemsList))

		for _, item := range itemsList {
			itemMap := item.(map[string]any)

			onboardingItem := jamfpro.SubsetOnboardingItemRequest{
				EntityID:              itemMap["entity_id"].(string),
				SelfServiceEntityType: itemMap["self_service_entity_type"].(string),
				Priority:              itemMap["priority"].(int),
			}

			if id, exists := itemMap["id"].(string); exists && id != "" {
				onboardingItem.ID = id
			}

			onboardingItems = append(onboardingItems, onboardingItem)
		}
	}

	resource := &jamfpro.ResourceUpdateOnboardingSettings{
		Enabled:         d.Get("enabled").(bool),
		OnboardingItems: onboardingItems,
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro macOS Onboarding Settings to JSON: %w", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro macOS Onboarding Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
