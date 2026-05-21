package macos_onboarding_settings

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpgradeMacOSOnboardingSettingsV0toV1 verifies that state produced by the
// V0 schema (where onboarding_items was a TypeList) can be re-loaded by the
// V1 schema (TypeSet) without losing data or producing spurious diffs.
//
// TypeList and TypeSet share the same wire format in raw state ([]any of
// map[string]any), so the upgrade itself is a pass-through. These tests pin
// that invariant down so a future change to either schema can't silently
// break existing state.
func TestUpgradeMacOSOnboardingSettingsV0toV1(t *testing.T) {
	tests := []struct {
		name     string
		rawState map[string]any
	}{
		{
			name: "single item",
			rawState: map[string]any{
				"enabled": true,
				"onboarding_items": []any{
					map[string]any{
						"id":                       "14",
						"entity_id":                "76",
						"entity_name":              "Recon Policy",
						"scope_description":        "All Computers",
						"site_description":         "None",
						"self_service_entity_type": "OS_X_POLICY",
						"priority":                 1,
					},
				},
			},
		},
		{
			name: "multiple items out of priority order",
			rawState: map[string]any{
				"enabled": true,
				"onboarding_items": []any{
					map[string]any{
						"id":                       "16",
						"entity_id":                "136",
						"entity_name":              "Item C",
						"scope_description":        "All Computers",
						"site_description":         "None",
						"self_service_entity_type": "OS_X_POLICY",
						"priority":                 3,
					},
					map[string]any{
						"id":                       "14",
						"entity_id":                "76",
						"entity_name":              "Item A",
						"scope_description":        "All Computers",
						"site_description":         "None",
						"self_service_entity_type": "OS_X_POLICY",
						"priority":                 1,
					},
					map[string]any{
						"id":                       "15",
						"entity_id":                "134",
						"entity_name":              "Item B",
						"scope_description":        "All Computers",
						"site_description":         "None",
						"self_service_entity_type": "OS_X_MAC_APP",
						"priority":                 2,
					},
				},
			},
		},
		{
			name: "empty onboarding_items",
			rawState: map[string]any{
				"enabled":          false,
				"onboarding_items": []any{},
			},
		},
		{
			name: "onboarding_items absent",
			rawState: map[string]any{
				"enabled": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sanity check: the rawState must be valid under the V0 schema we
			// are migrating from, otherwise the test premise is wrong.
			_ = schema.TestResourceDataRaw(
				t,
				resourceMacOSOnboardingSettingsV0().Schema,
				tt.rawState,
			)

			upgraded, err := upgradeMacOSOnboardingSettingsV0toV1(
				context.Background(),
				tt.rawState,
				nil,
			)
			require.NoError(t, err)
			require.NotNil(t, upgraded)

			// Pass-through migration: structure must be preserved verbatim so
			// the SDK can re-hash list entries into the new TypeSet without
			// losing data.
			assert.Equal(t, tt.rawState, upgraded)

			// The upgraded state must round-trip cleanly through the V1
			// (current) schema.
			rd := schema.TestResourceDataRaw(
				t,
				ResourceJamfProMacOSOnboardingSettings().Schema,
				upgraded,
			)
			assert.Equal(t, tt.rawState["enabled"], rd.Get("enabled"))

			gotItems := rd.Get("onboarding_items").(*schema.Set).List()
			wantItems, _ := tt.rawState["onboarding_items"].([]any)
			assert.Len(t, gotItems, len(wantItems))

			// Every input item must appear in the resulting set with all
			// fields intact (Set is order-independent, so compare as a
			// multiset keyed by id).
			byID := make(map[string]map[string]any, len(gotItems))
			for _, it := range gotItems {
				m := it.(map[string]any)
				byID[m["id"].(string)] = m
			}
			for _, raw := range wantItems {
				want := raw.(map[string]any)
				got, ok := byID[want["id"].(string)]
				if !assert.True(t, ok, "missing item id=%v in upgraded state", want["id"]) {
					continue
				}
				assert.Equal(t, want["entity_id"], got["entity_id"])
				assert.Equal(t, want["entity_name"], got["entity_name"])
				assert.Equal(t, want["scope_description"], got["scope_description"])
				assert.Equal(t, want["site_description"], got["site_description"])
				assert.Equal(t, want["self_service_entity_type"], got["self_service_entity_type"])
				assert.Equal(t, want["priority"], got["priority"])
			}
		})
	}
}

// TestStateUpgraderWiring verifies the resource's StateUpgraders slice is
// wired correctly so the V0->V1 path is actually exercised by Terraform.
func TestStateUpgraderWiring(t *testing.T) {
	r := ResourceJamfProMacOSOnboardingSettings()

	assert.Equal(t, 1, r.SchemaVersion, "current schema version must be 1")
	require.Len(t, r.StateUpgraders, 1)

	u := r.StateUpgraders[0]
	assert.Equal(t, 0, u.Version, "upgrader source version must be 0")
	assert.NotNil(t, u.Upgrade, "upgrader must have an Upgrade function")
	assert.NotNil(t, u.Type, "upgrader must declare the prior schema type")
}
