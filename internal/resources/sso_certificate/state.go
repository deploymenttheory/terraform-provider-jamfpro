package ssocertificate

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the state of the resource with the provided response
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceSSOKeystoreResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	var oldKeystore map[string]interface{}
	var oldDetails map[string]interface{}

	if keystoreList, ok := d.GetOk("keystore"); ok && len(keystoreList.([]interface{})) > 0 {
		oldKeystore = keystoreList.([]interface{})[0].(map[string]interface{})
	}

	if detailsList, ok := d.GetOk("keystore_details"); ok && len(detailsList.([]interface{})) > 0 {
		oldDetails = detailsList.([]interface{})[0].(map[string]interface{})
	}

	if err := setKeystoreData(d, resp.Keystore); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if resp.KeystoreDetails != nil {
		if err := setKeystoreDetails(d, resp.KeystoreDetails); err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}

	if hasSignificantChanges(oldKeystore, oldDetails, resp) {
		d.SetId("")
		return diags
	}

	d.SetId("jamfpro_sso_certificate_singleton")
	return diags
}

// hasSignificantChanges checks if there are significant changes in the keystore and details
func hasSignificantChanges(oldKeystore, oldDetails map[string]interface{}, resp *jamfpro.ResourceSSOKeystoreResponse) bool {
	if oldKeystore == nil || oldDetails == nil {
		return false
	}

	if oldKeystore["key"] != resp.Keystore.Key {
		return true
	}

	if resp.KeystoreDetails != nil {
		if oldDetails["serial_number"] != resp.KeystoreDetails.SerialNumber {
			return true
		}
		if oldDetails["expiration"] != resp.KeystoreDetails.Expiration {
			return true
		}
	}

	oldKeys, ok := oldKeystore["keys"].([]interface{})
	if ok && len(oldKeys) > 0 {
		for i, key := range resp.Keystore.Keys {
			if i >= len(oldKeys) {
				return true
			}
			oldKey := oldKeys[i].(map[string]interface{})
			if oldKey["id"] != key.ID || oldKey["valid"] != key.Valid {
				return true
			}
		}
	}

	return false
}
