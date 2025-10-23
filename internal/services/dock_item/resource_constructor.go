package dock_item

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// constructResource constructs a ResourceDockItem object from the provided framework resource model.
func constructResource(data *dockItemResourceModel) (*jamfpro.ResourceDockItem, diag.Diagnostics) {
	var diags diag.Diagnostics

	resource := &jamfpro.ResourceDockItem{
		Name:     data.Name.ValueString(),
		Type:     data.Type.ValueString(),
		Path:     data.Path.ValueString(),
		Contents: data.Contents.ValueString(),
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		diags.AddError(
			"Failed to marshal Jamf Pro Dock Item",
			fmt.Sprintf("Failed to marshal Jamf Pro Dock Item '%s' to XML: %v", resource.Name, err),
		)
		return nil, diags
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Dock Item XML:\n%s\n", string(resourceXML))

	return resource, diags
}
