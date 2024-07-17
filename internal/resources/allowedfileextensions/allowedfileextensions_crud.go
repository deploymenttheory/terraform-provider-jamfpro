package allowedfileextensions

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

// resourceJamfProAllowedFileExtensionCreate is responsible for creating a new AllowedFileExtension in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Allowed File Extension: %v", err))
	}

	var creationResponse *jamfpro.ResourceAllowedFileExtension
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateAllowedFileExtension(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Allowed File Extension '%s' after retries: %v", resource.Extension, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// resourceJamfProAllowedFileExtensionRead is responsible for reading the current state of an Allowed File Extension Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var response *jamfpro.ResourceAllowedFileExtension
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetAllowedFileExtensionByID(resourceIDInt)
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

// resourceJamfProAllowedFileExtensionReadWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// resourceJamfProAllowedFileExtensionReadNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// resourceJamfProAllowedFileExtensionUpdate handles the update operation for an AllowedFileExtension resource in Terraform.
// Since there is no direct update API endpoint, this function will delete the existing AllowedFileExtension and create a new one.
// This approach simulates an update operation from the user's perspective in Terraform.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	deleteDiags := delete(ctx, d, meta)
	if deleteDiags.HasError() {
		return deleteDiags
	}

	return create(ctx, d, meta)
}

// resourceJamfProAllowedFileExtensionDelete is responsible for deleting an Allowed File Extension in Jamf Pro.
// This function will delete the resource based on its ID from the Terraform state.
// If the resource cannot be found by ID, it will attempt to delete by the 'extension' attribute.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteAllowedFileExtensionByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("extension").(string)
			apiErrByName := client.DeleteAllowedFileExtensionByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Allowed File Extension '%s' (ID: %s) after retries: %v", d.Get("extension").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
