// accountgroups_data_source.go
package accountgroups

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProAccountGroup provides information about specific Jamf Pro Dock Items by their ID or Name.
func DataSourceJamfProAccountGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProAccountGroupsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the account group.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account group.",
			},
			"access_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The access level of the account group.",
			},
			"privilege_set": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The privilege set assigned to the account group.",
			},
			"site": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The site information associated with the account group if access_level is set to Site Access.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Jamf Pro Site ID. Value defaults to -1 aka not used.",
							Default:     -1,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Jamf Pro Site Name. Value defaults to 'None' aka not used",
							Computed:    true,
						},
					},
				},
			},
			"privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The privileges associated with the account group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"members": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Members of the account group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceJamfProDockItemsRead fetches the details of specific account group from Jamf Pro using either their unique Name or Id.
func dataSourceJamfProAccountGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var accountGroup *jamfpro.ResourceAccountGroup

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		accountGroupID, err := strconv.Atoi(d.Id())
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error converting id (%s) to integer: %s", d.Id(), err))
		}

		// Try fetching the account group using the ID
		accountGroup, err = conn.GetAccountGroupByID(accountGroupID)
		if err != nil {
			// Handle the APIError
			if apiError, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			return retry.RetryableError(err)
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Update the Terraform state with account group attributes
	d.Set("name", accountGroup.Name)
	d.Set("access_level", accountGroup.AccessLevel)
	d.Set("privilege_set", accountGroup.PrivilegeSet)

	// Update site information
	site := make(map[string]interface{})
	site["id"] = accountGroup.Site.ID
	site["name"] = accountGroup.Site.Name
	d.Set("site", []interface{}{site})

	// Update privileges
	privileges := make(map[string]interface{})
	privileges["jss_objects"] = accountGroup.Privileges.JSSObjects
	privileges["jss_settings"] = accountGroup.Privileges.JSSSettings
	privileges["jss_actions"] = accountGroup.Privileges.JSSActions
	privileges["recon"] = accountGroup.Privileges.Recon
	privileges["casper_admin"] = accountGroup.Privileges.CasperAdmin
	privileges["casper_remote"] = accountGroup.Privileges.CasperRemote
	privileges["casper_imaging"] = accountGroup.Privileges.CasperImaging
	d.Set("privileges", []interface{}{privileges})

	// Update members
	members := make([]interface{}, 0)
	for _, member := range accountGroup.Members {
		memberMap := map[string]interface{}{
			"id":   member.ID,
			"name": member.Name,
		}
		members = append(members, memberMap)
	}
	d.Set("members", members)

	return diags
}
