// macosconfigurationprofilesplist_resource.go
package macosconfigurationprofilesplist

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProMacOSConfigurationProfilesPlist defines the schema and CRUD operations for managing Jamf Pro macOS Configuration Profiles in Terraform.
func ResourceJamfProMacOSConfigurationProfilesPlist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProMacOSConfigurationProfilesPlistCreate,
		ReadContext:   resourceJamfProMacOSConfigurationProfilesPlistRead,
		UpdateContext: resourceJamfProMacOSConfigurationProfilesPlistUpdate,
		DeleteContext: resourceJamfProMacOSConfigurationProfilesPlistDelete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the macOS configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Jamf UI name for configuration profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the configuration profile.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The universally unique identifier for the profile.",
			},
			"site": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The site to which the configuration profile is scoped.",
				Optional:    true,
				Default:     nil,
				Elem:        sharedschemas.GetSharedSchemaSite(),
			},
			"category": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The category to which the configuration profile is scoped.",
				Optional:    true,
				Elem:        sharedschemas.GetSharedSchemaCategory(),
			},
			"distribution_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Install Automatically",
				Description:  "The distribution method for the configuration profile. ['Make Available in Self Service','Install Automatically']",
				ValidateFunc: validation.StringInSlice([]string{"Make Available in Self Service", "Install Automatically"}, false),
			},
			"user_removable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the configuration profile is user removeable or not.",
			},
			"level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "System",
				Description:  "The deployment level of the configuration profile. Available options are: 'User' or 'System'. Note: 'System' is mapped to 'Computer Level' in the Jamf Pro GUI.",
				ValidateFunc: validation.StringInSlice([]string{"User", "System"}, false),
			},
			"payloads": {
				Type:             schema.TypeString,
				Required:         true,
				StateFunc:        plist.NormalizePayloadState,
				ValidateFunc:     plist.ValidatePayload,
				DiffSuppressFunc: DiffSuppressPayloads,
				Description:      "A MacOS configuration profile as a plist-formatted XML string.",
			},
			"redeploy_on_update": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Newly Assigned", // This is always "Newly Assigned" on existing profile objects, but may be set "All" on profile update requests and in TF state.
				Description: "Defines the redeployment behaviour when a mobile device config profile update occurs.This is always 'Newly Assigned' on new profile objects, but may be set 'All' on profile update requests and in TF state",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					if v == "All" || v == "Newly Assigned" {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be either 'All' or 'Newly Assigned', got: %s", key, v))
					return warns, errs
				},
			},
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The scope of the configuration profile.",
				Required:    true,
				Elem:        sharedschemas.GetSharedmacOSComputerSchemaScope(),
			},
			"self_service": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Self Service Configuration",
				Optional:    true,
				Default:     nil,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"install_button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text shown on Self Service install button",
							Default:     "Install",
						},
						"self_service_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description shown in Self Service",
							Default:     nil,
						},
						"force_users_to_view_description": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Forces users to view the description",
						},
						"feature_on_main_page": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Shows Configuration Profile on Self Service main page",
						},
						"notification": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enables Notification for this profile in self service",
						},
						"notification_subject": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "no message subject set",
							Description: "Message Subject",
						},
						"notification_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Message body",
						},
						// "self_service_icon": {
						// 	Type:        schema.TypeList,
						// 	MaxItems:    1,
						// 	Description: "Self Service icon settings",
						// 	Optional:    true,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"id":       {},
						// 			"uri":      {},
						// 			"data":     {},
						// 			"filename": {},
						// 		},
						// 	},
						// }, // TODO fix this broken crap later
						"self_service_categories": {
							Type:        schema.TypeList,
							Description: "Self Service category options",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Description: "ID of category",
										Optional:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Name of category",
										Optional:    true,
									},
									"display_in": {
										Type:        schema.TypeBool,
										ForceNew:    true,
										Description: "Display this profile in this category?",
										Required:    true,
									},
									"feature_in": {
										Type:        schema.TypeBool,
										Description: "Feature this profile in this category?",
										ForceNew:    true,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
