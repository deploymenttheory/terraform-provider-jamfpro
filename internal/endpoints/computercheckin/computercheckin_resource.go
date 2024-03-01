// computercheckin_resource.go
package computercheckin

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProComputerCheckindefines the schema and RU operations for managing Jamf Pro computer checkin configuration in Terraform.
func ResourceJamfProComputerCheckin() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProComputerCheckinCreate,
		ReadContext:   ResourceJamfProComputerCheckinRead,
		UpdateContext: ResourceJamfProComputerCheckinUpdate,
		DeleteContext: ResourceJamfProComputerCheckinDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			return validateComputerCheckinDependencies(d)
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
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

// ResourceJamfProComputerCheckinCreate is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
// Since this resource is a configuration set and not a resource that is 'created' in the traditional sense,
// this function will simply set the initial state in Terraform.
// ResourceJamfProComputerCheckinCreate is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
func ResourceJamfProComputerCheckinCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	checkinConfig, err := constructJamfProComputerCheckin(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Check-In for update: %v", err))
	}

	// Update (or effectively create) the check-in configuration with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := conn.UpdateComputerCheckinInformation(checkinConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Computer Check-In configuration after retries: %v", err))
	}

	// Since this resource is a singleton, use a fixed ID to represent it in the Terraform state
	d.SetId("jamfpro_computer_checkin_singleton")

	// Read the site to ensure the Terraform state is up to date
	readDiags := ResourceJamfProComputerCheckinRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProComputerCheckinRead is responsible for reading the current state of the Jamf Pro computer check-in configuration.
func ResourceJamfProComputerCheckinRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	var resource *jamfpro.ResourceComputerCheckin

	// Read operation with retry
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetComputerCheckinInformation()
		if apiErr != nil {
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the site, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		d.SetId("") // Remove from Terraform state if unable to read after retries
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Computer Check-In with ID '%s' after retries: %v", resourceID, err))
	}

	// The constant ID "jamfpro_computer_checkin_singleton" is assigned to satisfy Terraform's requirement for an ID.
	d.SetId("jamfpro_computer_checkin_singleton")

	// Map the configuration fields from the API response to a structured map
	checkinData := map[string]interface{}{
		"check_in_frequency":                       resource.CheckInFrequency,
		"create_startup_script":                    resource.CreateStartupScript,
		"log_startup_event":                        resource.LogStartupEvent,
		"check_for_policies_at_startup":            resource.CheckForPoliciesAtStartup,
		"apply_computer_level_managed_preferences": resource.ApplyComputerLevelManagedPrefs,
		"ensure_ssh_is_enabled":                    resource.EnsureSSHIsEnabled,
		"create_login_logout_hooks":                resource.CreateLoginLogoutHooks,
		"log_username":                             resource.LogUsername,
		"check_for_policies_at_login_logout":       resource.CheckForPoliciesAtLoginLogout,
		"apply_user_level_managed_preferences":     resource.ApplyUserLevelManagedPreferences,
		"hide_restore_partition":                   resource.HideRestorePartition,
		"perform_login_actions_in_background":      resource.PerformLoginActionsInBackground,
	}

	// Set the structured map in the Terraform state
	for key, val := range checkinData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// ResourceJamfProComputerCheckinUpdate is responsible for updating the Jamf Pro computer check-in configuration.
func ResourceJamfProComputerCheckinUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	checkinConfig, err := constructJamfProComputerCheckin(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Check-In for update: %v", err))
	}

	// Update (or effectively create) the check-in configuration with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := conn.UpdateComputerCheckinInformation(checkinConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Computer Check-In configuration after retries: %v", err))
	}

	// Since this resource is a singleton, use a fixed ID to represent it in the Terraform state
	d.SetId("jamfpro_computer_checkin_singleton")

	// Read the site to ensure the Terraform state is up to date
	readDiags := ResourceJamfProComputerCheckinRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProComputerCheckinDelete is responsible for 'deleting' the Jamf Pro computer check-in configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func ResourceJamfProComputerCheckinDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Simply remove the resource from the Terraform state by setting the ID to an empty string.
	d.SetId("")

	return nil
}
