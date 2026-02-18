package smart_mobile_device_group

import (
	"fmt"
	"time"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	And                          string = "and"
	Or                           string = "or"
	SearchTypeIs                 string = "is"
	SearchTypeIsNot              string = "is not"
	SearchTypeHas                string = "has"
	SearchTypeDoesNotHave        string = "does not have"
	SearchTypeMemberOf           string = "member of"
	SearchTypeNotMemberOf        string = "not member of"
	SearchTypeBeforeYYYYMMDD     string = "before (yyyy-mm-dd)"
	SearchTypeAfterYYYYMMDD      string = "after (yyyy-mm-dd)"
	SearchTypeMoreThanXDaysAgo   string = "more than x days ago"
	SearchTypeLessThanXDaysAgo   string = "less than x days ago"
	SearchTypeLike               string = "like"
	SearchTypeNotLike            string = "not like"
	SearchTypeGreaterThan        string = "greater than"
	SearchTypeMoreThan           string = "more than"
	SearchTypeLessThan           string = "less than"
	SearchTypeGreaterThanOrEqual string = "greater than or equal"
	SearchTypeLessThanOrEqual    string = "less than or equal"
	SearchTypeMatchesRegex       string = "matches regex"
	SearchTypeDoesNotMatch       string = "does not match regex"
)

// ResourceJamfProSmartMobileGroups defines the schema and CRUD operations for managing Jamf Pro smart mobile device groups in Terraform.
func ResourceJamfProSmartMobileGroups() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is superseded by jamfpro_smart_mobile_device_group_v1, which uses newly available API endpoints. The endpoints backing jamfpro_smart_mobile_device_group are likely to be deprecated. It is recommended to migrate to jamfpro_smart_mobile_device_group_v1 as soon as possible. Review the migration guide for more details: https://github.com/deploymenttheory/terraform-provider-jamfpro/docs/resource-migration-guide.md",
		CreateContext:      create,
		ReadContext:        readWithCleanup,
		UpdateContext:      update,
		DeleteContext:      delete,
		CustomizeDiff:      mainCustomDiffFunc,
		Description:        "This resource is superseded by jamfpro_smart_mobile_device_group_v1, which uses newly available API endpoints. The endpoints backing jamfpro_smart_mobile_device_group are likely to be deprecated. It is recommended to migrate to jamfpro_smart_mobile_device_group_v1 as soon as possible. Review the migration guide for more details: https://github.com/deploymenttheory/terraform-provider-jamfpro/docs/resource-migration-guide.md",
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
				Description: "The unique identifier of the mobile group.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro mobile group.",
			},
			"site_id": sharedschemas.GetSharedSchemaSite(),
			"criteria": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the smart group search criteria. Can be from the Jamf built in enteries or can be an extension attribute.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "The priority of the criterion.",
						},
						"and_or": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Either 'and', 'or', or blank.",
							Default:      "and",
							ValidateFunc: validation.StringInSlice([]string{"", And, Or}, false),
						},
						"search_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "is",
							Description:  fmt.Sprintf("The type of smart group search operator. Allowed values are '%v'", getCriteriaOperators()),
							ValidateFunc: validation.StringInSlice(getCriteriaOperators(), false),
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Search value for the smart group criteria to match with.",
						},
						"opening_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Opening parenthesis flag used during smart group construction.",
						},
						"closing_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Closing parenthesis flag used during smart group construction.",
						},
					},
				},
			},
		},
	}
}
