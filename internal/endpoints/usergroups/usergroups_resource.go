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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProUserGroups defines the schema and CRUD operations for managing Jamf Pro Scripts in Terraform.
func ResourceJamfProUserGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProUserGroupCreate,
		ReadContext:   ResourceJamfProUserGroupRead,
		UpdateContext: ResourceJamfProUserGroupUpdate,
		DeleteContext: ResourceJamfProUserGroupDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
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
							Description: "Logical operator to use with the next criterion (AND/OR).",
						},
						"search_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of search to perform (e.g., equals, contains).",
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
				Computed:    true,
				Description: "The users belonging to the user group.",
				Elem: &schema.Resource{
					Schema: userGroupSubsetUserItemSchema(),
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
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The unique identifier of the user.",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The username of the user.",
		},
		"full_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The full name of the user.",
		},
		"phone_number": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The phone number of the user.",
		},
		"email_address": {
			Type:        schema.TypeString,
			Required:    true,
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
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Script: %v", err))
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
	d.SetId(creationResponse.ID)

	// Retry reading the resource to ensure the Terraform state is up to date
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProUserGroupRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		// Successfully read the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to synchronize Terraform state for Jamf Pro User Group '%s' after creation: %v", resource.Name, err))
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

	// Read operation with retry
	var userGroup *jamfpro.ResourceUserGroup

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		userGroup, apiErr = conn.GetUserGroupByID(resourceIDInt)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				// User Group not found or gone, remove from Terraform state
				return retry.NonRetryableError(fmt.Errorf("user group not found, marked for deletion"))
			}
			// Convert any other API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the user group, exit the retry loop
		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "user group not found, marked for deletion") {
			d.SetId("") // Remove from Terraform state
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "User Group not found or gone",
				Detail:   fmt.Sprintf("Jamf Pro User Group with ID '%s' was not found on the server and is marked for deletion from Terraform state.", resourceID),
			})
			return diags
		}
		// Handle the final error after all retries have been exhausted
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro User Group with ID '%s' after retries: %v", resourceID, err))
	}

	// Update Terraform state with the user group information
	if err := d.Set("id", userGroup.ID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	// Set the 'name' attribute
	if err := d.Set("name", userGroup.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'is_smart' attribute
	if err := d.Set("is_smart", userGroup.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'is_notify_on_change' attribute
	if err := d.Set("is_notify_on_change", userGroup.IsNotifyOnChange); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'site' attribute, ensuring to handle it properly as it's a nested structure
	site := []interface{}{map[string]interface{}{
		"id":   userGroup.Site.ID,
		"name": userGroup.Site.Name,
	}}
	if err := d.Set("site", site); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'criteria' attribute, also a nested structure
	criteria := make([]interface{}, len(userGroup.Criteria))
	for i, criterion := range userGroup.Criteria {
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
	if err := d.Set("criteria", criteria); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set the 'users', 'user_additions', and 'user_deletions' attributes, if applicable
	// with helper function to convert the user items
	if err := d.Set("users", convertUserItems(userGroup.Users)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("user_additions", convertUserItems(userGroup.UserAdditions)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("user_deletions", convertUserItems(userGroup.UserDeletions)); err != nil {
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
