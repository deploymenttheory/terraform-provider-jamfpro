// allowedfileextensions_resource.go
package allowedfileextensions

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"

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
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
	// Assert the meta interface to the expected client type
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProAllowedFileExtension(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Allowed File Extension: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResourceAllowedFileExtension
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateAllowedFileExtension(resource)
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

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return client.GetAllowedFileExtensionByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Allowed File Extension", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)

	if waitDiags.HasError() {
		return waitDiags
	}

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
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Attempt to fetch the resource by ID
	resource, err := client.GetAllowedFileExtensionByID(resourceIDInt)

	if err != nil {
		// Handle not found error or other errors
		return state.HandleResourceNotFoundError(err, d)
	}

	// Update the Terraform state with the fetched data from the resource
	diags = updateTerraformState(d, resource)

	// Handle any errors and return diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
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
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

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
		apiErr := client.DeleteAllowedFileExtensionByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("extension").(string)
			apiErrByName := client.DeleteAllowedFileExtensionByName(resourceName)
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
