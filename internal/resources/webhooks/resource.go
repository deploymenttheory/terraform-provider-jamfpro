// webhooks_resource.go
package webhooks

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProWebhooks defines the schema and CRUD operations for managing Jamf Pro Webhooks in Terraform.
func ResourceJamfProWebhooks() *schema.Resource {
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
				Description: "The unique identifier of the Jamf Pro webhook.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Jamf Pro webhook.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "The Jamf Pro webhooks enablement state.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL the webhook will post data to.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
					if !strings.HasPrefix(v, "http://") && !strings.HasPrefix(v, "https://") {
						errs = append(errs, fmt.Errorf("%q must start with 'http://' or 'https://', got: %s", key, v))
					}
					return warns, errs
				},
			},
			"content_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The content type of the webhook payload (e.g., 'application/json' or 'text/xml').",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validContentTypes := []string{"text/xml", "application/json"}

					for _, validType := range validContentTypes {
						if v == validType {
							return // Valid value found, return without error
						}
					}

					//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validContentTypes, v))
					return warns, errs
				},
			},
			"event": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The event type that triggers the webhook.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validEvents := []string{
						"ComputerAdded", "ComputerCheckIn", "ComputerInventoryCompleted",
						"ComputerPatchPolicyCompleted", "ComputerPolicyFinished", "ComputerPushCapabilityChanged",
						"DeviceAddedToDEP", "DeviceRateLimited", "JSSShutdown", "JSSStartup",
						"MobileDeviceCheckIn", "MobileDeviceCommandCompleted", "MobileDeviceEnrolled",
						"MobileDeviceInventoryCompleted", "MobileDevicePushSent", "MobileDeviceUnEnrolled",
						"PatchSoftwareTitleUpdated", "PushSent", "RestAPIOperation", "SCEPChallenge",
						"SmartGroupComputerMembershipChange", "SmartGroupMobileDeviceMembershipChange", "SmartGroupUserMembershipChange",
					}

					for _, validEvent := range validEvents {
						if v == validEvent {
							return // Valid value found, return without error
						}
					}

					//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validEvents, v))
					return warns, errs
				},
			},
			"connection_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "Amount of time to wait for a response from the webhook's host server after sending a request, in seconds.Value must be an integer between 1 and 15",
				ValidateFunc: func(val interface{}, key string) ([]string, []error) {
					v := val.(int)
					if v < 0 || v > 16 {
						errs := make([]error, 0)
						//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
						errs = append(errs, fmt.Errorf("%q must be between 1 and 15, inclusive, got: %d", key, v))
						return nil, errs
					}
					return nil, nil
				},
			},
			"read_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Amount of time to attempt to connect to the webhook's host server, in seconds.Value must be an integer between 1 and 15",
				ValidateFunc: func(val interface{}, key string) ([]string, []error) {
					v := val.(int)
					if v < 0 || v > 16 {
						errs := make([]error, 0)
						//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
						errs = append(errs, fmt.Errorf("%q must be between 1 and 15, inclusive, got: %d", key, v))
						return nil, errs
					}
					return nil, nil
				},
			},
			"authentication_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "None",
				Description: "The type of authentication required for the webhook (e.g., BASIC 'Basic Authentication', HEADER for 'Header Authentication').",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validAuthTypes := []string{"BASIC", "HEADER"}

					for _, validType := range validAuthTypes {
						if v == validType {
							return // Valid value found, return without error
						}
					}

					//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validAuthTypes, v))
					return warns, errs
				},
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username for authentication, if applicable.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for authentication, if applicable.",
			},
			"enable_display_fields_for_group": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable display fields for the group associated with the webhook.",
			},
			"smart_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the smart group associated with the webhook.",
			},
			"display_fields": {
				Type:        schema.TypeList,
				Description: "List of displayfields",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
