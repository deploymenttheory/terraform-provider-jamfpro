// apiroles_resource.go
package apiroles

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
)

// ResourceJamfProAPIRoles defines the schema for managing Jamf Pro API Roles in Terraform.
func ResourceJamfProAPIRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAPIRolesCreate,
		ReadContext:   ResourceJamfProAPIRolesRead,
		UpdateContext: ResourceJamfProAPIRolesUpdate,
		DeleteContext: ResourceJamfProAPIRolesDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Jamf API Role.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the Jamf API Role.",
			},
			"privileges": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of privileges associated with the Jamf API Role.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: func(val interface{}, key string) ([]string, []error) {
						return validateResourceApiRolesDataFields(val, key)
					},
				},
			},
		},
	}
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

// constructAPIRoleFromSchema constructs an APIRole object from the provided schema data.
func constructAPIRoleFromSchema(d *schema.ResourceData) *jamfpro.APIRole {
	// Extract the display name and privileges from the schema
	displayName := d.Get("display_name").(string)
	privilegesSet := d.Get("privileges").(*schema.Set)
	privilegesInterface := privilegesSet.List()

	// Convert privileges from []interface{} to []string
	privileges := make([]string, len(privilegesInterface))
	for i, v := range privilegesInterface {
		privileges[i] = v.(string)
	}

	// Construct the APIRole object
	apiRole := &jamfpro.APIRole{
		DisplayName: displayName,
		Privileges:  privileges,
	}

	// Log the successful construction of the API Role
	log.Printf("[INFO] Successfully constructed APIRole with display name: %s", displayName)

	return apiRole
}

// ResourceJamfProAPIRolesCreate handles the creation of a Jamf Pro API Role.
// The function:
// 1. Constructs the API role data using the provided Terraform configuration.
// 2. Calls the API to create the role in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created role.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProAPIRolesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	var createdRole *jamfpro.APIRole
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the API role
		role := constructAPIRoleFromSchema(d)

		// Log the details of the role that is about to be created
		log.Printf("[INFO] Attempting to create APIRole with display name: %s", role.DisplayName)

		// Directly call the API to create the resource
		createdRole, err = conn.CreateJamfApiRole(role)
		if err != nil {
			// Log the error from the API call
			log.Printf("[ERROR] Error creating APIRole with display name: %s. Error: %s", role.DisplayName, err)

			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		// Log the response from the API call
		log.Printf("[INFO] Successfully created APIRole with ID: %s and display name: %s", createdRole.ID, createdRole.DisplayName)

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(createdRole.ID)

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAPIRolesRead(ctx, d, meta)
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

// ResourceJamfProAPIRolesRead handles reading a Jamf Pro API Role from the remote system.
// The function:
// 1. Tries to fetch the API role based on the ID from the Terraform state.
// 2. If fetching by ID fails, attempts to fetch it by the display name.
// 3. Updates the Terraform state with the fetched data.
func ResourceJamfProAPIRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Retrieve the ID and display name of the API role from the Terraform state
	roleIDString := d.Id()
	displayName, _ := d.GetOk("display_name")

	// Convert roleIDString to an integer
	roleID, err := strconv.Atoi(roleIDString)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to convert role ID to integer: %v", err))
	}

	// Use the retry function for the read operation
	var fetchedRole *jamfpro.APIRole
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Try to fetch the role by ID
		fetchedRole, err = conn.GetJamfApiRolesByID(roleID)
		if err != nil {
			log.Printf("[WARN] Error reading APIRole with ID: %d. Error: %s. Trying by display name: %s", roleID, err, displayName)

			// If fetching by ID fails, try fetching by display name
			fetchedRole, err = conn.GetJamfApiRolesNameById(displayName.(string))
			if err != nil {
				// Log the error from the second API call
				log.Printf("[ERROR] Error reading APIRole with display name: %s. Error: %s", displayName, err)
				return retry.NonRetryableError(err)
			}
		}

		// Log the response from the successful API call
		log.Printf("[INFO] Successfully read APIRole with ID: %d and display name: %s", roleID, fetchedRole.DisplayName)

		return nil
	})

	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Update the Terraform state with the fetched data
	if err := d.Set("id", fetchedRole.ID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'id': %v", err))
	}
	if err := d.Set("display_name", fetchedRole.DisplayName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'display_name': %v", err))
	}
	if err := d.Set("privileges", fetchedRole.Privileges); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'privileges': %v", err))
	}

	return diags
}

// ResourceJamfProAPIRolesUpdate handles updating a Jamf Pro API Role.
// The function:
// 1. Constructs the updated API role data using the provided Terraform configuration.
// 2. Calls the API to update the role in Jamf Pro.
// 3. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProAPIRolesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		// Construct the updated API role
		role := constructAPIRoleFromSchema(d)

		// Convert the ID from the Terraform state into a string to be used for the API request
		roleID := d.Id()

		// Log the details of the role that is about to be updated
		log.Printf("[INFO] Attempting to update APIRole with ID: %s and display name: %s", roleID, role.DisplayName)

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateJamfApiRoleByID(roleID, role)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(apiErr)
		}

		// Log the successful update
		log.Printf("[INFO] Successfully updated APIRole with ID: %s and display name: %s", roleID, role.DisplayName)

		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAPIRolesRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to read the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the state for the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update state for")
	}

	return diags
}

// ResourceJamfProAPIRolesDelete handles the deletion of a Jamf Pro API Role.
func ResourceJamfProAPIRolesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the **DELETE** operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Retrieve the ID of the API role from the Terraform state
		roleID := d.Id()

		// Log the details of the role that is about to be deleted
		log.Printf("[INFO] Attempting to **DELETE** APIRole with ID: %s", roleID)

		// Directly call the API to **DELETE** the resource
		apiErr := conn.DeleteJamfApiRoleByID(roleID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(apiErr)
		}

		// Log the successful deletion
		log.Printf("[INFO] Successfully deleted APIRole with ID: %s", roleID)

		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while deleting the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Remove the resource from the Terraform state
	d.SetId("")

	return diags
}
