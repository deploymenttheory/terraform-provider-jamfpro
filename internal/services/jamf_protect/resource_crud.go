package jamf_protect

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create initializes the Jamf Protect integration by creating a new resource.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Protect integration: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.CreateJamfProtectIntegration(*resource.Registration, resource.Settings.AutoInstall)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Protect integration: %v", err))
	}

	d.SetId("jamfpro_jamf_protect_singleton")

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read retrieves the current state of the Jamf Protect integration.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	d.SetId("jamfpro_jamf_protect_singleton")

	var response *jamfpro.ResponseJamfProtectSettings
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetJamfProtectSettings()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateState(d, response)...)
}

// readWithCleanup is a wrapper for the read function that handles cleanup.
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup is a wrapper for the read function that does not handle cleanup.
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update modifies the existing Jamf Protect integration settings.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Protect settings: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateJamfProtectSettings(*resource.Settings)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Protect settings: %v", err))
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete removes the Jamf Protect integration.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteJamfProtectIntegration()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Protect integration: %v", err))
	}

	d.SetId("")
	return diags
}
