package sso_certificate

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro SSO Certificate in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.CreateSSOCertificate()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create SSO certificate: %v", err))
	}

	d.SetId("jamfpro_sso_certificate_singleton")

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of a Jamf Pro SSO Certificate from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	d.SetId("jamfpro_sso_certificate_singleton")

	var response *jamfpro.ResourceSSOKeystoreResponse
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetSSOCertificate()
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

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// delete is responsible for deleting a Jamf Pro SSO Certificate.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteSSOCertificate()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to delete SSO certificate: %v", err))...)
	}

	d.SetId("")
	return diags
}
