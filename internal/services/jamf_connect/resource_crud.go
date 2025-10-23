package jamf_connect

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for initializing the Jamf Pro Jamf Connect configuration in Terraform.
// create doesn't follow standard crud pattern as this config is pre-existing in Jamf Pro and we must
// reference the resource by it's pre-existing UUID to update the the jamf connect update config.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Connect config profile for create: %v", err))
	}

	configProfileUUID := d.Get("config_profile_uuid").(string)

	var updatedProfile *jamfpro.ResourceJamfConnectConfigProfile
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		updatedProfile, apiErr = client.UpdateJamfConnectConfigProfileByConfigProfileUUID(
			configProfileUUID,
			resource,
		)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Connect config profile after retries: %v", err))
	}

	d.SetId(updatedProfile.UUID)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro Jamf Connect configuration.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	uuid := d.Id()

	var targetProfile *jamfpro.ResourceJamfConnectConfigProfile
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr := client.GetJamfConnectConfigProfiles(nil)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		for _, profile := range response.Results {
			if profile.UUID == uuid {
				targetProfile = &profile
				break
			}
		}

		if targetProfile == nil {
			return retry.RetryableError(fmt.Errorf("jamf connect config profile with UUID %s not found", uuid))
		}

		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateState(d, targetProfile)...)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating the Jamf Pro Jamf Connect configuration.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Connect config profile for update: %v", err))
	}

	configProfileUUID := d.Get("config_profile_uuid").(string)

	var updatedProfile *jamfpro.ResourceJamfConnectConfigProfile
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		updatedProfile, apiErr = client.UpdateJamfConnectConfigProfileByConfigProfileUUID(
			configProfileUUID,
			resource,
		)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Connect config profile after retries: %v", err))
	}

	d.SetId(updatedProfile.UUID)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro Jamf Connect configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
