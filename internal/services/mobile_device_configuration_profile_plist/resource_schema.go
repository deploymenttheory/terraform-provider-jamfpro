// mobiledeviceconfigurationprofilesplist_resource.go
package mobile_device_configuration_profile_plist

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProMobileDeviceConfigurationProfilesPlist defines the schema for mobile device configuration profiles in Terraform.
func ResourceJamfProMobileDeviceConfigurationProfilesPlist() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
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
				Description: "The unique identifier for the mobile device configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the mobile device configuration profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the mobile device configuration profile.",
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The level at which the mobile device configuration profile is applied, can be either 'Device Level' or 'User Level'.",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "Device Level" || v == "User Level" {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be either 'Device Level' or 'User Level', got: %s", key, v))
					return warns, errs
				},
			},
			"site_id":     sharedschemas.GetSharedSchemaSite(),
			"category_id": sharedschemas.GetSharedSchemaCategory(),
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The universally unique identifier for the profile.",
			},
			"deployment_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The deployment method for the mobile device configuration profile, can be either 'Install Automatically' or 'Make Available in Self Service'.",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "Install Automatically" || v == "Make Available in Self Service" {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be either 'Install Automatically' or 'Make Available in Self Service', got: %s", key, v))
					return warns, errs
				},
			},
			"redeploy_on_update": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Defines the redeployment behaviour when an update to a mobile device config profile" +
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
			"redeploy_days_before_cert_expires": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of days before certificate expiration when the profile should be redeployed.",
			},
			"payloads": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: DiffSuppressPayloads,
				Description: "The iOS / iPadOS / tvOS configuration profile payload. Can be a file path to a .mobileconfig or a string with an embedded mobileconfig plist." +
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
				Description: "Controls validation of the Mobile device  configuration profile plist. When enabled (default), " +
					"performs the following validations:\n\n" +
					"1. Payload State Normalization (normalizePayloadState):\n" +
					"   - Normalizes the payload structure for consistent state management\n" +
					"   - Ensures profile format matches Jamf Pro's expected structure\n\n" +
					"Set to false when importing profiles from external sources that may not " +
					"strictly conform to Jamf Pro's plist requirements. Disabling validation " +
					"bypasses these checks but may result in deployment issues if the profile " +
					"structure is incompatible with Jamf Pro, or triggers jamf pro plist processing " +
					"not handled by 'payloads' diff suppression. Switch off at your own risk.",
			},
			// Scope
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The scope of the configuration profile.",
				Required:    true,
				Elem:        sharedschemas.GetSharedMobileDeviceSchemaScope(),
			},
		},
	}
}
