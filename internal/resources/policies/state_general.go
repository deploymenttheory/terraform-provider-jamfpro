package policies

// TODO remove log.prints, debug use only
// TODO maybe review error handling here too?

import (
	"reflect"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// stateGeneral Reads response and states general/root level item block
func stateGeneral(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	var err error

	err = d.Set("name", resp.General.Name)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("enabled", resp.General.Enabled)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_checkin", resp.General.TriggerCheckin)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_enrollment_complete", resp.General.TriggerEnrollmentComplete)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_login", resp.General.TriggerLogin)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_network_state_changed", resp.General.TriggerNetworkStateChanged)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_startup", resp.General.TriggerStartup)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_other", resp.General.TriggerOther)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("frequency", resp.General.Frequency)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("retry_event", resp.General.RetryEvent)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("retry_attempts", resp.General.RetryAttempts)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("target_drive", resp.General.OverrideDefaultSettings.TargetDrive)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("notify_on_each_failed_retry", resp.General.NotifyOnEachFailedRetry)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	err = d.Set("offline", resp.General.Offline)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	if resp.General.NetworkRequirements != "" {
		err = d.Set("network_requirements", resp.General.NetworkRequirements)
		if err != nil {
			*diags = append(*diags, diag.FromErr(err)...)
		}
	}

	// Site
	err = d.Set("site_id", resp.General.Site.ID)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	// Category
	err = d.Set("category_id", resp.General.Category.ID)
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}

	// Set DateTime Limitations
	setGeneralDateTimeLimitations(d, resp, diags)

	// Set Network Limitations
	setGeneralNetworkLimitations(d, resp, diags)

}

// setGeneralDateTimeLimitations is a helper function to set the date_time_limitations block under general
func setGeneralDateTimeLimitations(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	if resp.General.DateTimeLimitations == nil {
		return
	}

	// Check if all values are at their default (empty string or zero value)
	v := reflect.ValueOf(*resp.General.DateTimeLimitations)
	allDefault := true

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if (field.Kind() == reflect.String && field.String() != "") ||
			(field.Kind() == reflect.Int && field.Int() != 0) {
			allDefault = false
			break
		}
	}

	if allDefault {
		return
	}

	// Otherwise, proceed to set the date_time_limitations block
	dateTimeLimitations := make(map[string]interface{})
	dateTimeLimitations["activation_date"] = resp.General.DateTimeLimitations.ActivationDate
	dateTimeLimitations["activation_date_epoch"] = resp.General.DateTimeLimitations.ActivationDateEpoch
	dateTimeLimitations["activation_date_utc"] = resp.General.DateTimeLimitations.ActivationDateUTC
	dateTimeLimitations["expiration_date"] = resp.General.DateTimeLimitations.ExpirationDate
	dateTimeLimitations["expiration_date_epoch"] = resp.General.DateTimeLimitations.ExpirationDateEpoch
	dateTimeLimitations["expiration_date_utc"] = resp.General.DateTimeLimitations.ExpirationDateUTC
	dateTimeLimitations["no_execute_start"] = resp.General.DateTimeLimitations.NoExecuteStart
	dateTimeLimitations["no_execute_end"] = resp.General.DateTimeLimitations.NoExecuteEnd

	err := d.Set("date_time_limitations", []interface{}{dateTimeLimitations})
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}

// setGeneralNetworkLimitations is a helper function to set the network_limitations block under general
func setGeneralNetworkLimitations(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	if resp.General.NetworkLimitations == nil {
		return
	}

	// Check if all values are at their default (true, "No Minimum", or empty string)
	v := reflect.ValueOf(*resp.General.NetworkLimitations)
	allDefault := true

	defaults := map[string]interface{}{
		"MinimumNetworkConnection": "No Minimum",
		"AnyIPAddress":             true,
		"NetworkSegments":          "",
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name

		switch field.Kind() {
		case reflect.Bool:
			if field.Bool() != defaults[fieldName] {
				allDefault = false
			}
		case reflect.String:
			if field.String() != defaults[fieldName] {
				allDefault = false
			}
		}
		if !allDefault {
			break
		}
	}

	if allDefault {
		return
	}

	// Otherwise, proceed to set the network_limitations block
	networkLimitations := make(map[string]interface{})
	networkLimitations["minimum_network_connection"] = resp.General.NetworkLimitations.MinimumNetworkConnection
	networkLimitations["any_ip_address"] = resp.General.NetworkLimitations.AnyIPAddress
	//Appears to be removed from gui
	//networkLimitations["network_segments"] = resp.General.NetworkLimitations.NetworkSegments

	err := d.Set("network_limitations", []interface{}{networkLimitations})
	if err != nil {
		*diags = append(*diags, diag.FromErr(err)...)
	}
}
