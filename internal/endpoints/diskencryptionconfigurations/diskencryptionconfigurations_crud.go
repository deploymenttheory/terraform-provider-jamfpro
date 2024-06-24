package diskencryptionconfigurations

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

// resourceJamfProDiskEncryptionConfigurationsCreate is responsible for creating a new Jamf Pro Disk Encryption Configuration in the remote system.
// The function:
// 1. Constructs the disk encryption configuration data using the provided Terraform configuration.
// 2. Calls the API to create the disk encryption configuration in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created disk encryption configuration.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProDiskEncryptionConfigurationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProDiskEncryptionConfiguration(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Disk Encryption Configuration: %v", err))
	}

	var creationResponse *jamfpro.ResponseDiskEncryptionConfigurationCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateDiskEncryptionConfiguration(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Disk Encryption Configuration '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	return append(diags, resourceJamfProDiskEncryptionConfigurationsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProDiskEncryptionConfigurationsRead is responsible for reading the current state of a Jamf Pro Disk Encryption Configuration resource from Jamf Pro and updating the Terraform state with the retrieved data.
func resourceJamfProDiskEncryptionConfigurationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var response *jamfpro.ResourceDiskEncryptionConfiguration
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetDiskEncryptionConfigurationByID(resourceIDInt)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, state.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateTerraformState(d, response)...)
}

func resourceJamfProDiskEncryptionConfigurationsReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProDiskEncryptionConfigurationsRead(ctx, d, meta, true)
}

func resourceJamfProDiskEncryptionConfigurationsReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProDiskEncryptionConfigurationsRead(ctx, d, meta, false)
}

// resourceJamfProDiskEncryptionConfigurationsUpdate is responsible for updating an existing Jamf Pro Disk Encryption Configuration on the remote system.
func resourceJamfProDiskEncryptionConfigurationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProDiskEncryptionConfiguration(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Disk Encryption Configuration for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateDiskEncryptionConfigurationByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Disk Encryption Configuration '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	return append(diags, resourceJamfProDiskEncryptionConfigurationsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProDiskEncryptionConfigurationsDelete is responsible for deleting a Jamf Pro Disk Encryption Configuration.
func resourceJamfProDiskEncryptionConfigurationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteDiskEncryptionConfigurationByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteDiskEncryptionConfigurationByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Disk Encryption Configuration '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
