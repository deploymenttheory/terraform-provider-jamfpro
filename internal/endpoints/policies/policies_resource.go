// policies_resource.go
package policies

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProPolicies defines the schema and CRUD operations for managing Jamf Pro Policy in Terraform.
func ResourceJamfProPolicies() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProPoliciesCreate,
		ReadContext:   ResourceJamfProPoliciesRead,
		UpdateContext: ResourceJamfProPoliciesUpdate,
		DeleteContext: ResourceJamfProPoliciesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
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
			// "trigger": { // NOTE appears to be redundant when used with the below. Maybe this use to be a multiple choice option?
			// 	Type:         schema.TypeString,
			// 	Required:     true,
			// 	Description:  "Event(s) triggers to use to initiate the policy. Values can be 'USER_INITIATED' for self self trigger and 'EVENT' for an event based trigger",
			// 	ValidateFunc: validation.StringInSlice([]string{"EVENT", "USER_INITIATED"}, false),
			// },
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
			// "trigger_logout": { // NOTE appears to be redundant
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Description: "Trigger policy when a user logout.",
			// 	Default:     false,
			// },
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
				Default:     "",
				Description: "Any other trigger for the policy.",
				// TODO need a validation func here to make sure this cannot be provided as empty.
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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetInt(val)
					if v == -1 || (v > 0 && v <= 10) {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be -1 if not being set or between 1 and 10 if it is being set, got: %d", key, v))
					return warns, errs
				},
			},
			"notify_on_each_failed_retry": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Send notifications for each failed policy retry attempt. ",
				Default:     false,
			},
			// "location_user_only": { // NOTE Can't find in GUI
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Description: "Location-based policy for user only.",
			// 	Default:     false,
			// },
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
			"category": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Category to add the policy to.",
				Computed:    true,
				Elem:        sharedschemas.GetSharedSchemaCategory(),
			},
			"date_time_limitations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Server-side limitations use your Jamf Pro host server's time zone and settings. The Jamf Pro host service is in UTC time.",
				Computed:    true,
				Elem:        GetPolicySchemaDateTimeLimitations(),
			},
			"network_limitations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Network limitations for the policy.",
				Computed:    true,
				Elem:        &schema.Resource{},
			}, // END OF General UI
			"payloads": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "All payloads container",
				Elem:        getPolicySchemaPayloads(),
			}, // MOVING EVERYTHING BELOW INTO HERE
		},
	}
}
