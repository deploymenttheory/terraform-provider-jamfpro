package smart_computer_group_v2

import (
	"context"

	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/validation"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the resource implements the ResourceWithConfigValidators interface
var _ resource.ResourceWithConfigValidators = &smartComputerGroupV2FrameworkResource{}

// ConfigValidators returns a list of config validators for the resource
func (r *smartComputerGroupV2FrameworkResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		validation.IncrementingInt32SequenceValidator[smartComputerGroupV2ResourceModel]{},
	}
}

// GetInt32Sequence exposes the priority sequence for validation.
func (m smartComputerGroupV2ResourceModel) GetInt32Sequence() []int32 {
	if m.Criteria.IsNull() || m.Criteria.IsUnknown() {
		return nil
	}

	criteria, diags := schemahelpers.Expand[smartComputerGroupV2CriteriaDataModel](context.Background(), m.Criteria)
	if diags.HasError() {
		return nil
	}

	priorities := make([]int32, 0, len(criteria))
	for _, c := range criteria {
		if c.Priority.IsUnknown() {
			return nil
		}
		priorities = append(priorities, c.Priority.ValueInt32())
	}

	return priorities
}
