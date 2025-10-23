package jamf_protect_plan

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrProvideEitherIDOrName = fmt.Errorf("please provide either 'id' or 'name', not both")
	ErrPlanNotFound          = fmt.Errorf("jamf Protect Plan not found using the provided identifier")
	ErrFailedToReadPlan      = fmt.Errorf("failed to read Jamf Protect Plan after retries")
)

// DataSourceJamfProtectPlan provides information about a Jamf Protect Plan
func DataSourceJamfProtectPlan() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the Jamf Protect Plan",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Jamf Protect Plan",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Jamf Protect Plan",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the Jamf Protect Plan",
			},
			"profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Profile ID associated with the plan",
			},
			"profile_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Profile name associated with the plan",
			},
			"scope_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the plan's scope",
			},
		},
	}
}
