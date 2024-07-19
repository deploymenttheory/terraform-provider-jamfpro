// diskencryptionconfigurations_state.go
package diskencryptionconfigurations

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Disk EncryptionC onfiguration information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceDiskEncryptionConfiguration) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", strconv.Itoa(resp.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("key_type", resp.KeyType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("file_vault_enabled_users", resp.FileVaultEnabledUsers); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if resp.InstitutionalRecoveryKey != nil {
		irk := make(map[string]interface{})
		irk["certificate_type"] = resp.InstitutionalRecoveryKey.CertificateType
		irk["data"] = resp.InstitutionalRecoveryKey.Data

		if err := d.Set("institutional_recovery_key", []interface{}{irk}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

	} else {
		if err := d.Set("institutional_recovery_key", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
