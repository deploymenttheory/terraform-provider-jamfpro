// accountgroups_resource.go
package accountgroups

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

// ResourceJamfProAccountGroup defines the schema and CRUD operations for managing buildings in Terraform.
func ResourceJamfProAccountGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAccountGroupCreate,
		ReadContext:   ResourceJamfProAccountGroupRead,
		UpdateContext: ResourceJamfProAccountGroupUpdate,
		DeleteContext: ResourceJamfProAccountGroupDelete,
		CustomizeDiff: customDiffAccountGroups,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The access level of the account. This can be either Full Access, or scoped to a jamf pro site with Site Access",
				ValidateFunc: validation.StringInSlice([]string{"Full Access", "Site Access"}, false),
			},
			"privilege_set": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The privilege set assigned to the account group.",
				ValidateFunc: validation.StringInSlice([]string{"Administrator", "Auditor", "Enrollment Only", "Custom"}, false),
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
			"jss_objects_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to JSS Objects.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jss_settings_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to JSS Settings.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jss_actions_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to JSS Actions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"casper_admin_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to Casper Admin.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"casper_remote_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to Casper Remote.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"casper_imaging_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to Casper Imaging.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recon_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to Recon.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

// constructJamfProAccountGroup constructs an AccountGroup object from the provided schema data.
func constructJamfProAccountGroup(d *schema.ResourceData) (*jamfpro.ResourceAccountGroup, error) {
	accountGroup := &jamfpro.ResourceAccountGroup{
		Name:         util.GetStringFromInterface(d.Get("name")),
		AccessLevel:  util.GetStringFromInterface(d.Get("access_level")),
		PrivilegeSet: util.GetStringFromInterface(d.Get("privilege_set")),
	}

	// Construct Site
	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		siteMap := v.([]interface{})[0].(map[string]interface{})
		accountGroup.Site = jamfpro.SharedResourceSite{
			ID:   util.GetIntFromInterface(siteMap["id"]),
			Name: util.GetStringFromInterface(siteMap["name"]),
		}
	}

	// Construct Privileges
	accountGroup.Privileges = jamfpro.AccountSubsetPrivileges{
		JSSObjects:    util.GetStringSliceFromInterface(d.Get("jss_objects_privileges")),
		JSSSettings:   util.GetStringSliceFromInterface(d.Get("jss_settings_privileges")),
		JSSActions:    util.GetStringSliceFromInterface(d.Get("jss_actions_privileges")),
		CasperAdmin:   util.GetStringSliceFromInterface(d.Get("casper_admin_privileges")),
		CasperRemote:  util.GetStringSliceFromInterface(d.Get("casper_remote_privileges")),
		CasperImaging: util.GetStringSliceFromInterface(d.Get("casper_imaging_privileges")),
		Recon:         util.GetStringSliceFromInterface(d.Get("recon_privileges")),
	}

	// Construct Members
	if v, ok := d.GetOk("members"); ok {
		var members jamfpro.AccountGroupSubsetMembers
		for _, member := range v.([]interface{}) {
			memberMap := member.(map[string]interface{})
			memberStruct := struct {
				ID   int    `json:"id,omitempty" xml:"id,omitempty"`
				Name string `json:"name,omitempty" xml:"name,omitempty"`
			}{
				ID:   util.GetIntFromInterface(memberMap["id"]),
				Name: util.GetStringFromInterface(memberMap["name"]),
			}
			members = append(members, memberStruct)
		}
		accountGroup.Members = members
	}

	log.Printf("[INFO] Successfully constructed Account Group with name: %s", accountGroup.Name)
	return accountGroup, nil
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

// ResourceJamfProAccountGroupCreate is responsible for creating a new Jamf Pro Script in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
// ResourceJamfProAccountGroupCreate is responsible for creating a new Jamf Pro Account Group in the remote system.
func ResourceJamfProAccountGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	var createdAccountGroup *jamfpro.ResponseAccountGroupCreated
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the account group
		accountGroup, err := constructJamfProAccountGroup(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the account group for terraform create: %w", err))
		}

		// Directly call the API to create the resource
		createdAccountGroup, err = conn.CreateAccountGroup(accountGroup)
		if err != nil {
			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(strconv.Itoa(createdAccountGroup.ID))

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAccountGroupRead(ctx, d, meta)
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

// ResourceJamfProAccountGroupRead is responsible for reading the current state of a Jamf Pro Account Group Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProAccountGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			// If fetching by ID fails, try fetching by Name
			groupName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			accountGroup, err = conn.GetAccountGroupByName(groupName)
			if err != nil {
				// Handle the APIError
				if apiError, ok := err.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(err)
			}
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

// ResourceJamfProAccountGroupUpdate is responsible for updating an existing Jamf Pro Account Group on the remote system.
func ResourceJamfProAccountGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the update operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated account group
		accountGroup, err := constructJamfProAccountGroup(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the account group for terraform update: %w", err))
		}

		// Obtain the ID from the Terraform state to be used for the API request
		accountGroupID, err := strconv.Atoi(d.Id())
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error converting id (%s) to integer: %s", d.Id(), err))
		}

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateAccountGroupByID(accountGroupID, accountGroup)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			groupName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string in update"))
			}

			_, apiErr = conn.UpdateAccountGroupByName(groupName, accountGroup)
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
		readDiags := ResourceJamfProAccountGroupRead(ctx, d, meta)
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

// ResourceJamfProAccountGroupDelete is responsible for deleting a Jamf Pro account group.
func ResourceJamfProAccountGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the delete operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Obtain the ID from the Terraform state to be used for the API request
		accountGroupID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse dock item ID: %v", convertErr))
		}

		// Directly call the API to delete the resource
		apiErr := conn.DeleteAccountGroupByID(accountGroupID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by name
			accountGroupName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			apiErr = conn.DeleteAccountGroupByName(accountGroupName)
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
