package sso_certificate

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead reads the SSO certificate settings from Jamf Pro
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	var response *jamfpro.ResourceSSOKeystoreResponse
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetSSOCertificate()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read SSO certificate settings: %v", err))
	}

	d.SetId(fmt.Sprintf("jamfpro_sso_certificate_%s", response.Keystore.KeystoreFileName))

	if err := d.Set("keystore", []interface{}{
		map[string]interface{}{
			"key":                 response.Keystore.Key,
			"type":                response.Keystore.Type,
			"keystore_file_name":  response.Keystore.KeystoreFileName,
			"keystore_setup_type": response.Keystore.KeystoreSetupType,
			"keys":                flattenCertKeys(response.Keystore.Keys),
		},
	}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if response.KeystoreDetails != nil {
		if err := d.Set("keystore_details", []interface{}{
			map[string]interface{}{
				"keys":          response.KeystoreDetails.Keys,
				"issuer":        response.KeystoreDetails.Issuer,
				"subject":       response.KeystoreDetails.Subject,
				"expiration":    response.KeystoreDetails.Expiration,
				"serial_number": response.KeystoreDetails.SerialNumber,
			},
		}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
