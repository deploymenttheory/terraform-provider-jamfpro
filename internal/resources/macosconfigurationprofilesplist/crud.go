package macosconfigurationprofilesplist

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProMacOSConfigurationProfilesPlistCreate is responsible for creating a new Jamf Pro macOS Configuration Profile in the remote system.
// The function follows these steps:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Checks if a resource with the same name already exists in Jamf Pro.
//   - If it exists, it uses the existing resource and returns a warning.
//
// 3. If no existing resource is found, it attempts to create the resource in Jamf Pro.
// 4. Implements a retry mechanism with conflict detection:
//   - If a conflict error occurs during creation, it rechecks for the resource's existence.
//   - If found after a conflict, it uses the existing resource.
//
// 5. For successful creations, it updates the Terraform state with the ID of the newly created or found resource.
// 6. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
// This approach helps mitigate race conditions in concurrent operations and handles pre-existing resources gracefully.
func resourceJamfProMacOSConfigurationProfilesPlistCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProMacOSConfigurationProfilePlist(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro macOS Configuration Profile: %v", err))
	}

	existingResource, err := client.GetMacOSConfigurationProfileByName(resource.General.Name)
	if err == nil && existingResource != nil {
		d.SetId(strconv.Itoa(existingResource.General.ID))

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Resource already exists",
			Detail:   fmt.Sprintf("A macOS Configuration Profile with name '%s' already exists. Using existing resource.", resource.General.Name),
		})
		return append(diags, resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup(ctx, d, meta)...)
	}

	var creationResponse *jamfpro.ResponseMacOSConfigurationProfileCreationUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateMacOSConfigurationProfile(resource)
		if apiErr != nil {

			if strings.Contains(apiErr.Error(), "Conflict") || strings.Contains(apiErr.Error(), "Duplicate name") {

				existingResource, getErr := client.GetMacOSConfigurationProfileByName(resource.General.Name)
				if getErr == nil && existingResource != nil {

					d.SetId(strconv.Itoa(existingResource.General.ID))
					return nil
				}
			}
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro macOS Configuration Profile '%s' after retries: %v", resource.General.Name, err))
	}

	if creationResponse != nil {
		d.SetId(strconv.Itoa(creationResponse.ID))
	}

	return append(diags, resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistRead is responsible for reading the current state of a Jamf Pro config profile Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProMacOSConfigurationProfilesPlistRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	var response *jamfpro.ResourceMacOSConfigurationProfile
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetMacOSConfigurationProfileByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateState(d, response)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProMacOSConfigurationProfilesPlistReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProMacOSConfigurationProfilesPlistRead(ctx, d, meta, true)
}

// resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProMacOSConfigurationProfilesPlistRead(ctx, d, meta, false)
}

// resourceJamfProMacOSConfigurationProfilesPlistUpdate is responsible for updating an existing Jamf Pro config profile on the remote system.
func resourceJamfProMacOSConfigurationProfilesPlistUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProMacOSConfigurationProfilePlist(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro macOS Configuration Profile for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateMacOSConfigurationProfileByID(resourceID, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro macOS Configuration Profile '%s' (ID: %s) after retries: %v", resource.General.Name, resourceID, err))
	}

	return append(diags, resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistDelete is responsible for deleting a Jamf Pro config profile.
func resourceJamfProMacOSConfigurationProfilesPlistDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceName := d.Get("name").(string)

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteMacOSConfigurationProfileByID(resourceID)
		if apiErr != nil {
			apiErrByName := client.DeleteMacOSConfigurationProfileByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro macOS Configuration Profile '%s' (ID: %s) after retries: %v", resourceName, resourceID, err))
	}

	d.SetId("")

	return diags
}
