// diskencryptionconfigurations_object.go
package diskencryptionconfigurations

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/constructobject"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProDiskEncryptionConfiguration constructs a ResourceDiskEncryptionConfiguration object from the provided schema data.
func constructJamfProDiskEncryptionConfiguration(d *schema.ResourceData) (*jamfpro.ResourceDiskEncryptionConfiguration, error) {
	diskEncryptionConfig := &jamfpro.ResourceDiskEncryptionConfiguration{
		Name:                  d.Get("name").(string),
		KeyType:               d.Get("key_type").(string),
		FileVaultEnabledUsers: d.Get("file_vault_enabled_users").(string),
	}

	// Handle institutional_recovery_key
	if v, ok := d.GetOk("institutional_recovery_key"); ok && len(v.([]interface{})) > 0 {
		irkData := v.([]interface{})[0].(map[string]interface{})
		diskEncryptionConfig.InstitutionalRecoveryKey = &jamfpro.DiskEncryptionConfigurationInstitutionalRecoveryKey{
			Key:             irkData["key"].(string),
			CertificateType: irkData["certificate_type"].(string),
			Password:        irkData["password"].(string),
			Data:            irkData["data"].(string),
		}
	}

	// Print the constructed XML output to the log
	xmlOutput, err := constructobject.SerializeAndRedactXML(diskEncryptionConfig, []string{"Password"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Disk Encryption Configurations XML:\n%s\n", string(xmlOutput))

	return diskEncryptionConfig, nil
}
