// scripts_resource.go
package scripts

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProScripts defines the schema and CRUD operations for managing Jamf Pro Scripts in Terraform.
func ResourceJamfProScripts() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProScriptsCreate,
		ReadContext:   ResourceJamfProScriptsRead,
		UpdateContext: ResourceJamfProScriptsUpdate,
		DeleteContext: ResourceJamfProScriptsDelete,
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

// ResourceJamfProScriptsCreate is responsible for creating a new Jamf Pro Script in the remote system.
// The function:
// 1. Constructs the script data using the provided Terraform configuration.
// 2. Calls the API to create the script in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created script.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProScriptsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProScript(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Script: %v", err))
	}

	var creationResponse *jamfpro.ResponseScriptCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateScript(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Script '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(creationResponse.ID)

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	return client.GetScriptByID(id.(string))
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Script", creationResponse.ID, checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)
	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProScriptsRead(ctx, d, meta)...)
}

// ResourceJamfProScriptsRead is responsible for reading the current state of a Jamf Pro Script Resource from the remote system.
// The function:
// 1. Fetches the script's current state using its ID. If it fails then obtain script's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the script being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProScriptsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	resource, err := client.GetScriptByID(resourceID)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProScriptsUpdate is responsible for updating an existing Jamf Pro Department on the remote system.
func ResourceJamfProScriptsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProScript(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Script for update: %v", err))
	}
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateScriptByID(resourceID, resource)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			_, apiErrByName := client.UpdateScriptByName(resourceName, resource)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Script '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	return append(diags, ResourceJamfProScriptsRead(ctx, d, meta)...)
}

// ResourceJamfProScriptsDelete is responsible for deleting a Jamf Pro Department.
func ResourceJamfProScriptsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteScriptByID(resourceID)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteScriptByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Script '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
