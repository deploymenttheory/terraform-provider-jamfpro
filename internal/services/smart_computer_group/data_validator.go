package smart_computer_group

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the resource implements the ResourceWithConfigValidators interface
var _ resource.ResourceWithConfigValidators = &smartComputerGroupFrameworkResource{}

// ConfigValidators returns a list of config validators for the resource
func (r *smartComputerGroupFrameworkResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		schema.CriteriaPriorityValidator[smartComputerGroupResourceModel]{},
	}
}

// GetCriteria implements schema.CriteriaModel for validation
func (m smartComputerGroupResourceModel) GetCriteria() []schema.CriterionModel {
	criteria := make([]schema.CriterionModel, len(m.Criteria))
	for i, c := range m.Criteria {
		criteria[i] = c
	}
	return criteria
}

// GetPriority implements schema.CriterionModel for validation
func (c smartComputerGroupCriteriaDataModel) GetPriority() types.Int64 {
	return c.Priority
}
