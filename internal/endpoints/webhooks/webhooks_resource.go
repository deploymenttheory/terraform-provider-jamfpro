// webhooks_resource.go
package webhooks

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProWebhooks defines the schema and CRUD operations for managing Jamf Pro Webhooks in Terraform.
func ResourceJamfProWebhooks() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProWebhookCreate,
		ReadContext:   resourceJamfProWebhookRead,
		UpdateContext: resourceJamfProWebhookUpdate,
		DeleteContext: resourceJamfProWebhookDelete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
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
					v := util.GetString(val)
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
					v := util.GetString(val)
					validContentTypes := []string{"text/xml", "application/json"}

					for _, validType := range validContentTypes {
						if v == validType {
							return // Valid value found, return without error
						}
					}

					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validContentTypes, v))
					return warns, errs
				},
			},
			"event": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The event type that triggers the webhook.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
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
					v := util.GetInt(val)
					if v < 0 || v > 16 {
						errs := make([]error, 0)
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
					v := util.GetInt(val)
					if v < 0 || v > 16 {
						errs := make([]error, 0)
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
					v := util.GetString(val)
					validAuthTypes := []string{"BASIC", "HEADER"}

					for _, validType := range validAuthTypes {
						if v == validType {
							return // Valid value found, return without error
						}
					}

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
				Optional:    true,
				Description: "List of display fields associated with the webhook.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the display field.",
						},
					},
				},
			},
		},
	}
}

// resourceJamfProWebhooksCreate is responsible for creating a new Jamf Pro Webhook in the remote system.
// The function:
// 1. Constructs the Webhook data using the provided Terraform configuration.
// 2. Calls the API to create the Webhook in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created Webhook.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProWebhookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProWebhook(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Webhook: %v", err))
	}

	var creationResponse *jamfpro.ResourceWebhook
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateWebhook(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Webhook '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetWebhookByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Webhook", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)
	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, resourceJamfProWebhookRead(ctx, d, meta)...)
}

// resourceJamfProWebhookRead is responsible for reading the current state of a Jamf Pro Webhook Resource from the remote system.
// The function:
// 1. Fetches the user group's current state using its ID. If it fails, it tries to obtain the user group's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the user group being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProWebhookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetWebhookByID(resourceIDInt)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// resourceJamfProWebhookUpdate is responsible for updating an existing Jamf Pro Webhook on the remote system.
func resourceJamfProWebhookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProWebhook(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Webhook for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateWebhookByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Webhook '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	return append(diags, resourceJamfProWebhookRead(ctx, d, meta)...)
}

// resourceJamfProWebhookDelete is responsible for deleting a Jamf Pro Webhook.
func resourceJamfProWebhookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteWebhookByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteWebhookByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Webhook '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
