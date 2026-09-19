package account_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProAccountGroups defines the resource and upgrades legacy collection state.
func ResourceJamfProAccountGroups() *schema.Resource {
	r := resourceSchema()
	r.SchemaVersion = 1
	r.StateUpgraders = []schema.StateUpgrader{{
		Version: 0,
		Type:    resourceV0().CoreConfigSchema().ImpliedType(),
		Upgrade: upgradeV0ToV1,
	}}
	return r
}

// Lists and sets share the JSON array representation. The SDK decodes legacy
// flatmap state with the complete v0 schema and re-encodes every field as v1.
func upgradeV0ToV1(_ context.Context, state map[string]any, _ any) (map[string]any, error) {
	return state, nil
}
