package smart_mobile_device_group_v1

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read fetches the smart mobile device group data from Jamf Pro.
func (d *smartMobileDeviceGroupFrameworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data smartMobileDeviceGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := data.ID.ValueString()
	name := data.Name.ValueString()

	var resource *jamfpro.ResourceSmartMobileDeviceGroupV1

	switch {
	case resourceID != "":
		tflog.Debug(ctx, fmt.Sprintf("Reading Smart Mobile Device Group by ID: %s", resourceID))
		fetchedResource, err := d.client.GetSmartMobileDeviceGroupByIDV1(resourceID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Smart Mobile Device Group",
				fmt.Sprintf("Failed to read Smart Mobile Device Group with ID '%s': %v", resourceID, err),
			)
			return
		}
		if fetchedResource == nil {
			resp.Diagnostics.AddError(
				"Resource Not Found",
				fmt.Sprintf("Smart Mobile Device Group with ID '%s' was not found", resourceID),
			)
			return
		}
		resource = fetchedResource
	case name != "":
		tflog.Debug(ctx, fmt.Sprintf("Looking up Smart Mobile Device Group by name: %s", name))
		listItem, err := d.client.GetSmartMobileDeviceGroupByNameV1(name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Finding Smart Mobile Device Group",
				fmt.Sprintf("Failed to find Smart Mobile Device Group with name '%s': %v", name, err),
			)
			return
		}
		if listItem == nil {
			resp.Diagnostics.AddError(
				"Resource Not Found",
				fmt.Sprintf("Smart Mobile Device Group with name '%s' was not found", name),
			)
			return
		}
		resourceID = listItem.GroupID
		tflog.Debug(ctx, fmt.Sprintf("Reading Smart Mobile Device Group by ID: %s", resourceID))
		fetchedResource, err := d.client.GetSmartMobileDeviceGroupByIDV1(resourceID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Smart Mobile Device Group",
				fmt.Sprintf("Failed to read Smart Mobile Device Group with ID '%s': %v", resourceID, err),
			)
			return
		}
		if fetchedResource == nil {
			resp.Diagnostics.AddError(
				"Resource Not Found",
				fmt.Sprintf("Smart Mobile Device Group with ID '%s' was not found", resourceID),
			)
			return
		}
		resource = fetchedResource
	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be provided",
		)
		return
	}

	data.ID = types.StringValue(resourceID)

	data.Name = types.StringValue(resource.GroupName)

	if resource.GroupDescription != "" {
		data.Description = types.StringValue(resource.GroupDescription)
	} else {
		data.Description = types.StringNull()
	}

	if resource.SiteId != nil && *resource.SiteId != "" {
		data.SiteID = types.StringValue(*resource.SiteId)
	} else {
		data.SiteID = types.StringNull()
	}

	criteriaModels := make([]smartMobileDeviceGroupCriteriaDataModel, 0, len(resource.Criteria))
	for _, criterion := range resource.Criteria {
		criteriaModel := smartMobileDeviceGroupCriteriaDataModel{
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
	resp.Diagnostics.Append(criteriaDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Criteria = criteriaList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
