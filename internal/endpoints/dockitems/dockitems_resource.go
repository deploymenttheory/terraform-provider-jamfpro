// dockitems_resource.go
package dockitems

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProDockItems defines the schema and CRUD operations for managing Jamf Pro Dock Items in Terraform.
func ResourceJamfProDockItems() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProDockItemsCreate,
		ReadContext:   ResourceJamfProDockItemsRead,
		UpdateContext: ResourceJamfProDockItemsUpdate,
		DeleteContext: ResourceJamfProDockItemsDelete,
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
				Description: "The unique identifier of the dock item.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the dock item.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the dock item (App/File/Folder).",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v, ok := val.(string)
					if !ok {
						errs = append(errs, fmt.Errorf("expected a string for %q but got a different type", key))
						return
					}
					validTypes := map[string]bool{
						"App":    true,
						"File":   true,
						"Folder": true,
					}
					if !validTypes[v] {
						errs = append(errs, fmt.Errorf("%q must be one of 'App', 'File', or 'Folder', got: %s", key, v))
					}
					return
				},
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path of the dock item.",
			},
			"contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contents of the dock item.",
			},
		},
	}
}

const (
	JamfProResourceDockItem = "Dock Item"
)

// constructJamfProDockItem constructs a ResourceDockItem object from the provided schema data.
func constructJamfProDockItem(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceDockItem, error) {
	// Initialize the logging subsystem for the construction operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemConstruct, hclog.Debug)

	dockItem := &jamfpro.ResourceDockItem{
		Name:     util.GetStringFromInterface(d.Get("name")),
		Type:     util.GetStringFromInterface(d.Get("type")),
		Path:     util.GetStringFromInterface(d.Get("path")),
		Contents: util.GetStringFromInterface(d.Get("contents")),
	}

	// Serialize and pretty-print the dockitem object as XML
	resourceXML, err := xml.MarshalIndent(dockItem, "", "  ")
	if err != nil {
		logging.LogTFConstructResourceXMLMarshalFailure(subCtx, JamfProResourceDockItem, err.Error())
		return nil, err
	}

	// Log the successful construction and serialization to XML
	logging.LogTFConstructedXMLResource(subCtx, JamfProResourceDockItem, string(resourceXML))

	return dockItem, nil
}

// ResourceJamfProDockItemsCreate is responsible for creating a new Jamf Pro Dock Item in the remote system.
// The function:
// 1. Constructs the dock item data using the provided Terraform configuration.
// 2. Calls the API to create the dock item in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created dock item.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProDockItemsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	var creationResponse *jamfpro.ResourceDockItem
	var apiErrorCode int
	resourceName := d.Get("name").(string)

	// Initialize the logging subsystem with the create operation context
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Construct the dockitem object outside the retry loop to avoid reconstructing it on each retry
	dockItem, err := constructJamfProDockItem(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDockItem, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDockItem)

	// Retry the API call to create the dockitem in Jamf Pro
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateDockItem(dockItem)
		if apiErr != nil {
			// Extract and log the API error code if available
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailedAfterRetry(subCtx, JamfProResourceDockItem, resourceName, apiErr.Error(), apiErrorCode)
			// Return a non-retryable error to break out of the retry loop
			return retry.NonRetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final error and append it to the diagnostics
		logging.LogAPICreateFailure(subCtx, JamfProResourceDockItem, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Log successful creation of the dockitem and set the resource ID in Terraform state
	logging.LogAPICreateSuccess(subCtx, JamfProResourceDockItem, strconv.Itoa(creationResponse.ID))

	d.SetId(strconv.Itoa(creationResponse.ID))

	// Retry reading the dockitem to ensure the Terraform state is up to date
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDockItemsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			// Log any read errors and return a retryable error to retry the read operation
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceDockItem, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		// Successfully read the dockitem, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final state sync failure and append it to the diagnostics
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceDockItem, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {
		// Log successful state synchronization
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceDockItem, d.Id())
	}

	return diags
}

// ResourceJamfProDockItemsRead is responsible for reading the current state of a Jamf Pro Dock Item Resource from the remote system.
// The function:
// 1. Fetches the dock item's current state using its ID. If it fails then obtain dock item's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the dock item being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProDockItemsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var apiErrorCode int
	var dockItem *jamfpro.ResourceDockItem
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDockItem, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Read operation with retry
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		dockItem, apiErr = conn.GetDockItemByID(resourceIDInt)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceDockItem, resourceID, apiErr.Error(), apiErrorCode)
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the script, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		d.SetId("") // Remove from Terraform state if unable to read after retries
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceDockItem, resourceID)
		return diag.FromErr(err)
	}

	// Check if dockItem data exists
	if dockItem != nil {
		// Set the fields directly in the Terraform state
		if err := d.Set("id", strconv.Itoa(dockItem.ID)); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("name", dockItem.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("type", dockItem.Type); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("path", dockItem.Path); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("contents", dockItem.Contents); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// ResourceJamfProDockItemsUpdate is responsible for updating an existing Jamf Pro Site on the remote system.
func ResourceJamfProDockItemsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDockItem, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Construct the resource object
	dockItem, err := constructJamfProDockItem(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDockItem, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDockItem)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateDockItemByID(resourceIDInt, dockItem)
		if apiErr != nil {
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceDockItem, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateDockItemByName(resourceName, dockItem)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceDockItem, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceDockItem, resourceID, resourceName)
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceDockItem, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the Site to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProDockItemsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceDockItem, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceDockItem, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceDockItem, resourceID)
	}

	return nil
}

// ResourceJamfProDockItemsDelete is responsible for deleting a Jamf Pro Site.
func ResourceJamfProDockItemsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Construct the resource object
	resource, err := constructJamfProDockItem(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Dock Item for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateDockItemByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Dock Item '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProDiskEncryptionConfigurationsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProDiskEncryptionConfigurationsDelete is responsible for deleting a Jamf Pro Disk Encryption Configuration.
func ResourceJamfProDiskEncryptionConfigurationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteDockItemByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteDockItemByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Dock Item '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
