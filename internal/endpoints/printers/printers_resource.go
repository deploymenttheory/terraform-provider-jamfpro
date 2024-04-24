// printers_resource.go
package printers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"

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
		CustomizeDiff: validateJamfProResourcePrinterDataFields,
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
				Default:     "/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/PrintCore.framework/Resources/Generic.ppd",
			},
			"ppd_contents": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The contents of the PPD file.",
			},
		},
	}
}

// ResourceJamfProPrintersCreate is responsible for creating a new Jamf Pro Printer in the remote system.
// The function:
// 1. Constructs the printer data using the provided Terraform configuration.
// 2. Calls the API to create the printer in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created printer.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProPrintersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProPrinter(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Printer: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResponsePrinterCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreatePrinter(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Printer '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return apiclient.Conn.GetPrinterByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Printer", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second, apiclient.EnableCookieJar)
	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProPrintersRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProPrintersRead is responsible for reading the current state of a Jamf Pro Printer Resource from the remote system.
// The function:
// 1. Fetches the printer's current state using its ID. If it fails, then obtain the printer's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the printer being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProPrintersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Attempt to fetch the resource by ID
	resource, err := apiclient.Conn.GetPrinterByID(resourceIDInt)

	if err != nil {
		// Handle not found error or other errors
		return state.HandleResourceNotFoundError(err, d)
	}

	// Update the Terraform state with the fetched data from the resource
	diags = updateTerraformState(d, resource)

	// Handle any errors and return diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProPrintersUpdate is responsible for updating an existing Jamf Pro Printer on the remote system.
func ResourceJamfProPrintersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	resource, err := constructJamfProPrinter(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Printer for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdatePrinterByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Printer '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProPrintersRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProPrintersDelete is responsible for deleting a Jamf Pro Printer.
func ResourceJamfProPrintersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		apiErr := conn.DeletePrinterByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeletePrinterByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Printer '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
