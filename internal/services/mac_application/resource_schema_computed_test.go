package mac_application

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestResourceSchemaIsValid(t *testing.T) {
	require.NoError(t, ResourceJamfProMacApplication().InternalValidate(nil, true))
}

// TestAppleOwnedMetadataDoesNotDrift runs the real schema diff with state as
// Jamf Pro reports it. Jamf Pro updates version, description and icon for
// Mac App Store apps on its own, so omitting them from config must not produce
// a diff, while setting them explicitly must still be honored.
// See https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/1186.
func TestAppleOwnedMetadataDoesNotDrift(t *testing.T) {
	watched := []string{"version", "self_service.0.self_service_description", "self_service.0.self_service_icon"}

	tests := []struct {
		name        string
		config      map[string]any
		wantChanged []string
	}{
		{
			name:        "omitted fields keep the values Jamf Pro reports",
			config:      map[string]any{},
			wantChanged: nil,
		},
		{
			name:        "explicit version is still applied",
			config:      map[string]any{"version": "16.100"},
			wantChanged: []string{"version"},
		},
		{
			name:        "explicit description is still applied",
			config:      map[string]any{"self_service_description": "Our own description"},
			wantChanged: []string{"self_service.0.self_service_description"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResourceJamfProMacApplication()

			state := r.Data(nil)
			state.SetId("42")
			require.NoError(t, state.Set("name", "Microsoft Excel"))
			require.NoError(t, state.Set("bundle_id", "com.microsoft.Excel"))
			require.NoError(t, state.Set("version", "16.113"))
			require.NoError(t, state.Set("scope", []any{map[string]any{"all_computers": false}}))
			require.NoError(t, state.Set("self_service", []any{map[string]any{
				"install_button_text":      "Install",
				"self_service_description": "Description written by Apple",
				"self_service_icon":        []any{map[string]any{"id": 5, "uri": "https://example.com/icon.png"}},
			}}))

			selfService := map[string]any{"install_button_text": "Install"}
			if v, ok := tt.config["self_service_description"]; ok {
				selfService["self_service_description"] = v
			}
			config := map[string]any{
				"name":         "Microsoft Excel",
				"bundle_id":    "com.microsoft.Excel",
				"scope":        []any{map[string]any{"all_computers": false}},
				"self_service": []any{selfService},
			}
			if v, ok := tt.config["version"]; ok {
				config["version"] = v
			}

			// An empty client satisfies the scope validator, which only calls the API
			// when directory service group names are in scope.
			diff, err := r.Diff(context.Background(), state.State(), terraform.NewResourceConfigRaw(config), &jamfpro.Client{})
			require.NoError(t, err)

			var changed []string
			if diff != nil {
				for key, attr := range diff.Attributes {
					for _, w := range watched {
						if (key == w || strings.HasPrefix(key, w+".")) && attr.Old != attr.New {
							changed = append(changed, w)
						}
					}
				}
			}
			sort.Strings(changed)
			changed = compactStrings(changed)
			require.Equal(t, tt.wantChanged, changed)
		})
	}
}

func compactStrings(in []string) []string {
	var out []string
	for i, s := range in {
		if i == 0 || s != in[i-1] {
			out = append(out, s)
		}
	}
	return out
}
