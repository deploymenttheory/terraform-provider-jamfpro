package apiintegrations

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

// resourceJamfProApiIntegrationsCreate is responsible for creating a new Jamf Pro API Integration in the remote system.
func resourceJamfProApiIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro API Integration: %v", err))
	}

	var creationResponse *jamfpro.ResourceApiIntegration
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateApiIntegration(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro API Integration '%s' after retries: %v", resource.DisplayName, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	return append(diags, resourceJamfProApiIntegrationsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProApiIntegrationsRead is responsible for reading the current state of a Jamf Pro API Integration from the remote system.
func resourceJamfProApiIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var response *jamfpro.ResourceApiIntegration
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetApiIntegrationByID(resourceIDInt)
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

// resourceJamfProApiIntegrationsReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProApiIntegrationsReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProApiIntegrationsRead(ctx, d, meta, true)
}

// resourceJamfProApiIntegrationsReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProApiIntegrationsReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProApiIntegrationsRead(ctx, d, meta, false)
}

// resourceJamfProApiIntegrationsUpdate is responsible for updating an existing Jamf Pro API Integration on the remote system.
func resourceJamfProApiIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro API Integration for update: %v", err))
	}

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateApiIntegrationByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro API Integration '%s' (ID: %s) after retries: %v", resource.DisplayName, resourceID, err))
	}

	return append(diags, resourceJamfProApiIntegrationsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProApiIntegrationsDelete is responsible for deleting a Jamf Pro API Integration.
func resourceJamfProApiIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteApiIntegrationByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("display_name").(string)
			apiErrByName := client.DeleteApiIntegrationByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro API Integration '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
