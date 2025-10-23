package device_enrollments_public_key

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourcePublicKeyRead fetches the device enrollments public key from Jamf Pro.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
