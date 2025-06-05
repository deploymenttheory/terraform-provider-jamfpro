package sso_failover

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the values from the Jamf Pro SSO Failover response.
func updateState(d *schema.ResourceData, resp *jamfpro.ResponseSSOFailover) diag.Diagnostics {
	var diags diag.Diagnostics

	oldURL, hasOldURL := d.GetOk("failover_url")
	oldTime, hasOldTime := d.GetOk("generation_time")

	if err := d.Set("failover_url", resp.FailoverURL); err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("generation_time", resp.GenerationTime); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if hasOldURL && hasOldTime && d.Id() != "" {
		if oldURL.(string) != resp.FailoverURL &&
			oldTime.(int) != int(resp.GenerationTime) {
			d.SetId("")
			return diags
		}
	}

	d.SetId("jamfpro_sso_failover_singleton")
	return diags
}
