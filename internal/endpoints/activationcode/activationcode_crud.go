package activationcode

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProActivationCodeCreate is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
func resourceJamfProActivationCodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	activationCodeConfig, err := constructJamfProActivationCode(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Activation Code for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := client.UpdateActivationCode(activationCodeConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Activation Code configuration after retries: %v", err))...)
	}

	// TODO document why this is not an ID
	d.SetId("jamfpro_activation_code_singleton")

	return append(diags, resourceJamfProActivationCodeRead(ctx, d, meta)...)
}

// resourceJamfProActivationCodeRead is responsible for reading the current state of the Jamf Pro computer check-in configuration.
func resourceJamfProActivationCodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// TODO here too
	d.SetId("jamfpro_computer_checkin_singleton")

	var response *jamfpro.ResourceActivationCode
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetActivationCode()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return append(diags, updateTerraformState(d, response)...)
}

// resourceJamfProActivationCodeUpdate is responsible for updating the Jamf Pro computer check-in configuration.
func resourceJamfProActivationCodeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	activationCodeConfig, err := constructJamfProActivationCode(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Activation Code for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := client.UpdateActivationCode(activationCodeConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Activation Code configuration after retries: %v", err))
	}

	// TODO and here
	d.SetId("jamfpro_activation_code_singleton")

	return append(diags, resourceJamfProActivationCodeRead(ctx, d, meta)...)
}

// resourceJamfProActivationCodeDelete is responsible for 'deleting' the Jamf Pro computer check-in configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func resourceJamfProActivationCodeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
