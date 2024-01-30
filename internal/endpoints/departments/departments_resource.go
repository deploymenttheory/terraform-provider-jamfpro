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

// constructJamfProDepartment
func constructJamfProDepartment(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceDepartment, error) {
	department := &jamfpro.ResourceDepartment{
		Name: util.GetStringFromInterface(d.Get("name")),
	}

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemConstruct, hclog.Debug)

	// Serialize and pretty-print the department object as XML
	deptXML, err := xml.MarshalIndent(department, "", "  ")
	if err != nil {
		logging.LogTFConstructResourceXMLMarshalFailure(subCtx, JamfProResourceDepartment, err.Error())
		return nil, err
	}

	// Log the successful construction and serialization to XML
	logging.LogTFConstructedXMLResourceXML(subCtx, JamfProResourceDepartment, string(deptXML))

	return department, nil
}

// ResourceJamfProDepartmentsCreate is responsible for creating a new Jamf Pro Site in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProDepartmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
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

	// construct the resource object with retries

	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		department, err := constructJamfProDepartment(subCtx, d)
		if err != nil {
			logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDepartment, err.Error())
			return retry.NonRetryableError(err)
		}
		logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDepartment)

		createdAttribute, err := conn.CreateDepartment(department)
		if err != nil {
			var apiErrorCode int
			if apiError, ok := err.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailure(subCtx, JamfProResourceDepartment, err.Error(), apiErrorCode)
			return retry.NonRetryableError(err)
		}

		logging.LogAPICreateSuccess(subCtx, JamfProResourceDepartment, createdAttribute.ID)
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		var apiErrorCode int
		if apiError, ok := err.(*http_client.APIError); ok {
			apiErrorCode = apiError.StatusCode
		}

		logging.LogAPICreateFailure(subCtx, JamfProResourceDepartment, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// set newly created resource ID into terraform state
	d.SetId(createdAttribute.ID)

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDepartmentsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subCtx, JamfProResourceDepartment, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subCtx, JamfProResourceDepartment, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {
		logging.LogTFStateSyncSuccess(subCtx, JamfProResourceDepartment, d.Id())
	}

	return diags
}

/*
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
*/

// ResourceJamfProDepartmentsRead is responsible for reading the current state of a Jamf Pro Department Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
// ResourceJamfProDepartmentsRead reads a Jamf Pro Department resource from the remote system.
func ResourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	departmentID := d.Id()

	// read operation
	department, err := conn.GetDepartmentByID(departmentID)
	if err != nil {
		logging.LogFailedReadByID(subCtx, JamfProResourceDepartment, departmentID, err.Error())
		d.SetId("") // Remove from Terraform state
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceDepartment, departmentID)
		return diags
	}

	// Assuming successful read if no error
	logging.LogAPIReadSuccess(subCtx, JamfProResourceDepartment, departmentID)

	if err := d.Set("id", departmentID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", department.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProDepartmentsUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func ResourceJamfProDepartmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
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

	// Construct the resource object
	department, err := constructJamfProDepartment(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDepartment, err.Error())
		return diag.FromErr(fmt.Errorf("failed to construct department for update: %w", err))
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDepartment)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateDepartmentByID(departmentID, department)
		if apiErr != nil {
			var apiErrorCode int
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceDepartment, departmentID, departmentName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateDepartmentByName(departmentName, department)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceDepartment, departmentName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceDepartment, departmentID, departmentName)
		}
		return nil
	})

	if err != nil {
		logging.LogAPIUpdateFailedAfterRetry(subCtx, JamfProResourceDepartment, departmentID, departmentName, err.Error())
		return diag.FromErr(err)
	}

	// Retry reading the department to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDepartmentsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subCtx, JamfProResourceDepartment, departmentID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subCtx, JamfProResourceDepartment, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subCtx, JamfProResourceDepartment, departmentID)
	}

	return nil
}

/*
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
*/

// ResourceJamfProDepartmentsDelete is responsible for deleting a Jamf Pro Department.
func ResourceJamfProDepartmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the delete operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	// Initialize Delete operation variables
	resourceID := d.Id()
	resourceName := d.Get("name").(string)

	// Use the retry function for the delete operation with appropriate timeout
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Delete By ID
		apiErr := conn.DeleteDepartmentByID(resourceID)
		if apiErr != nil {
			var apiErrorCode int
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourceDepartment, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			// If Delete by ID fails then try Delete by Name
			apiErr = conn.DeleteDepartmentByName(resourceName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIDeleteFailureByName(subCtx, JamfProResourceDepartment, resourceName, apiErr.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceDepartment, resourceID, resourceName, err.Error())
		return diag.FromErr(err)
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourceDepartment, resourceID, resourceName)

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
