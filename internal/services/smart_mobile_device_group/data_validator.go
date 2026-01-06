package smart_mobile_device_group

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the resource implements the ResourceWithConfigValidators interface
var _ resource.ResourceWithConfigValidators = &smartMobileDeviceGroupFrameworkResource{}

// ConfigValidators returns a list of config validators for the resource
func (r *smartMobileDeviceGroupFrameworkResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		schema.CriteriaPriorityValidator[smartMobileDeviceGroupResourceModel]{},
	}
}

// GetCriteria implements schema.CriteriaModel for validation
func (m smartMobileDeviceGroupResourceModel) GetCriteria() []schema.CriterionModel {
	criteria := make([]schema.CriterionModel, len(m.Criteria))
	for i, c := range m.Criteria {
		criteria[i] = c
	}
	return criteria
}

// GetPriority implements schema.CriterionModel for validation
func (c smartMobileDeviceGroupCriteriaDataModel) GetPriority() types.Int32 {
	return c.Priority
}
