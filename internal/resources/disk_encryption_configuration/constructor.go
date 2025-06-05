// diskencryptionconfigurations_object.go
package diskencryptionconfigurations

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProDiskEncryptionConfiguration constructs a ResourceDiskEncryptionConfiguration object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceDiskEncryptionConfiguration, error) {
	resource := &jamfpro.ResourceDiskEncryptionConfiguration{
		Name:                  d.Get("name").(string),
		KeyType:               d.Get("key_type").(string),
		FileVaultEnabledUsers: d.Get("file_vault_enabled_users").(string),
	}

	if v, ok := d.GetOk("institutional_recovery_key"); ok && len(v.([]interface{})) > 0 {
		irkData := v.([]interface{})[0].(map[string]interface{})
		resource.InstitutionalRecoveryKey = &jamfpro.DiskEncryptionConfigurationInstitutionalRecoveryKey{
			Key:             irkData["key"].(string),
			CertificateType: irkData["certificate_type"].(string),
			Password:        irkData["password"].(string),
			Data:            irkData["data"].(string),
		}
	}

	xmlOutput, err := common.SerializeAndRedactXML(resource, []string{"Password"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Disk Encryption Configurations XML:\n%s\n", string(xmlOutput))

	return resource, nil
}
