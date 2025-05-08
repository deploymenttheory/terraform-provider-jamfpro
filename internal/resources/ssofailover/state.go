package ssofailover

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func updateState(d *schema.ResourceData, resp *jamfpro.ResponseSSOFailover) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := map[string]interface{}{
		"failover_url":    resp.FailoverURL,
		"generation_time": resp.GenerationTime,
	}

	// Set all settings
	for key, val := range settings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
