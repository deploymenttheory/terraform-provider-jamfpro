package sso_failover

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProSSOFailover provides information about Jamf Pro's SSO Failover configuration
func DataSourceJamfProSSOFailover() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"failover_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SSO failover URL for Jamf Pro",
			},
			"generation_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timestamp when the failover URL was generated",
			},
		},
	}
}

// dataSourceRead reads the SSO failover settings from Jamf Pro
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	var response *jamfpro.ResponseSSOFailover
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetSSOFailoverSettings()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read SSO failover settings: %v", err))
	}

	d.SetId(fmt.Sprintf("jamfpro_sso_failover_%d", response.GenerationTime))

	if err := d.Set("failover_url", response.FailoverURL); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("generation_time", response.GenerationTime); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
