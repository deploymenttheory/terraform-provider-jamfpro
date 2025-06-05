// Constructor file (device_communication_settings_constructor.go)
package device_communication_settings

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func construct(d *schema.ResourceData) (*jamfpro.ResourceDeviceCommunicationSettings, error) {
	resource := &jamfpro.ResourceDeviceCommunicationSettings{
		AutoRenewMobileDeviceMdmProfileWhenCaRenewed:                  d.Get("auto_renew_mobile_device_mdm_profile_when_ca_renewed").(bool),
		AutoRenewMobileDeviceMdmProfileWhenDeviceIdentityCertExpiring: d.Get("auto_renew_mobile_device_mdm_profile_when_device_identity_cert_expiring").(bool),
		AutoRenewComputerMdmProfileWhenCaRenewed:                      d.Get("auto_renew_computer_mdm_profile_when_ca_renewed").(bool),
		AutoRenewComputerMdmProfileWhenDeviceIdentityCertExpiring:     d.Get("auto_renew_computer_mdm_profile_when_device_identity_cert_expiring").(bool),
		MdmProfileMobileDeviceExpirationLimitInDays:                   d.Get("mdm_profile_mobile_device_expiration_limit_in_days").(int),
		MdmProfileComputerExpirationLimitInDays:                       d.Get("mdm_profile_computer_expiration_limit_in_days").(int),
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Device Communication Settings to JSON: %v", err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Device Communication Settings JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
