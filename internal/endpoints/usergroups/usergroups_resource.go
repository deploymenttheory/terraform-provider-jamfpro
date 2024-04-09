// usergroups_object.go
package usergroups

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
	And                    UserGroupAndOr = "and"
	Or                     UserGroupAndOr = "or"
	SearchTypeIs                          = "is"
	SearchTypeIsNot                       = "is not"
	SearchTypeLike                        = "like"
	SearchTypeNotLike                     = "not like"
	SearchTypeMatchesRegex                = "matches regex"
	SearchTypeDoesNotMatch                = "does not match regex"
	SearchTypeMemberOf                    = "member of"
	SearchTypeNotMemberOf                 = "not member of"
)

type UserGroupAndOr string

// ResourceJamfProUserGroups defines the schema and CRUD operations for managing Jamf Pro Scripts in Terraform.
func ResourceJamfProUserGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProUserGroupCreate,
		ReadContext:   ResourceJamfProUserGroupRead,
		UpdateContext: ResourceJamfProUserGroupUpdate,
		DeleteContext: ResourceJamfProUserGroupDelete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the user group.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user group.",
			},
			"is_smart": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the user group is a smart group.",
			},
			"is_notify_on_change": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if notifications are sent on change.",
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
							Description: "The unique identifier of the site.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the site.",
						},
					},
				},
				Description: "The site associated with the user group.",
			},
			"criteria": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The criteria used for defining the smart user group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the criterion.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
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
							Optional: true,
							Description: fmt.Sprintf("The type of user smart group search operator. Allowed values are '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'.",
								string(SearchTypeIs), string(SearchTypeIsNot), string(SearchTypeLike),
								string(SearchTypeNotLike), string(SearchTypeMatchesRegex), string(SearchTypeDoesNotMatch),
								string(SearchTypeMemberOf), string(SearchTypeNotMemberOf)),
							ValidateFunc: validation.StringInSlice([]string{
								string(SearchTypeIs), string(SearchTypeIsNot), string(SearchTypeLike),
								string(SearchTypeNotLike), string(SearchTypeMatchesRegex), string(SearchTypeDoesNotMatch),
								string(SearchTypeMemberOf), string(SearchTypeNotMemberOf),
							}, false),
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value to search for.",
						},
						"opening_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if there is an opening parenthesis before this criterion, denoting the start of a grouped expression.",
						},
						"closing_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if there is a closing parenthesis after this criterion, denoting the end of a grouped expression.",
						},
					},
				},
			},
			"users": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A block representing the users belonging to the user group.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of jamf pro user object ID's for use within a static group.",
						},
					},
				},
			},
			"user_additions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Users added to the user group.",
				Elem: &schema.Resource{
					Schema: userGroupSubsetUserItemSchema(),
				},
			},
			"user_deletions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Users removed from the user group.",
				Elem: &schema.Resource{
					Schema: userGroupSubsetUserItemSchema(),
				},
			},
		},
	}
}

func userGroupSubsetUserItemSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The unique identifier of the user.",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The username of the user.",
		},
		"full_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The full name of the user.",
		},
		"phone_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The phone number of the user.",
		},
		"email_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The email address of the user.",
		},
	}
}

// ResourceJamfProUserGroupCreate is responsible for creating a new Jamf Pro User Group in the remote system.
// The function:
// 1. Constructs the User Group data using the provided Terraform configuration.
// 2. Calls the API to create the User Group in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created User Group.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProUserGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProUserGroup(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro User Group: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResponseUserGroupCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateUserGroup(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro User Group '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return apiclient.Conn.GetUserGroupByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro User Group", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.JamfProPropagationDelay)*time.Second)
	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProUserGroupRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProUserGroupRead is responsible for reading the current state of a Jamf Pro User Group Resource from the remote system.
// The function:
// 1. Fetches the user group's current state using its ID. If it fails, it tries to obtain the user group's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the user group being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProUserGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	resource, err := conn.GetUserGroupByID(resourceIDInt)

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
						Detail:   fmt.Sprintf("Jamf Pro User Group resource with ID '%s' was not found and has been removed from the Terraform state.", resourceID),
					},
				}
			}
		}
		// For other errors, or if this is a create operation, return a diagnostic error
		return diag.FromErr(err)
	}

	// Update the Terraform state with the fetched data
	if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", resource.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_notify_on_change", resource.IsNotifyOnChange); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'site' attribute in the state only if it's not empty (i.e., not default values)
	site := []interface{}{}

	if resource.Site.ID != -1 || resource.Site.Name != "None" {
		site = append(site, map[string]interface{}{
			"id":   resource.Site.ID,
			"name": resource.Site.Name,
		})
	}
	if len(site) > 0 {
		if err := d.Set("site", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// 'criteria' attribute
	criteria := make([]interface{}, len(resource.Criteria))
	for i, criterion := range resource.Criteria {
		criteria[i] = map[string]interface{}{
			"name":          criterion.Name,
			"priority":      criterion.Priority,
			"and_or":        criterion.AndOr,
			"search_type":   criterion.SearchType,
			"value":         criterion.Value,
			"opening_paren": criterion.OpeningParen,
			"closing_paren": criterion.ClosingParen,
		}
	}
	d.Set("criteria", criteria)

	// Set the user id's only if the group is not smart
	if !resource.IsSmart {
		var userIDStrList []string
		for _, user := range resource.Users {
			userIDStrList = append(userIDStrList, strconv.Itoa(user.ID))
		}

		if err := d.Set("users", []interface{}{
			map[string]interface{}{
				"id": userIDStrList,
			},
		}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if err := d.Set("user_additions", convertUserItems(resource.UserAdditions)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("user_deletions", convertUserItems(resource.UserDeletions)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProUserGroupUpdate is responsible for updating an existing Jamf Pro Printer on the remote system.
func ResourceJamfProUserGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	resource, err := constructJamfProUserGroup(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro User Group for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateUserGroupByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro User Group '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProUserGroupRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProUserGroupDelete is responsible for deleting a Jamf Pro User Group.
func ResourceJamfProUserGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		apiErr := conn.DeleteUserGroupByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteUserGroupByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro User Group '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
