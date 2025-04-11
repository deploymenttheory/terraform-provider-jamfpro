// advancedmobiledevicesearches_resource.go
package advancedmobiledevicesearches

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProAdvancedMobileDeviceSearches defines the schema for managing advanced mobile device searches in Terraform.
func ResourceJamfProAdvancedMobileDeviceSearches() *schema.Resource {
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
				Description: "The unique identifier of the advanced mobile device search",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the advanced mobile device search",
			},
			"site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-1",
				Description: "The ID of the site to associate the search with",
			},
			"criteria": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of search criteria",
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
								"matches regex", "does not match regex", "member of", "not member of"
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
				Type:        schema.TypeSet, // Use TypeSet instead of TypeList as jamf uses a ranrom order for the fields
				Optional:    true,
				Description: "List of fields to display in the search results",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
