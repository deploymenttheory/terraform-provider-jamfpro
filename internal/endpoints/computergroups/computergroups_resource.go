// computergroup_resource.go
package computergroups

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	And                          DeviceGroupAndOr = "and"
	Or                           DeviceGroupAndOr = "or"
	SearchTypeIs                                  = "is"
	SearchTypeIsNot                               = "is not"
	SearchTypeHas                                 = "has"
	SearchTypeDoesNotHave                         = "does not have"
	SearchTypeMemberOf                            = "member of"
	SearchTypeNotMemberOf                         = "not member of"
	SearchTypeBeforeYYYYMMDD                      = "before (yyyy-mm-dd)"
	SearchTypeAfterYYYYMMDD                       = "after (yyyy-mm-dd)"
	SearchTypeMoreThanXDaysAgo                    = "more than x days ago"
	SearchTypeLessThanXDaysAgo                    = "less than x days ago"
	SearchTypeLike                                = "like"
	SearchTypeNotLike                             = "not like"
	SearchTypeGreaterThan                         = "greater than"
	SearchTypeMoreThan                            = "more than"
	SearchTypeLessThan                            = "less than"
	SearchTypeGreaterThanOrEqual                  = "greater than or equal"
	SearchTypeLessThanOrEqual                     = "less than or equal"
	SearchTypeMatchesRegex                        = "matches regex"
	SearchTypeDoesNotMatch                        = "does not match regex"
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
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
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
							Description: fmt.Sprintf("The type of smart group search operator. Allowed values are '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'. ",
								SearchTypeIs, SearchTypeIsNot, SearchTypeHas, SearchTypeDoesNotHave, SearchTypeMemberOf, SearchTypeNotMemberOf,
								SearchTypeBeforeYYYYMMDD, SearchTypeAfterYYYYMMDD, SearchTypeMoreThanXDaysAgo, SearchTypeLessThanXDaysAgo,
								SearchTypeLike, SearchTypeNotLike, SearchTypeGreaterThan, SearchTypeMoreThan, SearchTypeLessThan, SearchTypeGreaterThanOrEqual,
								SearchTypeLessThanOrEqual, SearchTypeMatchesRegex, SearchTypeDoesNotMatch),
							ValidateFunc: validation.StringInSlice([]string{
								SearchTypeIs, SearchTypeIsNot, SearchTypeHas, SearchTypeDoesNotHave, SearchTypeMemberOf, SearchTypeNotMemberOf,
								SearchTypeBeforeYYYYMMDD, SearchTypeAfterYYYYMMDD, SearchTypeMoreThanXDaysAgo, SearchTypeLessThanXDaysAgo,
								SearchTypeLike, SearchTypeNotLike, SearchTypeGreaterThan, SearchTypeMoreThan, SearchTypeLessThan, SearchTypeGreaterThanOrEqual,
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

// ResourceJamfProComputerGroupsCreate is responsible for creating a new Jamf Pro Computer Group in the remote system.
func ResourceJamfProComputerGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProComputerGroup(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Group: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResourceComputerGroup
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateComputerGroup(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Computer Group '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return apiclient.Conn.GetComputerGroupByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Computer Group", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second, apiclient.EnableCookieJar)

	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProComputerGroupsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProComputerGroupsRead is responsible for reading the current state of a Jamf Pro Computer Group from the remote system.
func ResourceJamfProComputerGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Attempt to fetch the resource by ID
	resource, err := conn.GetComputerGroupByID(resourceIDInt)

	if err != nil {
		// Skip resource state removal if this is a create operation
		if !d.IsNewResource() {
			// If the error is a "not found" error, remove the resource from the state
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "410") {
				d.SetId("") // Remove the resource from Terraform state
				return diag.Diagnostics{
					{
						Severity: diag.Warning,
						Summary:  "Resource not found",
						Detail:   fmt.Sprintf("Jamf Pro Computer Group resource with ID '%s' was not found and has been removed from the Terraform state.", resourceID),
					},
				}
			}
		}
		// For other errors, or if this is a create operation, return a diagnostic error
		return diag.FromErr(err)
	}

	// Update the Terraform state with the fetched data
	if resource != nil {
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("is_smart", resource.IsSmart); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		site := map[string]interface{}{
			"id":   resource.Site.ID,
			"name": resource.Site.Name,
		}
		if err := d.Set("site", []interface{}{site}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		// Set the criteria
		criteriaList := make([]interface{}, len(resource.Criteria.Criterion))
		for i, crit := range resource.Criteria.Criterion {
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

		// Set the computers only if the group is not smart
		if !resource.IsSmart {
			computersList := make([]interface{}, len(resource.Computers))
			for i, comp := range resource.Computers {
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
		}
	}

	return diags
}

// ResourceJamfProComputerGroupsUpdate is responsible for updating an existing Jamf Pro Computer Group on the remote system.
func ResourceJamfProComputerGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Construct the resource object
	resource, err := constructJamfProComputerGroup(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Group for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateComputerGroupByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Computer Group '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProComputerGroupsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProComputerGroupsDelete is responsible for deleting a Jamf Pro Computer Group.
func ResourceJamfProComputerGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteComputerGroupByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteComputerGroupByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Computer Group '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
