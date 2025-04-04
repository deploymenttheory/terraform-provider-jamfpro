package accountdrivenuserenrollmentsettings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceJamfProAccountDrivenUserEnrollmentSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable Account Driven User Enrollment Session Token Settings.",
			},
			"expiration_interval_days": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"expiration_interval_seconds"},
				Description:   "Configure the number of days to prompt users to re-authenticate on devices enrolled using account-driven User Enrollment.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// If seconds is explicitly set, suppress changes to days
					if v, ok := d.GetOk("expiration_interval_seconds"); ok && v.(int) > 0 {
						return true
					}
					return false
				},
			},
			"expiration_interval_seconds": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"expiration_interval_days"},
				Description:   "Configure the number of seconds to prompt users to re-authenticate on devices enrolled using account-driven User Enrollment. Must be in scienctific notation. e.g 'scientific' for 30 days.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// If days is explicitly set, suppress changes to seconds
					if v, ok := d.GetOk("expiration_interval_days"); ok && v.(int) > 0 {
						return true
					}
					return false
				},
			},
		},
	}
}
