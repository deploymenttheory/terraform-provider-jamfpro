// macosconfigurationprofilesplistgenerator_resource.go
package macosconfigurationprofilesplistgenerator

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProMacOSConfigurationProfilesPlistGenerator defines the schema and CRUD operations for managing Jamf Pro macOS Configuration Profiles in Terraform.
func ResourceJamfProMacOSConfigurationProfilesPlistGenerator() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProMacOSConfigurationProfilesPlistGeneratorCreate,
		ReadContext:   resourceJamfProMacOSConfigurationProfilesPlistGeneratorReadWithCleanup,
		UpdateContext: resourceJamfProMacOSConfigurationProfilesPlistGeneratorUpdate,
		DeleteContext: resourceJamfProMacOSConfigurationProfilesPlistGeneratorDelete,
		//CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(15 * time.Second),
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
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of payloads for the macOS configuration profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payload_root": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The root level of the payloads.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"payload_description_root": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of the payload.",
									},
									"payload_display_name_root": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the payload.",
									},
									"payload_enabled_root": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the payload is enabled.",
									},
									"payload_identifier_root": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifier for the payload.",
									},
									"payload_organization_root": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Organization associated with the payload.",
									},
									"payload_removal_disallowed_root": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the payload removal is disallowed.",
									},
									"payload_scope_root": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scope of the payload. Computed by what is set by level. 'System' or 'User'.",
									},
									"payload_type_root": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the config profile payload.",
									},
									"payload_uuid_root": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "UUID of the payload.",
									},
									"payload_version_root": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Version of the payload.",
									},
								},
							},
						},
						"payload_content": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Content of the payload.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"configuration": {
										Type:             schema.TypeList,
										Optional:         true,
										StateFunc:        plist.NormalizePayloadState,
										DiffSuppressFunc: DiffSuppressPayloads,
										Description:      "A list of key-value pairs for the macOS configuration profile payload.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The key for the plist entry.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The value for the plist entry.",
												},
											},
										},
									},
									"payload_description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of the payload.",
									},
									"payload_display_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Display name of the payload.",
									},
									"payload_enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the payload is enabled.",
									},
									"payload_identifier": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unique identifier for the payload within the mdm profile. Required for a valid request to be sent but then overwritten by the Jamf Pro server. This key changes every time a profile is updated.",
									},
									"payload_organization": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Organization associated with the payload.",
									},
									"payload_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the config profile payload.",
									},
									"payload_uuid": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "payload UUID for the payload within the mdm profile. Required for a valid request to be sent but then overwritten by the Jamf Pro server. This key changes every time a profile is updated.",
									},
									"payload_version": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Version of the payload.",
									},
									"payload_removal_disallowed": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the payload removal is disallowed.",
									},
									"payload_scope": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Scope of the payload. Computed by what is set by level. 'System' or 'User'.",
									},
								},
							},
						},
					},
				},
			},
			"redeploy_on_update": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Newly Assigned",
				Description: "Defines the redeployment behaviour when a mobile device config profile update occurs. This is always 'Newly Assigned' on new profile objects, but may be set 'All' on profile update requests and in TF state.",
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
