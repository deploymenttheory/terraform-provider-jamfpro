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
// Since this resource is a configuration set and not a resource that is 'created' in the traditional sense,
// this function will simply set the initial state in Terraform.
// ResourceJamfProActivationCodeCreate is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
func ResourceJamfProActivationCodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	activationCodeConfig, err := constructJamfProActivationCode(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Activation Code for update: %v", err))
	}

	// Update (or effectively create) the activation code configuration with retries
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

	// Since this resource is a singleton, use a fixed ID to represent it in the Terraform state
	d.SetId("jamfpro_activation_code_singleton")

	// Read the site to ensure the Terraform state is up to date
	readDiags := ResourceJamfProActivationCodeRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProActivationCodeRead is responsible for reading the current state of the Jamf Pro computer check-in configuration.
func ResourceJamfProActivationCodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics

	// Attempt to fetch the resource by ID
	resource, err := client.GetActivationCode()

	// The constant ID "jamfpro_computer_checkin_singleton" is assigned to satisfy Terraform's requirement for an ID.
	d.SetId("jamfpro_computer_checkin_singleton")

	if err != nil {
		// Handle not found error or other errors
		return state.HandleResourceNotFoundError(err, d)
	}

	// Update the Terraform state with the fetched data from the resource
	diags = updateTerraformState(d, resource)

	// Handle any errors and return diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProActivationCodeUpdate is responsible for updating the Jamf Pro computer check-in configuration.
func ResourceJamfProActivationCodeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	activationCodeConfig, err := constructJamfProActivationCode(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Activation Code for update: %v", err))
	}

	// Update (or effectively create) the activation code configuration with retries
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

	// Since this resource is a singleton, use a fixed ID to represent it in the Terraform state
	d.SetId("jamfpro_activation_code_singleton")

	// Read the site to ensure the Terraform state is up to date
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
	// Simply remove the resource from the Terraform state by setting the ID to an empty string.
	d.SetId("")

	return nil
}
