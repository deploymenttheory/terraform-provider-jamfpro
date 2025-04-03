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
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Second),
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
				Description: "The macOS configuration profile payload. Can be a file path to a .mobileconfig or a string with an embedded mobileconfig plist." +
					"Jamf Pro stores configuration profiles as XML property lists (plists). When profiles are uploaded, " +
					"Jamf Pro processes and reformats them for consistency. This means the XML that is considered valid " +
					"for an upload may look different from what Jamf Pro returns. To handle these differences, the provider " +
					"implements comprehensive diff suppression for the following cases:\n\n" +
					"Differences are suppressed in the following cases:\n\n" +
					"1. Base64 Content Normalization:\n" +
					"   - Removes whitespace, newlines, and tabs from base64 encoded strings\n" +
					"   - Example: 'SGVs bG8g V29y bGQ=' vs 'SGVsbG8gV29ybGQ='\n\n" +
					"2. XML Tag Formatting:\n" +
					"   - Standardizes self-closing tag formats\n" +
					"   - Examples: '<true/>' vs '< true/>' vs '<true />' vs '<true    />' vs '<true  \\t />'\n\n" +
					"3. Empty String Standardization:\n" +
					"   - Normalizes various representations of empty strings\n" +
					"   - Converts strings containing only whitespace to empty strings\n" +
					"   - Example: '' vs '    ' vs '\\n\\t'\n\n" +
					"4. HTML Entity Decoding:\n" +
					"   - Unescapes HTML entities for comparison\n" +
					"   - Example: '&lt;string&gt;' vs '<string>'\n" +
					"   - Example: '&quot;text&quot;' vs '\"text\"'\n\n" +
					"5. Key Ordering:\n" +
					"   - Sorts dictionary keys alphabetically for consistent comparison\n" +
					"   - Example: '{\"b\":1,\"a\":2}' vs '{\"a\":2,\"b\":1}'\n\n" +
					"6. Field Exclusions:\n" +
					"   - Ignores Jamf Pro-managed identifiers that may change between environments\n" +
					"   - Excluded fields: PayloadUUID, PayloadIdentifier, PayloadOrganization, PayloadDisplayName\n" +
					"   - These fields are removed from comparison as they are managed by Jamf Pro\n\n" +
					"7. Trailing Whitespace:\n" +
					"   - Removes trailing whitespace from each line\n" +
					"   - Example: 'value    ' vs 'value'\n\n" +
					"This normalization approach ensures that functionally identical profiles are " +
					"recognized as equivalent despite superficial formatting differences. " +
					"\n" +
					"NOTE - This provider only supports plists generated from Jamf Pro. It does not support " +
					"importing plists from other sources. If you need to import a plist from an external source," +
					"(e.g. iMazing, Apple Configurator, etc.) " +
					"you must first import it into Jamf Pro, then export it from Jamf Pro to generate a compatible plist. " +
					"This provider cannot diff suppress plists generated from external sources.",
			},
			"payload_validate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Controls validation of the MacOS configuration profile plist. When enabled (default), " +
					"performs the following validations:\n\n" +
					"1. Profile Structure Validation (validatePayload):\n" +
					"   - Verifies valid plist XML format\n" +
					"   - Validates PayloadIdentifier matches PayloadUUID\n" +
					"   - Checks required profile fields\n\n" +
					"2. Payload State Normalization (normalizePayloadState):\n" +
					"   - Normalizes the payload structure for consistent state management\n" +
					"   - Ensures profile format matches Jamf Pro's expected structure\n\n" +
					"3. Distribution Method Validation (validateDistributionMethod):\n" +
					"   - Verifies self-service configuration matches distribution method\n" +
					"   - Example: 'Make Available in Self Service' requires self_service block\n" +
					"   - Example: 'Install Automatically' must not have self_service block\n\n" +
					"4. Profile Level Validation (validateMacOSConfigurationProfileLevel):\n" +
					"   - Ensures PayloadScope in plist matches the 'level' attribute\n" +
					"   - Example: If level is 'System', PayloadScope must be 'System'\n\n" +
					"Set to false when importing profiles from external sources that may not " +
					"strictly conform to Jamf Pro's plist requirements. Disabling validation " +
					"bypasses these checks but may result in deployment issues if the profile " +
					"structure is incompatible with Jamf Pro, or triggers jamf pro plist processing " +
					"not handled by 'payloads' diff suppression. Switch off at your own risk.",
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
							Description: "Message Subject",
						},
						"notification_message": {
							Type:        schema.TypeString,
							Optional:    true,
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
