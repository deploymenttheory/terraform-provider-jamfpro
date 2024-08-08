package policies

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getPolicySchemaReboot() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The reboot message displayed to the user.",
				// Default:     "This computer will restart in 5 minutes. Please save anything you are working on and log out by choosing Log Out from the bottom of the Apple menu.",
			},
			"specify_startup": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Reboot Method",
				// Default:     "",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validMethods := []string{"", "Standard Restart", "MDM Restart with Kernel Cache Rebuild"}
					for _, method := range validMethods {
						if v == method {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validMethods, v))
					return
				},
			},
			"startup_disk": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Disk to boot computers to",
				Default:     "Current Startup Disk",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validDisks := []string{"Current Startup Disk", "Currently Selected Startup Disk (No Bless)", "macOS Installer", "Specify Local Startup Disk"}
					for _, disk := range validDisks {
						if v == disk {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validDisks, v))
					return
				},
			},
			"no_user_logged_in": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Action to take if no user is logged in to the computer",
				// Default:     "Do not restart",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validOptions := []string{"Restart if a package or update requires it", "Restart Immediately", "Do not restart"}
					for _, option := range validOptions {
						if v == option {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validOptions, v))
					return
				},
			},
			"user_logged_in": {
				Type:     schema.TypeString,
				Optional: true,
				// Default:     "Do not restart",
				Description: "Action to take if a user is logged in to the computer",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validOptions := []string{"Restart if a package or update requires it", "Restart Immediately", "Restart", "Do not restart"}
					for _, option := range validOptions {
						if v == option {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validOptions, v))
					return
				},
			},
			"minutes_until_reboot": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Amount of time to wait before the restart begins.",
				// Default:     5,
			},
			"start_reboot_timer_immediately": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defines if the reboot timer should start immediately once the policy applies to a macOS device.",
				// Default:     false,
			},
			"file_vault_2_reboot": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Perform authenticated restart on computers with FileVault 2 enabled. Restart FileVault 2-encrypted computers without requiring an unlock during the next startup",
				// Default:     false,
			},
		}}
	return out
}
