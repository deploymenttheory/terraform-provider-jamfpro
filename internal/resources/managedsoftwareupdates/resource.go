package managedsoftwareupdates

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProManagedSoftwareUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: mainCustomDiffFunc,
		Schema: map[string]*schema.Schema{
			"plan_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the managed software update plan.",
			},
			"group": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the Jamf Pro device group for the update plan.",
						},
						"object_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"COMPUTER_GROUP", "MOBILE_DEVICE_GROUP"}, false),
							Description:  "The type of the group (COMPUTER_GROUP or MOBILE_DEVICE_GROUP).",
						},
					},
				},
			},
			"device": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the individual device for the update plan.",
						},
						"object_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"COMPUTER", "MOBILE_DEVICE", "APPLE_TV"}, false),
							Description:  "The device type that the device_id refers to (COMPUTER, MOBILE_DEVICE, or APPLE_TV).",
						},
					},
				},
			},

			// Root level attributes, previously in the config block
			"update_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DOWNLOAD_ONLY", "DOWNLOAD_INSTALL", "DOWNLOAD_INSTALL_ALLOW_DEFERRAL", "DOWNLOAD_INSTALL_RESTART", "DOWNLOAD_INSTALL_SCHEDULE", "UNKNOWN"}, false),
				Description:  "The software update action to perform.",
			},
			"version_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LATEST_MAJOR", "LATEST_MINOR", "LATEST_ANY", "SPECIFIC_VERSION", "CUSTOM_VERSION", "UNKNOWN"}, false),
				Description:  "The type of version to update to.",
			},
			"specific_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"15.0", "14.7", "14.6.1", "14.6", "14.5", "13.7", "13.6.9", "13.6.8", "13.6.7", "12.7.6", "12.7.5", "11.7.10", "NO_SPECIFIC_VERSION"}, false),
				Description:  "Optional. Indicates the specific version to update to. Only available when the version type is set to specific version or custom version, otherwise defaults to NO_SPECIFIC_VERSION.",
			},
			"build_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional. Indicates the build version to update to. Only available when the version type is set to CUSTOM_VERSION.",
			},
			"max_deferrals": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  "Required when the provided update_action is DOWNLOAD_INSTALL_ALLOW_DEFERRAL, not applicable to all managed software update plans.",
			},
			"force_install_local_date_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional. Indicates the local date and time of the device to force update by.",
			},
		},
	}
}
