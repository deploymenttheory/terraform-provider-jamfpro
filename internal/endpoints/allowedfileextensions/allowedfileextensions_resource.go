// allowedfileextensions_resource.go
package allowedfileextensions

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

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
func ResourceJamfProAllowedFileExtensionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProAllowedFileExtension(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Allowed File Extension: %v", err))
	}

	var creationResponse *jamfpro.ResourceAllowedFileExtension
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateAllowedFileExtension(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Allowed File Extension '%s' after retries: %v", resource.Extension, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetAllowedFileExtensionByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Allowed File Extension", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)

	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProAllowedFileExtensionRead(ctx, d, meta)...)
}

// ResourceJamfProAllowedFileExtensionRead is responsible for reading the current state of an Allowed File Extension Resource from the remote system.
func ResourceJamfProAllowedFileExtensionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetAllowedFileExtensionByID(resourceIDInt)
	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProAllowedFileExtensionUpdate handles the update operation for an AllowedFileExtension resource in Terraform.
// Since there is no direct update API endpoint, this function will delete the existing AllowedFileExtension and create a new one.
// This approach simulates an update operation from the user's perspective in Terraform.
func ResourceJamfProAllowedFileExtensionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO come back to this
	deleteDiags := ResourceJamfProAllowedFileExtensionDelete(ctx, d, meta)
	if deleteDiags.HasError() {
		return deleteDiags
	}

	return ResourceJamfProAllowedFileExtensionCreate(ctx, d, meta)
}

// ResourceJamfProAllowedFileExtensionDelete is responsible for deleting an Allowed File Extension in Jamf Pro.
// This function will delete the resource based on its ID from the Terraform state.
// If the resource cannot be found by ID, it will attempt to delete by the 'extension' attribute.
func ResourceJamfProAllowedFileExtensionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteAllowedFileExtensionByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("extension").(string)
			apiErrByName := client.DeleteAllowedFileExtensionByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Allowed File Extension '%s' (ID: %s) after retries: %v", d.Get("extension").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
