package categories

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProCategoriesCreate is responsible for creating a new Jamf Pro Category in the remote system.
func resourceJamfProCategoriesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return append(diags, resourceJamfProCategoriesReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProCategoriesRead is responsible for reading the current state of a Jamf Pro Category Resource from the remote system.
func resourceJamfProCategoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	var response *jamfpro.ResourceCategory
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetCategoryByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, state.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateTerraformState(d, response)...)
}

func resourceJamfProCategoriesReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProCategoriesRead(ctx, d, meta, true)
}

func resourceJamfProCategoriesReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProCategoriesRead(ctx, d, meta, false)
}

// resourceJamfProCategoriesUpdate is responsible for updating an existing Jamf Pro Category on the remote system.
func resourceJamfProCategoriesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("final attempt to update Category '%s' failed: %v", resourceName, err))...)
	}

	return append(diags, resourceJamfProCategoriesReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProCategoriesDelete is responsible for deleting a Jamf Pro Category.
func resourceJamfProCategoriesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	d.SetId("")

	return diags
}
