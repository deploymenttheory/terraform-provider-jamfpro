package plist

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDictionarySetsAndOrderedArrays(t *testing.T) {
	child := &schema.Resource{Schema: map[string]*schema.Schema{
		"key":   {Type: schema.TypeString, Required: true},
		"value": {Type: schema.TypeString, Optional: true},
	}}
	entries := []any{map[string]any{"key": "B", "value": "two"}, map[string]any{"key": "A", "value": "three"}}
	for _, collection := range []any{entries, schema.NewSet(schema.HashResource(child), entries)} {
		value, err := settingValue(map[string]any{"dictionary": collection})
		require.NoError(t, err)
		require.Equal(t, map[string]any{"B": "two", "A": "three"}, value)
	}
	for _, raw := range []string{`["b","a","b"]`, `[{"key":"b"},{"key":"a"},{"key":"b"}]`, `[]`, `[["b","a","b"],true,2]`} {
		value, err := settingValue(map[string]any{"key": "array", "array_json": raw})
		require.NoError(t, err)
		var extracted []any
		extractNestedConfigurationSettings(map[string]any{"array": value}, &extracted)
		require.JSONEq(t, raw, extracted[0].(map[string]any)["array_json"].(string))
		restored, err := settingValue(extracted[0].(map[string]any))
		require.NoError(t, err)
		require.Equal(t, value, restored)
	}
}

func TestNestedScalarAndLeafMap(t *testing.T) {
	value, err := settingValue(map[string]any{"dictionary": []any{map[string]any{"key": "scalar", "value": "hello", "dictionary": []any{}}}})
	require.NoError(t, err)
	require.Equal(t, map[string]any{"scalar": "hello"}, value)
	value, err = settingValue(map[string]any{"dictionary": map[string]any{"leaf": "true"}})
	require.NoError(t, err)
	require.Equal(t, map[string]any{"leaf": true}, value)
	_, err = settingValue(map[string]any{"array_json": "{}"})
	require.Error(t, err)
}

func TestArraySurvivesPlistSerialization(t *testing.T) {
	raw := `["b","a","b",{"inside":[2,1,2]},true,1.5]`
	value, err := settingValue(map[string]any{"array_json": raw})
	require.NoError(t, err)
	profile := &ConfigurationProfile{PayloadContent: []PayloadContent{{ConfigurationItems: map[string]any{"Ordered": value}}}}
	xml, err := MarshalPayload(profile)
	require.NoError(t, err)
	hcl, err := ConvertPlistToHCL(xml)
	require.NoError(t, err)
	settings := hcl[0].(map[string]any)["payload_content"].([]any)[0].(map[string]any)["setting"].([]any)
	require.JSONEq(t, raw, settings[0].(map[string]any)["array_json"].(string))
}

func TestArrayDoesNotSilentlyConvertBinaryData(t *testing.T) {
	var output []any
	require.Error(t, extractNestedConfigurationSettings(map[string]any{"Data": []any{[]byte("data")}}, &output))
}
