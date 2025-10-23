package smtp_server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func construct(d *schema.ResourceData) (*jamfpro.ResourceSMTPServer, error) {
	resource := &jamfpro.ResourceSMTPServer{
		Enabled:            d.Get("enabled").(bool),
		AuthenticationType: d.Get("authentication_type").(string),
	}

	// Handle Connection Settings
	if v, ok := d.GetOk("connection_settings"); ok && len(v.([]any)) > 0 {
		connSettings := v.([]any)[0].(map[string]any)
		resource.ConnectionSettings = &jamfpro.ResourceSMTPServerConnectionSettings{
			Host:              connSettings["host"].(string),
			Port:              connSettings["port"].(int),
			EncryptionType:    connSettings["encryption_type"].(string),
			ConnectionTimeout: connSettings["connection_timeout"].(int),
		}
	}

	// Handle Sender Settings
	if v, ok := d.GetOk("sender_settings"); ok && len(v.([]any)) > 0 {
		senderSettings := v.([]any)[0].(map[string]any)
		resource.SenderSettings = &jamfpro.ResourceSMTPServerSenderSettings{
			DisplayName:  senderSettings["display_name"].(string),
			EmailAddress: senderSettings["email_address"].(string),
		}
	}

	// Handle Basic Auth Credentials
	if v, ok := d.GetOk("basic_auth_credentials"); ok && len(v.([]any)) > 0 {
		basicAuth := v.([]any)[0].(map[string]any)
		resource.BasicAuthCredentials = &jamfpro.ResourceSMTPServerBasicAuthCredentials{
			Username: basicAuth["username"].(string),
			Password: basicAuth["password"].(string),
		}
	}

	// Handle Graph API Credentials
	if v, ok := d.GetOk("graph_api_credentials"); ok && len(v.([]any)) > 0 {
		graphApi := v.([]any)[0].(map[string]any)
		resource.GraphApiCredentials = &jamfpro.ResourceSMTPServerGraphApiCredentials{
			TenantId:     graphApi["tenant_id"].(string),
			ClientId:     graphApi["client_id"].(string),
			ClientSecret: graphApi["client_secret"].(string),
		}
	}

	// Handle Google Mail Credentials
	if v, ok := d.GetOk("google_mail_credentials"); ok && len(v.([]any)) > 0 {
		googleMail := v.([]any)[0].(map[string]any)
		resource.GoogleMailCredentials = &jamfpro.ResourceSMTPServerGoogleMailCredentials{
			ClientId:     googleMail["client_id"].(string),
			ClientSecret: googleMail["client_secret"].(string),
		}
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro SMTP Server Settings to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro SMTP Server Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
