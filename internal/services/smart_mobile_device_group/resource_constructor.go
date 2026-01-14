package smart_mobile_device_group

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// constructResource constructs a ResourceSmartMobileDeviceGroupV2 object from the provided framework resource model.
func constructResource(ctx context.Context, data *smartMobileDeviceGroupResourceModel) (*jamfpro.ResourceSmartMobileDeviceGroupV1, diag.Diagnostics) {
	var diags diag.Diagnostics

	resource := &jamfpro.ResourceSmartMobileDeviceGroupV1{
		GroupName:        data.Name.ValueString(),
		GroupDescription: data.Description.ValueString(),
		SiteId:           data.SiteID.ValueStringPointer(),
	}

	criteriaModels, critDiags := schemahelpers.Expand[smartMobileDeviceGroupCriteriaDataModel](ctx, data.Criteria)
	diags.Append(critDiags...)
	if diags.HasError() {
		return nil, diags
	}

	if len(criteriaModels) > 0 {
		resource.Criteria = make([]jamfpro.SharedSubsetCriteriaJamfProAPI, len(criteriaModels))
		for i, criterion := range criteriaModels {
			apiCriterion := jamfpro.SharedSubsetCriteriaJamfProAPI{
				Name:       criterion.Name.ValueString(),
				Priority:   int(criterion.Priority.ValueInt32()),
				AndOr:      criterion.AndOr.ValueString(),
				SearchType: criterion.SearchType.ValueString(),
				Value:      criterion.Value.ValueString(),
			}

			if !criterion.OpeningParen.IsNull() && !criterion.OpeningParen.IsUnknown() {
				val := criterion.OpeningParen.ValueBool()
				apiCriterion.OpeningParen = &val
			}

			if !criterion.ClosingParen.IsNull() && !criterion.ClosingParen.IsUnknown() {
				val := criterion.ClosingParen.ValueBool()
				apiCriterion.ClosingParen = &val
			}

			resource.Criteria[i] = apiCriterion
		}
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		diags.AddError(
			"Failed to marshal Smart Mobile Device Group",
			fmt.Sprintf("Failed to marshal Smart Mobile Device Group '%s' to JSON: %v", resource.GroupName, err),
		)
		return nil, diags
	}

	log.Printf("[DEBUG] Constructed Smart Mobile Device Group JSON:\n%s\n", string(resourceJSON))

	return resource, diags
}
