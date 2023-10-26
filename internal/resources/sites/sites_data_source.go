// sites_date_source.go
package sites

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

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
				Optional:    true,
				Description: "The unique identifier of the Jamf Pro site.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique name of the Jamf Pro site.",
				Computed:    true,
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
	conn := meta.(*client.APIClient).Conn

	var site *jamfpro.ResponseSite
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		siteName := v.(string)
		site, err = conn.GetSiteByName(siteName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch Jamf Pro site by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		siteID, convertErr := strconv.Atoi(v.(string))
		if convertErr != nil {
			return diag.FromErr(fmt.Errorf("failed to convert Jamf Pro site ID to integer: %v", convertErr))
		}
		site, err = conn.GetSiteByID(siteID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch Jamf Pro site by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	// Set the data source attributes using the fetched data
	d.SetId(fmt.Sprintf("%d", site.ID))
	d.Set("name", site.Name)

	return nil
}
