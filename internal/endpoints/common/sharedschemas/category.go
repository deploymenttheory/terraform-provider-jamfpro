package sharedschemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetSharedSchemaCategory() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the category to which the configuration profile is scoped.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the category to which the configuration profile is scoped.",
			},
		},
	}

	return out
}
