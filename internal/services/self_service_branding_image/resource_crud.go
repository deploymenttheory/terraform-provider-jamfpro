// self_service_branding_image_crud.go
package self_service_branding_image

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrFailedConstruct       = fmt.Errorf("failed to construct self service branding image file path")
	ErrUploadFailed          = fmt.Errorf("failed to upload self service branding image")
	ErrCreateAfterRetries    = fmt.Errorf("failed to create Jamf Pro Self Service branding image after retries")
	ErrFailedConstructUpdate = fmt.Errorf("failed to construct self service branding image file path for update")
	ErrUpdateAfterRetries    = fmt.Errorf("failed to update Jamf Pro Self Service branding image after retries")
)

// create is responsible for initializing the Jamf Pro Self Service branding image configuration in Terraform.
// The API returns a URL for the uploaded branding image. Use the final path
// segment as the Terraform resource ID (the schema documents the ID is derived
// from the URL).
// Extract the final path segment from the returned URL (e.g. .../download/5 -> 5)
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	filePath, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrFailedConstruct.Error(), err))
	}

	var uploadResponse *jamfpro.ResponseSelfServiceBrandingImage
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		uploadResponse, apiErr = client.UploadSelfServiceBrandingImage(filePath)
		if apiErr != nil {
			return retry.RetryableError(fmt.Errorf("%s: %w", ErrUploadFailed.Error(), apiErr))
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrCreateAfterRetries.Error(), err))
	}

	u, err := url.Parse(uploadResponse.URL)
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid URL: %w", err))
	}
	seg := path.Base(u.Path)
	if unescaped, err := url.PathUnescape(seg); err == nil {
		seg = unescaped
	}

	d.SetId(seg)
	common.CleanupDownloadedIcon(d.Get("self_service_branding_image_file_web_source").(string), filePath)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro Self Service branding image configuration.
// This resource has no GET endpoint in the Jamf Pro API. It is created/updated by
// uploading an image and removed from state by Terraform delete. Therefore the
// Read will simply return the current state unless the ID is missing.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceID := d.Id()
	if resourceID == "" {
		return diags
	}

	return diags
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// delete is responsible for 'deleting' the Jamf Pro Self Service branding image configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
