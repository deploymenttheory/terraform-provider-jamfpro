// deviceenrollments_public_key_data_source.go
package device_enrollments_publickey

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDeviceEnrollmentsPublicKey provides the public key for device enrollments in Jamf Pro.
func DataSourceJamfProDeviceEnrollmentsPublicKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public key used for device enrollments.",
			},
		},
	}
}

// dataSourcePublicKeyRead fetches the device enrollments public key from Jamf Pro.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	publicKey, err := client.GetDeviceEnrollmentsPublicKey()
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch device enrollments public key: %v", err))
	}

	d.SetId("jamfpro_device_enrollments_public_key_singleton")

	if err := d.Set("public_key", publicKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting public_key: %v", err))
	}

	return diag.Diagnostics{}
}
