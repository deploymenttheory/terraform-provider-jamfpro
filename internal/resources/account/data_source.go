// accounts_data_source.go
package account

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProAccounts provides information about specific Jamf Pro Accounts by their ID or Name.
func DataSourceJamfProAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the jamf pro account.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the jamf pro account.",
			},
		},
	}
}

// dataSourceRead fetches the details of specific account from Jamf Pro using either their unique Name or Id.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)
	userName := d.Get("name").(string)

	if resourceID == "" && userName == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceAccount
	if userName != "" {
		var apiErr error
		resource, apiErr = client.GetAccountByName(userName)
		if apiErr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Could not find Jamf Pro Account with Name '%s': %v", userName, apiErr),
			})
			d.SetId("")
			return diags
		}
	} else if resourceID != "" {
		var apiErr error
		resource, apiErr = client.GetAccountByID(resourceID)
		if apiErr != nil {
			return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Account with ID '%s': %v", resourceID, apiErr))
		}
	}

	if resource != nil {
		d.SetId(fmt.Sprintf("%d", resource.ID)) // or resource.ID if it's a string
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Account with ID '%v': %v", resource.ID, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}
