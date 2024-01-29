// department_resource.go
package departments

import (
	"context"
	"encoding/xml"
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
			Create: schema.DefaultTimeout(10 * time.Second),
			Read:   schema.DefaultTimeout(10 * time.Second),
			Update: schema.DefaultTimeout(10 * time.Second),
			Delete: schema.DefaultTimeout(10 * time.Second),
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
	department := &jamfpro.ResourceDepartment{
		// Assuming your ResourceDepartment struct has a Name field
		Name: util.GetStringFromInterface(d.Get("name")),
	}

	// Initialize the logging subsystem for the construction operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemConstruct, hclog.Debug)

	// Serialize and pretty-print the department object as XML
	deptXML, err := xml.MarshalIndent(department, "", "  ")
	if err != nil {
		logging.Error(subCtx, logging.SubsystemConstruct, "Failed to marshal department to XML", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	// Log the pretty-printed XML
	logging.Debug(subCtx, logging.SubsystemConstruct, "Constructed Department XML", map[string]interface{}{"xml": string(deptXML)})

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

	// Initialize tflog
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)

	// Initialize variables
	var createdAttribute *jamfpro.ResponseDepartmentCreate

	// construct the resource object
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
			var apiErrorCode int
			if apiError, ok := err.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.Error(subCtx, logging.SubsystemAPI, "API Error during department creation", map[string]interface{}{
				"name":       department.Name,
				"error":      err.Error(),
				"error_code": apiErrorCode,
			})
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemCreate, "Failed to create department", map[string]interface{}{
			"error": err.Error(),
		})
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

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
// ResourceJamfProDepartmentsRead is responsible for reading the current state of a Jamf Pro Department Resource from the remote system.
func ResourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	attributeID := d.Id()
	attributeName := d.Get("name").(string)

	// Use the retry function for the read operation with appropriate timeout
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		attribute, apiErr := conn.GetDepartmentByID(attributeID)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				logging.Error(subCtx, logging.SubsystemRead, "API error encountered", map[string]interface{}{
					"id":          attributeID,
					"api_error":   apiError.Error(),
					"status_code": apiError.StatusCode,
				})

				if apiError.StatusCode == 404 {
					// Department not found, remove from Terraform state
					d.SetId("")
					logging.Info(subCtx, logging.SubsystemRead, "Department not found, removed from Terraform state", map[string]interface{}{
						"id":   attributeID,
						"name": attributeName,
					})
					return nil // Exit retry loop
				}
			} else {
				logging.Error(subCtx, logging.SubsystemRead, "Unknown error encountered", map[string]interface{}{
					"id":    attributeID,
					"error": apiErr.Error(),
				})
			}
			return retry.RetryableError(apiErr) // Retry for other API errors
		}

		// Successfully found the department, update Terraform state
		d.Set("id", attribute.ID)
		d.Set("name", attribute.Name)
		logging.Info(subCtx, logging.SubsystemRead, "Successfully read department", map[string]interface{}{
			"id":   attribute.ID,
			"name": attribute.Name,
		})

		return nil // Success, exit retry loop
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemRead, "Failed to read department after retries", map[string]interface{}{
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

	// Initialize the logging subsystem for the update operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)

	// Initialize variables
	departmentID := d.Id()
	departmentName := d.Get("name").(string)

	// construct the resource object
	department, err := constructJamfProDepartment(subCtx, d)
	if err != nil {
		logging.Error(subCtx, logging.SubsystemUpdate, "Failed to construct department for update", map[string]interface{}{
			"error": err.Error(),
			"id":    departmentID,
		})
		return diag.FromErr(fmt.Errorf("failed to construct department for update: %w", err))
	}

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateDepartmentByID(departmentID, department)
		if apiErr != nil {
			var apiErrorCode int
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.Error(subCtx, logging.SubsystemUpdate, "Failed to update department by ID, trying by name", map[string]interface{}{
				"error":      apiErr.Error(),
				"error_code": apiErrorCode,
				"id":         departmentID,
				"name":       departmentName,
			})

			_, apiErrByName := conn.UpdateDepartmentByName(departmentName, department)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.Error(subCtx, logging.SubsystemUpdate, "API error during department update by name", map[string]interface{}{
					"error":      apiErrByName.Error(),
					"error_code": apiErrByNameCode,
					"name":       departmentName,
				})
				return retry.RetryableError(apiErrByName)
			}
		}

		logging.Info(subCtx, logging.SubsystemUpdate, "Successfully updated department", map[string]interface{}{
			"name": department.Name,
			"id":   departmentID,
		})
		return nil
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemUpdate, "Failed to update department", map[string]interface{}{
			"error": err.Error(),
			"id":    departmentID,
			"name":  departmentName,
		})
		return diag.FromErr(err)
	}

	// Retry reading the department to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDepartmentsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.Error(subCtx, logging.SubsystemUpdate, "Failed to read department after update", map[string]interface{}{
				"summary": readDiags[0].Summary,
				"id":      departmentID,
			})
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemUpdate, "Failed to synchronize Terraform state after department update", map[string]interface{}{
			"error": err.Error(),
			"id":    departmentID,
		})
		return diag.FromErr(err)
	}

	return nil
}

// ResourceJamfProDepartmentsDelete is responsible for deleting a Jamf Pro Department.
func ResourceJamfProDepartmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, assertionOk := meta.(*client.APIClient)
	if !assertionOk {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the delete operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	// Initialize resource variables
	departmentID := d.Id()
	departmentName := d.Get("name").(string)

	// Use the retry function for the delete operation with appropriate timeout
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := conn.DeleteDepartmentByID(departmentID)
		if apiErr != nil {
			apiError, apiErrorOk := apiErr.(*http_client.APIError)
			if apiErrorOk && apiError.StatusCode == 404 {
				// Department not found by ID, try by name
				apiErrByName := conn.DeleteDepartmentByName(departmentName)
				if apiErrByName != nil {
					apiErrorByName, apiErrorByNameOk := apiErrByName.(*http_client.APIError)
					if apiErrorByNameOk && apiErrorByName.StatusCode == 404 {
						// Department not found by name, safe to remove from Terraform state
						d.SetId("")
						return nil // No need to retry
					}
					// Log error if deletion by name fails for reasons other than not found
					logging.Error(subCtx, logging.SubsystemDelete, "API error during department deletion by name", map[string]interface{}{
						"error": apiErrByName.Error(),
						"name":  departmentName,
					})
					return retry.RetryableError(apiErrByName)
				}
			} else {
				// Log error if deletion by ID fails for reasons other than not found
				logging.Error(subCtx, logging.SubsystemDelete, "Failed to delete department by ID", map[string]interface{}{
					"error": apiErr.Error(),
					"id":    departmentID,
				})
				return retry.RetryableError(apiErr)
			}
		}
		return nil // Department deleted successfully
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemDelete, "Failed to delete department", map[string]interface{}{
			"id":    departmentID,
			"name":  departmentName,
			"error": err.Error(),
		})
		return diag.FromErr(err)
	}

	// Log the successful removal of the department from the Terraform state
	logging.Info(subCtx, logging.SubsystemDelete, "Successfully removed department from Terraform state", map[string]interface{}{
		"id":   departmentID,
		"name": departmentName,
	})

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
