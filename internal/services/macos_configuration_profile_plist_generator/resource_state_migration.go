package macos_configuration_profile_plist_generator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProMacOSConfigurationProfilesPlistGenerator defines the resource and upgrades legacy collection state.
func ResourceJamfProMacOSConfigurationProfilesPlistGenerator() *schema.Resource {
	r := resourceSchema()
	r.SchemaVersion = 1
	r.StateUpgraders = []schema.StateUpgrader{{
		Version: 0,
		Type:    resourceV0().CoreConfigSchema().ImpliedType(),
		Upgrade: upgradeV0ToV1,
	}}
	return r
}

// The SDK decodes legacy lists with the complete v0 schema. Initialize the
// optional array representation while preserving every pre-existing value.
func upgradeV0ToV1(_ context.Context, state map[string]any, _ any) (map[string]any, error) {
	payloads, _ := state["payloads"].([]any)
	for _, rawPayload := range payloads {
		payload, _ := rawPayload.(map[string]any)
		contents, _ := payload["payload_content"].([]any)
		for _, rawContent := range contents {
			content, _ := rawContent.(map[string]any)
			initializeArrayState(content["setting"])
		}
	}
	return state, nil
}

func initializeArrayState(value any) {
	entries, _ := value.([]any)
	for _, item := range entries {
		entry, ok := item.(map[string]any)
		if !ok || entry == nil {
			continue
		}
		entry["array_json"] = ""
		initializeArrayState(entry["dictionary"])
	}
}
