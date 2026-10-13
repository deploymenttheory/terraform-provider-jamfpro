package enrollment_customization

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func setFixture() map[string]any {
	return map[string]any{"enrollment_customization_image_source": "", "branding_settings": []any{map[string]any{"text_color": "000000", "button_color": "000000", "button_text_color": "FFFFFF", "background_color": "FFFFFF"}}, "display_name": "set-test", "description": "Preserve description", "ldap_pane": []any{map[string]any{"display_name": "LDAP", "rank": 1, "title": "Sign in", "username_label": "Username", "password_label": "Password", "back_button_text": "Back", "continue_button_text": "Continue", "ldap_group_access": []any{map[string]any{"group_name": "bravo", "ldap_server_id": 2}, map[string]any{"group_name": "alpha", "ldap_server_id": 1}}}}}
}
func TestStateUpgradeProtocol(t *testing.T) {
	current := ResourceJamfProEnrollmentCustomization()
	require.Equal(t, 1, current.SchemaVersion)
	require.Len(t, current.StateUpgraders, 1)
	server := schema.NewGRPCProviderServer(&schema.Provider{ResourcesMap: map[string]*schema.Resource{"jamfpro_enrollment_customization": current}})
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
					response, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{TypeName: "jamfpro_enrollment_customization", Version: int64(version), RawState: raw})
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
					normalized, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{TypeName: "jamfpro_enrollment_customization", Version: 1, RawState: &tfprotov5.RawState{JSON: wantJSON}})
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
		{"legacy list order", resourceV0(), true}, {"set order", ResourceJamfProEnrollmentCustomization(), false},
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
	for _, key := range []string{"ldap_group_access"} {
		for _, action := range []string{"change", "add", "remove", "clear", "duplicate"} {
			t.Run(key+"/"+action, func(t *testing.T) {
				r := ResourceJamfProEnrollmentCustomization()
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

func TestConstructorUsesLDAPAccessSet(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ResourceJamfProEnrollmentCustomization().Schema, setFixture())
	pane := data.Get("ldap_pane").([]any)[0].(map[string]any)
	payload, err := constructLDAPPane(pane)
	require.NoError(t, err)
	require.Len(t, payload.LDAPGroupAccess, 2)
	require.Equal(t, 1, payload.Rank)
}

// TestExportOfflineFixture optionally writes synthetic state for CLI plan verification.
// It never contacts Jamf Pro; the fixture must only be planned with -refresh=false.
func TestExportOfflineFixture(t *testing.T) {
	dir := os.Getenv("JAMF_SET_OFFLINE_FIXTURE_DIR")
	if dir == "" {
		t.Skip("optional CLI fixture export")
	}
	require.NoError(t, os.MkdirAll(dir, 0700))
	old := resourceV0()
	input := setFixture()
	data := schema.TestResourceDataRaw(t, old.Schema, input)
	data.SetId("42")
	value, err := data.State().AttrsAsObjectValue(old.CoreConfigSchema().ImpliedType())
	require.NoError(t, err)
	attrs, err := ctyjson.Marshal(value, old.CoreConfigSchema().ImpliedType())
	require.NoError(t, err)
	state := map[string]any{"version": 4, "terraform_version": "1.14.6", "serial": 1, "lineage": "00000000-0000-4000-8000-000000000042", "outputs": map[string]any{}, "resources": []any{map[string]any{"mode": "managed", "type": "jamfpro_enrollment_customization", "name": "test", "provider": "provider[\"registry.terraform.io/deploymenttheory/jamfpro\"]", "instances": []any{map[string]any{"schema_version": 0, "attributes": json.RawMessage(attrs), "sensitive_attributes": []any{}}}}}}
	encoded, err := json.MarshalIndent(state, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "terraform.tfstate"), encoded, 0600))
	reverseCollections(input)
	config := map[string]any{"resource": map[string]any{"jamfpro_enrollment_customization": map[string]any{"test": input}}}
	encoded, err = json.MarshalIndent(config, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf.json"), encoded, 0600))
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
