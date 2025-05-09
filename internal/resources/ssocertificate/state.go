package ssocertificate

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func updateState(d *schema.ResourceData, resp *jamfpro.ResourceSSOKeystoreResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create the keystore map with nested keys
	keystoreMap := map[string]interface{}{
		"key":                 resp.Keystore.Key,
		"type":                resp.Keystore.Type,
		"keystore_file_name":  resp.Keystore.KeystoreFileName,
		"keystore_setup_type": resp.Keystore.KeystoreSetupType,
		"keys":                flattenCertKeys(resp.Keystore.Keys),
	}

	// Set the entire keystore structure
	if err := d.Set("keystore", []interface{}{keystoreMap}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set keystore details if available
	if resp.KeystoreDetails != nil {
		keystoreDetails := map[string]interface{}{
			"keys":          resp.KeystoreDetails.Keys,
			"issuer":        resp.KeystoreDetails.Issuer,
			"subject":       resp.KeystoreDetails.Subject,
			"expiration":    resp.KeystoreDetails.Expiration,
			"serial_number": resp.KeystoreDetails.SerialNumber,
		}
		if err := d.Set("keystore_details", []interface{}{keystoreDetails}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
