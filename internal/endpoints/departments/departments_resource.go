// department_resource.go
package departments

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	JamfProResourceDepartment = "Department"
)

// ResourceJamfProDepartments defines the schema and CRUD operations for managing Jamf Pro Departments in Terraform.
func ResourceJamfProDepartments() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProDepartmentsCreate,
		ReadContext:   ResourceJamfProDepartmentsRead,
		UpdateContext: ResourceJamfProDepartmentsUpdate,
		DeleteContext: ResourceJamfProDepartmentsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
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
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceName := d.Get("name").(string)

	// Construct the department object
	department, err := constructJamfProDepartment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Department '%s': %v", resourceName, err))
	}

	// Attempt to create the department in Jamf Pro
	creationResponse, err := conn.CreateDepartment(department)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Department '%s': %v", resourceName, err))
	}

	// Set the resource ID in the Terraform state
	d.SetId(creationResponse.ID)

	// Sync the Terraform state with the remote system
	readDiags := ResourceJamfProDepartmentsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProDepartmentsRead is responsible for reading the current state of a Jamf Pro Department Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
// ResourceJamfProDepartmentsRead is responsible for reading the current state of a Jamf Pro Department Resource from the remote system.
func ResourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	var resource *jamfpro.ResourceDepartment

	// Read operation with retry
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetDepartmentByID(resourceID)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				// Return non-retryable error with a message to avoid SDK issues
				return retry.NonRetryableError(fmt.Errorf("resource not found, marked for deletion"))
			}
			// Retry for other types of errors
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	// If err is not nil, check if it's due to the resource being not found
	if err != nil {
		if err.Error() == "resource not found, marked for deletion" {
			// Resource not found, remove from Terraform state
			d.SetId("")
			// Append a warning diagnostic and return
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("Jamf Pro Distribution Point with ID '%s' was not found on the server and is marked for deletion from terraform state.", resourceID),
			})
			return diags
		}

		// For other errors, return an error diagnostic
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Distribution Point with ID '%s' after retries: %v", resourceID, err))
	}

	// Update the Terraform state with the fetched data
	if resource != nil {
		// Set the fields directly in the Terraform state
		if err := d.Set("id", resourceID); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// ResourceJamfProDepartmentsUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func ResourceJamfProDepartmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)

	// Construct the resource object
	department, err := constructJamfProDepartment(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error constructing Jamf Pro Department '%s': %v", resourceName, err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateDepartmentByID(resourceID, department)
		if apiErr == nil {
			// Successfully updated the department, exit the retry loop
			return nil
		}

		// If update by ID fails, attempt to update by Name
		_, apiErrByName := conn.UpdateDepartmentByName(resourceName, department)
		if apiErrByName != nil {
			// Log the error and return a retryable error
			return retry.RetryableError(fmt.Errorf("failed to update department '%s' by ID '%s' and by name due to errors: %v, %v", resourceName, resourceID, apiErr, apiErrByName))
		}

		// Successfully updated the department by name, exit the retry loop
		return nil
	})

	// Handle error after all retries are exhausted
	if err != nil {
		diags = append(diags, diag.FromErr(fmt.Errorf("final attempt to update department '%s' failed: %v", resourceName, err))...)
		return diags
	}

	// Log the successful update
	hclog.FromContext(ctx).Info(fmt.Sprintf("Successfully updated department '%s' with ID '%s'", resourceName, resourceID))

	// Sync the Terraform state with the remote system
	readDiags := ResourceJamfProDepartmentsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProDepartmentsDelete is responsible for deleting a Jamf Pro Department.
func ResourceJamfProDepartmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)

	// Use the retry function for the delete operation with appropriate timeout
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteDepartmentByID(resourceID)
		if apiErr != nil {
			// If deletion by ID fails, attempt to delete by Name
			apiErrByName := conn.DeleteDepartmentByName(resourceName)
			if apiErrByName != nil {
				// Log the error and return a retryable error
				return retry.RetryableError(fmt.Errorf("failed to delete department '%s' by ID '%s' and by name due to errors: %v, %v", resourceName, resourceID, apiErr, apiErrByName))
			}
		}
		// Successfully deleted the department, exit the retry loop
		return nil
	})

	// Handle error after all retries are exhausted
	if err != nil {
		diags = append(diags, diag.FromErr(fmt.Errorf("final attempt to delete department '%s' failed: %v", resourceName, err))...)
		return diags
	}

	// Log the successful deletion
	hclog.FromContext(ctx).Info(fmt.Sprintf("Successfully deleted department '%s' with ID '%s'", resourceName, resourceID))

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
