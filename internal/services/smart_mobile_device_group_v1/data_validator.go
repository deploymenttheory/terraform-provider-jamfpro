package smart_mobile_device_group_v1

import (
	"context"

	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/validation"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the resource implements the ResourceWithConfigValidators interface
var _ resource.ResourceWithConfigValidators = &smartMobileDeviceGroupV1FrameworkResource{}

// ConfigValidators returns a list of config validators for the resource
func (r *smartMobileDeviceGroupV1FrameworkResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		validation.IncrementingInt32SequenceValidator[smartMobileDeviceGroupV1ResourceModel]{},
	}
}

// GetInt32Sequence exposes the priority sequence for validation.
func (m smartMobileDeviceGroupV1ResourceModel) GetInt32Sequence() []int32 {
	if m.Criteria.IsNull() || m.Criteria.IsUnknown() {
		return nil
	}

	criteria, diags := schemahelpers.Expand[smartMobileDeviceGroupV1CriteriaDataModel](context.Background(), m.Criteria)
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
