// diskencryptionconfigurations_object.go
package diskencryptionconfigurations

import (
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructDiskEncryptionConfiguration constructs a ResourceDiskEncryptionConfiguration object from the provided schema data.
func constructDiskEncryptionConfiguration(d *schema.ResourceData) (*jamfpro.ResourceDiskEncryptionConfiguration, error) {
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

	// Serialize and pretty-print the diskEncryptionConfig object as XML for logging
	resourceXML, err := xml.MarshalIndent(diskEncryptionConfig, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Disk Encryption Configuration '%s' to XML: %v", diskEncryptionConfig.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro Disk Encryption Configuration XML:\n%s\n", string(resourceXML))

	return diskEncryptionConfig, nil
}
