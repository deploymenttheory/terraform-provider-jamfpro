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
		Schema: map[string]*schema.Schema{
			"plan_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the managed software update plan.",
			},
			"devices": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the device for the update plan.",
						},
						"object_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"COMPUTER", "MOBILE_DEVICE", "APPLE_TV"}, false),
							Description:  "The type of the device (COMPUTER, MOBILE_DEVICE, or APPLE_TV).",
						},
					},
				},
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
							Description: "The ID of the group for the update plan.",
						},
						"object_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"COMPUTER", "MOBILE_DEVICE", "APPLE_TV"}, false),
							Description:  "The type of the group (COMPUTER, MOBILE_DEVICE, or APPLE_TV).",
						},
					},
				},
			},
			"config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"update_action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"DOWNLOAD_ONLY", "DOWNLOAD_INSTALL", "DOWNLOAD_INSTALL_ALLOW_DEFERRAL", "DOWNLOAD_INSTALL_RESTART", "DOWNLOAD_INSTALL_SCHEDULE", "UNKNOWN"}, false),
							Description:  "The update action to perform.",
						},
						"version_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"LATEST_MAJOR", "LATEST_MINOR", "LATEST_ANY", "SPECIFIC_VERSION", "CUSTOM_VERSION", "UNKNOWN"}, false),
							Description:  "The type of version to update to.",
						},
						"specific_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The specific version to update to, if applicable. Required when version_type is SPECIFIC_VERSION or CUSTOM_VERSION.",
						},
						"build_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The build version to update to. Only applicable when version_type is CUSTOM_VERSION.",
						},
						"max_deferrals": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Description:  "The maximum number of times the update can be deferred. Required when update_action is DOWNLOAD_INSTALL_ALLOW_DEFERRAL.",
						},
						"force_install_local_date_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The local date and time to force the installation.",
						},
					},
				},
			},
		},
	}
}
