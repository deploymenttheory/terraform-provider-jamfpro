// policies_data_validation.go
package policies

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateJamfProResourcePolicyDataFields ensures that a Jamf Pro policy meets specific criteria:
// 1. It cannot be set to have an ongoing frequency when the trigger_checkin is true and it applies to all computers.
// 2. The 'offline' field can only be set to true if the 'frequency' is "Ongoing".
// This function is used as a CustomizeDiff function in a Terraform resource schema to validate the policy configuration.
func validateJamfProResourcePolicyDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Validate the first scenario
	triggerCheckin, triggerCheckinOk := diff.GetOk("general.0.trigger_checkin")
	frequency, frequencyOk := diff.GetOk("general.0.frequency")
	allComputers, allComputersOk := diff.GetOk("scope.0.all_computers")

	if triggerCheckinOk && frequencyOk && allComputersOk {
		if triggerCheckin.(bool) && frequency.(string) == "Ongoing" && allComputers.(bool) {
			return fmt.Errorf("jamf pro policies that update inventory on all computers cannot be set to ongoing frequency at recurring check-in. Please update your terraform configuration")
		}
	}

	// Validate the second scenario
	offline, offlineOk := diff.GetOk("general.0.offline")
	if offlineOk && offline.(bool) && frequencyOk && frequency.(string) != "Ongoing" {
		return fmt.Errorf("jamf pro policy triggers can only be set to 'offline' if the policy frequency is set to 'Ongoing'. Please update your terraform configuration")
	}

	return nil
}
