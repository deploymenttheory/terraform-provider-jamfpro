package smart_mobile_device_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the resource implements the ResourceWithConfigValidators interface
var _ resource.ResourceWithConfigValidators = &smartMobileDeviceGroupFrameworkResource{}

// ConfigValidators returns a list of config validators for the resource
func (r *smartMobileDeviceGroupFrameworkResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		&criteriaPriorityValidator{},
	}
}

// criteriaPriorityValidator validates that criteria priorities start at 0 and increment by 1
type criteriaPriorityValidator struct{}

// Description returns a plain text description of the validator's behavior
func (v criteriaPriorityValidator) Description(ctx context.Context) string {
	return "Ensures criteria priorities start at 0 and increment by 1 for each subsequent criterion"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior
func (v criteriaPriorityValidator) MarkdownDescription(ctx context.Context) string {
	return "Ensures criteria priorities start at 0 and increment by 1 for each subsequent criterion"
}

// ValidateResource performs the validation
func (v criteriaPriorityValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data smartMobileDeviceGroupResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(data.Criteria) <= 1 {
		return
	}

	expectedPriority := int64(0)
	for index, criterion := range data.Criteria {
		priority := criterion.Priority.ValueInt64()

		if index == 0 && priority != 0 {
			resp.Diagnostics.AddError(
				"Invalid Criteria Priority",
				fmt.Sprintf("The first criterion must have a priority of 0, got %d", priority),
			)
			return
		}

		if index > 0 && priority != expectedPriority {
			resp.Diagnostics.AddError(
				"Invalid Criteria Priority",
				fmt.Sprintf("Criterion %d has an invalid priority %d, expected %d.", index, priority, expectedPriority),
			)
			return
		}

		expectedPriority++
	}
}
