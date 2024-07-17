package scripts

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProScriptsCreate is responsible for creating a new Jamf Pro Script in the remote system.
// The function:
// 1. Constructs the script data using the provided Terraform configuration.
// 2. Calls the API to create the script in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created script.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProScriptsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Script: %v", err))
	}

	var creationResponse *jamfpro.ResponseScriptCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateScript(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Script '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(creationResponse.ID)

	return append(diags, resourceJamfProScriptsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProScriptsRead is responsible for reading the current state of a Jamf Pro Script Resource from the remote system.
// The function:
// 1. Fetches the script's current state using its ID. If it fails then obtain script's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the script being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProScriptsRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	var response *jamfpro.ResourceScript
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetScriptByID(resourceID)
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

// resourceJamfProScriptsReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProScriptsReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProScriptsRead(ctx, d, meta, true)
}

// resourceJamfProScriptsReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProScriptsReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProScriptsRead(ctx, d, meta, false)
}

// resourceJamfProScriptsUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func resourceJamfProScriptsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Script for update: %v", err))
	}
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateScriptByID(resourceID, resource)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			_, apiErrByName := client.UpdateScriptByName(resourceName, resource)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Script '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	return append(diags, resourceJamfProScriptsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProScriptsDelete is responsible for deleting a Jamf Pro Department.
func resourceJamfProScriptsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteScriptByID(resourceID)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteScriptByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Script '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
