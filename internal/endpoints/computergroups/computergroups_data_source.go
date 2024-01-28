// computergroup_data_source.go
package computergroups

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProComputerGroups provides information about a specific computer group in Jamf Pro.
func DataSourceJamfProComputerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProComputerGroupsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the computer group.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the Jamf Pro computer group.",
			},
			"is_smart": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Smart or static group.",
			},
			"site": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Jamf Pro Site-related settings of the computer group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the site.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the site.",
						},
					},
				},
			},
			"criteria": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the criteria.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The priority of the criterion.",
						},
						"and_or": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Either 'and' or 'or'.",
						},
						"search_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of search operator.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Criteria search value.",
						},
						"opening_paren": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Opening parenthesis flag.",
						},
						"closing_paren": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Closing parenthesis flag.",
						},
					},
				},
			},
			"computers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the computer.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the computer.",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC Address of the computer.",
						},
						"alt_mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alternative MAC Address of the computer.",
						},
						"serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Serial number of the computer.",
						},
					},
				},
			},
		},
	}
}

// DataSourceJamfProComputerGroupsRead fetches the details of a specific computer group
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
func DataSourceJamfProComputerGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var group *jamfpro.ResourceComputerGroup
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok {
		groupName, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'name' as string")
		}
		if groupName != "" {
			group, err = conn.GetComputerGroupByName(groupName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to fetch computer group by name: %v", err))
			}
		}
	} else if v, ok := d.GetOk("id"); ok {
		groupIDStr, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'id' as string")
		}
		groupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse computer group ID: %v", err))
		}
		group, err = conn.GetComputerGroupByID(groupID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch computer group by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	// Set the data source attributes using the fetched data
	d.SetId(fmt.Sprintf("%d", group.ID))

	if err := d.Set("name", group.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'name': %v", err))
	}
	if err := d.Set("is_smart", group.IsSmart); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'is_smart': %v", err))
	}
	if err := d.Set("site", []interface{}{map[string]interface{}{
		"id":   group.Site.ID,
		"name": group.Site.Name,
	}}); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'site': %v", err))
	}

	// Set the criteria
	criteriaList := make([]interface{}, len(group.Criteria.Criterion))
	for i, crit := range group.Criteria.Criterion {
		criteriaMap := map[string]interface{}{
			"name":          crit.Name,
			"priority":      crit.Priority,
			"and_or":        crit.AndOr,
			"search_type":   crit.SearchType,
			"value":         crit.Value,
			"opening_paren": crit.OpeningParen,
			"closing_paren": crit.ClosingParen,
		}
		criteriaList[i] = criteriaMap
	}

	if err := d.Set("criteria", criteriaList); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'criteria': %v", err))
	}

	// Set the computers
	computersList := make([]interface{}, len(group.Computers))
	for i, computer := range group.Computers {
		computersList[i] = map[string]interface{}{
			"id":              computer.ID,
			"name":            computer.Name,
			"mac_address":     computer.MacAddress,
			"alt_mac_address": computer.AltMacAddress,
			"serial_number":   computer.SerialNumber,
		}
	}

	if err := d.Set("computers", computersList); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'computers': %v", err))
	}

	return nil
}
