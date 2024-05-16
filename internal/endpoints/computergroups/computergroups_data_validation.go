// computergroup_data_validation.go
package computergroups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customDiffComputeGroups is the top-level custom diff function.
func customDiffComputeGroups(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := validateComputersNotAllowedWithSmart(ctx, d, meta); err != nil {
		return err
	}

	if err := validateCriteriaFields(ctx, d, meta); err != nil {
		return err
	}

	return nil
}

// validateComputersNotAllowedWithSmart checks that 'computers' is not set when 'is_smart' is true.
func validateComputersNotAllowedWithSmart(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	isSmart := d.Get("is_smart").(bool)
	site := d.Get("site").([]interface{})[0].(map[string]interface{})
	if isSmart && site["id"] != -1 {
		if computers, exists := d.GetOk("computers"); exists && len(computers.([]interface{})) > 0 {
			return fmt.Errorf("'computers' field is not allowed when 'is_smart' is true %v", site)
		}
	}
	return nil
}

// validateCriteriaFields validates that 'name', 'and_or', and 'search_type' are set in each 'criteria' if 'criteria' is populated.
func validateCriteriaFields(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if criteria, ok := d.GetOk("criteria"); ok && len(criteria.([]interface{})) > 0 {
		for i, c := range criteria.([]interface{}) {
			criterion, ok := c.(map[string]interface{})
			if !ok {
				continue // Skip invalid structure.
			}

			// Validate 'name', 'and_or', and 'search_type' in each criterion.
			if criterion["name"] == nil || criterion["name"].(string) == "" {
				return fmt.Errorf("'name' field is required for 'criteria' at index %d", i)
			}
			if criterion["and_or"] == nil || criterion["and_or"].(string) == "" {
				return fmt.Errorf("'and_or' field is required for 'criteria' at index %d", i)
			}
			if criterion["search_type"] == nil || criterion["search_type"].(string) == "" {
				return fmt.Errorf("'search_type' field is required for 'criteria' at index %d", i)
			}
		}
	}
	return nil
}
