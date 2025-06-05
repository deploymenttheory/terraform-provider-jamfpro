package computercheckin

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
// Doesn't follow crud pattern as requires 2 api calls to complete operation.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructClientCheckInSettings(d)
	if err != nil {

		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Client Check-In for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.UpdateClientCheckinSettings(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {

		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Client Check-In configuration after retries: %v", err))
	}

	policyProperties, err := constructPolicyProperties(d)
	if err != nil {

		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Policy Properties for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.UpdatePolicyProperties(*policyProperties)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {

		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Policy Properties configuration after retries: %v", err))
	}

	d.SetId("jamfpro_client_checkin_singleton")

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro computer check-in configuration.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	d.SetId("jamfpro_client_checkin_singleton")

	var checkinResponse *jamfpro.ResourceClientCheckinSettings
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		checkinResponse, apiErr = client.GetClientCheckinSettings()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	diags = append(diags, updateState(d, checkinResponse)...)

	var policyResponse *jamfpro.ResourcePolicyProperties
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		policyResponse, apiErr = client.GetPolicyProperties()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	if err := d.Set("allow_network_state_change_triggers", policyResponse.AllowNetworkStateChangeTriggers); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating the Jamf Pro computer check-in configuration.
// Doesn't follow crud pattern as requires 2 api calls to complete operation.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	checkinConfig, err := constructClientCheckInSettings(d)
	if err != nil {

		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Client Check-In for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.UpdateClientCheckinSettings(*checkinConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {

		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Client Check-In configuration after retries: %v", err))
	}

	if d.HasChange("allow_network_state_change_triggers") {
		policyProperties, err := constructPolicyProperties(d)
		if err != nil {

			return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Policy Properties for update: %v", err))
		}

		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			_, apiErr := client.UpdatePolicyProperties(*policyProperties)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {

			return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Policy Properties configuration after retries: %v", err))
		}
	}

	d.SetId("jamfpro_client_checkin_singleton")

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro computer check-in configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
