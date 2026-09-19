package policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourcePolicyV1 preserves the complete schema before repeated blocks became sets.
// Keep these historical collection types when extending the current schema.
func resourcePolicyV1() *schema.Resource {
	r := resourcePolicy()
	payloads := r.Schema["payloads"].Elem.(*schema.Resource)
	for _, name := range []string{"scripts", "printers", "dock_items"} {
		payloads.Schema[name].Type = schema.TypeList
	}
	payloads.Schema["packages"].Elem.(*schema.Resource).Schema["package"].Type = schema.TypeList
	maintenance := payloads.Schema["account_maintenance"].Elem.(*schema.Resource)
	maintenance.Schema["local_accounts"].Elem.(*schema.Resource).Schema["account"].Type = schema.TypeList
	maintenance.Schema["directory_bindings"].Elem.(*schema.Resource).Schema["binding"].Type = schema.TypeList
	r.Schema["self_service"].Elem.(*schema.Resource).Schema["self_service_category"].Type = schema.TypeList
	return r
}

// resourcePolicyV0 preserves all policy fields while restoring the old deferral name.
func resourcePolicyV0() *schema.Resource {
	r := resourcePolicyV1()
	payloads := r.Schema["payloads"].Elem.(*schema.Resource)
	interaction := payloads.Schema["user_interaction"].Elem.(*schema.Resource)
	fields := make(map[string]*schema.Schema, len(interaction.Schema))
	for name, field := range interaction.Schema {
		if name == "allow_users_to_defer" {
			name = "allow_user_to_defer"
		}
		fields[name] = field
	}
	interaction.Schema = fields
	return r
}

// Lists and sets use the same JSON array representation. Terraform re-encodes
// these values with the V2 set schema without dropping existing policy fields.
func upgradePolicyV1toV2(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return rawState, nil
}

func upgradePolicyUserInteractionV0toV1(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if payloads, ok := rawState["payloads"].([]any); ok && len(payloads) > 0 {
		payload := payloads[0].(map[string]any)
		if userInteractions, ok := payload["user_interaction"].([]any); ok && len(userInteractions) > 0 {
			userInteraction := userInteractions[0].(map[string]any)

			// Create a new interaction block with the new field name
			newInteraction := map[string]any{
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
