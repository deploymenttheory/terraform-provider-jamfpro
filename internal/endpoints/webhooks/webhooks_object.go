// webhooks_object.go
package webhooks

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/constructobject"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProWebhook constructs a ResourceWebhook object from the provided schema data.
func constructJamfProWebhook(d *schema.ResourceData) (*jamfpro.ResourceWebhook, error) {
	webhook := &jamfpro.ResourceWebhook{
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

	// Handle Display Fields if provided
	if v, ok := d.GetOk("display_fields"); ok {
		displayFieldsData := v.([]interface{})
		for _, fieldData := range displayFieldsData {
			field := fieldData.(map[string]interface{})
			displayField := jamfpro.SharedAdvancedSearchContainerDisplayField{
				Size: field["size"].(int),
			}

			subFieldsData := field["display_field"].([]interface{})
			for _, subFieldData := range subFieldsData {
				subField := subFieldData.(map[string]interface{})
				displayField.DisplayField = append(displayField.DisplayField, jamfpro.SharedAdvancedSearchSubsetDisplayField{
					Name: subField["name"].(string),
				})
			}

			webhook.DisplayFields = append(webhook.DisplayFields, displayField)
		}
	}

	// Print the constructed XML output to the log
	xmlOutput, err := constructobject.SerializeAndRedactXML(webhook, []string{"password"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Webhook XML:\n%s\n", xmlOutput)

	return webhook, nil
}
