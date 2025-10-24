// State file (device_communication_settings_state.go)
package device_communication_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func updateState(d *schema.ResourceData, resp *jamfpro.ResourceDeviceCommunicationSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := map[string]any{
		"auto_renew_mobile_device_mdm_profile_when_ca_renewed":                    resp.AutoRenewMobileDeviceMdmProfileWhenCaRenewed,
		"auto_renew_mobile_device_mdm_profile_when_device_identity_cert_expiring": resp.AutoRenewMobileDeviceMdmProfileWhenDeviceIdentityCertExpiring,
		"auto_renew_computer_mdm_profile_when_ca_renewed":                         resp.AutoRenewComputerMdmProfileWhenCaRenewed,
		"auto_renew_computer_mdm_profile_when_device_identity_cert_expiring":      resp.AutoRenewComputerMdmProfileWhenDeviceIdentityCertExpiring,
		"mdm_profile_mobile_device_expiration_limit_in_days":                      resp.MdmProfileMobileDeviceExpirationLimitInDays,
		"mdm_profile_computer_expiration_limit_in_days":                           resp.MdmProfileComputerExpirationLimitInDays,
	}

	for key, val := range settings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
