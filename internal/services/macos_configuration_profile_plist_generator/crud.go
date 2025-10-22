// macosconfigurationprofilesplistgenerator_crud.go
package macos_configuration_profile_plist_generator

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProMacOSConfigurationProfilesPlistCreate is responsible for creating a new Jamf Pro macOS Configuration Profile in the remote system.
// The function:
// 1. Constructs the configuration profile data using the provided Terraform configuration.
// 2. Calls the API to create the configuration profile in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created configuration profile.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProMacOSConfigurationProfilesPlistGeneratorCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Lock the mutex to ensure only one profile plist create can run this function at a time

	resource, err := constructJamfProMacOSConfigurationProfilesPlistGenerator(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro macOS Configuration Profile: %v", err))
	}

	var creationResponse *jamfpro.ResponseMacOSConfigurationProfileCreationUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateMacOSConfigurationProfile(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro macOS Configuration Profile '%s' after retries: %v", resource.General.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	return append(diags, resourceJamfProMacOSConfigurationProfilesPlistGeneratorReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistGeneratorRead is responsible for reading the current state of a Jamf Pro macOS Configuration Profile Resource from the remote system.
// The function:
// 1. Fetches the configuration profile's current state using its ID. If it fails then obtain configuration profile's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the configuration profile being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProMacOSConfigurationProfilesPlistGeneratorRead(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
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
		return append(diags, sdkv2.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateState(d, response)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistGeneratorReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProMacOSConfigurationProfilesPlistGeneratorReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return resourceJamfProMacOSConfigurationProfilesPlistGeneratorRead(ctx, d, meta, true)
}

// resourceJamfProMacOSConfigurationProfilesPlistGeneratorReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProMacOSConfigurationProfilesPlistGeneratorReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return resourceJamfProMacOSConfigurationProfilesPlistGeneratorRead(ctx, d, meta, false)
}

// resourceJamfProMacOSConfigurationProfilesPlistGeneratorUpdate is responsible for updating an existing Jamf Pro config profile on the remote system.
func resourceJamfProMacOSConfigurationProfilesPlistGeneratorUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProMacOSConfigurationProfilesPlistGenerator(d)
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

	return append(diags, resourceJamfProMacOSConfigurationProfilesPlistGeneratorReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistGeneratorDelete is responsible for deleting a Jamf Pro config profile.
func resourceJamfProMacOSConfigurationProfilesPlistGeneratorDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()
	var err error

	resourceName := d.Get("name").(string)

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
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
