package static_computer_group

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// state updates the Terraform model with the latest Static Computer Group V2 information from the Jamf Pro API.
func state(data *staticComputerGroupResourceModel, resp *jamfpro.ResponseStaticComputerGroupListItemV2) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(resp.ID)
	data.Name = types.StringValue(resp.Name)
	data.SiteID = types.StringValue(resp.SiteID)

	if resp.Description == "" && (data.Description.IsNull() || data.Description.IsUnknown()) {
		data.Description = types.StringNull()
	} else {
		data.Description = types.StringValue(resp.Description)
	}

	return diags
}
