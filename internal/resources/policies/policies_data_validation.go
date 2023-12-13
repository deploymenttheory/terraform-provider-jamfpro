// policies_data_validation.go
package policies

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateJamfProResourcePolicyDataFields ensures that a Jamf Pro policy is not set
// to have an ongoing frequency when the trigger_checkin is true and it applies to all computers.
// This function is intended to be used as a CustomizeDiff function in a Terraform resource schema
// to validate the policy configuration during the plan phase.
func validateJamfProResourcePolicyDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	triggerCheckin, triggerCheckinOk := diff.GetOk("general.0.trigger_checkin")
	frequency, frequencyOk := diff.GetOk("general.0.frequency")
	allComputers, allComputersOk := diff.GetOk("scope.0.all_computers")

	if triggerCheckinOk && frequencyOk && allComputersOk {
		if triggerCheckin.(bool) && frequency.(string) == "Ongoing" && allComputers.(bool) {
			return fmt.Errorf("policies that update inventory on all computers cannot be set to ongoing frequency at recurring check-in")
		}
	}

	return nil
}
