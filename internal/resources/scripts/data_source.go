// scripts_data_source.go
package scripts

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProScripts provides information about a specific Jamf Pro script by its ID or Name.
func DataSourceJamfProScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The Jamf Pro unique identifier (ID) of the script.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Display name for the script.",
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
				Description: "Notes to display about the script (e.g., who created it and when it was created).",
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The script can only be run on computers with these operating system versions. Each version must be separated by a comma (e.g., 10.11, 15, 16.1).",
			},
			"priority": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Execution priority of the script (BEFORE, AFTER, AT_REBOOT).",
			},
			"script_contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contents of the script. Must be non-compiled and in an accepted format.",
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

// dataSourceRead fetches the details of a specific Jamf Pro script
// from Jamf Pro
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	searchID := d.Get("id").(string)
	searchName := d.Get("name").(string)

	if searchID == "" && searchName == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var scriptsList *jamfpro.ResponseScriptsList
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		scriptsList, apiErr = client.GetScripts("")
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch list of scripts: %v", err))
	}

	var matchedID string
	if searchID != "" {
		for _, script := range scriptsList.Results {
			if script.ID == searchID {
				matchedID = searchID
				break
			}
		}
	} else {
		for _, script := range scriptsList.Results {
			if script.Name == searchName {
				matchedID = script.ID
				break
			}
		}
	}

	if matchedID == "" {
		return diag.FromErr(fmt.Errorf("no script found matching the provided criteria"))
	}

	var resource *jamfpro.ResourceScript
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetScriptByID(matchedID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read script with ID '%s': %v", matchedID, err))
	}

	d.SetId(matchedID)
	return updateState(d, resource)
}
