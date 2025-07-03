package reenrollment

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructReenrollmentSettings constructs a ResourceReenrollment object from the provided schema data and logs its XML representation.
func construct(d *schema.ResourceData) (*jamfpro.ResourceReenrollmentSettings, error) {
	resource := &jamfpro.ResourceReenrollmentSettings{
		FlushLocationInformation:        d.Get("flush_location_information").(bool), 
		FlushLocationInformationHistory: d.Get("flush_location_information_history").(bool),
		FlushPolicyHistory:              d.Get("flush_policy_history").(bool),
		FlushExtensionAttributes:        d.Get("flush_extension_attributes").(bool),
		FlushSoftwareUpdatePlans:        d.Get("flush_software_update_plans").(bool),
		FlushMdmQueue:                   d.Get("flush_mdm_queue").(string),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Reenrollment to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Reenrollment JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// // constructPolicyProperties constructs a ResourcePolicyProperties object from the provided schema data.
// func constructPolicyProperties(d *schema.ResourceData) (*jamfpro.ResourcePolicyProperties, error) {
// 	allowNetworkStateChangeTriggers := d.Get("allow_network_state_change_triggers").(bool)

// 	resource := &jamfpro.ResourcePolicyProperties{
// 		AllowNetworkStateChangeTriggers: allowNetworkStateChangeTriggers,
// 	}

// 	// Log the constructed resource for debugging
// 	log.Printf("[DEBUG] Constructed Jamf Pro Policy Properties: %+v\n", resource)

// 	return resource, nil
// }
