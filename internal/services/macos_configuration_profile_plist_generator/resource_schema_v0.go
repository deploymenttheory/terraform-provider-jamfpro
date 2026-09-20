package macos_configuration_profile_plist_generator

import (
	"fmt"
	"time"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceV0 preserves the complete schema before collection sets were introduced.
func resourceV0() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProMacOSConfigurationProfilesPlistGeneratorCreate,
		ReadContext:   resourceJamfProMacOSConfigurationProfilesPlistGeneratorReadWithCleanup,
		UpdateContext: resourceJamfProMacOSConfigurationProfilesPlistGeneratorUpdate,
		DeleteContext: resourceJamfProMacOSConfigurationProfilesPlistGeneratorDelete,
		CustomizeDiff: mainCustomDiffFunc,
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return utils.NormalizeWhitespace(old) == utils.NormalizeWhitespace(new)
				},
				DiffSuppressOnRefresh: true,
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
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payload_description_header": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Description of the payload at the header level of the plist. This provides a human-readable explanation of what the overall profile is intended to do or configure.",
						},
						"payload_display_name_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The display name of the payload at the header level of the plist. This is shown in user interfaces to identify the overall profile to users and administrators. Jamf Pro matches this to the name of the configuation profile, 'name' at the top of the schema.",
						},
						"payload_enabled_header": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether the payload is enabled at the header level of the plist. If set to false, the overall profile will be disabled.",
						},
						"payload_identifier_header": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique identifier for the payload within the MDM profile at the header level of the plist. This identifier is used to track and reference the overall profile uniquely.",
						},
						"payload_organization_header": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The organization associated with the payload at the header level of the plist. This represents the entity that created or is responsible for the overall profile.",
						},
						"payload_type_header": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of the config profile payload at the header level of the plist. This indicates what kind of settings or configurations the overall profile applies.",
						},
						"payload_uuid_header": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID for the payload within the MDM profile at the header level of the plist. This ensures the uniqueness of the overall profile.",
						},
						"payload_version_header": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The version of the payload at the header level of the plist. This helps in identifying the version of the overall profile settings or configurations being applied.",
						},
						"payload_removal_disallowed_header": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates whether the removal of the payload is disallowed. If set to true, the MDM profile cannot be removed by users.",
						},
						"payload_scope_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The scope of the payload at the header level of the plist. This defines the context in which the overall profile settings are applied, can be either 'System' or 'User'.",
						},
						"payload_content": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The payload content of the macOS configuration profile plist. Multiple payloads can be defined as needed.Defined as key value pairs and supports nested dictionaries.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"setting": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The key and value setting items of the macOS configuration profile plist",
										Elem: &schema.Resource{
											Schema: payloadContentSchemaV0().Schema,
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
										Description: "Identifier for the payload.A GUID.",
									},
									"payload_organization": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Organization associated with the payload.",
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
									"payload_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the config profile payload.",
									},
									"payload_uuid": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "UUID of the payload.",
									},
									"payload_version": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Version of the payload.",
									},
								},
							},
						},
					},
				},
			},
			"redeploy_on_update": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Defines the redeployment behaviour when an update to a macOS config profile" +
					"occurs. This is always 'Newly Assigned' on new profile objects, but may be set to 'All'" +
					"on profile update requests once the configuration profile has been deployed to at least" +
					" one device.",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
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
						"self_service_display_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Display name for the profile in Self Service (Self Service 10.0.0 or later)",
						},
						"install_button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name for the button that users click to install the profile",
						},
						"self_service_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description to display for the profile in Self Service",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return utils.NormalizeWhitespace(old) == utils.NormalizeWhitespace(new)
							},
							DiffSuppressOnRefresh: true,
						},
						"force_users_to_view_description": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Force users to view the description before the profile installs",
						},
						"feature_on_main_page": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Shows Configuration Profile on Self Service main page",
						},
						"notification": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "TEMPORARILY DISABLED\nEnables Notification for this profile in self service",
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
						"self_service_icon_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Icon for policy to use in self-service. Can be used in conjection with the icons resource",
						},
						"self_service_category": {
							Type:        schema.TypeSet,
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

// Define a finite level of nested dictionaries
func nestedDictionarySchemaV0(level int) *schema.Schema {
	if level <= 0 {
		return &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "A nested dictionary structure for xml plist definition.",
			Elem:        schema.TypeString,
		}
	}
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "A nested dictionary structure.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The key for the dictionary entry.",
				},
				"value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The value for the dictionary entry.",
				},
				"dictionary": nestedDictionarySchemaV0(level - 1),
			},
		},
	}
}

// Define the payload content schema with limited depth
func payloadContentSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key for the xml plist entry.",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The value for the xml plist entry.",
			},
			"dictionary": nestedDictionarySchemaV0(6),
		},
	}
}
