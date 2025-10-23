package smtp_server

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func updateState(d *schema.ResourceData, resp *jamfpro.ResourceSMTPServer) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := map[string]any{
		"enabled":             resp.Enabled,
		"authentication_type": resp.AuthenticationType,
	}

	// Set Connection Settings if present
	if resp.ConnectionSettings != nil {
		connSettings := []map[string]any{
			{
				"host":               resp.ConnectionSettings.Host,
				"port":               resp.ConnectionSettings.Port,
				"encryption_type":    resp.ConnectionSettings.EncryptionType,
				"connection_timeout": resp.ConnectionSettings.ConnectionTimeout,
			},
		}
		if err := d.Set("connection_settings", connSettings); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set Sender Settings if present
	if resp.SenderSettings != nil {
		senderSettings := []map[string]any{
			{
				"display_name":  resp.SenderSettings.DisplayName,
				"email_address": resp.SenderSettings.EmailAddress,
			},
		}
		if err := d.Set("sender_settings", senderSettings); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set Basic Auth Credentials if present
	if resp.BasicAuthCredentials != nil {
		basicAuth := []map[string]any{
			{
				"username": resp.BasicAuthCredentials.Username,
				"password": d.Get("basic_auth_credentials.0.password").(string),
			},
		}
		if err := d.Set("basic_auth_credentials", basicAuth); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set Graph API Credentials if present
	if resp.GraphApiCredentials != nil {
		graphApi := []map[string]any{
			{
				"tenant_id":     resp.GraphApiCredentials.TenantId,
				"client_id":     resp.GraphApiCredentials.ClientId,
				"client_secret": d.Get("graph_api_credentials.0.client_secret").(string),
			},
		}
		if err := d.Set("graph_api_credentials", graphApi); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set Google Mail Credentials if present
	if resp.GoogleMailCredentials != nil {
		googleMail := []map[string]any{
			{
				"client_id":     resp.GoogleMailCredentials.ClientId,
				"client_secret": d.Get("google_mail_credentials.0.client_secret").(string),
			},
		}
		if err := d.Set("google_mail_credentials", googleMail); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		// Set authentications if present
		var auths []map[string]any
		for _, auth := range resp.GoogleMailCredentials.Authentications {
			auths = append(auths, map[string]any{
				"email_address": auth.EmailAddress,
				"status":        auth.Status,
			})
		}
		if err := d.Set("authentications", auths); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Set base settings
	for key, val := range settings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
