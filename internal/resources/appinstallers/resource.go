package appinstallers

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProAppCatalogDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the app catalog deployment.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the app catalog deployment.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the deployment is enabled.",
			},
			"app_title_id": { //TODO will probably handle this within the constructor
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the app title.",
			},
			"deployment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTOMATIC", "SELF_SERVICE"}, false),
				Description:  "The type of deployment (AUTOMATIC or SELF_SERVICE).",
			},
			"update_behavior": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTOMATIC", "MANUAL"}, false),
				Description:  "The update behavior (AUTOMATIC or MANUAL).",
			},
			"category_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the category.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the site.",
			},
			"smart_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the smart group.",
			},
			"install_predefined_config_profiles": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to install predefined configuration profiles.",
			},
			"title_available_in_ais": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the title is available in AIS.",
			},
			"trigger_admin_notifications": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to trigger admin notifications.",
			},
			"notification_settings": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Notification settings for the deployment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notification_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The notification message.",
						},
						"notification_interval": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The interval for notifications.",
						},
						"deadline_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The deadline message.",
						},
						"deadline": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The deadline in minutes.",
						},
						"quit_delay": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The quit delay in minutes.",
						},
						"complete_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The completion message.",
						},
						"relaunch": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to relaunch after installation.",
						},
						"suppress": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Suppression settings.",
						},
					},
				},
			},
			"self_service_settings": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Self-service settings for the deployment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_in_featured_category": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to include in the featured category.",
						},
						"include_in_compliance_category": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to include in the compliance category.",
						},
						"force_view_description": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to force viewing the description.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The self-service description.",
						},
						"categories": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Categories for the self-service deployment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The ID of the category.",
									},
									"featured": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether the category is featured.",
									},
								},
							},
						},
					},
				},
			},
			"selected_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The selected version of the app.",
			},
			"latest_available_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest available version of the app.",
			},
			"version_removed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the version has been removed.",
			},
		},
	}
}
