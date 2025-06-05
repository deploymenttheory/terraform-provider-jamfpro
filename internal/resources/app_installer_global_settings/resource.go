package app_installer_global_settings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProAppInstallerGlobalSettings() *schema.Resource {
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
			"notification_message": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				Description: "Message to display as the push notification (up to 150 characters) when an app installer update is available. " +
					"This message applies to all App Installers deployments unless overridden in a specific app deployment.",
				ValidateFunc: validation.StringLenBetween(0, 150),
			},
			"notification_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
				Description: "How often the notification message is displayed in hours. " +
					"This applies globally to all App Installers deployments unless overridden in an individual deployment.",
			},
			"deadline_message": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				Description: "Message to display as the push notification (up to 150 characters) when the app is forced to quit to complete the update. " +
					"This applies globally to all App Installers deployments unless overridden in an individual deployment.",
				ValidateFunc: validation.StringLenBetween(0, 150),
			},
			"deadline": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
				Description: "How long the user can defer the update before the app is forced to quit to complete the update in hours. " +
					"This applies globally to all App Installers deployments unless overridden in an individual deployment.",
			},
			"quit_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
				Description: "Additional time for users to save work and close the app in minutes. " +
					"This applies globally to all App Installers deployments unless overridden in an individual deployment.",
			},
			"complete_message": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				Description: "Message to display as the push notification (up to 150 characters) when an update is complete. " +
					"This applies globally to all App Installers deployments unless overridden in an individual deployment.",
				ValidateFunc: validation.StringLenBetween(0, 150),
			},
			"relaunch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  nil,
				Description: "If true, the app will be relaunched after installation. " +
					"This applies globally to all App Installers deployments unless overridden in an individual deployment.",
			},
			"suppress": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  nil,
				Description: "If true, suppresses all user-facing notifications globally. " +
					"This applies globally to all App Installers deployments unless overridden in an individual deployment.",
			},
		},
	}
}
