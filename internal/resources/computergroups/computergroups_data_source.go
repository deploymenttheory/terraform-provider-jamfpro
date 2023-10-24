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
				Required:    true,
				Description: "The unique name of the Jamf Pro computer group.",
			},
			"is_smart": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Smart or static group.",
			},
			"site": {
				Type:     schema.TypeList,
				Computed: true,
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
							Description: "Operator.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Search value.",
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

// DataSourceJamfProComputerGroupsRead retrieves details about a specific computer group in Jamf Pro
// using the provided name or ID. The function prioritizes fetching details using the computer group's
// unique Name if provided. Otherwise, it uses the Id.
//
// The function populates all the fields from the schema to the Terraform state, ensuring all the
// details of the computer group are available for use in the Terraform configuration.
//
// The Jamf Pro API client is used to interact with the Jamf Pro instance to fetch the necessary details.
//
// Params:
// - ctx: The current context.
// - d: The Terraform resource data which contains information about the resource's attributes.
// - meta: The provider meta object, which contains a pre-configured Jamf Pro API client.
//
// Returns:
// - diag.Diagnostics: A list of diagnostic messages that provide information, warnings, or errors
//   encountered during the read operation.

func DataSourceJamfProComputerGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn

	var computerGroup *jamfpro.ComputerGroup
	var err error

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		computerGroupName := v.(string)
		computerGroup, err = conn.GetComputerGroupByName(computerGroupName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch computer group by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok && v.(string) != "" {
		computerGroupID, convertErr := strconv.Atoi(v.(string))
		if convertErr != nil {
			return diag.FromErr(fmt.Errorf("failed to convert computer group ID to integer: %v", convertErr))
		}
		computerGroup, err = conn.GetComputerGroupByID(computerGroupID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch computer group by ID: %v", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("either 'name' or 'id' must be specified"))
	}

	if computerGroup == nil {
		return diag.FromErr(fmt.Errorf("computer group not found"))
	}

	// Set values to the state
	d.SetId(fmt.Sprintf("%d", computerGroup.ID))
	d.Set("name", computerGroup.Name)
	d.Set("is_smart", computerGroup.IsSmart)

	// Set site values
	site := []interface{}{
		map[string]interface{}{
			"id":   computerGroup.Site.ID,
			"name": computerGroup.Site.Name,
		},
	}
	d.Set("site", site)

	// Set criteria values
	var criteriaList []interface{}
	for _, crit := range computerGroup.Criteria {
		criteriaMap := map[string]interface{}{
			"name":          crit.Name,
			"priority":      crit.Priority,
			"and_or":        string(crit.AndOr),
			"search_type":   crit.SearchType,
			"value":         crit.SearchValue,
			"opening_paren": crit.OpeningParen,
			"closing_paren": crit.ClosingParen,
		}
		criteriaList = append(criteriaList, criteriaMap)
	}
	d.Set("criteria", criteriaList)

	// Set computer values
	var computerList []interface{}
	for _, comp := range computerGroup.Computers {
		computerMap := map[string]interface{}{
			"id":              comp.ID,
			"name":            comp.Name,
			"mac_address":     comp.MacAddress,
			"alt_mac_address": comp.AltMacAddress,
			"serial_number":   comp.SerialNumber,
		}
		computerList = append(computerList, computerMap)
	}
	d.Set("computers", computerList)

	return nil
}
