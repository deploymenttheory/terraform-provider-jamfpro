package advancedmobiledevicesearches

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro mobile device Search in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return diag.FromErr(fmt.Errorf("failed to construct Advanced Mobile Search: %v", err))
	}

	var creationResponse *jamfpro.ResponseAdvancedMobileDeviceSearchCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateAdvancedMobileDeviceSearch(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return diag.FromErr(fmt.Errorf("failed to create Advanced Mobile Search '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(creationResponse.ID)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of a Jamf Pro mobile device Search from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetAdvancedMobileDeviceSearchByID,
		updateState,
	)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// update is responsible for updating an existing Jamf Pro mobile device Search on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := construct(d)
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return diag.FromErr(fmt.Errorf("failed to construct Advanced Mobile Device Search for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateAdvancedMobileDeviceSearchByID(resourceID, *resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return diag.FromErr(fmt.Errorf("failed to update Advanced Mobile Device Search '%s' (ID: %s) after retries: %v", resource.Name, resourceID, err))
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for deleting a Jamf Pro AdvancedMobileDeviceSearch.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteAdvancedMobileDeviceSearchByID,
	)
}
