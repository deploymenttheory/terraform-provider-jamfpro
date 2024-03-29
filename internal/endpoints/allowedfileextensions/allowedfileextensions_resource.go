// allowedfileextensions_resource.go
package allowedfileextensions

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProAllowedFileExtensions defines the schema and CRUD operations for managing AllowedFileExtentionss in Terraform.
func ResourceJamfProAllowedFileExtensions() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAllowedFileExtensionCreate,
		ReadContext:   ResourceJamfProAllowedFileExtensionRead,
		UpdateContext: ResourceJamfProAllowedFileExtensionUpdate,
		DeleteContext: ResourceJamfProAllowedFileExtensionDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extension": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// ResourceJamfProAllowedFileExtensionCreate is responsible for creating a new AllowedFileExtension in the remote system.
// The function:
// 1. Constructs the AllowedFileExtension data using the provided Terraform configuration.
// 2. Calls the API to create the AllowedFileExtension in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created AllowedFileExtension.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProAllowedFileExtensionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProAllowedFileExtension(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Site: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResourceAllowedFileExtension
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateAllowedFileExtension(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Allowed File Extension '%s' after retries: %v", resource.Extension, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProAllowedFileExtensionRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProAllowedFileExtensionRead is responsible for reading the current state of an Allowed File Extension Resource from the remote system.
// The function:
// 1. Fetches the Allowed File Extension's current state using its ID. If it fails, then obtain the Allowed File Extension's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the Allowed File Extension being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProAllowedFileExtensionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	var resource *jamfpro.ResourceAllowedFileExtension

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetAllowedFileExtensionByID(resourceIDInt)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				// Return non-retryable error with a message to avoid SDK issues
				return retry.NonRetryableError(fmt.Errorf("resource not found, marked for deletion"))
			}
			// Retry for other types of errors
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	// If err is not nil, check if it's due to the resource being not found
	if err != nil {
		if err.Error() == "resource not found, marked for deletion" {
			// Resource not found, remove from Terraform state
			d.SetId("")
			// Append a warning diagnostic and return
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("Jamf Pro Allowed File Extension with ID '%s' was not found on the server and is marked for deletion from terraform state.", resourceID),
			})
			return diags
		}

		// For other errors, return an error diagnostic
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Allowed File Extension with ID '%s' after retries: %v", resourceID, err))
	}

	// Update the Terraform state with the fetched data
	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("extension", resource.Extension); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProAllowedFileExtensionUpdate handles the update operation for an AllowedFileExtension resource in Terraform.
// Since there is no direct update API endpoint, this function will delete the existing AllowedFileExtension and create a new one.
// This approach simulates an update operation from the user's perspective in Terraform.
func ResourceJamfProAllowedFileExtensionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Step 1: Delete the existing AllowedFileExtension
	deleteDiags := ResourceJamfProAllowedFileExtensionDelete(ctx, d, meta)
	if deleteDiags.HasError() {
		return deleteDiags
	}

	// Step 2: Create a new AllowedFileExtension with the updated details
	createDiags := ResourceJamfProAllowedFileExtensionCreate(ctx, d, meta)
	return createDiags
}

// ResourceJamfProAllowedFileExtensionDelete is responsible for deleting an Allowed File Extension in Jamf Pro.
// This function will delete the resource based on its ID from the Terraform state.
// If the resource cannot be found by ID, it will attempt to delete by the 'extension' attribute.
func ResourceJamfProAllowedFileExtensionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		apiErr := conn.DeleteAllowedFileExtensionByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("extension").(string)
			apiErrByName := conn.DeleteAllowedFileExtensionByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Allowed File Extension '%s' (ID: %s) after retries: %v", d.Get("extension").(string), resourceID, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
