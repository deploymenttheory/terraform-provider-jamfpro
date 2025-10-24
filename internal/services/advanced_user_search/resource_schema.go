// advancedusersearches_resource.go
package advanced_user_search

import (
	"time"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProAdvancedUserSearches defines the schema for managing advanced user Searches in Terraform.
func ResourceJamfProAdvancedUserSearches() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the advanced user search",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the advanced user search",
			},
			"criteria": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the search criteria field",
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 100),
							Description:  "Priority order of the criteria. Default is 0, 0 is always used for the first criterion.",
						},
						"and_or": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"and", "or"}, false),
							Description:  "Logical operator (and/or) for the search criteria",
						},
						"search_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"is", "is not", "like", "not like", "has", "does not have",
								"greater than", "less than", "greater than or equal", "less than or equal",
								"matches regex", "does not match regex", "member of", "not member of",
								"more than x days ago",
							}, false),
							Description: "Type of search to perform",
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							Description:  "Value to search for",
						},
						"opening_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether this criterion has an opening parenthesis",
						},
						"closing_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether this criterion has a closing parenthesis",
						},
					},
				},
			},
			"display_fields": {
				Type:        schema.TypeSet,
				Description: "Set of display fields",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"site_id": sharedschemas.GetSharedSchemaSite(),
		},
	}
}
