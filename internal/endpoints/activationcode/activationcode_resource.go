// activationcode_resource.go
package activationcode

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProActivationCode defines the schema and CRUD operations for managing Jamf Pro activation code configuration in Terraform.
func ResourceJamfProActivationCode() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProActivationCodeCreate,
		ReadContext:   ResourceJamfProActivationCodeRead,
		UpdateContext: ResourceJamfProActivationCodeUpdate,
		DeleteContext: ResourceJamfProActivationCodeDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"organization_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the organization associated with the activation code.",
			},
			"code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The activation code.",
			},
		},
	}
}

// ResourceJamfProActivationCodeCreate is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
func ResourceJamfProActivationCodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// TODO document why this is not an ID
	d.SetId("jamfpro_activation_code_singleton")

	readDiags := ResourceJamfProActivationCodeRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProActivationCodeRead is responsible for reading the current state of the Jamf Pro computer check-in configuration.
func ResourceJamfProActivationCodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics

	resource, err := client.GetActivationCode()

	// TODO here too
	d.SetId("jamfpro_computer_checkin_singleton")

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	diags = updateTerraformState(d, resource)

	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProActivationCodeUpdate is responsible for updating the Jamf Pro computer check-in configuration.
func ResourceJamfProActivationCodeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	d.SetId("jamfpro_activation_code_singleton")

	readDiags := ResourceJamfProActivationCodeRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProActivationCodeDelete is responsible for 'deleting' the Jamf Pro computer check-in configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func ResourceJamfProActivationCodeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
