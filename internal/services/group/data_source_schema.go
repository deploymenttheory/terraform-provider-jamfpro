package group

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	errMustProvideOne    = fmt.Errorf("one of 'name', 'group_platform_id', or 'group_jamfpro_id' must be provided")
	errNameAndIDConflict = fmt.Errorf("'name' and 'group_jamfpro_id' cannot both be specified")
	errGroupTypeRequired = fmt.Errorf("'group_type' must be specified when using 'name' or 'group_jamfpro_id'")
	errGroupTypeAllowed  = fmt.Errorf("must be either 'COMPUTER' or 'MOBILE'")
)

// DataSourceJamfProGroups provides information about a specific group in Jamf Pro.
func DataSourceJamfProGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"group_jamfpro_id"},
				Description:   "The name of the group. Mutually exclusive with group_jamfpro_id. Requires group_type.",
			},
			"group_platform_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Platform ID of the group.",
			},
			"group_jamfpro_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				Description:   "Jamf Pro ID of the group. Mutually exclusive with name. Requires group_type.",
			},
			"group_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Brief explanation of the content or purpose of the group.",
			},
			"group_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Type of the group. Required if name or group_jamfpro_id is specified. Must be either COMPUTER or MOBILE.",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "" && v != "COMPUTER" && v != "MOBILE" {
						errs = append(errs, errGroupTypeAllowed)
					}
					return
				},
			},
			"smart": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the group is a smart group.",
			},
			"membership_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of members in the group.",
			},
		},
	}
}
