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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProComputerCheckindefines the schema and RU operations for managing Jamf Pro computer checkin configuration in Terraform.
func ResourceJamfProComputerCheckin() *schema.Resource {
	return &schema.Resource{
		ReadContext:   ResourceJamfProComputerCheckinRead,
		UpdateContext: ResourceJamfProComputerCheckinUpdate,
		DeleteContext: ResourceJamfProComputerCheckinDelete,
		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
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
				Description: "Determines if policies should be checked at startup.",
			},
			"apply_computer_level_managed_preferences": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Applies computer level managed preferences.",
			},
			"ensure_ssh_is_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Ensures that SSH is enabled.",
			},
			"create_login_logout_hooks": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Determines if login/logout hooks should be created.",
			},
			"log_username": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Logs the username.",
			},
			"check_for_policies_at_login_logout": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Checks for policies at login and logout.",
			},
			"apply_user_level_managed_preferences": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Applies user level managed preferences.",
			},
			"hide_restore_partition": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Determines if the restore partition should be hidden.",
			},
			"perform_login_actions_in_background": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Performs login actions in the background.",
			},
			"display_status_to_user": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Displays status to the user.",
			},
		},
	}
}

// constructComputerCheckin constructs a ResponseComputerCheckin object from the provided schema data and returns any errors encountered.
func constructComputerCheckin(d *schema.ResourceData) (*jamfpro.ResponseComputerCheckin, error) {
	checkin := &jamfpro.ResponseComputerCheckin{}

	fields := map[string]interface{}{
		"check_in_frequency":                  &checkin.CheckInFrequency,
		"create_startup_script":               &checkin.CreateStartupScript,
		"log_startup_event":                   &checkin.LogStartupEvent,
		"check_for_policies_at_startup":       &checkin.CheckForPoliciesAtStartup,
		"apply_computer_level_managed_prefs":  &checkin.ApplyComputerLevelManagedPrefs,
		"ensure_ssh_is_enabled":               &checkin.EnsureSSHIsEnabled,
		"create_login_logout_hooks":           &checkin.CreateLoginLogoutHooks,
		"log_username":                        &checkin.LogUsername,
		"check_for_policies_at_login_logout":  &checkin.CheckForPoliciesAtLoginLogout,
		"apply_user_level_managed_prefs":      &checkin.ApplyUserLevelManagedPreferences,
		"hide_restore_partition":              &checkin.HideRestorePartition,
		"perform_login_actions_in_background": &checkin.PerformLoginActionsInBackground,
		"display_status_to_user":              &checkin.DisplayStatusToUser,
	}

	for key, ptr := range fields {
		if v, ok := d.GetOk(key); ok {
			switch ptr := ptr.(type) {
			case *int:
				*ptr = v.(int)
			case *bool:
				*ptr = v.(bool)
			default:
				return nil, fmt.Errorf("unsupported data type for key '%s'", key)
			}
		}
	}

	// Log the successful construction of the group
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
	var checkinConfig *jamfpro.ResponseComputerCheckin
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

	// Map the configuration fields from the API response to the Terraform schema with type assertion
	if err := d.Set("check_in_frequency", checkinConfig.CheckInFrequency); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("create_startup_script", checkinConfig.CreateStartupScript); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("log_startup_event", checkinConfig.LogStartupEvent); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("check_for_policies_at_startup", checkinConfig.CheckForPoliciesAtStartup); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("apply_computer_level_managed_prefs", checkinConfig.ApplyComputerLevelManagedPrefs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ensure_ssh_is_enabled", checkinConfig.EnsureSSHIsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("create_login_logout_hooks", checkinConfig.CreateLoginLogoutHooks); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("log_username", checkinConfig.LogUsername); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("check_for_policies_at_login_logout", checkinConfig.CheckForPoliciesAtLoginLogout); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("apply_user_level_managed_prefs", checkinConfig.ApplyUserLevelManagedPreferences); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hide_restore_partition", checkinConfig.HideRestorePartition); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("perform_login_actions_in_background", checkinConfig.PerformLoginActionsInBackground); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_status_to_user", checkinConfig.DisplayStatusToUser); err != nil {
		return diag.FromErr(err)
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
