package client_checkin

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProClientCheckin() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
			return validateComputerCheckinDependencies(d)
		},
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
			"check_in_frequency": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The frequency at which computers check in with Jamf Pro for available policies. Valid values are 5, 15, 30, or 60.",
				ValidateFunc: validation.IntInSlice([]int{60, 30, 15, 5}),
			},
			"create_hooks": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if login/logout hooks should be created.",
			},
			"hook_log": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if login/logout events should be logged.",
			},
			"hook_policies": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if policies should be checked at login/logout.",
			},
			"create_startup_script": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if a startup script should be created.",
			},
			"startup_log": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if startup events should be logged.",
			},
			"startup_policies": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if policies should be checked at startup.",
			},
			"startup_ssh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable SSH (Remote Login) on computers that have it disabled at startup.",
			},
			"enable_local_configuration_profiles": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable local configuration profiles management.",
			},
			"allow_network_state_change_triggers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Check for policies with a 'Network State Change' trigger when a network change occurs, such as a network connection change, a computer name change, or an IP address change.",
			},
		},
	}
}
