package computercheckin

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProComputerCheckinCreate is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
// Since this resource is a configuration set and not a resource that is 'created' in the traditional sense,
// this function will simply set the initial state in Terraform.
// resourceJamfProComputerCheckinCreate is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
func resourceJamfProComputerCheckinCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProComputerCheckin(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Check-In for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := client.UpdateComputerCheckinInformation(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Computer Check-In configuration after retries: %v", err))
	}

	d.SetId("jamfpro_computer_checkin_singleton")

	return append(diags, resourceJamfProComputerCheckinRead(ctx, d, meta)...)
}

// resourceJamfProComputerCheckinRead is responsible for reading the current state of the Jamf Pro computer check-in configuration.
func resourceJamfProComputerCheckinRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	var err error

	// TODO not an ID?
	d.SetId("jamfpro_computer_checkin_singleton")

	var response *jamfpro.ResourceComputerCheckin
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetComputerCheckinInformation()
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

// resourceJamfProComputerCheckinUpdate is responsible for updating the Jamf Pro computer check-in configuration.
func resourceJamfProComputerCheckinUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	checkinConfig, err := constructJamfProComputerCheckin(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Check-In for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := client.UpdateComputerCheckinInformation(checkinConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Computer Check-In configuration after retries: %v", err))
	}

	d.SetId("jamfpro_computer_checkin_singleton")

	return append(diags, resourceJamfProComputerCheckinRead(ctx, d, meta)...)
}

// resourceJamfProComputerCheckinDelete is responsible for 'deleting' the Jamf Pro computer check-in configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func resourceJamfProComputerCheckinDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
