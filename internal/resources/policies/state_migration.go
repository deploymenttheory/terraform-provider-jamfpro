package policies

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourcePolicyUserInteractionV0 defines the v0 schema for the user interaction block in the policy resource.
func resourcePolicyUserInteractionV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"message_start":            {Type: schema.TypeString, Optional: true},
			"allow_user_to_defer":      {Type: schema.TypeBool, Optional: true, Default: false},
			"allow_deferral_until_utc": {Type: schema.TypeString, Optional: true},
			"allow_deferral_minutes":   {Type: schema.TypeInt, Optional: true, Default: 0},
			"message_finish":           {Type: schema.TypeString, Optional: true},
		},
	}
}

// Migration function to handle the upgrade from V0 to V1
// copies the existing fields and updates the allow_user_to_defer field to allow_users_to_defer
func upgradePolicyUserInteractionV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	newState := make(map[string]interface{})

	newState["message_start"] = rawState["message_start"]
	newState["allow_deferral_until_utc"] = rawState["allow_deferral_until_utc"]
	newState["allow_deferral_minutes"] = rawState["allow_deferral_minutes"]
	newState["message_finish"] = rawState["message_finish"]

	if deferVal, ok := rawState["allow_user_to_defer"]; ok {
		newState["allow_users_to_defer"] = deferVal
	} else {
		newState["allow_users_to_defer"] = false
	}

	return newState, nil
}
