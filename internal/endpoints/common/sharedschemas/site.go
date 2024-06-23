package sharedschemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetSharedSchemaSite() *schema.Schema {
	out := &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Jamf Pro Site-related settings of the policy.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Jamf Pro Site ID. Value defaults to -1 aka not used.",
					Default:     -1,
				},
			},
		},
	}

	return out
}
