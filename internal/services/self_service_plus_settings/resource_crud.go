package self_service_plus_settings

import (
	"context"
	"errors"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrConstruct   = errors.New("failed to construct Jamf Pro Self Service Plus Settings for update")
	ErrApplyConfig = errors.New("failed to apply Jamf Pro Self Service Plus Settings configuration after retries")
)

// create is responsible for initializing the Jamf Pro SelfServicePlus settings configuration in Terraform.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructSelfServicePlusSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrConstruct.Error(), err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := client.UpdateSelfServicePlusSettings(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrApplyConfig.Error(), err))
	}

	d.SetId("jamfpro_self_service_plus_settings_singleton")

	return append(diags, read(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro Self Service Plus settings configuration.
func read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	d.SetId("jamfpro_self_service_plus_settings_singleton")

	var response *jamfpro.ResourceSelfServicePlusSettings
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetSelfServicePlusSettings()
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

// update is responsible for updating the Jamf Pro Self Service Plus settings configuration.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	selfServicePlusSettingsConfig, err := constructSelfServicePlusSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrConstruct.Error(), err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		apiErr := client.UpdateSelfServicePlusSettings(*selfServicePlusSettingsConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrApplyConfig.Error(), err))
	}

	d.SetId("jamfpro_self_service_plus_settings_singleton")

	return append(diags, read(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro Self Service Plus settings configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")

	return nil
}
