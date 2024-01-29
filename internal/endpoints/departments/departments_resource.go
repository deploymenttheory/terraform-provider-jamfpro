// department_resource.go
package departments

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/sdkv2"

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
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
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

// constructJamfProDepartment constructs a ResourceDepartment object from the provided schema data.
func constructJamfProDepartment(d *schema.ResourceData) (*jamfpro.ResourceDepartment, error) {
	department := &jamfpro.ResourceDepartment{}

	// Utilize type assertion helper functions for direct field extraction
	department.Name = util.GetStringFromInterface(d.Get("name"))

	// Log the successful construction of the department
	log.Printf("[INFO] Successfully constructed Department with name: %s", department.Name)

	return department, nil
}

// ResourceJamfProDepartmentsCreate is responsible for creating a new Jamf Pro Site in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProDepartmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Instantiate the centralized logger
	logger := sdkv2.ConsoleLogger{}

	var createdAttribute *jamfpro.ResponseDepartmentCreate

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		department, err := constructJamfProDepartment(d)
		if err != nil {
			logger.Error("Failed to construct department", err.Error(), "name", d.Get("name"))
			return retry.NonRetryableError(err)
		}

		createdAttribute, err = conn.CreateDepartment(department)
		if err != nil {
			if apiErr, ok := err.(*http_client.APIError); ok {
				logger.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message, "name", department.Name)
				return retry.NonRetryableError(err)
			}
			logger.Errorf("Failed to create department: %s", err.Error(), "name", department.Name)
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		// Log and return the error using the centralized logger
		logger.Errorf("Failed to create department: %s", err.Error(), "name", d.Get("name"))
		return logger.Diagnostics
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(createdAttribute.ID)

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDepartmentsRead(ctx, d, meta)
		if len(readDiags) > 0 {
			logger.Errorf("Failed to read the created department: %s", readDiags[0].Summary, "name", d.Get("name"))
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Failed to update the Terraform state for the created department: %s", err.Error(), "name", d.Get("name"))
		return logger.Diagnostics
	}

	return nil
}

// ResourceJamfProDepartmentsRead is responsible for reading the current state of a Jamf Pro Department Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the centralized logger
	logger := sdkv2.ConsoleLogger{}

	var attribute *jamfpro.ResourceDepartment

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {

		attributeID := d.Id()

		var apiErr error

		attribute, apiErr = conn.GetDepartmentByID(attributeID)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				logger.Errorf("API Error (Code: %d): %s while fetching by ID", apiError.StatusCode, apiError.Message, "id", attributeID)
				return retry.NonRetryableError(apiErr)
			}

			attributeName, ok := d.Get("name").(string)
			if !ok {
				logger.Error("Unable to assert 'name' as a string for fetching department", "", "id", attributeID)
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			attribute, apiErr = conn.GetDepartmentByName(attributeName)
			if apiErr != nil {
				logger.Errorf("Error fetching department by name: %s", apiErr.Error(), "name", attributeName)
				return retry.RetryableError(apiErr)
			}
		}

		logger.Infof("Successfully fetched department: %s", attribute.Name)
		return nil
	})

	if err != nil {
		logger.Errorf("Failed to read department: %s", err.Error(), "id", d.Id())
		return logger.Diagnostics
	}

	if err := d.Set("name", attribute.Name); err != nil {
		logger.Errorf("Failed to set 'name' attribute in Terraform state: %s", err.Error(), "name", attribute.Name)
		return logger.Diagnostics
	}

	return nil
}

// ResourceJamfProDepartmentsUpdate is responsible for updating an existing Jamf Pro Site on the remote system.
func ResourceJamfProDepartmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the centralized logger
	logger := sdkv2.ConsoleLogger{}

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		department, err := constructJamfProDepartment(d)
		if err != nil {
			logger.Error("Failed to construct department for update", err.Error(), "id", d.Id())
			return retry.NonRetryableError(err)
		}

		departmentID := d.Id()
		_, apiErr := conn.UpdateDepartmentByID(departmentID, department)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				logger.Errorf("API Error (Code: %d): %s during update by ID", apiError.StatusCode, apiError.Message, "id", departmentID)
				return retry.NonRetryableError(apiErr)
			}

			departmentName := d.Get("name").(string)
			_, apiErr = conn.UpdateDepartmentByName(departmentName, department)
			if apiErr != nil {
				logger.Errorf("Error updating department by name: %s", apiErr.Error(), "name", departmentName)
				return retry.RetryableError(apiErr)
			}
		}

		logger.Infof("Successfully updated department: %s", department.Name)
		return nil
	})

	if err != nil {
		logger.Errorf("Failed to update department: %s", err.Error(), "id", d.Id())
		return logger.Diagnostics
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDepartmentsRead(ctx, d, meta)
		if len(readDiags) > 0 {
			logger.Errorf("Failed to read department after update: %s", readDiags[0].Summary, "id", d.Id())
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Failed to synchronize Terraform state after department update: %s", err.Error(), "id", d.Id())
		return logger.Diagnostics
	}

	return nil
}

// ResourceJamfProDepartmentsDelete is responsible for deleting a Jamf Pro Department.
func ResourceJamfProDepartmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the centralized logger
	logger := sdkv2.ConsoleLogger{}

	departmentID := d.Id()

	// Extract the retry timeout from the schema
	retryTimeout := d.Timeout(schema.TimeoutDelete)

	err := retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		apiErr := conn.DeleteDepartmentByID(departmentID)
		if apiErr != nil {
			departmentName := d.Get("name").(string)
			apiErr = conn.DeleteDepartmentByName(departmentName)
			if apiErr != nil {
				// Log the error and continue retrying within the retry time window
				logger.Error(
					"Failed to delete department. Retrying...",
					apiErr.Error(),
					"id", departmentID, "name", departmentName,
				)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	if err != nil {
		// Log the error with the retry time window information
		logger.Error(
			fmt.Sprintf("Failed to delete department within the retry time window of %s.", formatDuration(retryTimeout)),
			err.Error(),
			"id", departmentID,
		)
		return logger.Diagnostics
	}

	d.SetId("") // Clear the ID from the Terraform state
	return nil
}

// formatDuration formats the duration in a human-readable form.
func formatDuration(d time.Duration) string {
	return d.Round(time.Second).String()
}
