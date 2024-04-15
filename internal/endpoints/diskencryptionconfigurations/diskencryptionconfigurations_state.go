// diskencryptionconfigurations_state.go
package diskencryptionconfigurations

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Disk EncryptionC onfiguration information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceDiskEncryptionConfiguration) diag.Diagnostics {
	var diags diag.Diagnostics

	// Assuming successful retrieval, proceed to set the resource attributes in Terraform state
	if resource != nil {
		// Set the fields directly in the Terraform state
		if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("key_type", resource.KeyType); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("file_vault_enabled_users", resource.FileVaultEnabledUsers); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		// Handle Institutional Recovery Key
		if resource.InstitutionalRecoveryKey != nil {
			irk := make(map[string]interface{})
			irk["certificate_type"] = resource.InstitutionalRecoveryKey.CertificateType
			//irk["password"] = resource.InstitutionalRecoveryKey.Password // Uncomment if password should be set
			irk["data"] = resource.InstitutionalRecoveryKey.Data

			if err := d.Set("institutional_recovery_key", []interface{}{irk}); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		} else {
			// Ensure institutional_recovery_key is not set in the Terraform state if nil or empty
			if err := d.Set("institutional_recovery_key", []interface{}{}); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	}
	return diags
}
