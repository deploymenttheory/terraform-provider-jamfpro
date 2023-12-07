// computergroup_data_validation.go
package computergroups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customDiffComputeGroups is a CustomDiff function that enforces conditional logic on the 'computers' and 'criteria' fields of the JamfProComputerGroups resource based on the value of 'is_smart'.
// When is_smart is true, the criteria block is valid, and the computers block should not be set.
// When is_smart is false, the computers block is valid, and the criteria block should not be set.
func customDiffComputeGroups(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	isSmart := diff.Get("is_smart").(bool)

	// When 'is_smart' is true, 'computers' should not be set.
	if isSmart {
		if computers, exists := diff.GetOk("computers"); exists && len(computers.([]interface{})) > 0 {
			return fmt.Errorf("'computers' field is not allowed when 'is_smart' is true")
		}
	} else {
		// If 'is_smart' is false, 'criteria' should not be set.
		if criteria, exists := diff.GetOk("criteria"); exists && len(criteria.([]interface{})) > 0 {
			return fmt.Errorf("'criteria' field is not allowed when 'is_smart' is false")
		}
	}

	// Additional validations for 'criteria' when 'is_smart' is true.
	if isSmart {
		criteria, ok := diff.GetOk("criteria")
		if !ok || len(criteria.([]interface{})) == 0 {
			return fmt.Errorf("'criteria' field must be set when 'is_smart' is true")
		}

		for i, c := range criteria.([]interface{}) {
			criterion, ok := c.(map[string]interface{})
			if !ok {
				continue // Skip invalid structure.
			}

			// Validate 'name', 'and_or', and 'search_type' in each criterion.
			if criterion["name"] == nil || criterion["name"].(string) == "" {
				return fmt.Errorf("'name' field is required for 'criteria' at index %d when 'is_smart' is true", i)
			}
			if criterion["and_or"] == nil || criterion["and_or"].(string) == "" {
				return fmt.Errorf("'and_or' field is required for 'criteria' at index %d when 'is_smart' is true", i)
			}
			if criterion["search_type"] == nil || criterion["search_type"].(string) == "" {
				return fmt.Errorf("'search_type' field is required for 'criteria' at index %d when 'is_smart' is true", i)
			}
		}
	}

	return nil
}
