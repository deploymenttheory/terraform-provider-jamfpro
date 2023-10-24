// computergroup_resource.go
package computergroups

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	And DeviceGroupAndOr = "and"
	Or  DeviceGroupAndOr = "or"
)

const (
	SearchTypeIs           = "is"
	SearchTypeIsNot        = "is not"
	SearchTypeLike         = "like"
	SearchTypeNotLike      = "not like"
	SearchTypeMatchesRegex = "matches regex"
	SearchTypeDoesNotMatch = "does not match regex"
)

type DeviceGroupAndOr string

func ResourceJamfProComputerGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProComputerGroupsCreate,
		ReadContext:   ResourceJamfProComputerGroupsRead,
		UpdateContext: ResourceJamfProComputerGroupsUpdate,
		DeleteContext: ResourceJamfProComputerGroupsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Optional:    true,
				Default:     false,
				Description: "Smart or static group.",
			},
			"site": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID of the site.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the site.",
						},
					},
				},
			},
			"criteria": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the smart group search criteria.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The priority of the criterion.",
						},
						"and_or": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Either 'and' or 'or'.",
							ValidateFunc: validation.StringInSlice([]string{
								string(And),
								string(Or),
							}, false),
						},
						"search_type": {
							Type:     schema.TypeString,
							Required: true,
							Description: fmt.Sprintf("The type of search operator. Allowed values are '%s', '%s', '%s', '%s', '%s', and '%s'.",
								SearchTypeIs, SearchTypeIsNot, SearchTypeLike, SearchTypeNotLike, SearchTypeMatchesRegex, SearchTypeDoesNotMatch),
							ValidateFunc: validation.StringInSlice([]string{
								SearchTypeIs,
								SearchTypeIsNot,
								SearchTypeLike,
								SearchTypeNotLike,
								SearchTypeMatchesRegex,
								SearchTypeDoesNotMatch,
							}, false),
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Search value.",
						},
						"opening_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Opening parenthesis flag.",
						},
						"closing_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Closing parenthesis flag.",
						},
					},
				},
			},
			"computers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID of the computer.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the computer.",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "MAC Address of the computer.",
						},
						"alt_mac_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alternative MAC Address of the computer.",
						},
						"serial_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Serial number of the computer.",
						},
					},
				},
			},
		},
	}
}

// ResourceJamfProComputerGroupsCreate performs tf creation operations for jamf pro computer group resources
func ResourceJamfProComputerGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn

	computerGroupName := d.Get("name").(string)
	isSmart := d.Get("is_smart").(bool)

	// Extracting site details
	var siteID int
	var siteName string

	siteList, exists := d.GetOk("site")
	if exists {
		site := siteList.([]interface{})[0].(map[string]interface{})
		siteID = site["id"].(int)
		siteName = site["name"].(string)
	}

	// Extracting criteria
	criteriaData := d.Get("criteria").([]interface{})
	var criteria []jamfpro.ComputerGroupCriterion
	for _, item := range criteriaData {
		data := item.(map[string]interface{})
		criterion := jamfpro.ComputerGroupCriterion{
			Name:         data["name"].(string),
			Priority:     data["priority"].(int),
			AndOr:        jamfpro.DeviceGroupAndOr(data["and_or"].(string)),
			SearchType:   data["search_type"].(string),
			SearchValue:  data["value"].(string),
			OpeningParen: data["opening_paren"].(bool),
			ClosingParen: data["closing_paren"].(bool),
		}
		criteria = append(criteria, criterion)
	}

	// Extracting computers
	computersData := d.Get("computers").([]interface{})
	var computers []jamfpro.ComputerGroupComputerItem
	for _, item := range computersData {
		data := item.(map[string]interface{})
		computer := jamfpro.ComputerGroupComputerItem{
			ID:            data["id"].(int),
			Name:          data["name"].(string),
			SerialNumber:  data["serial_number"].(string),
			MacAddress:    data["mac_address"].(string),
			AltMacAddress: data["alt_mac_address"].(string),
		}
		computers = append(computers, computer)
	}

	groupRequest := &jamfpro.ComputerGroupRequest{
		Name:      computerGroupName,
		IsSmart:   isSmart,
		Site:      jamfpro.Site{ID: siteID, Name: siteName},
		Criteria:  criteria,
		Computers: computers,
	}

	group, err := conn.CreateComputerGroup(groupRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the ID of the computer group in the Terraform state
	d.SetId(fmt.Sprintf("%d", group.ID))

	return ResourceJamfProComputerGroupsRead(ctx, d, meta)
}

// ResourceJamfProComputerGroupsRead uses the jamf pro sdk to read a computer group object
func ResourceJamfProComputerGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn

	var diags diag.Diagnostics

	// Initially attempt to get the computer group by its ID.
	computerGroupID, err := strconv.Atoi(d.Id())
	if err == nil {
		computerGroup, err := conn.GetComputerGroupByID(computerGroupID)
		if err != nil {
			// If there's an error fetching by ID, log a warning and continue.
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to fetch computer group by ID",
				Detail:   fmt.Sprintf("Failed to fetch computer group with ID %d: %v", computerGroupID, err),
			})
		} else {
			// If successfully fetched the computer group by ID, set the details into the state and return.
			d.Set("name", computerGroup.Name)
			d.Set("is_smart", computerGroup.IsSmart)
			d.Set("site", []interface{}{
				map[string]interface{}{
					"id":   computerGroup.Site.ID,
					"name": computerGroup.Site.Name,
				},
			})
			criteriaList := make([]interface{}, len(computerGroup.Criteria))
			for i, crit := range computerGroup.Criteria {
				criteriaList[i] = map[string]interface{}{
					"name":          crit.Name,
					"priority":      crit.Priority,
					"and_or":        crit.AndOr,
					"search_type":   crit.SearchType,
					"value":         crit.SearchValue,
					"opening_paren": crit.OpeningParen,
					"closing_paren": crit.ClosingParen,
				}
			}
			d.Set("criteria", criteriaList)

			computerList := make([]interface{}, len(computerGroup.Computers))
			for i, comp := range computerGroup.Computers {
				computerList[i] = map[string]interface{}{
					"id":              comp.ID,
					"name":            comp.Name,
					"serial_number":   comp.SerialNumber,
					"mac_address":     comp.MacAddress,
					"alt_mac_address": comp.AltMacAddress,
				}
			}
			d.Set("computers", computerList)
			return diags
		}
	}

	// If fetching by ID failed or wasn't possible, try to fetch by the name.
	computerGroupName := d.Get("name").(string)
	computerGroup, err := conn.GetComputerGroupByName(computerGroupName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to fetch computer group by name",
			Detail:   fmt.Sprintf("Failed to fetch computer group with name %s: %v", computerGroupName, err),
		})
		return diags
	}

	// Set the fetched computer group details into the state.
	d.Set("name", computerGroup.Name)
	d.Set("is_smart", computerGroup.IsSmart)
	d.Set("site", []interface{}{
		map[string]interface{}{
			"id":   computerGroup.Site.ID,
			"name": computerGroup.Site.Name,
		},
	})
	criteriaList := make([]interface{}, len(computerGroup.Criteria))
	for i, crit := range computerGroup.Criteria {
		criteriaList[i] = map[string]interface{}{
			"name":          crit.Name,
			"priority":      crit.Priority,
			"and_or":        crit.AndOr,
			"search_type":   crit.SearchType,
			"value":         crit.SearchValue,
			"opening_paren": crit.OpeningParen,
			"closing_paren": crit.ClosingParen,
		}
	}
	d.Set("criteria", criteriaList)

	computerList := make([]interface{}, len(computerGroup.Computers))
	for i, comp := range computerGroup.Computers {
		computerList[i] = map[string]interface{}{
			"id":              comp.ID,
			"name":            comp.Name,
			"serial_number":   comp.SerialNumber,
			"mac_address":     comp.MacAddress,
			"alt_mac_address": comp.AltMacAddress,
		}
	}
	d.Set("computers", computerList)
	d.Set("id", fmt.Sprintf("%d", computerGroup.ID))

	return diags
}

// ResourceJamfProComputerGroupsUpdate performs tf update operations upon jamf pro computer group resources
func ResourceJamfProComputerGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn

	var diags diag.Diagnostics

	// Create the ComputerGroupRequest object for the update.
	groupRequest := &jamfpro.ComputerGroupRequest{
		Name:    d.Get("name").(string),
		IsSmart: d.Get("is_smart").(bool),
	}

	// Extracting site details
	var siteID int
	var siteName string

	if siteList, exists := d.GetOk("site"); exists {
		site := siteList.([]interface{})[0].(map[string]interface{})
		siteID = site["id"].(int)
		siteName = site["name"].(string)
	}
	groupRequest.Site = jamfpro.Site{ID: siteID, Name: siteName}

	// Extracting criteria
	criteriaData := d.Get("criteria").([]interface{})
	for _, item := range criteriaData {
		data := item.(map[string]interface{})
		criterion := jamfpro.ComputerGroupCriterion{
			Name:         data["name"].(string),
			Priority:     data["priority"].(int),
			AndOr:        jamfpro.DeviceGroupAndOr(data["and_or"].(string)),
			SearchType:   data["search_type"].(string),
			SearchValue:  data["value"].(string),
			OpeningParen: data["opening_paren"].(bool),
			ClosingParen: data["closing_paren"].(bool),
		}
		groupRequest.Criteria = append(groupRequest.Criteria, criterion)
	}

	// Extracting computers
	computersData := d.Get("computers").([]interface{})
	for _, item := range computersData {
		data := item.(map[string]interface{})
		computer := jamfpro.ComputerGroupComputerItem{
			ID:            data["id"].(int),
			Name:          data["name"].(string),
			SerialNumber:  data["serial_number"].(string),
			MacAddress:    data["mac_address"].(string),
			AltMacAddress: data["alt_mac_address"].(string),
		}
		groupRequest.Computers = append(groupRequest.Computers, computer)
	}

	// Check if the name has changed and update accordingly
	if d.HasChange("name") {
		oldName, _ := d.GetChange("name")

		// Initially attempt to update the computer group by its ID.
		computerGroupID, err := strconv.Atoi(d.Id())
		if err == nil {
			_, err = conn.UpdateComputerGroupByID(computerGroupID, groupRequest)
		}

		// If updating by ID failed or wasn't possible, try to update by the old name.
		if err != nil {
			_, err = conn.UpdateComputerGroupByName(oldName.(string), groupRequest)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Failed to update computer group with old name %s", oldName.(string)),
					Detail:   err.Error(),
				})
				return diags
			}
		}
	}

	// Even if the update was successful, run the Read function to get the latest state and verify the update.
	readDiags := ResourceJamfProComputerGroupsRead(ctx, d, meta)
	diags = append(diags, readDiags...)

	return diags
}

// ResourceJamfProComputerGroupsDelete performs tf delete operations upon jamf pro computer group resources
func ResourceJamfProComputerGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn

	var diags diag.Diagnostics

	// Initially attempt to delete the computer group by its ID.
	computerGroupID, err := strconv.Atoi(d.Id())
	if err == nil {
		err := conn.DeleteComputerGroupByID(computerGroupID)
		if err == nil {
			// Successfully deleted the computer group by ID.
			return diags
		}

		// If there's an error deleting by ID, log a warning and continue.
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to delete computer group by ID",
			Detail:   fmt.Sprintf("Failed to delete computer group with ID %d: %v", computerGroupID, err),
		})
	}

	// If deleting by ID failed or wasn't possible, try to delete by the name.
	computerGroupName := d.Get("name").(string)
	err = conn.DeleteComputerGroupByName(computerGroupName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete computer group by name",
			Detail:   fmt.Sprintf("Failed to delete computer group with name %s: %v", computerGroupName, err),
		})
	}

	return diags
}
