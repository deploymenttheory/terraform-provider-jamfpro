// Resource file (device_communication_settings_resource.go)
package devicecommunicationsettings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProDeviceCommunicationSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew_mobile_device_mdm_profile_when_ca_renewed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Automatically renew mobile device MDM profile when CA is renewed",
			},
			"auto_renew_mobile_device_mdm_profile_when_device_identity_cert_expiring": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Automatically renew mobile device MDM profile when device identity certificate is expiring",
			},
			"auto_renew_computer_mdm_profile_when_ca_renewed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Automatically renew computer MDM profile when CA is renewed",
			},
			"auto_renew_computer_mdm_profile_when_device_identity_cert_expiring": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Automatically renew computer MDM profile when device identity certificate is expiring",
			},
			"mdm_profile_mobile_device_expiration_limit_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      90,
				Description:  "MDM profile expiration limit in days for mobile devices. Valid values are 90, 120, or 180 days.",
				ValidateFunc: validation.IntInSlice([]int{90, 120, 180}),
			},
			"mdm_profile_computer_expiration_limit_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      90,
				Description:  "MDM profile expiration limit in days for computers. Valid values are 90, 120, or 180 days.",
				ValidateFunc: validation.IntInSlice([]int{90, 120, 180}),
			},
		},
	}
}
