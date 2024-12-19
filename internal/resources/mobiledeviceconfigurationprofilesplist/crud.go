package mobiledeviceconfigurationprofilesplist

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

// resourceJamfProMobileDeviceConfigurationProfileCreate is responsible for creating a new Jamf Pro Mobile Device Configuration Profile in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProMobileDeviceConfigurationProfilePlistCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProMobileDeviceConfigurationProfilePlist(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Mobile Device Configuration Profile: %v", err))
	}

	var creationResponse *jamfpro.ResponseMobileDeviceConfigurationProfileCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateMobileDeviceConfigurationProfile(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Mobile Device Configuration Profile '%s' after retries: %v", resource.General.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	return append(diags, resourceJamfProMobileDeviceConfigurationProfilePlistReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProMobileDeviceConfigurationProfilePlistRead is responsible for reading the current state of a Jamf Pro Mobile Device Configuration Profile Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProMobileDeviceConfigurationProfilePlistRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	var response *jamfpro.ResourceMobileDeviceConfigurationProfile
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetMobileDeviceConfigurationProfileByID(resourceID)
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

// resourceJamfProMobileDeviceConfigurationProfilePlistReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProMobileDeviceConfigurationProfilePlistReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProMobileDeviceConfigurationProfilePlistRead(ctx, d, meta, true)
}

// resourceJamfProMobileDeviceConfigurationProfilePlistReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProMobileDeviceConfigurationProfilePlistReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProMobileDeviceConfigurationProfilePlistRead(ctx, d, meta, false)
}

// resourceJamfProMobileDeviceConfigurationProfilePlistUpdate is responsible for updating an existing Jamf Pro Mobile Device Configuration Profile on the remote system.
func resourceJamfProMobileDeviceConfigurationProfilePlistUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProMobileDeviceConfigurationProfilePlist(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Mobile Device Configuration Profile for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateMobileDeviceConfigurationProfileByID(resourceID, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Mobile Device Configuration Profile '%s' (ID: %s) after retries: %v", resource.General.Name, resourceID, err))
	}

	return append(diags, resourceJamfProMobileDeviceConfigurationProfilePlistReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProMobileDeviceConfigurationProfilePlistDelete is responsible for deleting a Jamf Pro Mobile Device Configuration Profile.
func resourceJamfProMobileDeviceConfigurationProfilePlistDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteMobileDeviceConfigurationProfileByID(resourceID)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteMobileDeviceConfigurationProfileByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Mobile Device Configuration Profile '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
