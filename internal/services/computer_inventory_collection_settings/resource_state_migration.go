package computer_inventory_collection_settings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Old schema for state upgrade (version 0)
func resourceComputerInventoryCollectionSettingsV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"computer_inventory_collection_preferences": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_fonts":   {Type: schema.TypeBool, Optional: true},
						"include_plugins": {Type: schema.TypeBool, Optional: true},
					},
				},
			},
			"font_paths":   pathSchema(""),
			"plugin_paths": pathSchema(""),
		},
	}
}

// upgrader for v0 -> v1: remove deprecated fields
func upgradeComputerInventoryCollectionSettingsV0toV1(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	newState := make(map[string]any)
	for k, v := range rawState {
		if k == "font_paths" || k == "plugin_paths" {
			continue
		}
		newState[k] = v
	}

	if v, ok := newState["computer_inventory_collection_preferences"]; ok {
		if list, ok := v.([]any); ok && len(list) > 0 {
			if prefs, ok := list[0].(map[string]any); ok {
				cleaned := make(map[string]any)
				for pk, pv := range prefs {
					if pk == "include_fonts" || pk == "include_plugins" {
						continue
					}
					cleaned[pk] = pv
				}
				list[0] = cleaned
				newState["computer_inventory_collection_preferences"] = list
			}
		}
	}

	return newState, nil
}
