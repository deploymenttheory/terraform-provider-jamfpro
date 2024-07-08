package departments

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProDepartmentsCreate is responsible for creating a new Jamf Pro Department in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
// resourceJamfProDepartmentsCreate is responsible for creating a new Jamf Pro Department in the remote system.
func resourceJamfProDepartmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProDepartment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Department: %v", err))
	}

	var creationResponse *jamfpro.ResponseDepartmentCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateDepartment(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Department '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(creationResponse.ID)

	return append(diags, resourceJamfProDepartmentsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProDepartmentsRead is responsible for reading the current state of a Jamf Pro Department Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
// resourceJamfProDepartmentsRead is responsible for reading the current state of a Jamf Pro Department Resource from the remote system.
func resourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	var response *jamfpro.ResourceDepartment
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetDepartmentByID(resourceID)
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

// resourceJamfProDepartmentsReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProDepartmentsReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProDepartmentsRead(ctx, d, meta, true)
}

// resourceJamfProDepartmentsReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProDepartmentsReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProDepartmentsRead(ctx, d, meta, false)
}

// resourceJamfProDepartmentsUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func resourceJamfProDepartmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)

	department, err := constructJamfProDepartment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error constructing Jamf Pro Department '%s': %v", resourceName, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateDepartmentByID(resourceID, department)
		if apiErr == nil {
			return nil
		}

		_, apiErrByName := client.UpdateDepartmentByName(resourceName, department)
		if apiErrByName != nil {
			return retry.RetryableError(fmt.Errorf("failed to update department '%s' by ID '%s' and by name due to errors: %v, %v", resourceName, resourceID, apiErr, apiErrByName))
		}

		return nil
	})

	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("final attempt to update department '%s' failed: %v", resourceName, err))...)
	}

	return append(diags, resourceJamfProDepartmentsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProDepartmentsDelete is responsible for deleting a Jamf Pro Department.
func resourceJamfProDepartmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteDepartmentByID(resourceID)
		if apiErr != nil {
			apiErrByName := client.DeleteDepartmentByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(fmt.Errorf("failed to delete department '%s' by ID '%s' and by name due to errors: %v, %v", resourceName, resourceID, apiErr, apiErrByName))
			}
		}
		return nil
	})

	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("final attempt to delete department '%s' failed: %v", resourceName, err))...)
	}

	d.SetId("")

	return diags
}
