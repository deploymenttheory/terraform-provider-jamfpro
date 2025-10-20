package managed_software_update_feature_toggle

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for initializing the Jamf Pro Managed Software Update Feature Toggle configuration in Terraform.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructManagedSoftwareUpdateFeatureToggle(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Managed Software Update Feature Toggle for create: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.UpdateManagedSoftwareUpdateFeatureToggle(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Managed Software Update Feature Toggle configuration: %v", err))
	}

	d.SetId("jamfpro_managed_software_update_feature_toggle_singleton")

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro Managed Software Update Feature Toggle configuration.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	var response *jamfpro.ResourceManagedSoftwareUpdateFeatureToggle
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetManagedSoftwareUpdateFeatureToggle()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateState(d, response)...)
}

// readWithCleanup reads a resources and states with cleanup
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads a resource without cleanup
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating the Jamf Pro Managed Software Update Feature Toggle configuration.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	response, err := constructManagedSoftwareUpdateFeatureToggle(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Managed Software Update Feature Toggle for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateManagedSoftwareUpdateFeatureToggle(response)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Managed Software Update Feature Toggle configuration: %v", err))
	}

	d.SetId("jamfpro_managed_software_update_feature_toggle_singleton")

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro Managed Software Update Feature Toggle configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
