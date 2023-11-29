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
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
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
				Description: "The category of the printer.",
			},
			"uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URI of the printer.",
			},
			"CUPS_name": {
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

// constructJamfProPrinter constructs a ResponsePrinters object from the provided schema data and returns any errors encountered.
func constructJamfProPrinter(d *schema.ResourceData) (*jamfpro.ResponsePrinters, error) {
	printer := &jamfpro.ResponsePrinters{}

	fields := map[string]interface{}{
		"name":         &printer.Name,
		"category":     &printer.Category,
		"uri":          &printer.URI,
		"CUPS_name":    &printer.CUPSName,
		"location":     &printer.Location,
		"model":        &printer.Model,
		"info":         &printer.Info,
		"notes":        &printer.Notes,
		"make_default": &printer.MakeDefault,
		"use_generic":  &printer.UseGeneric,
		"ppd":          &printer.PPD,
		"ppd_path":     &printer.PPDPath,
		"ppd_contents": &printer.PPDContents,
	}

	for key, ptr := range fields {
		if v, ok := d.GetOk(key); ok {
			switch ptr := ptr.(type) {
			case *string:
				*ptr = v.(string)
			case *bool:
				*ptr = v.(bool)
			default:
				return nil, fmt.Errorf("unsupported data type for key '%s'", key)
			}
		}
	}

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
func ResourceJamfProPrintersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation.
	var createdPrinter *jamfpro.ResponsePrinters
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the printer.
		printer, err := constructJamfProPrinter(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the printer for terraform create: %w", err))
		}

		// Directly call the API to create the resource.
		createdPrinter, err = conn.CreatePrinters(printer)
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
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var printer *jamfpro.ResponsePrinters

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		printerID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse printer ID: %v", convertErr))
		}

		// Try fetching the printer using the ID
		var apiErr error
		printer, apiErr = conn.GetPrinterByID(printerID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			printerName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			printer, apiErr = conn.GetPrinterByName(printerName)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Safely set attributes in the Terraform state
	if err := d.Set("name", printer.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("category", printer.Category); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("uri", printer.URI); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("CUPS_name", printer.CUPSName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("location", printer.Location); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("model", printer.Model); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("info", printer.Info); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("notes", printer.Notes); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("make_default", printer.MakeDefault); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("use_generic", printer.UseGeneric); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("ppd", printer.PPD); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("ppd_path", printer.PPDPath); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("ppd_contents", printer.PPDContents); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags

}

// ResourceJamfProPrintersUpdate is responsible for updating an existing Jamf Pro Printer on the remote system.
func ResourceJamfProPrintersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the printer
		printer, err := constructJamfProPrinter(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the printer for terraform update: %w", err))
		}

		// Convert the ID from the Terraform state into an integer to be used for the API request
		printerID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse printer ID: %v", convertErr))
		}

		// Directly call the API to update the resource by ID
		_, apiErr := conn.UpdatePrinterByID(printerID, printer)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			printerName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			_, apiErr = conn.UpdatePrinterByName(printerName, printer)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProPrintersRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to update the Terraform state for the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	return diags
}

// ResourceJamfProPrintersDelete is responsible for deleting a Jamf Pro Printer.
func ResourceJamfProPrintersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the DELETE operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		printerID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse printer ID: %v", convertErr))
		}

		// Directly call the API to DELETE the resource
		apiErr := conn.DeletePrinterByID(printerID)
		if apiErr != nil {
			// If the DELETE by ID fails, try deleting by name
			printerName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			apiErr = conn.DeletePrinterByName(printerName)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while deleting the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
