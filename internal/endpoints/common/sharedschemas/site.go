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
				Computed:    true,
				Description: "The name of the Jamf Pro Site. Computed value based on the utilised site ID.",
			},
		},
	}

	return out
}
