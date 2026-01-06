package smart_computer_group

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/validation"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the resource implements the ResourceWithConfigValidators interface
var _ resource.ResourceWithConfigValidators = &smartComputerGroupFrameworkResource{}

// ConfigValidators returns a list of config validators for the resource
func (r *smartComputerGroupFrameworkResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		validation.IncrementingInt32SequenceValidator[smartComputerGroupResourceModel]{},
	}
}

// GetInt32Sequence exposes the priority sequence for validation.
func (m smartComputerGroupResourceModel) GetInt32Sequence() []int32 {
	priorities := make([]int32, len(m.Criteria))
	for i, c := range m.Criteria {
		priorities[i] = c.Priority.ValueInt32()
	}
	return priorities
}
