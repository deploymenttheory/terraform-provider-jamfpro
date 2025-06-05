// staticmobiledevicegroup_object.go
package static_mobile_device_group

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProStaticMobileDeviceGroup constructs a ResourceMobileDeviceGroup object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceMobileDeviceGroup, error) {
	resource := &jamfpro.ResourceMobileDeviceGroup{
		Name:    d.Get("name").(string),
		IsSmart: false,
	}

	resource.Site = *sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))

	assignedMobileDevices := d.Get("assigned_mobile_device_ids").([]interface{})
	if len(assignedMobileDevices) > 0 {
		mobile_devices := []jamfpro.MobileDeviceGroupSubsetDeviceItem{}
		for _, v := range assignedMobileDevices {
			mobile_devices = append(mobile_devices, jamfpro.MobileDeviceGroupSubsetDeviceItem{
				ID: v.(int),
			})
		}
		resource.MobileDevices = mobile_devices
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Mobile Device Group '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Mobile Device Group XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// Helper functions for nested structures
