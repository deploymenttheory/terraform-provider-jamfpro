package advancedcomputersearches

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProAdvancedComputerSearchCreate is responsible for creating a new Jamf Pro Advanced Computer Search in the remote system.
func resourceJamfProAdvancedComputerSearchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics

	resource, err := constructJamfProAdvancedComputerSearch(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Advanced Computer Search: %v", err))
	}

	var creationResponse *jamfpro.ResponseAdvancedComputerSearchCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateAdvancedComputerSearch(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Advanced Computer Search '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	return append(diags, resourceJamfProAdvancedComputerSearchRead(ctx, d, meta)...)
}

// resourceJamfProAdvancedComputerSearchRead is responsible for reading the current state of a Jamf Pro Advanced Computer Search from the remote system.
func resourceJamfProAdvancedComputerSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetAdvancedComputerSearchByID(resourceIDInt)

	// TODO come back to this
	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// resourceJamfProAdvancedComputerSearchUpdate is responsible for updating an existing Jamf Pro Advanced Computer Search on the remote system.
func resourceJamfProAdvancedComputerSearchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProAdvancedComputerSearch(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Advanced Computer Search for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateAdvancedComputerSearchByID(resourceIDInt, resource)
		if apiErr != nil {

			return retry.RetryableError(apiErr)
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Advanced Computer Search '%s' (ID: %s) after retries: %v", resource.Name, resourceID, err))
	}

	return append(diags, resourceJamfProAdvancedComputerSearchRead(ctx, d, meta)...)
}

// resourceJamfProAdvancedComputerSearchDelete is responsible for deleting a Jamf Pro AdvancedComputerSearch.
func resourceJamfProAdvancedComputerSearchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resourceName := d.Get("name").(string)

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteAdvancedComputerSearchByID(resourceIDInt)
		if apiErr != nil {
			apiErrByName := client.DeleteAdvancedComputerSearchByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Advanced Computer Search '%s' (ID: %s) after retries: %v", resourceName, resourceID, err))
	}

	d.SetId("")

	return diags
}
