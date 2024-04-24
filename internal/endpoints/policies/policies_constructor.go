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
	// Non computed values first

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
			// date time limitations processed
			// network limitations processed
		},
		Scope: &jamfpro.PolicySubsetScope{
			// Processed
		},
		SelfService: &jamfpro.PolicySubsetSelfService{
			UseForSelfService: d.Get("self_service.0.use_for_self_service").(bool),
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
	log.Println("STRUCT END")
	json, _ := json.MarshalIndent(out, " ", "    ")
	log.Println(string(json))

	log.Println("PROCESS START")

	// Processed Fields

	// General

	log.Println("LOG-CATEGORY")
	// Category
	if len(d.Get("category").([]interface{})) > 0 {
		out.General.Category = &jamfpro.SharedResourceCategory{
			ID:   d.Get("category.0.id").(int),
			Name: d.Get("category.0.name").(string),
		}
	}

	log.Println("LOG-SITE")
	// Site
	if len(d.Get("site").([]interface{})) > 0 {
		out.General.Site = &jamfpro.SharedResourceSite{
			ID:   d.Get("site.0.id").(int),
			Name: d.Get("site.0.name").(string),
		}
	}

	log.Println("LOG-DATETIME")
	// Date time Limitations
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

	log.Println("LOG-NETWORK")
	// Network Limitations
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

	policyXML, _ := xml.MarshalIndent(out, "", "  ")
	log.Println("LOGEND")
	log.Println(string(policyXML))
	return out, nil
}
