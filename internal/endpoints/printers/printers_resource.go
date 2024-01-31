// printers_resource.go
package printers

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

// ResourceJamfProPrinters defines the schema and CRUD operations for managing Jamf Pro Printers in Terraform.
func ResourceJamfProPrinters() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProPrintersCreate,
		ReadContext:   ResourceJamfProPrintersRead,
		UpdateContext: ResourceJamfProPrintersUpdate,
		DeleteContext: ResourceJamfProPrintersDelete,
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
				Description: "The unique identifier of the printer.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the printer.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "No category assigned",
				Description: "The jamf pro category of the printer.",
			},
			"uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URI of the printer.",
			},
			"cups_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CUPS name of the printer.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of the printer.",
			},
			"model": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The model of the printer.",
			},
			"info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the printer.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notes about the printer.",
			},
			"make_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the printer is the default printer.",
			},
			"use_generic": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the printer uses a generic driver.",
			},
			"ppd": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The PPD file name of the printer.",
			},
			"ppd_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to the PPD file of the printer.",
			},
			"ppd_contents": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The contents of the PPD file.",
			},
		},
	}
}

const (
	JamfProResourcePrinter = "Printer"
)

// constructJamfProPrinter constructs a ResourcePrinter object from the provided schema data.
func constructJamfProPrinter(d *schema.ResourceData) (*jamfpro.ResourcePrinter, error) {
	printer := &jamfpro.ResourcePrinter{}

	// Utilize type assertion helper functions for direct field extraction
	printer.Name = util.GetStringFromInterface(d.Get("name"))
	printer.Category = util.GetStringFromInterface(d.Get("category"))
	printer.URI = util.GetStringFromInterface(d.Get("uri"))
	printer.CUPSName = util.GetStringFromInterface(d.Get("cups_name"))
	printer.Location = util.GetStringFromInterface(d.Get("location"))
	printer.Model = util.GetStringFromInterface(d.Get("model"))
	printer.Info = util.GetStringFromInterface(d.Get("info"))
	printer.Notes = util.GetStringFromInterface(d.Get("notes"))
	printer.MakeDefault = util.GetBoolFromInterface(d.Get("make_default"))
	printer.UseGeneric = util.GetBoolFromInterface(d.Get("use_generic"))
	printer.PPD = util.GetStringFromInterface(d.Get("ppd"))
	printer.PPDPath = util.GetStringFromInterface(d.Get("ppd_path"))
	printer.PPDContents = util.GetStringFromInterface(d.Get("ppd_contents"))

	// Log the successful construction of the printer
	log.Printf("[INFO] Successfully constructed Printer with name: %s", printer.Name)

	return printer, nil
}

// Helper function to generate diagnostics based on the error type.
func generateTFDiagsFromHTTPError(err error, d *schema.ResourceData, action string) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceName, exists := d.GetOk("name")
	if !exists {
		resourceName = "unknown"
	}

	// Handle the APIError in the diagnostic
	if apiErr, ok := err.(*http_client.APIError); ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   fmt.Sprintf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message),
		})
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   err.Error(),
		})
	}
	return diags
}

// Further CRUD function definitions would go here...

// ResourceJamfProPrintersCreate is responsible for creating a new Jamf Pro Printer in the remote system.
// The function:
// 1. Constructs the printer data using the provided Terraform configuration.
// 2. Calls the API to create the printer in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created printer.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProprintersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	var creationResponse *jamfpro.ResponsePr
	var apiErrorCode int

	// Initialize the logging subsystem with the create operation context
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Construct the printer object outside the retry loop to avoid reconstructing it on each retry
	printer, err := constructJamfProPrinter(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourcePrinter, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourcePrinter)

	// Retry the API call to create the printer in Jamf Pro
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.Createprinter(printer)
		if apiErr != nil {
			// Extract and log the API error code if available
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailure(subCtx, JamfProResourcePrinter, apiErr.Error(), apiErrorCode)
			// Return a non-retryable error to break out of the retry loop
			return retry.NonRetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final error and append it to the diagnostics
		logging.LogAPICreateFailure(subCtx, JamfProResourcePrinter, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Log successful creation of the printer and set the resource ID in Terraform state
	logging.LogAPICreateSuccess(subCtx, JamfProResourcePrinter, strconv.Itoa(creationResponse.ID))

	d.SetId(strconv.Itoa(creationResponse.ID))

	// Retry reading the printer to ensure the Terraform state is up to date
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProPrintersRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			// Log any read errors and return a retryable error to retry the read operation
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourcePrinter, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		// Successfully read the printer, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final state sync failure and append it to the diagnostics
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourcePrinter, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {
		// Log successful state synchronization
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourcePrinter, d.Id())
	}

	return diags
}

func ResourceJamfProPrintersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation.
	var createdPrinter *jamfpro.ResourcePrinter
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the printer.
		printer, err := constructJamfProPrinter(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the printer for terraform create: %w", err))
		}

		// Directly call the API to create the resource.
		createdPrinter, err = conn.CreatePrinter(printer)
		if err != nil {
			// Check if the error is an APIError.
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable.
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(strconv.Itoa(createdPrinter.ID))

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProPrintersRead(ctx, d, meta)
		if len(readDiags) > 0 {
			// If readDiags is not empty, it means there's an error, so we retry
			return retry.RetryableError(fmt.Errorf("failed to read the created resource"))
		}
		return nil
	})

	if err != nil {
		// If there's an error while updating the state for the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update state for")
	}

	return diags
}

// ResourceJamfProPrintersRead is responsible for reading the current state of a Jamf Pro Printer Resource from the remote system.
// The function:
// 1. Fetches the printer's current state using its ID. If it fails, then obtain the printer's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the printer being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProPrintersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var printer *jamfpro.ResourcePrinter
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error
		logging.LogFailedReadByID(subCtx, JamfProResourcePrinter, resourceID, "Invalid resource ID format", 0)
		return diag.FromErr(err)
	}

	// read operation

	printer, err = conn.GetPrinterByID(resourceIDInt)
	if err != nil {
		if apiError, ok := err.(*http_client.APIError); ok {
			apiErrorCode = apiError.StatusCode
		}
		logging.LogFailedReadByID(subCtx, JamfProResourcePrinter, resourceID, err.Error(), apiErrorCode)
		d.SetId("") // Remove from Terraform state
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourcePrinter, resourceID)
		return diags
	}

	// Assuming successful read if no error
	logging.LogAPIReadSuccess(subCtx, JamfProResourcePrinter, resourceID)

	// Construct the printer attributes for Terraform state
	printerAttributes := map[string]interface{}{
		"id":           printer.ID,
		"name":         printer.Name,
		"category":     printer.Category,
		"uri":          printer.URI,
		"cups_name":    printer.CUPSName,
		"location":     printer.Location,
		"model":        printer.Model,
		"info":         printer.Info,
		"notes":        printer.Notes,
		"make_default": printer.MakeDefault,
		"use_generic":  printer.UseGeneric,
		"ppd":          printer.PPD,
		"ppd_path":     printer.PPDPath,
		"ppd_contents": printer.PPDContents,
	}

	// Safely set attributes in the Terraform state
	if err := d.Set("printer", []interface{}{printerAttributes}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProPrintersUpdate is responsible for updating an existing Jamf Pro Printer on the remote system.
func ResourceJamfProPrintersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		// Handle conversion error
		logging.LogFailedReadByID(subCtx, JamfProResourcePrinter, resourceID, "Invalid resource ID format", 0)
		return diag.FromErr(err)
	}

	// Construct the resource object
	printer, err := constructJamfProPrinter(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourcePrinter, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourcePrinter)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdatePrinterByID(resourceIDInt, printer)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourcePrinter, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdatePrinterByName(resourceName, printer)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourcePrinter, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourcePrinter, resourceID, resourceName)
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourcePrinter, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the printer to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProPrintersRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourcePrinter, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourcePrinter, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourcePrinter, resourceID)
	}

	return nil
}

// ResourceJamfProPrintersDelete is responsible for deleting a Jamf Pro Printer.
func ResourceJamfProPrintersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the delete operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error
		logging.LogFailedReadByID(subCtx, JamfProResourcePrinter, resourceID, "Invalid resource ID format", 0)
		return diag.FromErr(err)
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Delete By ID
		apiErr := conn.DeletePrinterByID(resourceIDInt)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourcePrinter, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			// If Delete by ID fails then try Delete by Name
			apiErr = conn.DeletePrinterByName(resourceName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIDeleteFailureByName(subCtx, JamfProResourcePrinter, resourceName, apiErr.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourcePrinter, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourcePrinter, resourceID, resourceName)

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
