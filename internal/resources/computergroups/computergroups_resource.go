// computergroup_resource.go
package computergroups

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
				Computed: true,
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
							Optional:    true,
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
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Name of the smart group search criteria.",
							ValidateFunc: validateSmartGroupCriteriaNameWrapper,
						},
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The priority of the criterion.",
						},
						"and_or": {
							Type:        schema.TypeString,
							Optional:    true,
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

// constructComputerGroup constructs a ResponseComputerGroup object from the provided schema data.
func constructComputerGroup(d *schema.ResourceData) *jamfpro.ResponseComputerGroup {
	// Construct the Site struct
	siteList := d.Get("site").([]interface{})
	var site jamfpro.ComputerGroupSite
	if len(siteList) > 0 {
		// Handle the provided site data.
		siteMap := siteList[0].(map[string]interface{})
		site = jamfpro.ComputerGroupSite{
			ID:   siteMap["id"].(int),
			Name: siteMap["name"].(string),
		}
	}

	// Construct the slice of ComputerGroupCriterion
	var criteria []jamfpro.CriterionContainer
	if v, ok := d.GetOk("criteria"); ok {
		for _, crit := range v.([]interface{}) {
			criterionMap := crit.(map[string]interface{})

			var name, searchType, searchValue string
			var priority int
			var andOr jamfpro.DeviceGroupAndOr
			var openingParen, closingParen bool

			if nameValue, ok := criterionMap["name"].(string); ok {
				name = nameValue
			}
			if priorityValue, ok := criterionMap["priority"].(int); ok {
				priority = priorityValue
			}
			if andOrValue, ok := criterionMap["and_or"].(string); ok {
				andOr = jamfpro.DeviceGroupAndOr(andOrValue)
			}
			if searchTypeValue, ok := criterionMap["search_type"].(string); ok {
				searchType = searchTypeValue
			}
			if searchValueValue, ok := criterionMap["search_value"].(string); ok {
				searchValue = searchValueValue
			}
			if openingParenValue, ok := criterionMap["opening_paren"].(bool); ok {
				openingParen = openingParenValue
			}
			if closingParenValue, ok := criterionMap["closing_paren"].(bool); ok {
				closingParen = closingParenValue
			}

			criteria = append(criteria, jamfpro.CriterionContainer{
				Size: 1, // You may need to adjust this as per your implementation
				Criterion: jamfpro.ComputerGroupCriterion{
					Name:         name,
					Priority:     priority,
					AndOr:        andOr,
					SearchType:   searchType,
					SearchValue:  searchValue,
					OpeningParen: openingParen,
					ClosingParen: closingParen,
				},
			})
		}
	}

	// Construct the slice of ComputerGroupComputerItem
	var computers []jamfpro.ComputerContainer
	if v, ok := d.GetOk("computers"); ok {
		for _, comp := range v.([]interface{}) {
			computerMap := comp.(map[string]interface{})
			computers = append(computers, jamfpro.ComputerContainer{
				Size: 1, // You may need to adjust this as per your implementation
				Computer: jamfpro.ComputerGroupComputerItem{
					ID:            computerMap["id"].(int),
					Name:          computerMap["name"].(string),
					SerialNumber:  computerMap["serial_number"].(string),
					MacAddress:    computerMap["mac_address"].(string),
					AltMacAddress: computerMap["alt_mac_address"].(string),
				},
			})
		}
	}

	// Construct the ResponseComputerGroup object
	group := &jamfpro.ResponseComputerGroup{
		Name:      d.Get("name").(string),
		IsSmart:   d.Get("is_smart").(bool),
		Site:      site,
		Criteria:  criteria,
		Computers: computers,
	}

	// Log the successful construction of the group
	log.Printf("[INFO] Successfully constructed ComputerGroup with name: %s", group.Name)

	return group
}

// Helper function to generate diagnostics based on the error type
func generateTFDiagsFromHTTPError(err error, d *schema.ResourceData, action string) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceName, exists := d.GetOk("name")
	if !exists {
		resourceName = "unknown"
	}

	// Handle the APIError in the diagnostic
	if apiErr, ok := err.(*http_client.APIError); ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   fmt.Sprintf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message),
		})
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   err.Error(),
		})
	}
	return diags
}

// ResourceJamfProComputerGroupsCreate is responsible for creating a new Jamf Pro Computer Group in the remote system.
func ResourceJamfProComputerGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the create operation
	var createdGroup *jamfpro.ResponseComputerGroup
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the computer group
		group := constructComputerGroup(d)

		// Log the details of the group that is about to be created
		log.Printf("[INFO] Attempting to create ComputerGroup with name: %s", group.Name)

		// Directly call the API to create the resource
		createdGroup, err = conn.CreateComputerGroup(group)
		if err != nil {
			// Log the error from the API call
			log.Printf("[ERROR] Error creating ComputerGroup with name: %s. Error: %s", group.Name, err)

			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		// Log the response from the API call
		log.Printf("[INFO] Successfully created ComputerGroup with ID: %d and name: %s", createdGroup.ID, createdGroup.Name)

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(strconv.Itoa(createdGroup.ID))

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProComputerGroupsRead(ctx, d, meta)
		if len(readDiags) > 0 {
			// If readDiags is not empty, it means there's an error, so we retry
			return retry.RetryableError(fmt.Errorf("failed to read the created resource"))
		}
		return nil
	})

	if err != nil {
		// If there's an error while updating the state for the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update state for")
	}

	return diags
}

// ResourceJamfProComputerGroupsRead is responsible for reading the current state of a Jamf Pro Computer Group from the remote system.
func ResourceJamfProComputerGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	var group *jamfpro.ResponseComputerGroup

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		groupID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse group ID: %v", convertErr))
		}

		// Try fetching the computer group using the ID
		var apiErr error
		group, apiErr = conn.GetComputerGroupByID(groupID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			groupName := d.Get("name").(string)
			group, apiErr = conn.GetComputerGroupByName(groupName)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Safely set attributes in the Terraform state
	if err := d.Set("name", group.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", group.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("site", []interface{}{map[string]interface{}{
		"id":   group.Site.ID,
		"name": group.Site.Name,
	}}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the criteria
	criteriaList := make([]interface{}, len(group.Criteria))
	for i, crit := range group.Criteria {
		criteriaList[i] = map[string]interface{}{
			"name":          crit.Criterion.Name,
			"priority":      crit.Criterion.Priority,
			"and_or":        string(crit.Criterion.AndOr),
			"search_type":   crit.Criterion.SearchType,
			"value":         crit.Criterion.SearchValue,
			"opening_paren": crit.Criterion.OpeningParen,
			"closing_paren": crit.Criterion.ClosingParen,
		}
	}
	if err := d.Set("criteria", criteriaList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the computers
	computersList := make([]interface{}, len(group.Computers))
	for i, comp := range group.Computers {
		computersList[i] = map[string]interface{}{
			"id":              comp.Computer.ID,
			"name":            comp.Computer.Name,
			"mac_address":     comp.Computer.MacAddress,
			"alt_mac_address": comp.Computer.AltMacAddress,
			"serial_number":   comp.Computer.SerialNumber,
		}
	}
	if err := d.Set("computers", computersList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags

}

// ResourceJamfProComputerGroupsUpdate is responsible for updating an existing Jamf Pro Computer Group on the remote system.
func ResourceJamfProComputerGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated computer group
		group := constructComputerGroup(d)

		// Convert the ID from the Terraform state into an integer to be used for the API request
		groupID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse group ID: %v", convertErr))
		}

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateComputerGroupByID(groupID, group)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			groupName := d.Get("name").(string)
			_, apiErr = conn.UpdateComputerGroupByName(groupName, group)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProComputerGroupsRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to update the Terraform state for the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	return diags
}

// ResourceJamfProComputerGroupsDelete is responsible for deleting a Jamf Pro Computer Group.
func ResourceJamfProComputerGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the delete operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		groupID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse group ID: %v", convertErr))
		}

		// Directly call the API to delete the resource
		apiErr := conn.DeleteComputerGroupByID(groupID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by name
			groupName := d.Get("name").(string)
			apiErr = conn.DeleteComputerGroupByName(groupName)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while deleting the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
