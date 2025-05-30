package accountdrivenuserenrollmentsettings

import (
	"context"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for initializing the Jamf Pro Account Driven User Enrollment Settings in Terraform.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	settings, err := construct(d)
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Account Driven User Enrollment Settings: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.UpdateADUESessionTokenSettings(*settings)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Account Driven User Enrollment Settings: %v", err))
	}

	d.SetId("jamfpro_account_driven_user_enrollment_settings_singleton")

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro Account Driven User Enrollment Settings.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	var err error

	var response *jamfpro.ResourceADUETokenSettings
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetADUESessionTokenSettings()
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

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating the Jamf Pro Account Driven User Enrollment Settings.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	if d.HasChanges("enabled", "expiration_interval_days", "expiration_interval_seconds") {
		settings, err := construct(d)
		if err != nil {
			//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
			return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Account Driven User Enrollment Settings: %v", err))
		}

		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			_, apiErr := client.UpdateADUESessionTokenSettings(*settings)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
			return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Account Driven User Enrollment Settings: %v", err))
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro Account Driven User Enrollment Settings.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// No actual deletion operation needed as this is a singleton configuration resource
	log.Printf("[DEBUG] Removing Account Driven User Enrollment Settings from state")

	// Remove from state (settings aren't actually deletable)
	d.SetId("")
	return nil
}
