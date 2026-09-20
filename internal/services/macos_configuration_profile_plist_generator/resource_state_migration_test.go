package macos_configuration_profile_plist_generator

import (
	"context"
	"fmt"
	"testing"

	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func setFixture() map[string]any {
	return map[string]any{"name": "set-test", "distribution_method": "Make Available in Self Service", "scope": []any{map[string]any{"all_computers": false}}, "payloads": []any{map[string]any{"payload_description_header": "Test", "payload_enabled_header": true, "payload_organization_header": "Test", "payload_type_header": "Configuration", "payload_version_header": 1, "payload_content": []any{map[string]any{"payload_enabled": true, "payload_type": "com.apple.test", "payload_version": 1, "setting": []any{map[string]any{"key": "B", "value": "two"}, map[string]any{"key": "A", "dictionary": []any{map[string]any{"key": "Y", "value": "four"}, map[string]any{"key": "X", "value": "three"}}}}}}}}}
}
func TestStateUpgradeProtocol(t *testing.T) {
	current := ResourceJamfProMacOSConfigurationProfilesPlistGenerator()
	require.Equal(t, 1, current.SchemaVersion)
	require.Len(t, current.StateUpgraders, 1)
	server := schema.NewGRPCProviderServer(&schema.Provider{ResourcesMap: map[string]*schema.Resource{"jamfpro_macos_configuration_profile_plist_generator": current}})
	for version, old := range []*schema.Resource{resourceV0()} {
		for _, empty := range []bool{false, true} {
			for _, flatmap := range []bool{false, true} {
				t.Run(fmt.Sprintf("v%d/flatmap=%t/empty=%t", version, flatmap, empty), func(t *testing.T) {
					input := setFixture()
					if empty {
						clearCollections(input)
					}
					oldData := schema.TestResourceDataRaw(t, old.Schema, input)
					oldData.SetId("42")
					oldValue, err := oldData.State().AttrsAsObjectValue(old.CoreConfigSchema().ImpliedType())
					require.NoError(t, err)
					require.True(t, oldValue.IsWhollyKnown())
					raw := &tfprotov5.RawState{}
					if flatmap {
						raw.Flatmap = oldData.State().Attributes
					} else {
						// The SDK's object conversion validates that this state uses the old schema.
						raw.JSON, err = ctyjson.Marshal(oldValue, old.CoreConfigSchema().ImpliedType())
						require.NoError(t, err)
					}
					response, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{TypeName: "jamfpro_macos_configuration_profile_plist_generator", Version: int64(version), RawState: raw})
					require.NoError(t, err)
					require.Empty(t, response.Diagnostics)
					require.NotNil(t, response.UpgradedState)
					upgraded, err := msgpack.Unmarshal(response.UpgradedState.MsgPack, current.CoreConfigSchema().ImpliedType())
					require.NoError(t, err)
					wantData := schema.TestResourceDataRaw(t, current.Schema, input)
					wantData.SetId("42")
					want, err := wantData.State().AttrsAsObjectValue(current.CoreConfigSchema().ImpliedType())
					require.NoError(t, err)
					wantJSON, err := ctyjson.Marshal(want, current.CoreConfigSchema().ImpliedType())
					require.NoError(t, err)
					normalized, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{TypeName: "jamfpro_macos_configuration_profile_plist_generator", Version: 1, RawState: &tfprotov5.RawState{JSON: wantJSON}})
					require.NoError(t, err)
					require.Empty(t, normalized.Diagnostics)
					want, err = msgpack.Unmarshal(normalized.UpgradedState.MsgPack, current.CoreConfigSchema().ImpliedType())
					require.NoError(t, err)
					for key := range want.Type().AttributeTypes() {
						require.True(t, want.GetAttr(key).RawEquals(upgraded.GetAttr(key)), "field %s was not preserved", key)
					}
				})
			}
		}
	}
}

func reverseCollections(value any) {
	switch v := value.(type) {
	case map[string]any:
		for _, child := range v {
			reverseCollections(child)
		}
	case []any:
		for _, child := range v {
			reverseCollections(child)
		}
		for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
			v[left], v[right] = v[right], v[left]
		}
	}
}

func TestSetReorderingAndChanges(t *testing.T) {
	for _, tc := range []struct {
		name    string
		r       *schema.Resource
		changed bool
	}{
		{"legacy list order", resourceV0(), true}, {"set order", ResourceJamfProMacOSConfigurationProfilesPlistGenerator(), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.r.CustomizeDiff = nil
			input := setFixture()
			data := schema.TestResourceDataRaw(t, tc.r.Schema, input)
			data.SetId("42")
			reverseCollections(input)
			diff, err := tc.r.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(input), nil)
			require.NoError(t, err)
			require.Equal(t, tc.changed, diff != nil && !diff.Empty())
		})
	}
}

func mutateCollection(value any, key, action string) bool {
	switch v := value.(type) {
	case map[string]any:
		if items, ok := v[key].([]any); ok && len(items) > 0 {
			switch action {
			case "remove":
				v[key] = items[1:]
			case "clear":
				v[key] = []any{}
			case "duplicate":
				v[key] = append(items, items[0])
			default:
				var changed any
				switch item := items[0].(type) {
				case int:
					changed = item + 100
				case map[string]any:
					m := map[string]any{}
					for k, val := range item {
						m[k] = val
					}
					if _, ok := m["group_name"]; ok {
						m["group_name"] = "changed"
					} else if _, ok := m["key"]; ok {
						m["value"] = "changed"
					} else if id, ok := m["id"].(int); ok {
						m["id"] = id + 100
					} else {
						m["username"] = "changed"
					}
					changed = m
				}
				if action == "add" {
					v[key] = append(items, changed)
				} else {
					items[0] = changed
				}
			}
			return true
		}
		for _, child := range v {
			if mutateCollection(child, key, action) {
				return true
			}
		}
	case []any:
		for _, child := range v {
			if mutateCollection(child, key, action) {
				return true
			}
		}
	}
	return false
}
func TestSetMembershipChanges(t *testing.T) {
	for _, key := range []string{"setting", "dictionary"} {
		for _, action := range []string{"change", "add", "remove", "clear", "duplicate"} {
			t.Run(key+"/"+action, func(t *testing.T) {
				r := ResourceJamfProMacOSConfigurationProfilesPlistGenerator()
				r.CustomizeDiff = nil
				input := setFixture()
				data := schema.TestResourceDataRaw(t, r.Schema, input)
				data.SetId("42")
				require.True(t, mutateCollection(input, key, action))
				diff, err := r.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(input), nil)
				require.NoError(t, err)
				require.Equal(t, action != "duplicate", diff != nil && !diff.Empty())
			})
		}
	}
}

func TestOrderedArrayValueChanges(t *testing.T) {
	r := ResourceJamfProMacOSConfigurationProfilesPlistGenerator()
	r.CustomizeDiff = nil
	input := setFixture()
	setting := input["payloads"].([]any)[0].(map[string]any)["payload_content"].([]any)[0].(map[string]any)["setting"].([]any)[0].(map[string]any)
	setting["value"] = ""
	setting["array_json"] = `["b","a","b"]`
	data := schema.TestResourceDataRaw(t, r.Schema, input)
	data.SetId("42")
	setting["array_json"] = `["a","b","b"]`
	diff, err := r.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(input), nil)
	require.NoError(t, err)
	require.NotNil(t, diff)
	require.False(t, diff.Empty())
}

func clearCollections(value any) {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			switch key {
			case "assigned_computer_ids", "assigned_mobile_device_ids", "member_ids", "assigned_user_ids", "user_additions", "user_deletions", "self_service_category", "ldap_group_access", "setting":
				v[key] = []any{}
			default:
				clearCollections(child)
			}
		}
	case []any:
		for _, child := range v {
			clearCollections(child)
		}
	}
}

func TestEntryIdentityNormalizesAbsentFieldsButPreservesArrayOrder(t *testing.T) {
	require.Equal(t, hashPlistEntry(map[string]any{"key": "A", "value": "one"}), hashPlistEntry(map[string]any{"key": "A", "value": "one", "array_json": "", "dictionary": schema.NewSet(hashPlistEntry, nil)}))
	require.Equal(t, hashPlistEntry(map[string]any{"key": "A", "array_json": `[ "b", "a", "b" ]`}), hashPlistEntry(map[string]any{"key": "A", "array_json": `["b","a","b"]`}))
	require.NotEqual(t, hashPlistEntry(map[string]any{"key": "A", "array_json": `["b","a","b"]`}), hashPlistEntry(map[string]any{"key": "A", "array_json": `["a","b","b"]`}))
}
