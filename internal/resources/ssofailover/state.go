package ssofailover

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the state of the SSO failover resource with the provided response data.
func updateState(d *schema.ResourceData, resp *jamfpro.ResponseSSOFailover) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := map[string]interface{}{
		"failover_url":    resp.FailoverURL,
		"generation_time": resp.GenerationTime,
	}

	for key, val := range settings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
