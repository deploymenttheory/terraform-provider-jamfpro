package apiroles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProAPIRolesCreate handles the creation of a Jamf Pro API Role.
// The function:
// 1. Constructs the API role data using the provided Terraform configuration.
// 2. Calls the API to create the role in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created role.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProAPIRolesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProApiRole(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Account Group: %v", err))
	}

	var creationResponse *jamfpro.ResourceAPIRole
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateJamfApiRole(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro API Role '%s' after retries: %v", resource.DisplayName, err))
	}

	d.SetId(creationResponse.ID)

	return append(diags, resourceJamfProAPIRolesRead(ctx, d, meta)...)
}

// resourceJamfProAPIRolesRead handles reading a Jamf Pro API Role from the remote system.
// The function:
// 1. Tries to fetch the API role based on the ID from the Terraform state.
// 2. If fetching by ID fails, attempts to fetch it by the display name.
// 3. Updates the Terraform state with the fetched data.
func resourceJamfProAPIRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	var err error
	resourceID := d.Id()

	var response *jamfpro.ResourceAPIRole
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetJamfApiRoleByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return append(diags, updateTerraformState(d, response)...)
}

// resourceJamfProAPIRolesUpdate handles updating a Jamf Pro API Role.
// The function:
// 1. Constructs the updated API role data using the provided Terraform configuration.
// 2. Calls the API to update the role in Jamf Pro.
// 3. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProAPIRolesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProApiRole(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro API Role for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateJamfApiRoleByID(resourceID, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro API Role '%s' (ID: %s) after retries: %v", resource.DisplayName, resourceID, err))
	}

	return append(diags, resourceJamfProAPIRolesRead(ctx, d, meta)...)
}

// resourceJamfProAPIRolesDelete handles the deletion of a Jamf Pro API Role.
func resourceJamfProAPIRolesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteJamfApiRoleByID(resourceID)
		if apiErr != nil {
			resourceName := d.Get("display_name").(string)
			apiErrByName := client.DeleteJamfApiRoleByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro API role '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
