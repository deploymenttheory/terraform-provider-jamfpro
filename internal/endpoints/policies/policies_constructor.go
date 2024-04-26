package policies

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func constructPolicy(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {
	log.Println(("LOGHERE"))

	// Main Object Definition with primitive values assigned.

	log.Println("STRUCT START")

	out := &jamfpro.ResourcePolicy{
		General: jamfpro.PolicySubsetGeneral{
			// ID computed
			Name:                       d.Get("name").(string),
			Enabled:                    d.Get("enabled").(bool),
			TriggerCheckin:             d.Get("trigger_checkin").(bool),
			TriggerEnrollmentComplete:  d.Get("trigger_enrollment_complete").(bool),
			TriggerLogin:               d.Get("trigger_login").(bool),
			TriggerNetworkStateChanged: d.Get("trigger_network_state_changed").(bool),
			TriggerStartup:             d.Get("trigger_startup").(bool),
			TriggerOther:               d.Get("trigger_other").(string),
			Frequency:                  d.Get("frequency").(string),
			RetryEvent:                 d.Get("retry_event").(string),
			RetryAttempts:              d.Get("retry_attempts").(int),
			NotifyOnEachFailedRetry:    d.Get("notify_on_each_failed_retry").(bool),
			TargetDrive:                d.Get("target_drive").(string),
			Offline:                    d.Get("offline").(bool),
			// Category processed
			// site processed
			// Date time limitations
			// network limitations processed
		},
		Scope: &jamfpro.PolicySubsetScope{
			AllComputers: d.Get("scope.0.all_computers").(bool),
			AllJSSUsers:  d.Get("scope.0.all_jss_users").(bool),
			// computer ids
			// computer group ids
			// jss_user_ids
			// jss_user_group_ids
			// building_ids
			// department_ids
			// Limitations :
			/// user_names
			/// network_segment_ids
			/// ibeacon_ids
			/// user_group_ids
			// Exclusions
			/// computer_ids
			/// computer_group_ids
			/// user_ids
			/// user_group_ids
			/// department_ids
			/// network_segment_ids
			/// jss_user_ids
			/// jss_user_group_ids
			/// ibeacon_ids
		},
		SelfService: &jamfpro.PolicySubsetSelfService{
			// UseForSelfService:           d.Get("self_service.0.use_for_self_service").(bool),
			// SelfServiceDisplayName:      d.Get("self_service_display_name").(string),
			// InstallButtonText:           d.Get("install_button_text").(string),
			// ReinstallButtonText:         d.Get("reinstall_button_text").(string),
			// SelfServiceDescription:      d.Get("self_service_description").(string),
			// ForceUsersToViewDescription: d.Get("force_users_to_view_description").(bool),
			// TODO self service icon later
			// FeatureOnMainPage: d.Get("feature_on_main_page").(bool),
			// TODO Self service categories later
		},
		// Package Configuration
		// Scripts
		// Printers
		// DockItems
		// Account Maintenance
		// FilesProcesses
		// UserInteraction
		// DiskEncryption
		// Reboot
	}

	// DEBUG
	log.Println("STRUCT END")
	json, _ := json.MarshalIndent(out, " ", "    ")
	log.Println(string(json))
	log.Println("PROCESS START")

	// Processed Fields

	// General

	// Category
	suppliedCategory := d.Get("category").([]interface{})
	if len(suppliedCategory) > 0 {
		// construct category if provided
		outCat := &jamfpro.SharedResourceCategory{
			ID: suppliedCategory[0].(map[string]interface{})["id"].(int),
		}
		out.General.Category = outCat
	} else {
		// if no category, supply empty cat to remove it.
		out.General.Category = &jamfpro.SharedResourceCategory{
			ID: 0,
		}
	}

	// Site
	suppliedSite := d.Get("site").([]interface{})
	if len(suppliedSite) > 0 {
		// If site provided, construct
		outSite := &jamfpro.SharedResourceSite{
			ID: suppliedSite[0].(map[string]interface{})["id"].(int),
		}
		out.General.Site = outSite
	} else {
		// If no site, construct no site obj. We have to do this for the site to be removed.
		out.General.Site = &jamfpro.SharedResourceSite{
			ID: 0,
		}
	}

	// Date time Limitations
	log.Println("LOG-DATETIME")
	if len(d.Get("date_time_limitations").([]interface{})) > 0 {
		pathRoot := "date_time_limitations.0."
		out.General.DateTimeLimitations = &jamfpro.PolicySubsetGeneralDateTimeLimitations{
			ActivationDate:      d.Get(fmt.Sprintf("%s%s", pathRoot, "activation_date")).(string),
			ActivationDateEpoch: d.Get(fmt.Sprintf("%s%s", pathRoot, "activation_date_epoch")).(int),
			ActivationDateUTC:   d.Get(fmt.Sprintf("%s%s", pathRoot, "activation_date_utc")).(string),
			ExpirationDate:      d.Get(fmt.Sprintf("%s%s", pathRoot, "expiration_date")).(string),
			ExpirationDateEpoch: d.Get(fmt.Sprintf("%s%s", pathRoot, "expiration_date_epoch")).(int),
			ExpirationDateUTC:   d.Get(fmt.Sprintf("%s%s", pathRoot, "expiration_date_utc")).(string),
			// no execute on // TODO
			NoExecuteStart: d.Get(fmt.Sprintf("%s%s", pathRoot, "no_execute_start")).(string),
			NoExecuteEnd:   d.Get(fmt.Sprintf("%s%s", pathRoot, "no_execute_end")).(string),
		}
	}

	// Network Limitations
	log.Println("LOG-NETWORK")
	if len(d.Get("network_limitations").([]interface{})) > 0 {
		log.Println("FLAG 1")
		pathRoot := "network_limitations.0."
		out.General.NetworkLimitations = &jamfpro.PolicySubsetGeneralNetworkLimitations{}
		log.Println("FLAG 2")
		// out.General.NetworkLimitations.MinimumNetworkConnection = d.Get(fmt.Sprintf("%s%s", pathRoot, "minimum_network_connection")).(string)
		out.General.NetworkRequirements = d.Get(fmt.Sprintf("%s%s", pathRoot, "minimum_network_connection")).(string)
		log.Println("FLAG 3")
		out.General.NetworkLimitations.AnyIPAddress = d.Get(fmt.Sprintf("%s%s", pathRoot, "any_ip_address")).(bool)
		log.Println("FLAG 4")
		out.General.NetworkLimitations.NetworkSegments = d.Get(fmt.Sprintf("%s%s", pathRoot, "network_segments")).(string)

		log.Println("FLAG 5")
	}

	// Scope

	// DEBUG
	policyXML, _ := xml.MarshalIndent(out, "", "  ")
	log.Println("LOGEND")
	log.Println(string(policyXML))

	// END
	return out, nil
}
