package policies

import (
	"encoding/xml"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func constructPolicy(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {

	// Non computed values first
	policy := &jamfpro.ResourcePolicy{
		General: jamfpro.PolicySubsetGeneral{
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
		},
		SelfService: &jamfpro.PolicySubsetSelfService{
			UseForSelfService: d.Get("self_service.0.use_for_self_service").(bool),
		},
	}

	policyXML, _ := xml.MarshalIndent(policy, "", "  ")
	log.Println("LOGHERE")
	log.Println(string(policyXML))
	return policy, nil
}
