package mobile_device_extension_attribute

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct builds a ResourceMobileDeviceExtensionAttribute object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceMobileDeviceExtensionAttribute, error) {
	resource := &jamfpro.ResourceMobileDeviceExtensionAttribute{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		DataType:             d.Get("data_type").(string),
		InventoryDisplayType: d.Get("inventory_display_type").(string),
		InputType:            d.Get("input_type").(string),
	}

	if v, ok := d.GetOk("popup_menu_choices"); ok {
		set := v.(*schema.Set)
		for _, choice := range set.List() {
			resource.PopupMenuChoices = append(resource.PopupMenuChoices, choice.(string))
		}
	}

	if v, ok := d.GetOk("ldap_attribute_mapping"); ok {
		resource.LDAPAttributeMapping = v.(string)
	}

	if v, ok := d.GetOk("ldap_extension_attribute_allowed"); ok {
		resource.LDAPExtensionAttributeAllowed = jamfpro.BoolPtr(v.(bool))
	}

	if err := validateInputType(resource); err != nil {
		return nil, fmt.Errorf("failed to construct: %w", err)
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Mobile Device Extension Attribute to JSON: %w", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Mobile Device Extension Attribute JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
