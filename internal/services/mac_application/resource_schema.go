package mac_application

import (
	"time"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProMacApplication defines the schema and CRUD operations for managing Jamf Pro Mac Applications in Terraform
func ResourceJamfProMacApplication() *schema.Resource {
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
				Description: "The unique identifier of the mac application.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the mac application.",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version of the application.",
			},
			"bundle_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bundle identifier of the application.",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the application.",
			},
			"is_free": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if the application is free.",
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Install Automatically/Prompt Users to Install",
					"Make Available in Self Service",
				}, false),
				Description: "The deployment type for the application. Valid values are 'Install Automatically/Prompt Users to Install' or 'Make Available in Self Service'."},
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
						"install_button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text displayed on the install button in self service.",
						},
						"force_users_to_view_description": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Force users to view the description.",
						},
						"feature_on_main_page": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Feature this application on the main page.",
						},
						"notification": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The notification setting for this application.",
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
						"self_service_category": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The ID of the self service category.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the self service category.",
									},
									"display_in": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Display in this category.",
									},
									"feature_in": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Feature in this category.",
									},
								},
							},
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
									"data": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The data of the self service icon.",
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
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The scope of the mac application.",
				Required:    true,
				Elem:        sharedschemas.GetSharedmacOSComputerSchemaScope(),
			},
		},
	}
}
