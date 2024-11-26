// computercheckin_resource.go
package computercheckin

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProComputerCheckindefines the schema and RU operations for managing Jamf Pro computer checkin configuration in Terraform.
func ResourceJamfProComputerCheckin() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			return validateComputerCheckinDependencies(d)
		},
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
			"check_in_frequency": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The frequency in minutes for computer check-in.",
				ValidateFunc: validation.IntInSlice([]int{60, 30, 15, 5}),
			},
			"create_startup_script": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines if a startup script should be created.",
			},
			"log_startup_event": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines if startup events should be logged.",
			},
			"check_for_policies_at_startup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, ensure that computers check for policies triggered by startup",
			},
			"apply_computer_level_managed_preferences": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Applies computer level managed preferences. Setting appears to be hard coded to false and cannot be changed. Thus field is set to computed.",
			},
			"ensure_ssh_is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable SSH (Remote Login) on computers that have it disabled.",
			},
			"create_login_logout_hooks": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines if login/logout hooks should be created. Create events that trigger each time a user logs in",
			},
			"log_username": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Log Computer Usage information at login. Log the username and date/time at login.",
			},
			"check_for_policies_at_login_logout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Checks for policies at login and logout.",
			},
			"apply_user_level_managed_preferences": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Applies user level managed preferences. Setting appears to be hard coded to false and cannot be changed. Thus field is set to computed.",
			},
			"hide_restore_partition": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines if the restore partition should be hidden. Setting appears to be hard coded to false and cannot be changed. Thus field is set to computed.",
			},
			"perform_login_actions_in_background": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Performs login actions in the background. Setting appears to be hard coded to false and cannot be changed. Thus field is set to computed.",
			},
		},
	}
}
