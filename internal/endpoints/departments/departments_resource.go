// department_resource.go
package departments

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProDepartments defines the schema and CRUD operations for managing Jamf Pro Departments in Terraform.
func ResourceJamfProDepartments() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProDepartmentsCreate,
		ReadContext:   ResourceJamfProDepartmentsRead,
		UpdateContext: ResourceJamfProDepartmentsUpdate,
		DeleteContext: ResourceJamfProDepartmentsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the department.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro department.",
			},
		},
	}
}

// ResourceJamfProDepartmentsCreate is responsible for creating a new Jamf Pro Department in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
// ResourceJamfProDepartmentsCreate is responsible for creating a new Jamf Pro Department in the remote system.
func ResourceJamfProDepartmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	return client.GetDepartmentByID(id.(string))
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Department", creationResponse.ID, checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)

	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProDepartmentsRead(ctx, d, meta)...)
}

// ResourceJamfProDepartmentsRead is responsible for reading the current state of a Jamf Pro Department Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
// ResourceJamfProDepartmentsRead is responsible for reading the current state of a Jamf Pro Department Resource from the remote system.
func ResourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := client.GetDepartmentByID(resourceID)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProDepartmentsUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func ResourceJamfProDepartmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	hclog.FromContext(ctx).Info(fmt.Sprintf("Successfully updated department '%s' with ID '%s'", resourceName, resourceID))

	return append(diags, ResourceJamfProDepartmentsRead(ctx, d, meta)...)
}

// ResourceJamfProDepartmentsDelete is responsible for deleting a Jamf Pro Department.
func ResourceJamfProDepartmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	hclog.FromContext(ctx).Info(fmt.Sprintf("Successfully deleted department '%s' with ID '%s'", resourceName, resourceID))

	d.SetId("")

	return diags
}
