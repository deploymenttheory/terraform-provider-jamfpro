// webhooks_object.go
package webhooks

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/constructobject"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProWebhook constructs a ResourceWebhook object from the provided schema data.
func constructJamfProWebhook(d *schema.ResourceData) (*jamfpro.ResourceWebhook, error) {
	var resource *jamfpro.ResourceWebhook

	resource = &jamfpro.ResourceWebhook{
		Name:                        d.Get("name").(string),
		Enabled:                     d.Get("enabled").(bool),
		URL:                         d.Get("url").(string),
		ContentType:                 d.Get("content_type").(string),
		Event:                       d.Get("event").(string),
		ConnectionTimeout:           d.Get("connection_timeout").(int),
		ReadTimeout:                 d.Get("read_timeout").(int),
		AuthenticationType:          d.Get("authentication_type").(string),
		Username:                    d.Get("username").(string),
		Password:                    d.Get("password").(string),
		EnableDisplayFieldsForGroup: d.Get("enable_display_fields_for_group").(bool),
		SmartGroupID:                d.Get("smart_group_id").(int),
	}

	if v, ok := d.GetOk("display_fields"); ok {
		displayFieldsData := v.([]interface{})
		for _, fieldData := range displayFieldsData {
			field, ok := fieldData.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("display_field is not a valid map")
			}
			var displayFields []jamfpro.SharedAdvancedSearchSubsetDisplayField
			if subFieldsData, ok := field["display_field"].([]interface{}); ok {
				for _, subFieldData := range subFieldsData {
					subField, ok := subFieldData.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("sub_display_field is not a valid map")
					}
					if name, ok := subField["name"].(string); ok {
						displayFields = append(displayFields, jamfpro.SharedAdvancedSearchSubsetDisplayField{
							Name: name,
						})
					}
				}
			}
			resource.DisplayFields = append(resource.DisplayFields, jamfpro.SharedAdvancedSearchContainerDisplayField{
				DisplayField: displayFields,
			})
		}
	}

	// Serialize and log the XML output for debugging
	xmlOutput, err := constructobject.SerializeAndRedactXML(resource, []string{"Password"})
	if err != nil {
		log.Fatalf("Error serializing webhook to XML: %v", err)
	}
	log.Printf("[DEBUG] Constructed Jamf Pro Webhook XML:\n%s\n", xmlOutput)

	return resource, nil
}
