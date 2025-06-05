package sso_certificate

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setKeystoreData sets the keystore data in the resource data
func setKeystoreData(d *schema.ResourceData, keystore jamfpro.ResourceSSOCertKeystore) error {
	if err := d.Set("keystore", []interface{}{
		map[string]interface{}{
			"key":                 keystore.Key,
			"type":                keystore.Type,
			"keystore_file_name":  keystore.KeystoreFileName,
			"keystore_setup_type": keystore.KeystoreSetupType,
			"keys":                flattenCertKeys(keystore.Keys),
		},
	}); err != nil {
		return fmt.Errorf("error setting keystore: %v", err)
	}
	return nil
}

// setKeystoreDetails sets the keystore details in the resource data
func setKeystoreDetails(d *schema.ResourceData, details *jamfpro.ResourceSSOKeystoreDetails) error {
	if err := d.Set("keystore_details", []interface{}{
		map[string]interface{}{
			"keys":          details.Keys,
			"issuer":        details.Issuer,
			"subject":       details.Subject,
			"expiration":    details.Expiration,
			"serial_number": details.SerialNumber,
		},
	}); err != nil {
		return fmt.Errorf("error setting keystore details: %v", err)
	}
	return nil
}

// flattenCertKeys flattens the cert keys into a slice of interfaces
func flattenCertKeys(keys []jamfpro.ResourceCertKey) []interface{} {
	var result []interface{}
	for _, key := range keys {
		result = append(result, map[string]interface{}{
			"id":    key.ID,
			"valid": key.Valid,
		})
	}
	return result
}
