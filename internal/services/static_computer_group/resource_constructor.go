package static_computer_group

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// constructResource constructs a ResourceStaticComputerGroupV2 object from the provided framework resource model.
func constructResource(data *staticComputerGroupResourceModel) (*jamfpro.ResourceStaticComputerGroupV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	resource := &jamfpro.ResourceStaticComputerGroupV2{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		SiteId:      data.SiteID.ValueStringPointer(),
	}

	if !data.AssignedComputerIDs.IsNull() && !data.AssignedComputerIDs.IsUnknown() {
		elements := data.AssignedComputerIDs.Elements()
		computerIDs := make([]string, 0, len(elements))
		for _, elem := range elements {
			computerIDs = append(computerIDs, elem.(types.String).ValueString())
		}
		resource.Assignments = computerIDs
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		diags.AddError(
			"Failed to marshal Static Computer Group",
			fmt.Sprintf("Failed to marshal Static Computer Group static_computer_group'%s' to JSON: %v", resource.Name, err),
		)
		return nil, diags
	}

	log.Printf("[DEBUG] Constructed Static Computer Group static_computer_groupJSON:\n%s\n", string(resourceJSON))

	return resource, diags
}
