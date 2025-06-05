// webhooks_object.go
package webhook

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProWebhook constructs a ResourceWebhook object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceWebhook, error) {
	resource := &jamfpro.ResourceWebhook{
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
		Header:                      d.Get("header").(string),
		SmartGroupID:                d.Get("smart_group_id").(int),
	}

	displayFieldsHcl := d.Get("display_fields").([]interface{})
	if len(displayFieldsHcl) > 0 {
		for _, v := range displayFieldsHcl {
			resource.DisplayFields = append(resource.DisplayFields, jamfpro.DisplayField{Name: v.(string)})
		}
	}

	// Serialize and log the XML output for debugging
	xmlOutput, err := common.SerializeAndRedactXML(resource, []string{"Password"})
	if err != nil {
		log.Fatalf("Error serializing webhook to XML: %v", err)
	}
	log.Printf("[DEBUG] Constructed Jamf Pro Webhook XML:\n%s\n", xmlOutput)

	return resource, nil
}
