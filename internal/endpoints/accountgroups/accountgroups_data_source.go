// accountgroups_data_source.go
package accountgroups

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

// DataSourceJamfProAccountGroup provides information about specific Jamf Pro Dock Items by their ID or Name.
func DataSourceJamfProAccountGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProAccountGroupsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
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
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceAccountGroup, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// read operation
	accountGroup, err := conn.GetAccountGroupByID(resourceIDInt)
	if err != nil {
		if apiError, ok := err.(*http_client.APIError); ok {
			apiErrorCode = apiError.StatusCode
		}
		logging.LogFailedReadByID(subCtx, JamfProResourceAccountGroup, resourceID, err.Error(), apiErrorCode)

		return diags
	}
	/*
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
		for _, memberStruct := range accountGroup.Members {
			member := memberStruct.User // Access the User field
			memberMap := map[string]interface{}{
				"id":   member.ID,
				"name": member.Name,
			}
			members = append(members, memberMap)
		}
		d.Set("members", members)
	*/
	if err := d.Set("name", accountGroup.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("access_level", accountGroup.AccessLevel); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("privilege_set", accountGroup.PrivilegeSet); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update site information
	site := make(map[string]interface{})
	site["id"] = accountGroup.Site.ID
	site["name"] = accountGroup.Site.Name
	if err := d.Set("site", []interface{}{site}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update privileges
	privileges := make(map[string]interface{})
	privileges["jss_objects"] = accountGroup.Privileges.JSSObjects
	privileges["jss_settings"] = accountGroup.Privileges.JSSSettings
	privileges["jss_actions"] = accountGroup.Privileges.JSSActions
	privileges["recon"] = accountGroup.Privileges.Recon
	privileges["casper_admin"] = accountGroup.Privileges.CasperAdmin
	privileges["casper_remote"] = accountGroup.Privileges.CasperRemote
	privileges["casper_imaging"] = accountGroup.Privileges.CasperImaging
	if err := d.Set("privileges", []interface{}{privileges}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update members
	members := make([]interface{}, 0)
	for _, memberStruct := range accountGroup.Members {
		member := memberStruct.User // Access the User field
		memberMap := map[string]interface{}{
			"id":   member.ID,
			"name": member.Name,
		}
		members = append(members, memberMap)
	}
	if err := d.Set("members", members); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Check if there were any errors and return the diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
}
