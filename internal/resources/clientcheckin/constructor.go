package computercheckin

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructClientCheckInSettings constructs a ResourceComputerCheckin object from the provided schema data and logs its XML representation.
func constructClientCheckInSettings(d *schema.ResourceData) (*jamfpro.ResourceClientCheckinSettings, error) {
	resource := &jamfpro.ResourceClientCheckinSettings{
		CheckInFrequency:                 d.Get("check_in_frequency").(int),
		CreateHooks:                      d.Get("create_hooks").(bool),
		HookLog:                          d.Get("hook_log").(bool),
		HookPolicies:                     d.Get("hook_policies").(bool),
		CreateStartupScript:              d.Get("create_startup_script").(bool),
		StartupLog:                       d.Get("startup_log").(bool),
		StartupPolicies:                  d.Get("startup_policies").(bool),
		StartupSsh:                       d.Get("startup_ssh").(bool),
		EnableLocalConfigurationProfiles: d.Get("enable_local_configuration_profiles").(bool),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {

		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Checkin to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Checkin JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// constructPolicyProperties constructs a ResourcePolicyProperties object from the provided schema data.
func constructPolicyProperties(d *schema.ResourceData) (*jamfpro.ResourcePolicyProperties, error) {
	allowNetworkStateChangeTriggers := d.Get("allow_network_state_change_triggers").(bool)

	resource := &jamfpro.ResourcePolicyProperties{
		AllowNetworkStateChangeTriggers: allowNetworkStateChangeTriggers,
	}

	// Log the constructed resource for debugging
	log.Printf("[DEBUG] Constructed Jamf Pro Policy Properties: %+v\n", resource)

	return resource, nil
}
