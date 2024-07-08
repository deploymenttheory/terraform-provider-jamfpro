package sharedschemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetSharedSchemaSite() *schema.Schema {
	out := &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     -1,
		Description: "Jamf Pro Site-related settings of the policy.",
	}

	return out
}
