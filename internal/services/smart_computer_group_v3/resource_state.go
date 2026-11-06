package smart_computer_group_v3

import (
	"context"
	"slices"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func normalizedSiteID(id *string) string {
	if id == nil || *id == "" {
		return "-1"
	}
	return *id
}

func orderedCriteria(criteria []jamfpro.SubsetComputerSmartGroupCriteriaV3) []jamfpro.SubsetComputerSmartGroupCriteriaV3 {
	result := slices.Clone(criteria)
	slices.SortFunc(result, func(a, b jamfpro.SubsetComputerSmartGroupCriteriaV3) int {
		if a.Priority < b.Priority {
			return -1
		}
		if a.Priority > b.Priority {
			return 1
		}
		return 0
	})
	for i := range result {
		result[i].AndOr = strings.ToLower(result[i].AndOr)
	}
	return result
}

func state(
	ctx context.Context, data *smartComputerGroupV3ResourceModel, group *jamfpro.ResourceSmartComputerGroupV3,
) diag.Diagnostics {
	data.Name = types.StringValue(group.Name)
	data.Description = types.StringValue(group.Description)
	data.SiteID = types.StringValue(normalizedSiteID(group.SiteID))
	criteria := orderedCriteria(group.Criteria)
	// An omitted criteria block must keep its null representation on read-back.
	if len(criteria) == 0 && data.Criteria.IsNull() {
		data.Criteria = types.ListNull(schemahelpers.ObjectType)
		return nil
	}
	models := make([]criterionModel, 0, len(criteria))
	for _, criterion := range criteria {
		models = append(models, criterionModel{
			Name: types.StringValue(criterion.Name), Priority: types.Int64Value(int64(criterion.Priority)),
			AndOr: types.StringValue(criterion.AndOr), SearchType: types.StringValue(criterion.SearchType),
			Value: types.StringValue(criterion.Value), OpeningParen: types.BoolValue(criterion.OpeningParen),
			ClosingParen: types.BoolValue(criterion.ClosingParen),
		})
	}
	list, diags := schemahelpers.Flatten(ctx, models)
	data.Criteria = list
	return diags
}

func groupReflectsWrite(got, desired *jamfpro.ResourceSmartComputerGroupV3, description bool) bool {
	if got == nil || got.Name != desired.Name || normalizedSiteID(got.SiteID) != normalizedSiteID(desired.SiteID) {
		return false
	}
	if description && got.Description != desired.Description {
		return false
	}
	return slices.Equal(orderedCriteria(got.Criteria), orderedCriteria(desired.Criteria))
}
