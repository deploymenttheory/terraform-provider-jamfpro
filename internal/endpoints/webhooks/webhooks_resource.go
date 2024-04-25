// webhooks_resource.go
package webhooks

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProWebhooks defines the schema and CRUD operations for managing Jamf Pro Webhooks in Terraform.
func ResourceJamfProWebhooks() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProWebhookCreate,
		ReadContext:   ResourceJamfProWebhookRead,
		UpdateContext: ResourceJamfProWebhookUpdate,
		DeleteContext: ResourceJamfProWebhookDelete,
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

// ResourceJamfProWebhooksCreate is responsible for creating a new Jamf Pro Webhook in the remote system.
// The function:
// 1. Constructs the Webhook data using the provided Terraform configuration.
// 2. Calls the API to create the Webhook in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created Webhook.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProWebhookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProWebhook(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Webhook: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResourceWebhook
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateWebhook(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Webhook '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return apiclient.Conn.GetWebhookByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Webhook", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second, apiclient.EnableCookieJar)
	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProWebhookRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProWebhookRead is responsible for reading the current state of a Jamf Pro Webhook Resource from the remote system.
// The function:
// 1. Fetches the user group's current state using its ID. If it fails, it tries to obtain the user group's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the user group being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProWebhookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Attempt to fetch the resource by ID
	resource, err := conn.GetWebhookByID(resourceIDInt)

	if err != nil {
		// Handle not found error or other errors
		return state.HandleResourceNotFoundError(err, d)
	}

	// Update the Terraform state with the fetched data from the resource
	diags = updateTerraformState(d, resource)

	// Handle any errors and return diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProWebhookUpdate is responsible for updating an existing Jamf Pro Webhook on the remote system.
func ResourceJamfProWebhookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Construct the resource object
	resource, err := constructJamfProWebhook(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Webhook for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateWebhookByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Webhook '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProWebhookRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProWebhookDelete is responsible for deleting a Jamf Pro Webhook.
func ResourceJamfProWebhookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteWebhookByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteWebhookByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Webhook '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
