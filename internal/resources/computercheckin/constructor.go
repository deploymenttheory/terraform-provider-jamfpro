package computercheckin

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProComputerCheckin constructs a ResourceComputerCheckin object from the provided schema data and logs its XML representation.
func construct(d *schema.ResourceData) (*jamfpro.ResourceComputerCheckin, error) {
	resource := &jamfpro.ResourceComputerCheckin{
		CheckInFrequency:          d.Get("check_in_frequency").(int),
		CreateStartupScript:       d.Get("create_startup_script").(bool),
		LogStartupEvent:           d.Get("log_startup_event").(bool),
		CheckForPoliciesAtStartup: d.Get("check_for_policies_at_startup").(bool),
		// Note: "apply_computer_level_managed_preferences" is computed, not set directly
		EnsureSSHIsEnabled:            d.Get("ensure_ssh_is_enabled").(bool),
		CreateLoginLogoutHooks:        d.Get("create_login_logout_hooks").(bool),
		LogUsername:                   d.Get("log_username").(bool),
		CheckForPoliciesAtLoginLogout: d.Get("check_for_policies_at_login_logout").(bool),
		// Note: "apply_user_level_managed_preferences", "hide_restore_partition", and "perform_login_actions_in_background" are computed, not set directly
	}

	// Serialize and pretty-print the Category object as XML for logging
	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Checkin to XML: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Computer Checkin XML:\n%s\n", string(resourceXML))

	return resource, nil
}
