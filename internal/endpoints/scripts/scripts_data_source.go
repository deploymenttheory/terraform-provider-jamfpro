// scripts_date_source.go
package scripts

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProScripts provides information about a specific Jamf Pro script by its ID or Name.
func DataSourceJamfProScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProScriptsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Jamf Pro unique identifier (ID) of the script.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name for the script.",
			},
		},
	}
}

// dataSourceJamfProScriptsRead fetches the details of a specific Jamf Pro script
// from Jamf Pro using either its unique Name or its Id.
func DataSourceJamfProScriptsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var script *jamfpro.ResourceScript

	// Get the distribution point ID from the data source's arguments
	resourceID := d.Get("id").(string)

	// Read operation with retry
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		script, apiErr = conn.GetScriptByID(resourceID)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceScript, resourceID, apiErr.Error(), apiErrorCode)
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the data, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		return diag.FromErr(err)
	}

	// Check if resource data exists and set the Terraform state
	if script != nil {
		d.SetId(resourceID) // Set the id in the Terraform state
		if err := d.Set("name", script.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		d.SetId("") // Data not found, unset the id in the Terraform state
	}

	return diags
}
