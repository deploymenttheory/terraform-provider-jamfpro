package static_computer_group

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// updateState updates the Terraform state with the latest Computer Group information from the Jamf Pro API.
func updateState(ctx context.Context, data *staticComputerGroupDataSourceModel, resp *jamfpro.ResponseStaticComputerGroupListItemV2) diag.Diagnostics {
	var diags diag.Diagnostics

	data.Name = types.StringValue(resp.Name)

	if resp.Description != "" {
		data.Description = types.StringValue(resp.Description)
	} else {
		data.Description = types.StringNull()
	}

	data.SiteID = types.StringValue(resp.SiteID)

	return diags
}
