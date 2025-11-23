package static_computer_group

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read fetches the static computer group data from Jamf Pro.
func (d *staticComputerGroupFrameworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data staticComputerGroupDataSourceModel

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

	var resource *jamfpro.ResponseStaticComputerGroupListItemV2
	var err error

	if name != "" {
		tflog.Debug(ctx, fmt.Sprintf("Reading Static Computer Group by name: %s", name))
		resource, err = d.client.GetStaticComputerGroupByNameV2(name)
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Reading Static Computer Group by ID: %s", resourceID))
		resource, err = d.client.GetStaticComputerGroupByIDV2(resourceID)
	}

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if name != "" {
			lookupMethod = "name"
			lookupValue = name
		}
		resp.Diagnostics.AddError(
			"Error Reading Static Computer Group",
			fmt.Sprintf("Failed to read Static Computer Group with %s '%s': %v", lookupMethod, lookupValue, err),
		)
		return
	}

	if resource == nil {
		resp.Diagnostics.AddError(
			"Resource Not Found",
			"The Jamf Pro Static Computer Group was not found",
		)
		return
	}

	data.ID = types.StringValue(resource.ID)

	resp.Diagnostics.Append(updateState(ctx, &data, resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
