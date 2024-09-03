package appinstallers

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProAppInstallers() *schema.Resource {
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
		CustomizeDiff: validateAppCatalogDeploymentName,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the app installer deployment.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the app installer deployment. This name cannot be freeform text like in the gui as the name is used to infur the automatically appTitleId field.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the deployment is enabled.",
			},
			"app_title_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The jamf pro app installer ID of the app title.",
			},
			"deployment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"INSTALL_AUTOMATICALLY", "SELF_SERVICE"}, false),
				Description:  "Initial distribution method to use for distributing the app to a computer for the initial installation. (INSTALL_AUTOMATICALLY or SELF_SERVICE).",
			},
			"update_behavior": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTOMATIC", "MANUAL"}, false),
				Description:  "The method to use for all future app updates, regardless of the initial distribution method. (AUTOMATIC or MANUAL).",
			},
			"category_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the category to assign to the app installer. Use -1 if not required.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the site. Use -1 if not required.",
			},
			"smart_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the smart group to scope the Jamf Pro App installer. Default is '1' - All Managed Clients. -1 is not an option.",
			},
			"install_predefined_config_profiles": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allows Jamf to automatically install necessary configuration profiles to support this App Installer. When unselected, you may need to create configuration profiles for some software titles. https://learn.jamf.com/en-US/bundle/technical-articles/page/Configuration_Profiles_for_Additional_App_Installers_Settings.html",
			},
			"title_available_in_ais": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the title is available in AIS.",
			},
			"trigger_admin_notifications": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Log event notifications for this app. Opt in to receiving notifications for certain events including app updates and installation failures. Account Notification settings. https://your-instance.jamfcloud.com/notifications.html?id=0&o=r",
			},
			"notification_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "End User Experience notification settings for the deployment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notification_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The notification message.",
						},
						"notification_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The interval for notifications.",
						},
						"deadline_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The deadline message.",
						},
						"deadline": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The Update deadline in hours.",
						},
						"quit_delay": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Force quit grace period in minutes.",
						},
						"complete_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Update complete message.",
						},
						"relaunch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to Automatically open app after installation.",
						},
						"suppress": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Suppressing notifications.",
						},
					},
				},
			},
			"self_service_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Self-service settings for the deployment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_in_featured_category": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to include in the featured category.",
						},
						"include_in_compliance_category": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Include the app in the Featured category.Jamf Pro must be integrated with Microsoft Intune to include the app in the Compliance category. Confirm the integration is enabled. If you previously integrated Microsoft Intune using Conditional Access, disregard this alert.",
						},
						"force_view_description": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Force users to view the description before installing the app.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description (up to 4000 characters) to display for the app in Self Service.",
						},
						"categories": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Categories in which to display or feature the app in Self Service.",
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
				Computed:    true,
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
