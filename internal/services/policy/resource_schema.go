// policies_resource.go
package policy

import (
	"fmt"
	"time"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProPolicies defines the schema and CRUD operations for managing Jamf Pro Policy in Terraform.
func ResourceJamfProPolicies() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: validateSelfServiceConfig,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourcePolicyV0().CoreConfigSchema().ImpliedType(),
				Upgrade: upgradePolicyUserInteractionV0toV1,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Jamf Pro policy.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Define whether the policy is enabled.",
			},
			"trigger_checkin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when device performs recurring check-in against the frequency configured in Jamf Pro",
			},
			"trigger_enrollment_complete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when device enrollment is complete.",
			},
			"trigger_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when a user logs in to a computer. A login event that checks for policies must be configured in Jamf Pro for this to work",
			},
			"trigger_network_state_changed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when it's network state changes. When a computer's network state changes (e.g., when the network connection changes, when the computer name changes, when the IP address changes)",
			},
			"trigger_startup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Trigger policy when a computer starts up. A startup script that checks for policies must be configured in Jamf Pro for this to work",
			},
			"trigger_other": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any other trigger for the policy.",
				Default:     "",
			},
			"frequency": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Frequency of policy execution.",
				Default:     "Once per computer",
				ValidateFunc: validation.StringInSlice([]string{
					"Once per computer",
					"Once per user per computer",
					"Once per user",
					"Once every day",
					"Once every week",
					"Once every month",
					"Ongoing",
				}, false),
			},
			"retry_event": { // Retry only relevant if frequency is Once Per Computer
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Event on which to retry policy execution.",
				Default:     "none",
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"trigger",
					"check-in",
				}, false),
			},
			"retry_attempts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of retry attempts for the jamf pro policy. Valid values are -1 (not configured) and 1 through 10.",
				Default:     -1,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					vInt := val.(int)
					if vInt == -1 || (vInt > 0 && vInt <= 10) {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be -1 if not being set or between 1 and 10 if it is being set, got: %d", key, val))
					return warns, errs
				},
			},
			"notify_on_each_failed_retry": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Send notifications for each failed policy retry attempt. ",
				Default:     false,
			},
			"target_drive": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The drive on which to run the policy (e.g. /Volumes/Restore/ ). The policy runs on the boot drive by default",
				Default:     "/",
			},
			"offline": { // Only avaible if frequency set to continuous else not needed
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Make policy available offline by caching the policy to the macOS device to ensure it runs when Jamf Pro is unavailable. Only used when execution policy is set to 'ongoing'. ",
				Default:     false,
			},
			"network_requirements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network requirements for the policy.",
				Default:     "Any",
				ValidateFunc: validation.StringInSlice([]string{
					"Any",
					"Ethernet",
				}, false),
			},
			"category_id": sharedschemas.GetSharedSchemaCategory(),
			"site_id":     sharedschemas.GetSharedSchemaSite(),
			"date_time_limitations": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Server-side limitations use your Jamf Pro host server's time zone and settings. The Jamf Pro host service is in UTC time.",
				Elem:        getPolicySchemaDateTimeLimitations(),
			},
			"network_limitations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Network limitations for the policy.",
				MaxItems:    1,
				Elem:        getPolicySchemaNetworkLimitations(),
			}, // END OF General UI
			"payloads": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "All payloads container",
				Elem:        getPolicySchemaPayloads(),
			},
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Scope configuration for the profile.",
				Elem:        sharedschemas.GetSharedmacOSComputerSchemaScope(),
			},
			"self_service": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				MaxItems:    1,
				Description: "Self-service settings of the policy.",
				Elem:        getPolicySchemaSelfService(),
			},
			"package_distribution_point": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "repository of which packages are collected from",
			},
		},
	}
}
