// macosconfigurationprofilesplist_resource.go
package macosconfigurationprofilesplist

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProMacOSConfigurationProfilesPlist defines the schema and CRUD operations for managing Jamf Pro macOS Configuration Profiles in Terraform.
func ResourceJamfProMacOSConfigurationProfilesPlist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProMacOSConfigurationProfilesPlistCreate,
		ReadContext:   resourceJamfProMacOSConfigurationProfilesPlistReadWithCleanup,
		UpdateContext: resourceJamfProMacOSConfigurationProfilesPlistUpdate,
		DeleteContext: resourceJamfProMacOSConfigurationProfilesPlistDelete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
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
			"site_id":     sharedschemas.GetSharedSchemaSite(),
			"category_id": sharedschemas.GetSharedSchemaCategory(),
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
				DiffSuppressFunc: DiffSuppressPayloads,
				Description:      "A MacOS configuration profile as a plist-formatted XML string.",
			},
			"payload_validate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Validates plist payload XML. Turn off to force malformed XML confguration." +
					"Required when the configuration profile is a non Jamf Pro source, e.g iMazing. Removing" +
					"this may cause unexpected stating behaviour.",
			},
			"redeploy_on_update": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Defines the redeployment behaviour when an update to a macOS config profile" +
					"occurs. This is always 'Newly Assigned' on new profile objects, but may be set to 'All'" +
					"on profile update requests once the configuration profile has been deployed to at least" +
					" one device.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v, ok := val.(string)
					if !ok {
						errs = append(errs, fmt.Errorf("%q must be a string, got: %T", key, val))
						return warns, errs
					}
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
						},
						"self_service_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description shown in Self Service",
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
							Description: "Message Subject",
						},
						"notification_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Message body",
						},
						// TODO
						"self_service_icon_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Icon for policy to use in self-service",
						},
						"self_service_category": {
							Type:        schema.TypeList,
							Description: "Self Service category options",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Description: "ID of category. Both ID and Name are required",
										Required:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Name of category. Both ID and Name are required",
										Computed:    true,
									},
									"display_in": {
										Type:        schema.TypeBool,
										Description: "Display this profile in this category?",
										Required:    true,
									},
									"feature_in": {
										Type:        schema.TypeBool,
										Description: "Feature this profile in this category?",
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
