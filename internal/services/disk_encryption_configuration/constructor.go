// diskencryptionconfigurations_object.go
package disk_encryption_configuration

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProDiskEncryptionConfiguration constructs a ResourceDiskEncryptionConfiguration object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceDiskEncryptionConfiguration, error) {
	resource := &jamfpro.ResourceDiskEncryptionConfiguration{
		Name:                  d.Get("name").(string),
		KeyType:               d.Get("key_type").(string),
		FileVaultEnabledUsers: d.Get("file_vault_enabled_users").(string),
	}

	if v, ok := d.GetOk("institutional_recovery_key"); ok && len(v.([]any)) > 0 {
		irkData := v.([]any)[0].(map[string]any)
		resource.InstitutionalRecoveryKey = &jamfpro.DiskEncryptionConfigurationInstitutionalRecoveryKey{
			Key:             irkData["key"].(string),
			CertificateType: irkData["certificate_type"].(string),
			Password:        irkData["password"].(string),
			Data:            irkData["data"].(string),
		}
	}

	xmlOutput, err := sdkv2.SerializeAndRedactXML(resource, []string{"Password"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Disk Encryption Configurations XML:\n%s\n", string(xmlOutput))

	return resource, nil
}
