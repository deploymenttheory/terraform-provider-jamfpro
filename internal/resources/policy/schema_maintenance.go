package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getPolicySchemaMaintenance() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recon": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to run recon (inventory update) as part of the maintenance. Forces computers to submit updated inventory information to Jamf Pro",
				Default:     false,
			},
			"reset_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to reset the computer name to the name stored in Jamf Pro. Changes the computer name on computers to match the computer name in Jamf Pro",
				Default:     false,
			},
			"install_all_cached_packages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to install all cached packages. Installs packages cached by Jamf Pro",
				Default:     false,
			},
			"heal": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to heal the policy.",
				Default:     false,
			},
			"prebindings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to update prebindings.",
				Default:     false,
			},
			"permissions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to fix Disk Permissions (Not compatible with macOS v10.12 or later)",
				Default:     false,
			},
			"byhost": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to fix ByHost files andnpreferences.",
				Default:     false,
			},
			"system_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to flush caches from /Library/Caches/ and /System/Library/Caches/, except for any com.apple.LaunchServices caches",
				Default:     false,
			},
			"user_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to flush caches from ~/Library/Caches/, ~/.jpi_cache/, and ~/Library/Preferences/Microsoft/Office version #/Office Font Cache. Enabling this may cause problems with system fonts displaying unless a restart option is configured.",
				Default:     false,
			},
			"verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to verify system files and structure on the Startup Disk",
				Default:     false,
			},
		},
	}

	return out
}
