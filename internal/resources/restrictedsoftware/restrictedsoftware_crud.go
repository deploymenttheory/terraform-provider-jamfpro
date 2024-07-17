package restrictedsoftware

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

// resourceJamfProRestrictedSoftwareCreate is responsible for creating a new Jamf Pro Restricted Software in the remote system.
// The function:
// 1. Constructs the User Group data using the provided Terraform configuration.
// 2. Calls the API to create the User Group in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created User Group.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProRestrictedSoftwareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Restricted Software: %v", err))
	}

	var creationResponse *jamfpro.ResponseRestrictedSoftwareCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateRestrictedSoftware(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Restricted Software '%s' after retries: %v", resource.General.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	return append(diags, resourceJamfProRestrictedSoftwareReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProRestrictedSoftwareRead is responsible for reading the current state of a Jamf Pro Restricted Software Resource from the remote system.
// The function:
// 1. Fetches the user group's current state using its ID. If it fails, it tries to obtain the user group's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the user group being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProRestrictedSoftwareRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var response *jamfpro.ResourceRestrictedSoftware
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetRestrictedSoftwareByID(resourceIDInt)
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

// resourceJamfProRestrictedSoftwareReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProRestrictedSoftwareReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProRestrictedSoftwareRead(ctx, d, meta, true)
}

// resourceJamfProRestrictedSoftwareReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProRestrictedSoftwareReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProRestrictedSoftwareRead(ctx, d, meta, false)
}

// resourceJamfProRestrictedSoftwareUpdate is responsible for updating an existing Jamf Pro Printer on the remote system.
func resourceJamfProRestrictedSoftwareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Restricted Software for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateRestrictedSoftwareByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Restricted Software '%s' (ID: %d) after retries: %v", resource.General.Name, resourceIDInt, err))
	}

	return append(diags, resourceJamfProRestrictedSoftwareReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProRestrictedSoftwareDelete is responsible for deleting a Jamf Pro Restricted Software.
func resourceJamfProRestrictedSoftwareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteRestrictedSoftwareByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteRestrictedSoftwareByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Restricted Software '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
