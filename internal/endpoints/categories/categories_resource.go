// categories_resource.go
package categories

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProCategories defines the schema and CRUD operations for managing Jamf Pro Categories in Terraform.
func ResourceJamfProCategories() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProCategoriesCreate,
		ReadContext:   ResourceJamfProCategoriesRead,
		UpdateContext: ResourceJamfProCategoriesUpdate,
		DeleteContext: ResourceJamfProCategoriesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(15 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the category.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro category.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     9,
				Description: "The priority of the Jamf Pro category.",
			},
		},
	}
}

// ResourceJamfProCategoriesCreate is responsible for creating a new Jamf Pro Category in the remote system.
func ResourceJamfProCategoriesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProCategory(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Category Group: %v", err))
	}

	var creationResponse *jamfpro.ResponseCategoryCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateCategory(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Category '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(creationResponse.ID)

	// region watifor
	// TODO remove waitfor I think?
	// // Wait for the resource to be fully available before reading it
	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	return client.GetCategoryByID(id.(string))
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Category", creationResponse.ID, checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)

	// if waitDiags.HasError() {
	// 	return waitDiags
	// }
	// endregion

	return append(diags, ResourceJamfProCategoriesRead(ctx, d, meta)...)
}

// ResourceJamfProCategoriesRead is responsible for reading the current state of a Jamf Pro Category Resource from the remote system.
func ResourceJamfProCategoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()
	resource, err := client.GetCategoryByID(resourceID)

	// TODO review this logic
	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProCategoriesUpdate is responsible for updating an existing Jamf Pro Category on the remote system.
func ResourceJamfProCategoriesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)

	Category, err := constructJamfProCategory(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error constructing Jamf Pro Category '%s': %v", resourceName, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateCategoryByID(resourceID, Category)
		if apiErr == nil {
			return nil
		}

		// TODO rid of this by name stuff
		_, apiErrByName := client.UpdateCategoryByName(resourceName, Category)
		if apiErrByName != nil {
			return retry.RetryableError(fmt.Errorf("failed to update Category '%s' by ID '%s' and by name due to errors: %v, %v", resourceName, resourceID, apiErr, apiErrByName))
		}

		return nil
	})

	// TODO move this up?
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("final attempt to update Category '%s' failed: %v", resourceName, err))...)
	}

	// TODO what is this log?
	hclog.FromContext(ctx).Info(fmt.Sprintf("Successfully updated Category '%s' with ID '%s'", resourceName, resourceID))

	return append(diags, ResourceJamfProCategoriesRead(ctx, d, meta)...)
}

// ResourceJamfProCategoriesDelete is responsible for deleting a Jamf Pro Category.
func ResourceJamfProCategoriesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteCategoryByID(resourceID)
		if apiErr != nil {
			apiErrByName := client.DeleteCategoryByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(fmt.Errorf("failed to delete Category '%s' by ID '%s' and by name due to errors: %v, %v", resourceName, resourceID, apiErr, apiErrByName))
			}
		}
		return nil
	})

	// TODO move this up?
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("final attempt to delete Category '%s' failed: %v", resourceName, err))...)
	}

	// TODO what's this log?
	hclog.FromContext(ctx).Info(fmt.Sprintf("Successfully deleted Category '%s' with ID '%s'", resourceName, resourceID))

	d.SetId("")

	return diags
}
