package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getPolicySchemaPackages() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the package.",
			},
			// "name": {
			// 	Type:        schema.TypeString,
			// 	Description: "Name of the package.",
			// 	Computed:    true,
			// }, // No scoping by name
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Action to be performed for the package.",
				ValidateFunc: validation.StringInSlice([]string{"Install", "Cache", "Install Cached"}, false),
				Default:      "Install",
			},
			"fill_user_template": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fill User Template (FUT).",
			},
			"fill_existing_user_template": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fill Existing Users (FEU).",
			},
			// "update_autorun": { // NOT IN THE UI or in response
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Description: "Update auto-run status of the package.",
			// },
		},
	}

	return out
}
