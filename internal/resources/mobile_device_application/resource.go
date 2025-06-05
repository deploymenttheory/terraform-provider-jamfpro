package mobile_device_application

import (
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProMobileDeviceApplication defines the schema and CRUD operations for managing Jamf Pro Mobile Device Applications in Terraform
func ResourceJamfProMobileDeviceApplication() *schema.Resource {
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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the mobile device application.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the mobile device application.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the mobile device application.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the mobile device application.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bundle identifier of the application.",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version of the application.",
			},
			"internal_app": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if this is an internal application.",
			},
			"ipa": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the IPA file.",
						},
						"uri": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URI of the IPA file.",
						},
						"data": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The base64 encoded IPA data.",
						},
					},
				},
			},
			"icon": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the icon.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the icon file.",
						},
						"uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URI of the icon.",
						},
					},
				},
			},
			"mobile_device_provisioning_profile": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The mobile device provisioning profile ID.",
			},
			"itunes_store_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The iTunes Store URL for the application.",
			},
			"make_available_after_install": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Make the application available after installation.",
			},
			"itunes_country_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The iTunes country/region for the application.",
			},
			"itunes_sync_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The iTunes sync time.",
			},
			"deployment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deployment type for the application.",
			},
			"deploy_automatically": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Deploy the application automatically.",
			},
			"deploy_as_managed_app": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Deploy as a managed application.",
			},
			"remove_app_when_mdm_profile_is_removed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Remove the application when the MDM profile is removed.",
			},
			"prevent_backup_of_app_data": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Prevent backup of application data.",
			},
			"allow_user_to_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow users to delete the application.",
			},
			"require_network_tethered": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Require network tethering for the application.",
			},
			"keep_description_and_icon_up_to_date": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Keep the description and icon up to date.",
			},
			"keep_app_updated_on_devices": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Keep the application updated on devices.",
			},
			"free": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if the application is free.",
			},
			"take_over_management": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Take over management of the application.",
			},
			"host_externally": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Host the application externally.",
			},
			"external_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The external URL for the application.",
			},
			"site_id":     sharedschemas.GetSharedSchemaSite(),
			"category_id": sharedschemas.GetSharedSchemaCategory(),
			"self_service": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"self_service_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The self service description.",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return normalizeWhitespace(old) == normalizeWhitespace(new)
							},
						},
						"self_service_install_button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text displayed on the install button in self service.",
						},
						"feature_on_main_page": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Feature this application on the main page.",
						},
						"notification": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enable notifications for this application.",
						},
						"notification_subject": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The subject of the notification.",
						},
						"notification_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The message of the notification.",
						},
						"self_service_icon": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The ID of the self service icon.",
									},
									"filename": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The filename of the self service icon.",
									},
									"uri": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URI of the self service icon.",
									},
								},
							},
						},
					},
				},
			},
			"vpp": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assign_vpp_device_based_licenses": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Assign VPP device-based licenses.",
						},
						"vpp_admin_account_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The VPP admin account ID.",
						},
					},
				},
			},
			"app_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"preferences": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The preferences configuration for the application.",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return normalizeWhitespace(old) == normalizeWhitespace(new)
							},
						},
					},
				},
			},
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The scope of the mobile device application.",
				Required:    true,
				Elem:        sharedschemas.GetSharedMobileDeviceSchemaScope(),
			},
		},
	}
}
