package appinstallers

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceJamfProAppInstallers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				Description:   "The unique identifier of the app installer deployment.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Description:   "The name of the app installer deployment.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the deployment is enabled.",
			},
			"app_title_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The jamf pro app installer ID of the app title.",
			},
			"deployment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Initial distribution method for distributing the app.",
			},
			"update_behavior": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The method to use for all future app updates.",
			},
			"category_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the category assigned to the app installer.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the site.",
			},
			"smart_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the smart group to scope the Jamf Pro App installer.",
			},
			"install_predefined_config_profiles": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether Jamf automatically installs necessary configuration profiles.",
			},
			"title_available_in_ais": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the title is available in AIS.",
			},
			"trigger_admin_notifications": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether log event notifications are enabled for this app.",
			},
			"notification_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "End User Experience notification settings for the deployment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notification_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The notification message.",
						},
						"notification_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The interval for notifications.",
						},
						"deadline_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The deadline message.",
						},
						"deadline": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Update deadline in hours.",
						},
						"quit_delay": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Force quit grace period in minutes.",
						},
						"complete_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Update complete message.",
						},
						"relaunch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to Automatically open app after installation.",
						},
						"suppress": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Suppressing notifications.",
						},
					},
				},
			},
			"self_service_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Self-service settings for the deployment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_in_featured_category": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to include in the featured category.",
						},
						"include_in_compliance_category": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Include the app in the Featured category.",
						},
						"force_view_description": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Force users to view the description before installing.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description to display for the app in Self Service.",
						},
						"categories": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Categories in which to display or feature the app.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the category.",
									},
									"featured": {
										Type:        schema.TypeBool,
										Computed:    true,
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Check if either ID or name is provided
	resourceID := d.Get("id").(string)
	name := d.Get("name").(string)

	if resourceID == "" && name == "" {

		return diag.FromErr(fmt.Errorf("either id or name must be provided"))
	}

	var appInstaller *jamfpro.ResourceJamfAppCatalogDeployment
	var err error

	if resourceID != "" {
		// Get deployment directly by ID
		appInstaller, err = client.GetJamfAppCatalogAppInstallerByID(resourceID)
		if err != nil {

			return diag.FromErr(fmt.Errorf("failed to fetch Jamf Pro App Installer Deployment by ID %s: %v", resourceID, err))
		}
	} else {
		// Get deployment by name
		appInstaller, err = client.GetJamfAppCatalogAppInstallerByName(name)
		if err != nil {

			return diag.FromErr(fmt.Errorf("failed to fetch Jamf Pro App Installer Deployment by name %s: %v", name, err))
		}
	}

	// Set ID from the deployment
	d.SetId(appInstaller.ID)

	// Update state using the same function as the resource
	return append(diags, updateState(d, appInstaller)...)
}
