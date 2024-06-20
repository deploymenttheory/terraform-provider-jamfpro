// computercheckin_resource.go
package computercheckin

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProComputerCheckindefines the schema and RU operations for managing Jamf Pro computer checkin configuration in Terraform.
func ResourceJamfProComputerCheckin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProComputerCheckinCreate,
		ReadContext:   resourceJamfProComputerCheckinRead,
		UpdateContext: resourceJamfProComputerCheckinUpdate,
		DeleteContext: resourceJamfProComputerCheckinDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			return validateComputerCheckinDependencies(d)
		},
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
			"check_in_frequency": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The frequency in minutes for computer check-in.",
				ValidateFunc: validation.IntInSlice([]int{60, 30, 15, 5}),
			},
			"create_startup_script": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines if a startup script should be created.",
			},
			"log_startup_event": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines if startup events should be logged.",
			},
			"check_for_policies_at_startup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, ensure that computers check for policies triggered by startup",
			},
			"apply_computer_level_managed_preferences": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Applies computer level managed preferences. Setting appears to be hard coded to false and cannot be changed. Thus field is set to computed.",
			},
			"ensure_ssh_is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable SSH (Remote Login) on computers that have it disabled.",
			},
			"create_login_logout_hooks": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines if login/logout hooks should be created. Create events that trigger each time a user logs in",
			},
			"log_username": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Log Computer Usage information at login. Log the username and date/time at login.",
			},
			"check_for_policies_at_login_logout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Checks for policies at login and logout.",
			},
			"apply_user_level_managed_preferences": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Applies user level managed preferences. Setting appears to be hard coded to false and cannot be changed. Thus field is set to computed.",
			},
			"hide_restore_partition": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines if the restore partition should be hidden. Setting appears to be hard coded to false and cannot be changed. Thus field is set to computed.",
			},
			"perform_login_actions_in_background": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Performs login actions in the background. Setting appears to be hard coded to false and cannot be changed. Thus field is set to computed.",
			},
		},
	}
}

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

	resource, err := client.GetComputerCheckinInformation()

	// TODO not an ID?
	d.SetId("jamfpro_computer_checkin_singleton")

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
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
