package smart_mobile_device_group_v1

import (
	"context"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// state updates the Terraform model with the latest Smart Mobile Device Group V1 information from the Jamf Pro API.
func state(ctx context.Context, data *smartMobileDeviceGroupV1ResourceModel, resourceID string, resp *jamfpro.ResourceSmartMobileDeviceGroupV1) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(resourceID)
	data.Name = types.StringValue(resp.GroupName)

	if resp.GroupDescription == "" && (data.Description.IsNull() || data.Description.IsUnknown()) {
		data.Description = types.StringNull()
	} else {
		data.Description = types.StringValue(resp.GroupDescription)
	}

	if resp.SiteId != nil && *resp.SiteId != "" {
		data.SiteID = types.StringValue(*resp.SiteId)
	} else {
		data.SiteID = types.StringNull()
	}

	criteriaModels := make([]smartMobileDeviceGroupV1CriteriaDataModel, 0, len(resp.Criteria))
	for _, criterion := range resp.Criteria {
		criteriaModel := smartMobileDeviceGroupV1CriteriaDataModel{
			Name:       types.StringValue(criterion.Name),
			Priority:   types.Int32Value(int32(criterion.Priority)),
			AndOr:      types.StringValue(strings.ToLower(criterion.AndOr)),
			SearchType: types.StringValue(criterion.SearchType),
			Value:      types.StringValue(criterion.Value),
		}

		if criterion.OpeningParen != nil {
			criteriaModel.OpeningParen = types.BoolValue(*criterion.OpeningParen)
		} else {
			criteriaModel.OpeningParen = types.BoolValue(false)
		}

		if criterion.ClosingParen != nil {
			criteriaModel.ClosingParen = types.BoolValue(*criterion.ClosingParen)
		} else {
			criteriaModel.ClosingParen = types.BoolValue(false)
		}

		criteriaModels = append(criteriaModels, criteriaModel)
	}

	criteriaList, criteriaDiags := schemahelpers.Flatten(ctx, criteriaModels)
	diags.Append(criteriaDiags...)
	if diags.HasError() {
		return diags
	}

	data.Criteria = criteriaList

	return diags
}
