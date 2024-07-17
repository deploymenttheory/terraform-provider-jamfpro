package webhooks

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProWebhooksCreate is responsible for creating a new Jamf Pro Webhook in the remote system.
// The function:
// 1. Constructs the Webhook data using the provided Terraform configuration.
// 2. Calls the API to create the Webhook in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created Webhook.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProWebhookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
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

	return append(diags, resourceJamfProWebhookReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProWebhookRead is responsible for reading the current state of a Jamf Pro Webhook Resource from the remote system.
// The function:
// 1. Fetches the user group's current state using its ID. If it fails, it tries to obtain the user group's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the user group being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProWebhookRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var response *jamfpro.ResourceWebhook
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetWebhookByID(resourceIDInt)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateTerraformState(d, response)...)
}

// resourceJamfProWebhookReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProWebhookReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProWebhookRead(ctx, d, meta, true)
}

// resourceJamfProWebhookReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProWebhookReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProWebhookRead(ctx, d, meta, false)
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

	resource, err := construct(d)
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

	return append(diags, resourceJamfProWebhookReadNoCleanup(ctx, d, meta)...)
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
