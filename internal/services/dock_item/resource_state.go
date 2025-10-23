package dock_item

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// state updates the resource state with the latest Dock Item information from the Jamf Pro API.
func state(data *dockItemResourceModel, resp *jamfpro.ResourceDockItem) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(resp.ID))
	data.Name = types.StringValue(resp.Name)
	data.Type = types.StringValue(resp.Type)
	data.Path = types.StringValue(resp.Path)
	data.Contents = types.StringValue(resp.Contents)

	return diags
}
