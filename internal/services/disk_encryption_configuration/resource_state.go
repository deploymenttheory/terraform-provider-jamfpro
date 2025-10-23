// diskencryptionconfigurations_state.go
package disk_encryption_configuration

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Disk Encryption Configuration information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceDiskEncryptionConfiguration) diag.Diagnostics {
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
		irk := make(map[string]any)
		irk["certificate_type"] = resp.InstitutionalRecoveryKey.CertificateType
		irk["data"] = resp.InstitutionalRecoveryKey.Data

		if err := d.Set("institutional_recovery_key", []any{irk}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

	} else {
		if err := d.Set("institutional_recovery_key", []any{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
