package smart_mobile_device_group

import (
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// state updates the Terraform model with the latest Smart Mobile Device Group V1 information from the Jamf Pro API.
func state(data *smartMobileDeviceGroupResourceModel, resourceID string, resp *jamfpro.ResourceSmartMobileDeviceGroupV1) diag.Diagnostics {
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

	data.Criteria = make([]smartMobileDeviceGroupCriteriaDataModel, 0, len(resp.Criteria))
	for _, criterion := range resp.Criteria {
		criteriaModel := smartMobileDeviceGroupCriteriaDataModel{
			Name:       types.StringValue(criterion.Name),
			Priority:   types.Int64Value(int64(criterion.Priority)),
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

		data.Criteria = append(data.Criteria, criteriaModel)
	}

	return diags
}
