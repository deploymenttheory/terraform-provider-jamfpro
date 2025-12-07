package smart_computer_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read fetches the static computer group data from Jamf Pro.
func (d *smartComputerGroupFrameworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data smartComputerGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := data.ID.ValueString()
	name := data.Name.ValueString()

	if resourceID == "" && name == "" {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be provided",
		)
		return
	}

	if name != "" {
		tflog.Debug(ctx, fmt.Sprintf("Looking up Smart Computer Group by name: %s", name))
		listItem, err := d.client.GetSmartComputerGroupByNameV2(name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Finding Smart Computer Group",
				fmt.Sprintf("Failed to find Smart Computer Group with name '%s': %v", name, err),
			)
			return
		}
		if listItem == nil {
			resp.Diagnostics.AddError(
				"Resource Not Found",
				fmt.Sprintf("Smart Computer Group with name '%s' was not found", name),
			)
			return
		}
		resourceID = listItem.ID
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading Smart Computer Group by ID: %s", resourceID))
	resource, err := d.client.GetSmartComputerGroupByIDV2(resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Smart Computer Group",
			fmt.Sprintf("Failed to read Smart Computer Group with ID '%s': %v", resourceID, err),
		)
		return
	}

	if resource == nil {
		resp.Diagnostics.AddError(
			"Resource Not Found",
			fmt.Sprintf("Smart Computer Group with ID '%s' was not found", resourceID),
		)
		return
	}

	data.ID = types.StringValue(resourceID)

	data.Name = types.StringValue(resource.Name)

	if resource.Description != "" {
		data.Description = types.StringValue(resource.Description)
	} else {
		data.Description = types.StringNull()
	}

	if resource.SiteId != nil && *resource.SiteId != "" {
		data.SiteID = types.StringValue(*resource.SiteId)
	} else {
		data.SiteID = types.StringNull()
	}

	data.Criteria = make([]smartComputerGroupCriteriaDataModel, 0, len(resource.Criteria))
	for _, criterion := range resource.Criteria {
		criteriaModel := smartComputerGroupCriteriaDataModel{
			Name:       types.StringValue(criterion.Name),
			Priority:   types.Int64Value(int64(criterion.Priority)),
			AndOr:      types.StringValue(criterion.AndOr),
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
