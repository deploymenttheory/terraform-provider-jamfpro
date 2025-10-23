package static_mobile_device_group

import (
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// ResourceJamfProStaticMobileDeviceGroups defines the schema and CRUD operations for managing Jamf Pro static Mobile Device Groups in Terraform.
func ResourceJamfProStaticMobileDeviceGroups() *schema.Resource {
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
				Description: "The unique identifier of the Jamf Pro static mobile device group.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro static mobile device group.",
			},
			"is_smart": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Computed value indicating whether the mobile device group is smart or static.",
			},
			"site_id": sharedschemas.GetSharedSchemaSite(),
			"assigned_mobile_device_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "assigned mobile device by ids",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}
