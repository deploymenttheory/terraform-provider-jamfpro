// smartcomputergroup_data_validator.go
package smart_computer_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	// Validate criteria priorities
	if err := validateCriteriaPriority(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateCriteriaPriority ensures the first criterion has a priority of 0 and each subsequent criterion has a priority incremented by 1.
func validateCriteriaPriority(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	criteria, ok := diff.Get("criteria").([]interface{})
	if !ok {
		return nil
	}

	resourceName, ok := diff.Get("name").(string)
	if !ok {
		return fmt.Errorf("unable to retrieve resource name for validation")
	}

	expectedPriority := 0
	for index, criterion := range criteria {
		priority := criterion.(map[string]interface{})["priority"].(int)
		if index == 0 && priority != 0 {
			return fmt.Errorf("in 'jamfpro_smart_computer_group.%s': the first criterion must have a priority of 0, got %d", resourceName, priority)
		} else if index > 0 && priority != expectedPriority {
			return fmt.Errorf("in 'jamfpro_smart_computer_group.%s': criterion %d has an invalid priority %d, expected %d", resourceName, index, priority, expectedPriority)
		}
		expectedPriority++
	}

	return nil
}

// getCriteriaOperators returns a list of criteria operators for Smart Computer Groups.
func getCriteriaOperators() []string {
	out := []string{
		And,
		Or,
		SearchTypeIs,
		SearchTypeIsNot,
		SearchTypeHas,
		SearchTypeDoesNotHave,
		SearchTypeMemberOf,
		SearchTypeNotMemberOf,
		SearchTypeBeforeYYYYMMDD,
		SearchTypeAfterYYYYMMDD,
		SearchTypeMoreThanXDaysAgo,
		SearchTypeLessThanXDaysAgo,
		SearchTypeLike,
		SearchTypeNotLike,
		SearchTypeGreaterThan,
		SearchTypeMoreThan,
		SearchTypeLessThan,
		SearchTypeGreaterThanOrEqual,
		SearchTypeLessThanOrEqual,
		SearchTypeMatchesRegex,
		SearchTypeDoesNotMatch,
	}
	return out
}
