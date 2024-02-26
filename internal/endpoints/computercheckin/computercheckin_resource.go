// computercheckin_resource.go
package computercheckin

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
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

const (
	JamfProResourceComputerCheckin = "Computer Checkin"
)

// constructComputerCheckin constructs a ResourceComputerCheckin object from the provided schema data.
func constructComputerCheckin(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceComputerCheckin, error) {
	// Initialize the logging subsystem for the construction operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemConstruct, hclog.Debug)

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

	// Serialize and pretty-print the site object as XML
	resourceXML, err := xml.MarshalIndent(checkin, "", "  ")
	if err != nil {
		logging.LogTFConstructResourceXMLMarshalFailure(subCtx, JamfProResourceComputerCheckin, err.Error())
		return nil, err
	}

	// Log the successful construction and serialization to XML
	logging.LogTFConstructedXMLResource(subCtx, JamfProResourceComputerCheckin, string(resourceXML))

	return checkin, nil
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

	// Initialize the logging subsystem for the update operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	var apiErrorCode int
	resourceID := d.Id()
	resourceName := "jamfpro_computer_checkin_singleton"

	// Construct the resource object
	checkinConfig, err := constructComputerCheckin(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceComputerCheckin, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceComputerCheckin)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		apiErr := conn.UpdateComputerCheckinInformation(checkinConfig)
		if apiErr != nil {
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceComputerCheckin, resourceID, resourceName, apiErr.Error(), apiErrorCode)
			return retry.RetryableError(apiErr)
		}
		logging.LogAPIUpdateSuccess(subCtx, JamfProResourceComputerCheckin, resourceID, resourceName)
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceComputerCheckin, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the resource to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProComputerCheckinRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceComputerCheckin, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceComputerCheckin, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceComputerCheckin, resourceID)
	}

	return diags
}

// ResourceJamfProComputerCheckinRead is responsible for reading the current state of the Jamf Pro computer check-in configuration.
func ResourceJamfProComputerCheckinRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	var apiErrorCode int
	var checkinConfig *jamfpro.ResourceComputerCheckin

	var apiErr error
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		checkinConfig, apiErr = conn.GetComputerCheckinInformation()
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceComputerCheckin, resourceID, apiErr.Error(), apiErrorCode)
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the account group, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		d.SetId("") // Remove from Terraform state if unable to read after retries
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceComputerCheckin, resourceID)
		return diag.FromErr(err)
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

	// Initialize the logging subsystem for the update operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	var apiErrorCode int
	resourceID := d.Id()
	resourceName := "jamfpro_computer_checkin_singleton"

	// Construct the resource object
	checkinConfig, err := constructComputerCheckin(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceComputerCheckin, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceComputerCheckin)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		apiErr := conn.UpdateComputerCheckinInformation(checkinConfig)
		if apiErr != nil {
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceComputerCheckin, resourceID, resourceName, apiErr.Error(), apiErrorCode)
			return retry.RetryableError(apiErr)
		}
		logging.LogAPIUpdateSuccess(subCtx, JamfProResourceComputerCheckin, resourceID, resourceName)
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceComputerCheckin, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the resource to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProComputerCheckinRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceComputerCheckin, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceComputerCheckin, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceComputerCheckin, resourceID)
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
