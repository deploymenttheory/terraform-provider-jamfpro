// macosconfigurationprofilesplistgenerator_resource.go
package macos_configuration_profile_plist_generator

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/* --------- A mapping of how terraform schema correlates to plist structure ---------

<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <!-- PayloadContent corresponds to payload_content in the schema -->
        <key>PayloadContent</key>
        <array>
            <dict>
                <!-- Example payload 1 -->
                <key>AllowUserOverrides</key> <!-- hcl schema: payload_content.setting.key -->
                <true/> <!-- hcl schema: payload_content.setting.value -->
                <key>AllowedSystemExtensions</key> <!-- hcl schema: payload_content.setting.key -->
                <dict> <!-- hcl schema: payload_content.setting.dictionary -->
                    <key>H8P3P53Q9W</key> <!-- hcl schema: payload_content.setting.key -->
                    <array>
                        <string>com.axissecurity.client.com-axissecurity-client-SystemNetworkExtension</string> <!-- hcl schema: payload_content.setting.value -->
                    </array>
                </dict>
                <!-- Payload-level metadata fields in the plist -->
                <key>PayloadDescription</key> <!-- hcl schema: payload_description -->
                <string/>
                <key>PayloadDisplayName</key> <!-- hcl schema: payload_display_name -->
                <string>System Extension Policy</string>
                <key>PayloadEnabled</key> <!-- hcl schema: payload_enabled -->
                <true/>
                <key>PayloadIdentifier</key> <!-- hcl schema: payload_identifier -->
                <string>com.apple.system-extension-policy.70B93937-265D-431B-9DF6-A7E031A368EF</string>
                <key>PayloadOrganization</key> <!-- hcl schema: payload_organization -->
                <string>Deployment Theory</string>
                <key>PayloadType</key> <!-- hcl schema: payload_type -->
                <string>com.apple.system-extension-policy</string>
                <key>PayloadUUID</key> <!-- hcl schema: payload_uuid -->
                <string>EF513BF0-9C22-4FBE-9559-7EE838CE7AFC</string>
                <key>PayloadVersion</key> <!-- hcl schema: payload_version -->
                <integer>1</integer>
            </dict>
            <dict>
                <!-- Example payload 2 -->
                <key>NotificationSettings</key> <!-- hcl schema: payload_content.setting.key -->
                <array>
                    <dict>
                        <key>AlertType</key> <!-- hcl schema: payload_content.setting.key -->
                        <integer>2</integer> <!-- hcl schema: payload_content.setting.value -->
                        <key>BadgesEnabled</key> <!-- hcl schema: payload_content.setting.key -->
                        <true/> <!-- hcl schema: payload_content.setting.value -->
                        <key>BundleIdentifier</key> <!-- hcl schema: payload_content.setting.key -->
                        <string>com.axissecurity.client.ui</string> <!-- hcl schema: payload_content.setting.value -->
                        <key>CriticalAlertEnabled</key> <!-- hcl schema: payload_content.setting.key -->
                        <true/> <!-- hcl schema: payload_content.setting.value -->
                        <key>NotificationsEnabled</key> <!-- hcl schema: payload_content.setting.key -->
                        <true/> <!-- hcl schema: payload_content.setting.value -->
                        <key>ShowInLockScreen</key> <!-- hcl schema: payload_content.setting.key -->
                        <true/> <!-- hcl schema: payload_content.setting.value -->
                        <key>ShowInNotificationCenter</key> <!-- hcl schema: payload_content.setting.key -->
                        <true/> <!-- hcl schema: payload_content.setting.value -->
                        <key>SoundsEnabled</key> <!-- hcl schema: payload_content.setting.key -->
                        <true/> <!-- hcl schema: payload_content.setting.value -->
                    </dict>
                </array>
                <!-- Payload-level metadata fields in the plist -->
                <key>PayloadDisplayName</key> <!-- hcl schema: payload_display_name -->
                <string>Notifications Payload</string>
                <key>PayloadIdentifier</key> <!-- hcl schema: payload_identifier -->
                <string>BFA5BB51-886B-4DB9-9A3C-AF67FB627F7A</string>
                <key>PayloadOrganization</key> <!-- hcl schema: payload_organization -->
                <string>JAMF Software</string>
                <key>PayloadType</key> <!-- hcl schema: payload_type -->
                <string>com.apple.notificationsettings</string>
                <key>PayloadUUID</key> <!-- hcl schema: payload_uuid -->
                <string>531AC0A1-87CE-498B-8AFC-14898BEC84B3</string>
                <key>PayloadVersion</key> <!-- hcl schema: payload_version -->
                <integer>1</integer>
            </dict>
        </array>
        <!-- Root-level 'header' metadata fields in the plist -->
        <key>PayloadDescription</key> <!-- hcl schema: payload_description_header -->
        <string/>
        <key>PayloadDisplayName</key> <!-- hcl schema: payload_display_name_header -->
        <string>dt-mcp-axis_security_ext-0.0.1-prod-eu-0-0</string>
        <key>PayloadEnabled</key> <!-- hcl schema: payload_enabled_header -->
        <true/>
        <key>PayloadIdentifier</key> <!-- hcl schema: payload_identifier_header -->
        <string>com.axissecurity.client.profile</string>
        <key>PayloadOrganization</key> <!-- hcl schema: payload_organization_header -->
        <string>Deployment Theory</string>
        <key>PayloadRemovalDisallowed</key> <!-- hcl schema: payload_removal_disallowed_header -->
        <true/>
        <key>PayloadScope</key> <!-- hcl schema: payload_scope_header -->
        <string>System</string>
        <key>PayloadType</key> <!-- hcl schema: payload_type_header -->
        <string>Configuration</string>
        <key>PayloadUUID</key> <!-- hcl schema: payload_uuid_header -->
        <string>1A803CC7-58DB-43DC-A783-D20C4D9A033A</string>
        <key>PayloadVersion</key> <!-- hcl schema: payload_version_header -->
        <integer>1</integer>
    </dict>
</plist>
*/

// resourceJamfProMacOSConfigurationProfilesPlistGenerator defines the schema and CRUD operations for managing Jamf Pro macOS Configuration Profiles in Terraform.
func ResourceJamfProMacOSConfigurationProfilesPlistGenerator() *schema.Resource {
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
											Schema: payloadContentSchema().Schema,
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
func nestedDictionarySchema(level int) *schema.Schema {
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
				"dictionary": nestedDictionarySchema(level - 1),
			},
		},
	}
}

// Define the payload content schema with limited depth
func payloadContentSchema() *schema.Resource {
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
			"dictionary": nestedDictionarySchema(6),
		},
	}
}
