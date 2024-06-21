// scripts_resource.go
package scripts

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProScripts defines the schema and CRUD operations for managing Jamf Pro Scripts in Terraform.
func ResourceJamfProScripts() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProScriptsCreate,
		ReadContext:   resourceJamfProScriptsRead,
		UpdateContext: resourceJamfProScriptsUpdate,
		DeleteContext: resourceJamfProScriptsDelete,
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
				Description: "The Jamf Pro unique identifier (ID) of the script.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name for the script.",
			},
			"category_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Jamf Pro unique identifier (ID) of the category. Optional. Category ID can be used in isolation or in tandem with category_name.",
			},
			"category_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the category to add the script to. Optional. Category name can be used with category_id or not at all.",
			},
			"info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Information to display to the administrator when the script is run.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notes to display about the script (e.g., who created it and when it was created).",
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The script can only be run on computers with these operating system versions. Each version must be separated by a comma (e.g., 10.11, 15, 16.1).",
			},
			"priority": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Execution priority of the script (BEFORE, AFTER, AT_REBOOT).",
				ValidateFunc: validation.StringInSlice([]string{"BEFORE", "AFTER", "AT_REBOOT"}, false),
			},
			"script_contents": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Contents of the script. Must be non-compiled and in an accepted format.",
			},
			"parameter4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script parameter label 4",
			},
			"parameter5": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script parameter label 5",
			},
			"parameter6": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script parameter label 6",
			},
			"parameter7": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script parameter label 7",
			},
			"parameter8": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script parameter label 8",
			},
			"parameter9": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script parameter label 9",
			},
			"parameter10": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script parameter label 10",
			},
			"parameter11": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script parameter label 11",
			},
		},
	}
}
