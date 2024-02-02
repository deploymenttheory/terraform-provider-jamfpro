// computergroup_resource.go
package computergroups

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

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
	SearchTypeIs                 = "is"
	SearchTypeIsNot              = "is not"
	SearchTypeHas                = "has"
	SearchTypeDoesNotHave        = "does not have"
	SearchTypeMemberOf           = "member of"
	SearchTypeNotMemberOf        = "not member of"
	SearchTypeBeforeYYYYMMDD     = "before (yyyy-mm-dd)"
	SearchTypeAfterYYYYMMDD      = "after (yyyy-mm-dd)"
	SearchTypeMoreThanXDaysAgo   = "more than x days ago"
	SearchTypeLessThanXDaysAgo   = "less than x days ago"
	SearchTypeLike               = "like"
	SearchTypeNotLike            = "not like"
	SearchTypeGreaterThan        = "greater than"
	SearchTypeLessThan           = "less than"
	SearchTypeGreaterThanOrEqual = "greater than or equal"
	SearchTypeLessThanOrEqual    = "less than or equal"
	SearchTypeMatchesRegex       = "matches regex"
	SearchTypeDoesNotMatch       = "does not match regex"
)

type DeviceGroupAndOr string

// ResourceJamfProComputerGroups defines the schema and CRUD operations for managing Jamf Pro Computer Groups in Terraform.
func ResourceJamfProComputerGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProComputerGroupsCreate,
		ReadContext:   ResourceJamfProComputerGroupsRead,
		UpdateContext: ResourceJamfProComputerGroupsUpdate,
		DeleteContext: ResourceJamfProComputerGroupsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		CustomizeDiff: customDiffComputeGroups,
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
				Description: "Boolean selection to state if the group is a Smart group or not. If false then the group is a static group.",
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
							Description: "The ID of the site assigned to the computer group.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the site assigned to the computer group.",
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
							Description: "Name of the smart group search criteria. Can be from the Jamf built in enteries or can be an extension attribute.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The priority of the criterion.",
						},
						"and_or": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Either 'and', 'or', or blank.",
							Default:     "and",
							ValidateFunc: validation.StringInSlice([]string{
								"",
								string(And),
								string(Or),
							}, false),
						},
						"search_type": {
							Type:     schema.TypeString,
							Required: true,
							Description: fmt.Sprintf("The type of smart group search operator. Allowed values are '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'. ",
								SearchTypeIs, SearchTypeIsNot, SearchTypeHas, SearchTypeDoesNotHave, SearchTypeMemberOf, SearchTypeNotMemberOf,
								SearchTypeBeforeYYYYMMDD, SearchTypeAfterYYYYMMDD, SearchTypeMoreThanXDaysAgo, SearchTypeLessThanXDaysAgo,
								SearchTypeLike, SearchTypeNotLike, SearchTypeGreaterThan, SearchTypeLessThan, SearchTypeGreaterThanOrEqual,
								SearchTypeLessThanOrEqual, SearchTypeMatchesRegex, SearchTypeDoesNotMatch),
							ValidateFunc: validation.StringInSlice([]string{
								SearchTypeIs, SearchTypeIsNot, SearchTypeHas, SearchTypeDoesNotHave, SearchTypeMemberOf, SearchTypeNotMemberOf,
								SearchTypeBeforeYYYYMMDD, SearchTypeAfterYYYYMMDD, SearchTypeMoreThanXDaysAgo, SearchTypeLessThanXDaysAgo,
								SearchTypeLike, SearchTypeNotLike, SearchTypeGreaterThan, SearchTypeLessThan, SearchTypeGreaterThanOrEqual,
								SearchTypeLessThanOrEqual, SearchTypeMatchesRegex, SearchTypeDoesNotMatch,
							}, false),
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Search value for the smart group criteria to match with.",
						},
						"opening_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Opening parenthesis flag used during smart group construction.",
						},
						"closing_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Closing parenthesis flag used during smart group construction.",
						},
					},
				},
			},
			"computers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The ID of the computer used during static computer group construction.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Name of the computer used during static computer group construction.",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "MAC Address of the computer used during static computer group construction.",
						},
						"alt_mac_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Alternative MAC Address of the computer used during static computer group construction.",
						},
						"serial_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Serial number of the computer used during static computer group construction.",
						},
					},
				},
			},
		},
	}
}

// constructJamfProComputerGroup constructs a ResourceComputerGroup object from the provided schema data.
func constructJamfProComputerGroup(d *schema.ResourceData) (*jamfpro.ResourceComputerGroup, error) {
	group := &jamfpro.ResourceComputerGroup{}

	// Utilize type assertion helper functions for direct field extraction
	group.Name = util.GetStringFromInterface(d.Get("name"))
	group.IsSmart = util.GetBoolFromInterface(d.Get("is_smart"))

	// Handle nested "site" field
	if siteList, ok := d.GetOk("site"); ok {
		siteData, ok := siteList.([]interface{})
		if ok && len(siteData) > 0 {
			siteMap, ok := siteData[0].(map[string]interface{})
			if ok {
				group.Site.ID = util.GetIntFromInterface(siteMap["id"])
				group.Site.Name = util.GetStringFromInterface(siteMap["name"])
			}
		}
	}

	// Handle "criteria" field
	if criteria, ok := d.GetOk("criteria"); ok {
		for _, crit := range criteria.([]interface{}) {
			criterionMap := crit.(map[string]interface{})
			newCriterion := jamfpro.SharedSubsetCriteria{
				Name:         util.GetStringFromInterface(criterionMap["name"]),
				Priority:     util.GetIntFromInterface(criterionMap["priority"]),
				AndOr:        util.GetStringFromInterface(criterionMap["and_or"]),
				SearchType:   util.GetStringFromInterface(criterionMap["search_type"]),
				Value:        util.GetStringFromInterface(criterionMap["value"]),
				OpeningParen: util.GetBoolFromInterface(criterionMap["opening_paren"]),
				ClosingParen: util.GetBoolFromInterface(criterionMap["closing_paren"]),
			}
			group.Criteria.Criterion = append(group.Criteria.Criterion, newCriterion)
		}
	}

	// Handle "computers" field
	if computers, ok := d.GetOk("computers"); ok {
		for _, comp := range computers.([]interface{}) {
			computerMap := comp.(map[string]interface{})
			group.Computers = append(group.Computers, jamfpro.ComputerGroupSubsetComputer{
				ID:            util.GetIntFromInterface(computerMap["id"]),
				Name:          util.GetStringFromInterface(computerMap["name"]),
				SerialNumber:  util.GetStringFromInterface(computerMap["serial_number"]),
				MacAddress:    util.GetStringFromInterface(computerMap["mac_address"]),
				AltMacAddress: util.GetStringFromInterface(computerMap["alt_mac_address"]),
			})
		}
	}

	// Log the successful construction of the group
	log.Printf("[INFO] Successfully constructed ComputerGroup with name: %s", group.Name)

	return group, nil
}

// Helper function to generate diagnostics based on the error type.
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
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	var createdGroup *jamfpro.ResourceComputerGroup
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the computer group
		group, err := constructJamfProComputerGroup(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the computer group for terraform create: %w", err))
		}

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
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var group *jamfpro.ResourceComputerGroup

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
			groupName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

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

	// Set attributes in the Terraform state
	if err := d.Set("name", group.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", group.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	site := map[string]interface{}{
		"id":   group.Site.ID,
		"name": group.Site.Name,
	}
	if err := d.Set("site", []interface{}{site}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
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
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the computers
	computersList := make([]interface{}, len(group.Computers))
	for i, comp := range group.Computers {
		computerMap := map[string]interface{}{
			"id":              comp.ID,
			"name":            comp.Name,
			"mac_address":     comp.MacAddress,
			"alt_mac_address": comp.AltMacAddress,
			"serial_number":   comp.SerialNumber,
		}
		computersList[i] = computerMap
	}
	if err := d.Set("computers", computersList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProComputerGroupsUpdate is responsible for updating an existing Jamf Pro Computer Group on the remote system.
func ResourceJamfProComputerGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated computer group
		group, err := constructJamfProComputerGroup(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the computer group for terraform update: %w", err))
		}

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
			groupName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}
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
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

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
			groupName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}
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
