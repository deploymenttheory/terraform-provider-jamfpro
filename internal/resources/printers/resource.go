// printers_resource.go
package printers

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProPrinters defines the schema and CRUD operations for managing Jamf Pro Printers in Terraform.
func ResourceJamfProPrinters() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
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
			"category_name": {
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
