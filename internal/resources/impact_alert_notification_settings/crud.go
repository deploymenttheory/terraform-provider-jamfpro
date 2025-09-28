package impact_alert_notification_settings

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
	ErrConstruct   = errors.New("failed to construct Jamf Pro Impact Alert Notification Settings for update")
	ErrApplyConfig = errors.New("failed to apply Jamf Pro Impact Alert Notification Settings configuration after retries")
)

// create is responsible for initializing the Jamf Pro Impact Alert Notification Settings configuration in Terraform.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructImpactAlertNotificationSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrConstruct.Error(), err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := client.UpdateImpactAlertNotificationSettings(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrApplyConfig.Error(), err))
	}

	d.SetId("jamfpro_impact_alert_notification_settings_singleton")

	return append(diags, read(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro Impact Alert Notification Settings configuration.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	d.SetId("jamfpro_impact_alert_notification_settings_singleton")

	var response *jamfpro.ResourceImpactAlertNotificationSettings
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetImpactAlertNotificationSettings()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return append(diags, updateState(d, response)...)
}

// update is responsible for updating the Jamf Pro Impact Alert Notification Settings configuration.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	impactAlertSettingsConfig, err := constructImpactAlertNotificationSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrConstruct.Error(), err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		apiErr := client.UpdateImpactAlertNotificationSettings(*impactAlertSettingsConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrApplyConfig.Error(), err))
	}

	d.SetId("jamfpro_impact_alert_notification_settings_singleton")

	return append(diags, read(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro Impact Alert Notification Settings configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
