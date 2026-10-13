package macos_configuration_profile_plist_generator

import (
	"context"
	"testing"

	plistutil "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestArrayNumericRoundTrip(t *testing.T) {
	for _, raw := range []string{`[1.0,1e0,1.00]`, `[9007199254740993,18446744073709551615]`} {
		t.Run(raw, func(t *testing.T) {
			input := setFixture()
			payload := input["payloads"].([]any)[0].(map[string]any)
			content := payload["payload_content"].([]any)[0].(map[string]any)
			content["setting"] = []any{map[string]any{"key": "Ordered", "array_json": raw}}
			d := schema.TestResourceDataRaw(t, resourceSchema().Schema, input)
			xml, err := plistutil.ConvertHCLToPlist(d)
			require.NoError(t, err)
			hcl, err := plistutil.ConvertPlistToHCL(xml)
			require.NoError(t, err)
			settings := hcl[0].(map[string]any)["payload_content"].([]any)[0].(map[string]any)["setting"].([]any)
			got := settings[0].(map[string]any)["array_json"]
			require.Equal(t, arrayJSONSchema().StateFunc(raw), arrayJSONSchema().StateFunc(got))
		})
	}
}

func TestArrayRejectsUnsupportedValuesBeforeApply(t *testing.T) {
	for _, raw := range []string{`[null,"a"]`, `[{"null":null}]`, `[18446744073709551616]`, `[-9223372036854775809]`, `[1e400]`, `[1e-400]`} {
		t.Run(raw, func(t *testing.T) {
			_, errors := arrayJSONSchema().ValidateFunc(raw, "array_json")
			require.NotEmpty(t, errors)
			_, err := plistutil.ParseJSONArray(raw)
			require.Error(t, err)
		})
	}
}

func TestMigrationPreservesNullPayloadContent(t *testing.T) {
	state := map[string]any{"id": "42", "payloads": []any{map[string]any{"payload_content": nil}}}
	require.NotPanics(t, func() {
		upgraded, err := upgradeV0ToV1(context.Background(), state, nil)
		require.NoError(t, err)
		require.Equal(t, state, upgraded)
	})
}

func TestArrayHashPreservesNumericTypes(t *testing.T) {
	hash := func(raw string) int { return hashPlistEntry(map[string]any{"key": "Numbers", "array_json": raw}) }
	require.Equal(t, hash(`[1.0,{"nested":[2e0,2.0]}]`), hash(`[1e0,{"nested":[2.0,2.00]}]`))
	require.NotEqual(t, hash(`[1.0]`), hash(`[1]`))
	require.NotEqual(t, hash(`[1,2,1]`), hash(`[1,1,2]`))
	require.NotEqual(t, hash(`[1,1]`), hash(`[1]`))
}
