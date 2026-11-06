package disk_encryption_configuration

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUpdateStateOmitsEmptyInstitutionalRecoveryKey(t *testing.T) {
	for _, test := range []struct {
		name  string
		key   *jamfpro.DiskEncryptionConfigurationInstitutionalRecoveryKey
		count int
	}{
		{name: "omitted"},
		{name: "empty response object", key: &jamfpro.DiskEncryptionConfigurationInstitutionalRecoveryKey{}},
		{name: "certificate", key: &jamfpro.DiskEncryptionConfigurationInstitutionalRecoveryKey{CertificateType: "PK12", Data: "dGVzdA=="}, count: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, ResourceJamfProDiskEncryptionConfigurations().Schema, map[string]any{
				"name": "tf-testing-filevault", "key_type": "Individual", "file_vault_enabled_users": "Current or Next User",
			})
			diags := updateState(data, &jamfpro.ResourceDiskEncryptionConfiguration{
				ID: 1, Name: "tf-testing-filevault", KeyType: "Individual", FileVaultEnabledUsers: "Current or Next User", InstitutionalRecoveryKey: test.key,
			})
			require.False(t, diags.HasError(), "%v", diags)
			require.Len(t, data.Get("institutional_recovery_key").([]any), test.count)
		})
	}
}
