package policy

// TODO remove log.prints, debug use only
// TODO maybe review error handling here too?

import (
	"fmt"
	"log"
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

// setGeneralDateTimeLimitations updates the Terraform state for date_time_limitations during Read.
// it supports two scenarios with or without the hcl block defineds. if the block is not in hcl it will
// ignore entirely. in scenario 2 it will state the block but, as usual, since there's an issue with the GET
// on the api, for the fields "no_execute_start" , and "no_execute_end" , we have to extract these values from
// the HCL directly and state those.
func setGeneralDateTimeLimitations(d *schema.ResourceData, resp *jamfpro.ResourcePolicy, diags *diag.Diagnostics) {
	// Check if the block is defined in the HCL configuration (current state)
	hclBlockRaw, hclBlockExists := d.GetOk("date_time_limitations")

	// --- Scenario 1: Block NOT defined in HCL ---
	if !hclBlockExists {
		log.Printf("[DEBUG] setGeneralDateTimeLimitations: Block 'date_time_limitations' not configured in HCL. Ensuring state is nil.")
		if err := d.Set("date_time_limitations", nil); err != nil {
			*diags = append(*diags, diag.FromErr(fmt.Errorf("failed to unset date_time_limitations in state: %w", err))...)
		}
		return
	}

	// --- Scenario 2: Block IS defined in HCL ---
	log.Printf("[DEBUG] setGeneralDateTimeLimitations: Block 'date_time_limitations' is configured in HCL. Populating state using API and HCL overrides.")

	newStateMap := make(map[string]interface{})
	apiBlock := resp.General.DateTimeLimitations

	if apiBlock != nil {
		newStateMap["activation_date"] = apiBlock.ActivationDate
		newStateMap["activation_date_epoch"] = int(apiBlock.ActivationDateEpoch)
		newStateMap["activation_date_utc"] = apiBlock.ActivationDateUTC
		newStateMap["expiration_date"] = apiBlock.ExpirationDate
		newStateMap["expiration_date_epoch"] = int(apiBlock.ExpirationDateEpoch)
		newStateMap["expiration_date_utc"] = apiBlock.ExpirationDateUTC

		var noExecuteOnItems []interface{}
		if apiBlock.NoExecuteOn != nil {
			noExecuteOnItems = make([]interface{}, len(apiBlock.NoExecuteOn))
			for i, day := range apiBlock.NoExecuteOn {
				noExecuteOnItems[i] = day
			}
		} else {
			noExecuteOnItems = []interface{}{}
		}

		newStateMap["no_execute_on"] = schema.NewSet(schema.HashString, noExecuteOnItems)

		// Set start/end from API initially (will be overwritten by HCL values next)
		newStateMap["no_execute_start"] = apiBlock.NoExecuteStart
		newStateMap["no_execute_end"] = apiBlock.NoExecuteEnd
		log.Printf("[DEBUG] setGeneralDateTimeLimitations: Populated map with API data: %+v", newStateMap)

	} else {
		newStateMap["activation_date"] = ""
		newStateMap["activation_date_epoch"] = 0
		newStateMap["activation_date_utc"] = ""
		newStateMap["expiration_date"] = ""
		newStateMap["expiration_date_epoch"] = 0
		newStateMap["expiration_date_utc"] = ""
		newStateMap["no_execute_start"] = ""
		newStateMap["no_execute_end"] = ""
		newStateMap["no_execute_on"] = schema.NewSet(schema.HashString, []interface{}{})
		log.Printf("[DEBUG] setGeneralDateTimeLimitations: API did not return date_time_limitations block. Initialized state map with defaults.")
	}

	var hclStartValue string = ""
	var hclEndValue string = ""

	hclList, listOk := hclBlockRaw.([]interface{})
	if listOk && len(hclList) > 0 && hclList[0] != nil {
		hclMap, mapOk := hclList[0].(map[string]interface{})
		if mapOk {
			if val, ok := hclMap["no_execute_start"].(string); ok {
				hclStartValue = val
			}
			if val, ok := hclMap["no_execute_end"].(string); ok {
				hclEndValue = val
			}
			log.Printf("[DEBUG] Extracted from HCL state: no_execute_start='%s', no_execute_end='%s'", hclStartValue, hclEndValue)
		} else {
			log.Printf("[WARN] Could not read HCL date_time_limitations block as map during state setting.")
		}
	} else {
		log.Printf("[WARN] HCL date_time_limitations block exists but is not a valid list or is empty during state setting.")
	}

	newStateMap["no_execute_start"] = hclStartValue
	newStateMap["no_execute_end"] = hclEndValue

	log.Printf("[DEBUG] Setting final date_time_limitations state: %+v", newStateMap)
	err := d.Set("date_time_limitations", []interface{}{newStateMap})
	if err != nil {
		*diags = append(*diags, diag.Errorf("Failed to set date_time_limitations in state: %s", err)...)
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
