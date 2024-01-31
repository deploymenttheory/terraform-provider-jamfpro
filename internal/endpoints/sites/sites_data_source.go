// sites_date_source.go
package sites

import (
	"context"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProSites provides information about a specific Jamf Pro site by its ID or Name.
func DataSourceJamfProSites() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProSitesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the Jamf Pro site.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the Jamf Pro site.",
			},
		},
	}
}

// dataSourceJamfProSitesRead fetches the details of a specific Jamf Pro site
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
func dataSourceJamfProSitesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	resourceID := d.Id()
	var apiErrorCode int

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error
		logging.LogFailedReadByID(subCtx, JamfProResourceSite, resourceID, "Invalid resource ID format", 0)
		return diag.FromErr(err)
	}

	// read operation

	site, err := conn.GetSiteByID(resourceIDInt)
	if err != nil {
		if apiError, ok := err.(*http_client.APIError); ok {
			apiErrorCode = apiError.StatusCode
		}
		logging.LogFailedReadByID(subCtx, JamfProResourceSite, resourceID, err.Error(), apiErrorCode)
		return diags
	}

	// Assuming successful read if no error
	logging.LogAPIReadSuccess(subCtx, JamfProResourceSite, resourceID)

	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", site.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
