package icon

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	commonerrors "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/errors"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/files"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrConstructIconPath = errors.New("failed to construct icon file path")
	ErrUploadIcon        = errors.New("failed to upload icon")
	ErrCreateIcon        = errors.New("failed to create Jamf Pro Icon after retries")
	ErrUpdateIcon        = errors.New("failed to update Jamf Pro Icon after retries")
	ErrParseIconID       = errors.New("failed to parse icon ID")
)

// create is responsible for initializing the Jamf Pro Icon configuration in Terraform.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	filePath, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %w", ErrConstructIconPath, err))
	}

	var uploadResponse *jamfpro.ResponseIconUpload
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		uploadResponse, apiErr = client.UploadIcon(filePath)
		if apiErr != nil {
			return retry.RetryableError(fmt.Errorf("%w: %w", ErrUploadIcon, apiErr))
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %w", ErrCreateIcon, err))
	}

	d.SetId(fmt.Sprintf("%d", uploadResponse.ID))

	files.CleanupDownloadedIcon(d.Get("icon_file_web_source").(string), filePath)
	if d.Get("icon_file_base64").(string) != "" {
		files.CleanupDownloadedIcon("base64", filePath)
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro Icon configuration.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resourceID := d.Id()
	iconID, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w %s: %w", ErrParseIconID, resourceID, err))
	}

	var response *jamfpro.ResponseIconUpload
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetIconByID(iconID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, commonerrors.HandleResourceNotFoundError(err, d, cleanup)...)
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

// update is responsible for updating the Jamf Pro Icon configuration.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	if d.HasChange("icon_file_path") || d.HasChange("icon_file_web_source") || d.HasChange("icon_file_base64") {
		filePath, err := construct(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %w", ErrConstructIconPath, err))
		}

		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			uploadResponse, apiErr := client.UploadIcon(filePath)
			if apiErr != nil {
				return retry.RetryableError(fmt.Errorf("%w: %w", ErrUploadIcon, apiErr))
			}

			d.SetId(fmt.Sprintf("%d", uploadResponse.ID))
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %w", ErrUpdateIcon, err))
		}

		files.CleanupDownloadedIcon(d.Get("icon_file_web_source").(string), filePath)
		if d.Get("icon_file_base64").(string) != "" {
			files.CleanupDownloadedIcon("base64", filePath)
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro Icon configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")

	return nil
}
