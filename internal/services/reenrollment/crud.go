package reenrollment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for initializing the Jamf Pro re-enrollment settings configuration in Terraform.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro re-enrollment Settings for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.UpdateReenrollmentSettings(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro re-enrollment Settings configuration after retries: %v", err))
	}

	d.SetId("jamfpro_reenrollment_settings_singleton")

	return append(diags, read(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro re-enrollment settings configuration.
func read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	d.SetId("jamfpro_reenrollment_settings_singleton")

	var response *jamfpro.ResourceReenrollmentSettings
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetReenrollmentSettings()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		diag.FromErr(err)
	}

	return append(diags, updateState(d, response)...)
}

// update is responsible for updating the Jamf Pro re-enrollment settings configuration.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	reenrollmentSettingsConfig, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro re-enrollment Settings for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateReenrollmentSettings(*reenrollmentSettingsConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro re-enrollment Settings configuration after retries: %v", err))
	}

	d.SetId("jamfpro_reenrollment_settings_singleton")

	return append(diags, read(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro re-enrollment settings configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")

	return nil
}
