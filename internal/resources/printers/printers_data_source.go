// printers_data_source.go
package printers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

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
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var printer *jamfpro.ResourcePrinter
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		printerName := v.(string)
		printer, err = conn.GetPrinterByName(printerName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch printer by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		printerID, err := strconv.Atoi(v.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse printer ID: %v", err))
		}
		printer, err = conn.GetPrinterByID(printerID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch printer by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	if printer == nil {
		return diag.FromErr(fmt.Errorf("printer not found"))
	}

	// Set the data source attributes using the fetched data
	if err := d.Set("name", printer.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'name': %v", err))
	}
	if err := d.Set("category", printer.Category); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'category': %v", err))
	}
	if err := d.Set("uri", printer.URI); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'uri': %v", err))
	}
	if err := d.Set("cups_name", printer.CUPSName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'cups_name': %v", err))
	}
	if err := d.Set("location", printer.Location); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'location': %v", err))
	}
	if err := d.Set("model", printer.Model); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'model': %v", err))
	}
	if err := d.Set("info", printer.Info); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'info': %v", err))
	}
	if err := d.Set("notes", printer.Notes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'notes': %v", err))
	}
	if err := d.Set("make_default", printer.MakeDefault); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'make_default': %v", err))
	}
	if err := d.Set("use_generic", printer.UseGeneric); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'use_generic': %v", err))
	}
	if err := d.Set("ppd", printer.PPD); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'ppd': %v", err))
	}
	if err := d.Set("ppd_path", printer.PPDPath); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'ppd_path': %v", err))
	}
	if err := d.Set("ppd_contents", printer.PPDContents); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'ppd_contents': %v", err))
	}

	d.SetId(fmt.Sprintf("%d", printer.ID))

	return nil

}
