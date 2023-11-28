// scripts_date_source.go
package scripts

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProScripts provides information about a specific Jamf Pro script by its ID or Name.
func DataSourceJamfProScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProScriptsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "The unique identifier of the script.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The unique name of the script.",
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "Category of the script.",
				Computed:    true,
			},
			"filename": {
				Type:        schema.TypeString,
				Description: "Filename of the script.",
				Computed:    true,
			},
			"info": {
				Type:        schema.TypeString,
				Description: "Information to display to the administrator when the script is run.",
				Computed:    true,
			},
			"notes": {
				Type:        schema.TypeString,
				Description: "Notes to display about the script (e.g., who created it and when it was created).",
				Computed:    true,
			},
			"priority": {
				Type:        schema.TypeString,
				Description: "Execution priority of the script (Before, After, At Reboot).",
				Computed:    true,
			},
			"parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Script parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter4": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter5": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter7": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter8": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter9": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter10": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter11": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS requirements for the script.",
			},
			"script_contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contents of the script.",
			},
			"script_contents_encoded": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Encoded contents of the script.",
			},
		},
	}
}

// dataSourceJamfProScriptsRead fetches the details of a specific Jamf Pro script
// from Jamf Pro using either its unique Name or its Id. The function prioritizes the 'name' attribute over the 'id'
// attribute for fetching details. If neither 'name' nor 'id' is provided, it returns an error.
// Once the details are fetched, they are set in the data source's state.
//
// Parameters:
// - ctx: The context within which the function is called. It's used for timeouts and cancellation.
// - d: The current state of the data source.
// - meta: The meta object that can be used to retrieve the API client connection.
//
// Returns:
// - diag.Diagnostics: Returns any diagnostics (errors or warnings) encountered during the function's execution.
func DataSourceJamfProScriptsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var script *jamfpro.ResponseScript
	var err error

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		scriptName := v.(string)
		script, err = conn.GetScriptsByName(scriptName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch Jamf Pro script by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		scriptID, convertErr := strconv.Atoi(v.(string))
		if convertErr != nil {
			return diag.FromErr(fmt.Errorf("failed to convert Jamf Pro script ID to integer: %v", convertErr))
		}
		script, err = conn.GetScriptsByID(scriptID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch Jamf Pro script by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	d.SetId(fmt.Sprintf("%d", script.ID))
	d.Set("name", script.Name)
	// Set other script attributes to the data source state

	return nil
}
