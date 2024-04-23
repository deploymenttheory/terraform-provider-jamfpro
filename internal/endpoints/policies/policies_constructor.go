package policies

import (
	"encoding/xml"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func constructJamfProPolicy(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {

	// Non computed values first
	policy := &jamfpro.ResourcePolicy{
		General: jamfpro.PolicySubsetGeneral{
			Name:           d.Get("name").(string),
			Enabled:        d.Get("enabled").(bool),
			TriggerCheckin: d.Get("trigger_checkin").(bool),
		},
	}

	policyXML, _ := xml.MarshalIndent(policy, "", "  ")
	log.Println("LOGHERE")
	log.Println(string(policyXML))
	return policy, nil
}
