// allowedfileextensions_resource.go
package allowedfileextensions

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
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
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"extension": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// constructAllowedFileExtension creates a new ResourceAllowedFileExtension instance from Terraform data.
func constructAllowedFileExtension(d *schema.ResourceData) (*jamfpro.ResourceAllowedFileExtension, error) {
	extension := &jamfpro.ResourceAllowedFileExtension{}

	// Utilize type assertion helper functions for direct field extraction
	extension.Extension = util.GetStringFromInterface(d.Get("extension"))

	// Log the successful construction of the file extension
	log.Printf("[INFO] Successfully constructed Allowed File Extension: %s", extension.Extension)

	return extension, nil
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

// ResourceJamfProAllowedFileExtensionCreate is responsible for creating a new AllowedFileExtension in the remote system.
// The function:
// 1. Constructs the AllowedFileExtension data using the provided Terraform configuration.
// 2. Calls the API to create the AllowedFileExtension in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created AllowedFileExtension.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProAllowedFileExtensionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	var createdAllowedFileExtension *jamfpro.ResourceAllowedFileExtension
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the AllowedFileExtension
		AllowedFileExtension, err := constructAllowedFileExtension(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the Allowed File Extension for terraform create: %w", err))
		}

		// Directly call the API to create the resource
		createdAllowedFileExtension, err = conn.CreateAllowedFileExtension(AllowedFileExtension)
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
	d.SetId(strconv.Itoa(createdAllowedFileExtension.ID))

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAllowedFileExtensionRead(ctx, d, meta)
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

// ResourceJamfProAllowedFileExtensionRead is responsible for reading the current state of an Allowed File Extension Resource from the remote system.
// The function:
// 1. Fetches the Allowed File Extension's current state using its ID. If it fails, then obtain the Allowed File Extension's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the Allowed File Extension being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProAllowedFileExtensionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var allowedFileExtension *jamfpro.ResourceAllowedFileExtension

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {

		// Convert the ID from the Terraform state into an integer to be used for the API request
		allowedFileExtensionID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse allowed file extension ID: %v", convertErr))
		}

		// Try fetching the AllowedFileExtension using the ID
		var apiErr error
		allowedFileExtension, apiErr = conn.GetAllowedFileExtensionByID(allowedFileExtensionID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			allowedFileExtensionName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			allowedFileExtension, apiErr = conn.GetAllowedFileExtensionByName(allowedFileExtensionName)
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

	// Update the Terraform state with the fetched data
	if err := d.Set("extension", allowedFileExtension.Extension); err != nil {
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

func ResourceJamfProAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		accountID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse dock item ID: %v", convertErr))
		}

		// Directly call the API to delete the resource
		apiErr := conn.DeleteAccountByID(accountID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by name
			accountName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			apiErr = conn.DeleteAccountByName(accountName)
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

// ResourceJamfProAllowedFileExtensionDelete is responsible for deleting an Allowed File Extension in Jamf Pro.
// This function will delete the resource based on its ID from the Terraform state.
// If the resource cannot be found by ID, it will attempt to delete by the 'extension' attribute.
func ResourceJamfProAllowedFileExtensionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		allowedFileExtensionID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse allowed file extension ID: %v", convertErr))
		}

		// Directly call the API to delete the resource
		apiErr := conn.DeleteAllowedFileExtensionByID(allowedFileExtensionID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by extension name
			allowedFileExtensionName, ok := d.Get("extension").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'extension' as a string"))
			}

			apiErr = conn.DeleteAllowedFileExtensionByName(allowedFileExtensionName)
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
