package accountdrivenuserenrollmentsettings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Account Driven User Enrollment Settings from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceADUETokenSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("enabled", resp.Enabled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("expiration_interval_days", resp.ExpirationIntervalDays); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("expiration_interval_seconds", resp.ExpirationIntervalSeconds); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
