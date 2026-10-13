package policy

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

// policySetFixture includes every migrated collection and unrelated policy fields.
func policySetFixture() map[string]any {
	return map[string]any{
		"name": "migration-test", "enabled": false, "frequency": "Ongoing",
		"scope": []any{map[string]any{"all_computers": false}},
		"self_service": []any{map[string]any{"self_service_display_name": "Migration", "self_service_category": []any{
			map[string]any{"id": 2, "display_in": true}, map[string]any{"id": 1, "feature_in": true},
		}}},
		"payloads": []any{map[string]any{
			"scripts":    []any{map[string]any{"id": "2", "priority": "Before", "parameter4": "second"}, map[string]any{"id": "1", "priority": "After", "parameter11": "first"}},
			"printers":   []any{map[string]any{"id": 2, "name": "Printer B", "action": "install"}, map[string]any{"id": 1, "name": "Printer A", "action": "uninstall"}},
			"dock_items": []any{map[string]any{"id": 2, "name": "Dock B", "action": "Remove"}, map[string]any{"id": 1, "name": "Dock A", "action": "Add To End"}},
			"packages":   []any{map[string]any{"distribution_point": "default", "package": []any{map[string]any{"id": 2, "action": "Cache"}, map[string]any{"id": 1, "action": "Install"}}}},
			"account_maintenance": []any{map[string]any{
				"local_accounts":     []any{map[string]any{"account": []any{map[string]any{"username": "bravo", "action": "Create", "password": "test-bravo"}, map[string]any{"username": "alpha", "action": "Create", "password": "test-alpha"}}}},
				"directory_bindings": []any{map[string]any{"binding": []any{map[string]any{"name": "Directory B"}, map[string]any{"name": "Directory A"}}}},
			}},
			"user_interaction": []any{map[string]any{"allow_users_to_defer": true, "message_start": "Migration test"}},
		}},
	}
}

func TestPolicyStateUpgradeProtocol(t *testing.T) {
	current := ResourceJamfProPolicies()
	require.Equal(t, 2, current.SchemaVersion)
	require.Len(t, current.StateUpgraders, 2)
	server := schema.NewGRPCProviderServer(&schema.Provider{ResourcesMap: map[string]*schema.Resource{"jamfpro_policy": current}})
	for version, old := range []*schema.Resource{resourcePolicyV0(), resourcePolicyV1()} {
		for _, flatmap := range []bool{false, true} {
			t.Run(fmt.Sprintf("v%d/flatmap=%t", version, flatmap), func(t *testing.T) {
				input := policySetFixture()
				if version == 0 {
					input["payloads"].([]any)[0].(map[string]any)["user_interaction"] = []any{map[string]any{"allow_user_to_defer": true, "message_start": "Migration test"}}
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
				response, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{TypeName: "jamfpro_policy", Version: int64(version), RawState: raw})
				require.NoError(t, err)
				require.Empty(t, response.Diagnostics)
				require.NotNil(t, response.UpgradedState)
				upgraded, err := msgpack.Unmarshal(response.UpgradedState.MsgPack, current.CoreConfigSchema().ImpliedType())
				require.NoError(t, err)
				wantData := schema.TestResourceDataRaw(t, current.Schema, policySetFixture())
				wantData.SetId("42")
				want, err := wantData.State().AttrsAsObjectValue(current.CoreConfigSchema().ImpliedType())
				require.NoError(t, err)
				wantJSON, err := ctyjson.Marshal(want, current.CoreConfigSchema().ImpliedType())
				require.NoError(t, err)
				normalized, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{TypeName: "jamfpro_policy", Version: 2, RawState: &tfprotov5.RawState{JSON: wantJSON}})
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

func reversePolicyCollections(value any) {
	switch v := value.(type) {
	case map[string]any:
		for _, child := range v {
			reversePolicyCollections(child)
		}
	case []any:
		for _, child := range v {
			reversePolicyCollections(child)
		}
		for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
			v[left], v[right] = v[right], v[left]
		}
	}
}

func TestPolicySetReorderingAndChanges(t *testing.T) {
	for _, tc := range []struct {
		name     string
		resource *schema.Resource
		wantDiff bool
	}{
		{"old lists reproduce ordering diff", resourcePolicyV1(), true},
		{"sets ignore ordering", ResourceJamfProPolicies(), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// LDAP validation needs a tenant client; this test exercises schema diffing.
			tc.resource.CustomizeDiff = nil
			input := policySetFixture()
			data := schema.TestResourceDataRaw(t, tc.resource.Schema, input)
			data.SetId("42")
			reversePolicyCollections(input)
			diff, err := tc.resource.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(input), nil)
			require.NoError(t, err)
			require.Equal(t, tc.wantDiff, diff != nil && !diff.Empty())
			// Actual script parameter changes must remain visible.
			input["payloads"].([]any)[0].(map[string]any)["scripts"].([]any)[0].(map[string]any)["parameter4"] = "changed"
			diff, err = tc.resource.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(input), nil)
			require.NoError(t, err)
			require.NotNil(t, diff)
			require.False(t, diff.Empty())
		})
	}
}

func TestPolicySetConstructAndPasswordRefresh(t *testing.T) {
	r := ResourceJamfProPolicies()
	data := schema.TestResourceDataRaw(t, r.Schema, policySetFixture())
	var payload jamfpro.ResourcePolicy
	constructSelfService(data, &payload)
	constructPayloadPackages(data, &payload)
	constructPayloadScripts(data, &payload)
	constructPayloadPrinters(data, &payload)
	constructPayloadDockItems(data, &payload)
	constructPayloadAccountMaintenance(data, &payload)
	require.Len(t, payload.Scripts, 2)
	require.Len(t, payload.PackageConfiguration.Packages, 2)
	require.Len(t, payload.Printers.Printer, 2)
	require.Len(t, payload.DockItems, 2)
	require.Len(t, payload.SelfService.SelfServiceCategories, 2)
	require.Len(t, *payload.AccountMaintenance.DirectoryBindings, 2)
	require.Len(t, *payload.AccountMaintenance.Accounts, 2)
	accounts := *payload.AccountMaintenance.Accounts
	accounts[0], accounts[1] = accounts[1], accounts[0]
	for i := range accounts {
		accounts[i].Password = ""
	}
	payload.AccountMaintenance.Accounts = &accounts
	out := []map[string]any{{}}
	prepStatePayloadAccountMaintenance(&out, &payload, data)
	restored := out[0]["account_maintenance"].([]map[string]any)[0]["local_accounts"].([]map[string]any)[0]["account"].([]map[string]any)
	for _, account := range restored {
		require.Equal(t, "test-"+account["username"].(string), account["password"])
	}
}

func TestPolicySetMembershipChanges(t *testing.T) {
	paths := [][]string{
		{"payloads", "scripts"}, {"payloads", "printers"}, {"payloads", "dock_items"},
		{"payloads", "packages", "package"},
		{"payloads", "account_maintenance", "local_accounts", "account"},
		{"payloads", "account_maintenance", "directory_bindings", "binding"},
		{"self_service", "self_service_category"},
	}
	for _, path := range paths {
		for _, operation := range []string{"remove", "change", "add"} {
			t.Run(fmt.Sprintf("%v/%s", path, operation), func(t *testing.T) {
				r := ResourceJamfProPolicies()
				r.CustomizeDiff = nil // Isolate schema diffing from live LDAP validation.
				input := policySetFixture()
				data := schema.TestResourceDataRaw(t, r.Schema, input)
				data.SetId("42")
				parent := input
				for _, name := range path[:len(path)-1] {
					parent = parent[name].([]any)[0].(map[string]any)
				}
				name := path[len(path)-1]
				items := parent[name].([]any)
				switch operation {
				case "remove":
					parent[name] = items[1:]
				case "change":
					item := items[0].(map[string]any)
					switch name {
					case "scripts":
						item["parameter11"] = "updated"
					case "printers":
						item["make_default"] = true
					case "dock_items":
						item["action"] = "Add To Beginning"
					case "package":
						item["fill_user_template"] = true
					case "account":
						item["password"] = "updated-password"
					case "binding":
						item["name"] = "Updated directory"
					case "self_service_category":
						item["feature_in"] = true
					}
				case "add":
					item := make(map[string]any)
					for key, value := range items[0].(map[string]any) {
						item[key] = value
					}
					switch name {
					case "scripts":
						item["id"] = "3"
					case "account":
						item["username"] = "charlie"
					case "binding":
						item["name"] = "Directory C"
					default:
						item["id"] = 3
					}
					parent[name] = append(items, item)
				}
				diff, err := r.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(input), nil)
				require.NoError(t, err)
				require.NotNil(t, diff)
				require.False(t, diff.Empty())
				require.False(t, diff.RequiresNew(), "nested collection edits must update the existing policy")
			})
		}
	}
}

func TestPolicySetEmptyStateUpgrade(t *testing.T) {
	current := ResourceJamfProPolicies()
	require.NoError(t, current.InternalValidate(nil, true))
	server := schema.NewGRPCProviderServer(&schema.Provider{ResourcesMap: map[string]*schema.Resource{"jamfpro_policy": current}})
	for _, state := range []string{
		`{"id":"42","name":"empty","payloads":[],"self_service":[]}`,
		`{"id":"42","name":"empty","payloads":[{"scripts":[],"printers":[],"dock_items":[]}],"self_service":[]}`,
		`{"id":"42","name":"empty","payloads":[{"scripts":null,"printers":null,"dock_items":null}],"self_service":null}`,
		`{"id":"42","name":"empty"}`,
	} {
		for _, version := range []int64{0, 1} {
			response, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{TypeName: "jamfpro_policy", Version: version, RawState: &tfprotov5.RawState{JSON: []byte(state)}})
			require.NoError(t, err)
			require.Empty(t, response.Diagnostics)
			require.NotNil(t, response.UpgradedState)
		}
	}
}

func TestPolicyScriptParameterChangeDoesNotAddEmptyElement(t *testing.T) {
	r := ResourceJamfProPolicies()
	r.CustomizeDiff = nil
	input := policySetFixture()
	data := schema.TestResourceDataRaw(t, r.Schema, input)
	data.SetId("42")
	for _, item := range input["payloads"].([]any)[0].(map[string]any)["scripts"].([]any) {
		item.(map[string]any)["parameter4"] = "updated"
	}
	diff, err := r.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(input), nil)
	require.NoError(t, err)
	before, err := data.State().AttrsAsObjectValue(r.CoreConfigSchema().ImpliedType())
	require.NoError(t, err)
	after, err := diff.ApplyToValue(before, r.CoreConfigSchema())
	require.NoError(t, err)
	payload := after.GetAttr("payloads").AsValueSlice()[0]
	scripts := payload.GetAttr("scripts").AsValueSlice()
	require.Len(t, scripts, 2, "changing parameters must not resurrect removed set hashes as empty scripts")
	for _, item := range scripts {
		require.False(t, item.GetAttr("id").IsNull())
		require.Equal(t, "updated", item.GetAttr("parameter4").AsString())
	}
}
