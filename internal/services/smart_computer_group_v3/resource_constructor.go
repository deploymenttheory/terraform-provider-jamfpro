package smart_computer_group_v3

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func constructResource(
	ctx context.Context, data *smartComputerGroupV3ResourceModel,
) (*jamfpro.ResourceSmartComputerGroupV3, diag.Diagnostics) {
	criteria, diags := schemahelpers.Expand[criterionModel](ctx, data.Criteria)
	if diags.HasError() {
		return nil, diags
	}
	siteID := data.SiteID.ValueString()
	if siteID == "" {
		siteID = "-1"
	}
	result := &jamfpro.ResourceSmartComputerGroupV3{
		Name: data.Name.ValueString(), Description: data.Description.ValueString(), SiteID: &siteID,
		Criteria: make([]jamfpro.SubsetComputerSmartGroupCriteriaV3, 0, len(criteria)),
	}
	for _, criterion := range criteria {
		result.Criteria = append(result.Criteria, jamfpro.SubsetComputerSmartGroupCriteriaV3{
			Name: criterion.Name.ValueString(), Priority: int(criterion.Priority.ValueInt64()),
			AndOr: criterion.AndOr.ValueString(), SearchType: criterion.SearchType.ValueString(),
			Value: criterion.Value.ValueString(), OpeningParen: criterion.OpeningParen.ValueBool(),
			ClosingParen: criterion.ClosingParen.ValueBool(),
		})
	}
	return result, diags
}

func constructClassic(input *jamfpro.ResourceSmartComputerGroupV3) (*jamfpro.ResourceComputerGroup, error) {
	siteID, err := strconv.Atoi(normalizedSiteID(input.SiteID))
	if err != nil {
		return nil, fmt.Errorf("parse site ID: %w", err)
	}
	criteria := make([]jamfpro.SharedSubsetCriteria, 0, len(input.Criteria))
	for _, criterion := range input.Criteria {
		criteria = append(criteria, jamfpro.SharedSubsetCriteria(criterion))
	}
	// A single Classic request must carry all criteria, name, and site: splitting
	// this write would expose an intermediate membership to scoped policies.
	return &jamfpro.ResourceComputerGroup{
		Name: input.Name, IsSmart: true, Site: &jamfpro.SharedResourceSite{ID: siteID},
		Criteria: &jamfpro.ComputerGroupSubsetContainerCriteria{Size: len(criteria), Criterion: &criteria},
	}, nil
}
