// computercheckin_resource.go
package computercheckin

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

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
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
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
				Required:    true,
				Description: "Determines if a startup script should be created.",
			},
			"log_startup_event": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Determines if startup events should be logged.",
			},
			"check_for_policies_at_startup": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "If set to true, ensure that computers check for policies triggered by startup",
			},
			"apply_computer_level_managed_preferences": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Applies computer level managed preferences. Setting appears to be hard coded to false and cannot be changed. Thus field is set to computed.",
			},
			"ensure_ssh_is_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable SSH (Remote Login) on computers that have it disabled.",
			},
			"create_login_logout_hooks": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Determines if login/logout hooks should be created. Create events that trigger each time a user logs in",
			},
			"log_username": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Log Computer Usage information at login. Log the username and date/time at login.",
			},
			"check_for_policies_at_login_logout": {
				Type:        schema.TypeBool,
				Required:    true,
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
			"display_status_to_user": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Displays status to the user.",
			},
		},
	}
}

// constructComputerCheckin constructs a ResourceComputerCheckin object from the provided schema data.
func constructComputerCheckin(d *schema.ResourceData) (*jamfpro.ResourceComputerCheckin, error) {
	checkin := &jamfpro.ResourceComputerCheckin{}

	// Utilize type assertion helper functions for direct field extraction
	checkin.CheckInFrequency = util.GetIntFromInterface(d.Get("check_in_frequency"))
	checkin.CreateStartupScript = util.GetBoolFromInterface(d.Get("create_startup_script"))
	checkin.LogStartupEvent = util.GetBoolFromInterface(d.Get("log_startup_event"))
	checkin.CheckForPoliciesAtStartup = util.GetBoolFromInterface(d.Get("check_for_policies_at_startup"))
	// Note: "apply_computer_level_managed_preferences" is computed, not set directly
	checkin.EnsureSSHIsEnabled = util.GetBoolFromInterface(d.Get("ensure_ssh_is_enabled"))
	checkin.CreateLoginLogoutHooks = util.GetBoolFromInterface(d.Get("create_login_logout_hooks"))
	checkin.LogUsername = util.GetBoolFromInterface(d.Get("log_username"))
	checkin.CheckForPoliciesAtLoginLogout = util.GetBoolFromInterface(d.Get("check_for_policies_at_login_logout"))
	// Note: "apply_user_level_managed_preferences", "hide_restore_partition", and "perform_login_actions_in_background" are computed, not set directly
	checkin.DisplayStatusToUser = util.GetBoolFromInterface(d.Get("display_status_to_user"))

	// Log the successful construction of the checkin configuration
	log.Printf("[INFO] Successfully constructed ComputerCheckin")

	return checkin, nil
}

// Helper function to generate diagnostics based on the error type.
func generateTFDiagsFromHTTPError(err error, d *schema.ResourceData, action string) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceName, exists := d.GetOk("name")
	if !exists {
		resourceName = "unknown"
	}

	// Handle the APIError in the diagnostic
	if apiErr, ok := err.(*http_client.APIError); ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   fmt.Sprintf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message),
		})
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   err.Error(),
		})
	}
	return diags
}

// ResourceJamfProComputerCheckinCreate is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
// Since this resource is a configuration set and not a resource that is 'created' in the traditional sense,
// this function will simply set the initial state in Terraform.
// ResourceJamfProComputerCheckinCreate is responsible for initializing the Jamf Pro computer check-in configuration in Terraform.
func ResourceJamfProComputerCheckinCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Construct the computer check-in configuration
	checkinConfig, err := constructComputerCheckin(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Log the attempt to initialize the configuration
	log.Printf("[INFO] Initializing Jamf Pro Computer Checkin Configuration")

	// Call the API to set the initial state (using update func as the only api option available)
	apiErr := conn.UpdateComputerCheckinInformation(checkinConfig)
	if apiErr != nil {
		return diag.FromErr(fmt.Errorf("failed to initialize Computer Checkin Configuration in Jamf Pro: %w", apiErr))
	}

	// Set a constant ID to satisfy Terraform's requirement for a resource ID
	d.SetId("jamfpro_computer_checkin_singleton")

	// Perform a read operation to sync the current state from Jamf Pro
	readDiags := ResourceJamfProComputerCheckinRead(ctx, d, meta)
	if len(readDiags) > 0 {
		return readDiags
	}

	return diags
}

// ResourceJamfProComputerCheckinRead is responsible for reading the current state of the Jamf Pro computer check-in configuration.
func ResourceJamfProComputerCheckinRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Fetch the computer check-in configuration using the API client
	var checkinConfig *jamfpro.ResourceComputerCheckin
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		checkinConfig, err = conn.GetComputerCheckinInformation()
		if err != nil {
			// Handle the APIError
			if apiError, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			return retry.RetryableError(err)
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// The constant ID "jamfpro_computer_checkin_singleton" is assigned to satisfy Terraform's requirement for an ID.
	d.SetId("jamfpro_computer_checkin_singleton")

	// Map the configuration fields from the API response to a structured map
	checkinData := map[string]interface{}{
		"check_in_frequency":                       checkinConfig.CheckInFrequency,
		"create_startup_script":                    checkinConfig.CreateStartupScript,
		"log_startup_event":                        checkinConfig.LogStartupEvent,
		"check_for_policies_at_startup":            checkinConfig.CheckForPoliciesAtStartup,
		"apply_computer_level_managed_preferences": checkinConfig.ApplyComputerLevelManagedPrefs,
		"ensure_ssh_is_enabled":                    checkinConfig.EnsureSSHIsEnabled,
		"create_login_logout_hooks":                checkinConfig.CreateLoginLogoutHooks,
		"log_username":                             checkinConfig.LogUsername,
		"check_for_policies_at_login_logout":       checkinConfig.CheckForPoliciesAtLoginLogout,
		"apply_user_level_managed_preferences":     checkinConfig.ApplyUserLevelManagedPreferences,
		"hide_restore_partition":                   checkinConfig.HideRestorePartition,
		"perform_login_actions_in_background":      checkinConfig.PerformLoginActionsInBackground,
		"display_status_to_user":                   checkinConfig.DisplayStatusToUser,
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
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated computer check-in configuration
		checkinConfig, err := constructComputerCheckin(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the computer check-in configuration for terraform update: %w", err))
		}

		// Directly call the SDK API to update the configuration
		apiErr := conn.UpdateComputerCheckinInformation(checkinConfig)
		if apiErr != nil {
			// Wrap the error with additional context for clarity
			return retry.NonRetryableError(fmt.Errorf("failed to update Computer Checkin settings in Jamf Pro: %w", apiErr))
		}

		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProComputerCheckinRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to update the Terraform state for the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
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
