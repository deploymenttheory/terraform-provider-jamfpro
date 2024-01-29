// department_resource.go
package departments

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

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
func constructJamfProDepartment(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceDepartment, error) {
	department := &jamfpro.ResourceDepartment{}

	// Utilize type assertion helper functions for direct field extraction
	department.Name = util.GetStringFromInterface(d.Get("name"))

	// Log the successful construction of the department using tflogger
	logging.Info(ctx, logging.SubsystemCreate, "Successfully constructed Department with name", map[string]interface{}{"name": department.Name})

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

	// Initialize diagnostics for collecting any issues to report back to Terraform
	var diags diag.Diagnostics

	// Use tflogging subsystem
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)

	var createdAttribute *jamfpro.ResponseDepartmentCreate

	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		department, err := constructJamfProDepartment(subCtx, d)
		if err != nil {
			logging.Error(subCtx, logging.SubsystemCreate, "Failed to construct department", map[string]interface{}{
				"name":  d.Get("name"),
				"error": err.Error(),
			})
			return retry.NonRetryableError(err)
		}

		createdAttribute, err = conn.CreateDepartment(department)
		if err != nil {
			if apiErr, ok := err.(*http_client.APIError); ok {
				logging.Error(subCtx, logging.SubsystemAPI, fmt.Sprintf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message), map[string]interface{}{
					"name": department.Name,
				})
				return retry.NonRetryableError(err)
			}
			logging.Error(subCtx, logging.SubsystemCreate, "Failed to create department", map[string]interface{}{
				"name":  department.Name,
				"error": err.Error(),
			})
			return retry.RetryableError(err)
		}
		return nil
	})

	// Log any errors to tf diagnostics
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(createdAttribute.ID)

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDepartmentsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.Error(subCtx, logging.SubsystemRead, "Failed to read the created department", map[string]interface{}{
				"name":    d.Get("name"),
				"summary": readDiags[0].Summary,
			})
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemCreate, "Failed to update the Terraform state for the created department", map[string]interface{}{
			"name":  d.Get("name"),
			"error": err.Error(),
		})
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
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

	attributeID := d.Id()
	attributeName := d.Get("name").(string)

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		attribute, apiErr := conn.GetDepartmentByID(attributeID)
		if apiErr != nil {
			attribute, apiErr = conn.GetDepartmentByName(attributeName)
			if apiErr != nil {
				// Log the error using tflog for internal logging
				logging.Error(ctx, logging.SubsystemRead, "Error fetching department", map[string]interface{}{
					"id":    attributeID,
					"name":  attributeName,
					"error": apiErr.Error(),
				})

				return retry.RetryableError(apiErr)
			}
		}

		// Log the successful fetch using tflog
		logging.Info(ctx, logging.SubsystemRead, "Successfully fetched department", map[string]interface{}{
			"id":   attributeID,
			"name": attribute.Name,
		})

		// Check if attribute is not nil
		if attribute != nil {
			// Set the fields directly in the Terraform state
			if err := d.Set("id", attribute.ID); err != nil {
				return retry.RetryableError(err)
			}
			if err := d.Set("name", attribute.Name); err != nil {
				return retry.RetryableError(err)
			}
			// Add more attributes here as needed
		}

		return nil
	})

	if err != nil {
		// Log the final error using tflog
		logging.Error(ctx, logging.SubsystemRead, "Failed to read department", map[string]interface{}{
			"id":    attributeID,
			"error": err.Error(),
		})

		return diag.FromErr(err)
	}

	return nil
}

// ResourceJamfProDepartmentsUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func ResourceJamfProDepartmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	departmentID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		department, err := constructJamfProDepartment(ctx, d)
		if err != nil {
			// Log the construction error using tflog and the update subsystem
			logging.Error(ctx, logging.SubsystemUpdate, "Failed to construct department for update", map[string]interface{}{
				"error": err.Error(),
				"id":    departmentID,
			})
			return retry.NonRetryableError(err)
		}

		_, apiErr := conn.UpdateDepartmentByID(departmentID, department)
		if apiErr != nil {
			// Log the API error using tflog and the update subsystem
			logging.Error(ctx, logging.SubsystemUpdate, "API error during department update", map[string]interface{}{
				"error": apiErr.Error(),
				"id":    departmentID,
			})
			return retry.RetryableError(apiErr)
		}

		// Log the successful update using tflog and the update subsystem
		logging.Info(ctx, logging.SubsystemUpdate, "Successfully updated department", map[string]interface{}{
			"name": department.Name,
			"id":   departmentID,
		})
		return nil
	})

	if err != nil {
		// Log the final error using tflog and the update subsystem
		logging.Error(ctx, logging.SubsystemUpdate, "Failed to update department", map[string]interface{}{
			"error": err.Error(),
			"id":    departmentID,
		})
		return diag.FromErr(err)
	}

	// Retry reading the department to synchronize the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDepartmentsRead(ctx, d, meta)
		if len(readDiags) > 0 {
			// Log the read error using tflog and the update subsystem
			logging.Error(ctx, logging.SubsystemUpdate, "Failed to read department after update", map[string]interface{}{
				"summary": readDiags[0].Summary,
				"id":      departmentID,
			})
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		// Log the synchronization error using tflog and the update subsystem
		logging.Error(ctx, logging.SubsystemUpdate, "Failed to synchronize Terraform state after department update", map[string]interface{}{
			"error": err.Error(),
			"id":    departmentID,
		})
		return diag.FromErr(err)
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

	departmentID := d.Id()
	departmentName := d.Get("name").(string)

	// Extract the retry timeout from the schema
	retryTimeout := d.Timeout(schema.TimeoutDelete)

	err := retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		apiErr := conn.DeleteDepartmentByID(departmentID)
		if apiErr != nil {
			apiErr = conn.DeleteDepartmentByName(departmentName)
			if apiErr != nil {
				// Log the error using helper function for internal logging
				logging.Error(ctx, logging.SubsystemDelete, "Failed to delete department, retrying...",
					map[string]interface{}{
						"id":   departmentID,
						"name": departmentName,
					},
				)

				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	if err != nil {
		// Log the final error using helper function
		logging.Error(ctx, logging.SubsystemDelete, fmt.Sprintf("Failed to delete department within the retry time window of %s.", formatDuration(retryTimeout)),
			map[string]interface{}{
				"id":    departmentID,
				"error": err.Error(),
			},
		)

		return diag.FromErr(err)
	}

	d.SetId("") // Clear the ID from the Terraform state
	return nil
}

// formatDuration formats the duration in a human-readable form.
func formatDuration(d time.Duration) string {
	return d.Round(time.Second).String()
}
