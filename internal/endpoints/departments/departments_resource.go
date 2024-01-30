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
	logging.LogTFConstructedXMLResource(subCtx, JamfProResourceDepartment, string(deptXML))

	return department, nil
}

// ResourceJamfProDepartmentsCreate is responsible for creating a new Jamf Pro Site in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProDepartmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	var createdAttribute *jamfpro.ResponseDepartmentCreate
	var apiErrorCode int

	// Initialize the logging subsystem with the create operation context
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Construct the department object outside the retry loop to avoid reconstructing it on each retry
	department, err := constructJamfProDepartment(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDepartment, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDepartment)

	// Retry the API call to create the department in Jamf Pro
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		createdAttribute, apiErr = conn.CreateDepartment(department)
		if apiErr != nil {
			// Extract and log the API error code if available
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailure(subCtx, JamfProResourceDepartment, apiErr.Error(), apiErrorCode)
			// Return a non-retryable error to break out of the retry loop
			return retry.NonRetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final error and append it to the diagnostics
		logging.LogAPICreateFailure(subCtx, JamfProResourceDepartment, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Log successful creation of the department and set the resource ID in Terraform state
	logging.LogAPICreateSuccess(subCtx, JamfProResourceDepartment, createdAttribute.ID)
	d.SetId(createdAttribute.ID)

	// Retry reading the department to ensure the Terraform state is up to date
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDepartmentsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			// Log any read errors and return a retryable error to retry the read operation
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceDepartment, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		// Successfully read the department, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final state sync failure and append it to the diagnostics
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceDepartment, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {
		// Log successful state synchronization
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceDepartment, d.Id())
	}

	return diags
}

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
	resourceID := d.Id()
	var apiErrorCode int

	// read operation

	department, err := conn.GetDepartmentByID(resourceID)
	if err != nil {
		if apiError, ok := err.(*http_client.APIError); ok {
			apiErrorCode = apiError.StatusCode
		}
		logging.LogFailedReadByID(subCtx, JamfProResourceDepartment, resourceID, err.Error(), apiErrorCode)
		d.SetId("") // Remove from Terraform state
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceDepartment, resourceID)
		return diags
	}

	// Assuming successful read if no error
	logging.LogAPIReadSuccess(subCtx, JamfProResourceDepartment, resourceID)

	if err := d.Set("id", resourceID); err != nil {
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
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Construct the resource object
	department, err := constructJamfProDepartment(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDepartment, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDepartment)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateDepartmentByID(resourceID, department)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceDepartment, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateDepartmentByName(resourceName, department)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceDepartment, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceDepartment, resourceID, resourceName)
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceDepartment, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the department to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDepartmentsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceDepartment, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceDepartment, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceDepartment, resourceID)
	}

	return nil
}

// ResourceJamfProDepartmentsDelete is responsible for deleting a Jamf Pro Department.
func ResourceJamfProDepartmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Initialize the logging subsystem for the delete operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	// Use the retry function for the delete operation with appropriate timeout
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Delete By ID
		apiErr := conn.DeleteDepartmentByID(resourceID)
		if apiErr != nil {
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

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceDepartment, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourceDepartment, resourceID, resourceName)

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
