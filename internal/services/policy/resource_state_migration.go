package policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourcePolicyV0 defines the v0 schema for the user interaction block in the policy resource.
func resourcePolicyV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"payloads": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_interaction": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"message_start": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allow_user_to_defer": { // Old field name
										Type:     schema.TypeBool,
										Optional: true,
									},
									"allow_deferral_until_utc": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allow_deferral_minutes": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"message_finish": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func upgradePolicyUserInteractionV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if payloads, ok := rawState["payloads"].([]interface{}); ok && len(payloads) > 0 {
		payload := payloads[0].(map[string]interface{})
		if userInteractions, ok := payload["user_interaction"].([]interface{}); ok && len(userInteractions) > 0 {
			userInteraction := userInteractions[0].(map[string]interface{})

			// Create a new interaction block with the new field name
			newInteraction := map[string]interface{}{
				"message_start":            userInteraction["message_start"],
				"allow_users_to_defer":     userInteraction["allow_user_to_defer"],
				"allow_deferral_until_utc": userInteraction["allow_deferral_until_utc"],
				"allow_deferral_minutes":   userInteraction["allow_deferral_minutes"],
				"message_finish":           userInteraction["message_finish"],
			}

			userInteractions[0] = newInteraction
			payload["user_interaction"] = userInteractions
			payloads[0] = payload
			rawState["payloads"] = payloads
		}
	}
	return rawState, nil
}
