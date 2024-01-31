// printers_data_source.go
package printers

import (
	"context"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProPrinters provides information about a specific Jamf Pro printer by its ID or Name.
func DataSourceJamfProPrinters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProPrintersRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the printer.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the printer.",
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The category of the printer.",
			},
			"uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URI of the printer.",
			},
			"cups_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CUPS name of the printer.",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location of the printer.",
			},
			"model": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The model of the printer.",
			},
			"info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the printer.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Notes about the printer.",
			},
			"make_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the printer is the default printer.",
			},
			"use_generic": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the printer uses a generic driver.",
			},
			"ppd": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The PPD file name of the printer.",
			},
			"ppd_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path to the PPD file of the printer.",
			},
			"ppd_contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The contents of the PPD file.",
			},
		},
	}
}

// dataSourceJamfProPrintersRead fetches the details of a specific printer from Jamf Pro using either its unique Name or its Id.
func dataSourceJamfProPrintersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diags
	}

	// Assuming successful read if no error
	logging.LogAPIReadSuccess(subCtx, JamfProResourcePrinter, resourceID)

	// Set individual attributes in the Terraform state with error handling
	if err := d.Set("id", strconv.Itoa(printer.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", printer.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("category", printer.Category); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("uri", printer.URI); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("cups_name", printer.CUPSName); err != nil {
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
