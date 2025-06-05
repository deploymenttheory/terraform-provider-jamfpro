// scripts_data_source.go
package scripts

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProScripts provides information about a specific script in Jamf Pro.
func DataSourceJamfProScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the script.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the script.",
			},
			"category_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script Category",
			},
			"info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Information to display to the administrator when the script is run.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Notes to display about the script.",
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The script can only be run on computers with these operating system versions.",
			},
			"priority": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Execution priority of the script (BEFORE, AFTER, AT_REBOOT).",
			},
			"script_contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contents of the script.",
			},
			"parameter4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script parameter label 4",
			},
			"parameter5": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script parameter label 5",
			},
			"parameter6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script parameter label 6",
			},
			"parameter7": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script parameter label 7",
			},
			"parameter8": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script parameter label 8",
			},
			"parameter9": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script parameter label 9",
			},
			"parameter10": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script parameter label 10",
			},
			"parameter11": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Script parameter label 11",
			},
		},
	}
}

// dataSourceRead fetches the details of a specific script from Jamf Pro using either its unique Name or its Id.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	resourceID := d.Get("id").(string)
	name := d.Get("name").(string)

	if resourceID == "" && name == "" {

		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceScript
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if name != "" {
			resource, apiErr = client.GetScriptByName(name)
		} else {
			resource, apiErr = client.GetScriptByID(resourceID)
		}

		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if name != "" {
			lookupMethod = "name"
			lookupValue = name
		}

		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Script with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")

		return diag.FromErr(fmt.Errorf("the Jamf Pro Script was not found"))
	}

	d.SetId(resource.ID)
	return updateState(d, resource)
}
