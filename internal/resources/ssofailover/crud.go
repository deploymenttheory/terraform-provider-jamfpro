package ssofailover

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro SSO Failover URL in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, apiErr := client.UpdateFailoverUrl()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		d.SetId("jamfpro_sso_failover_singleton")
		d.Set("failover_url", response.FailoverURL)
		d.Set("generation_time", response.GenerationTime)
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create SSO failover settings: %v", err))
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of a Jamf Pro SSO Failover URL from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	d.SetId("jamfpro_sso_failover_singleton")

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
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	d.Set("failover_url", response.FailoverURL)
	d.Set("generation_time", response.GenerationTime)

	return append(diags, updateState(d, response)...)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating an existing Jamf Pro SSO Failover URL on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	if d.Get("regenerate").(bool) {
		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, apiErr := client.UpdateFailoverUrl()
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}

			d.Set("failover_url", response.FailoverURL)
			d.Set("generation_time", response.GenerationTime)
			d.Set("regenerate", false)
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update SSO failover settings: %v", err))
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for deleting a Jamf Pro SSO Failover URL.
// SSO Failover URL cannot be deleted, only regenerated
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
