package user_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i any) error {
	if err := validateIsSmartAttribute(ctx, diff, i); err != nil {
		return err
	}

	if err := validateCriteriaPrioritySequence(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateIsSmartAttribute checks the conditions related to the 'is_smart' attribute.
func validateIsSmartAttribute(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)
	isSmart, ok := diff.GetOkExists("is_smart")

	if !ok {
		return nil
	}

	usersBlockExists := len(diff.Get("assigned_user_ids").([]any)) > 0
	criteriaBlockExists := len(diff.Get("criteria").([]any)) > 0

	if isSmart.(bool) && usersBlockExists {
		return fmt.Errorf("in 'jamfpro_user_group.%s': 'users' block is not allowed when 'is_smart' is set to true", resourceName)
	}

	if !isSmart.(bool) && criteriaBlockExists {
		return fmt.Errorf("in 'jamfpro_user_group.%s': 'criteria' block is not allowed when 'is_smart' is set to false", resourceName)
	}

	if isSmart.(bool) && !criteriaBlockExists {
		return fmt.Errorf("in 'jamfpro_user_group.%s': 'criteria' block is required when 'is_smart' is set to true", resourceName)
	}

	if !isSmart.(bool) && !usersBlockExists {
		return fmt.Errorf("in 'jamfpro_user_group.%s': 'users' block is required when 'is_smart' is set to false", resourceName)
	}

	return nil
}

// validateCriteriaPrioritySequence checks that the 'priority' fields in the 'criteria' blocks are sequential starting from 0.
func validateCriteriaPrioritySequence(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	resourceName := diff.Get("name").(string)

	if criteriaBlocks, ok := diff.GetOk("criteria"); ok {
		criteriaList := criteriaBlocks.([]any)

		for expectedPriority, criteria := range criteriaList {
			criteriaMap := criteria.(map[string]any)

			if actualPriority, ok := criteriaMap["priority"].(int); !ok || actualPriority != expectedPriority {
				return fmt.Errorf("in 'jamfpro_user_group.%s': 'priority' value in 'criteria' block must be sequential starting from 0, found priority '%d' at position '%d'", resourceName, actualPriority, expectedPriority)
			}
		}
	}

	return nil
}
