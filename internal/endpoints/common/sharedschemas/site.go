package sharedschemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetSharedSchemaSite() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Jamf Pro Site ID. Value defaults to -1 aka not used.",
				Default:     -1,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Jamf Pro Site Name. Value defaults to 'None' aka not used",
				Default:     "None",
			},
		},
	}

	return out
}
