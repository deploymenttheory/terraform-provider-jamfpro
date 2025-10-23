package icon

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for initializing the Jamf Pro Icon configuration in Terraform.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	filePath, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct icon file path: %v", err))
	}

	var uploadResponse *jamfpro.ResponseIconUpload
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		uploadResponse, apiErr = client.UploadIcon(filePath)
		if apiErr != nil {
			return retry.RetryableError(fmt.Errorf("failed to upload icon: %v", apiErr))
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Icon after retries: %v", err))
	}

	d.SetId(fmt.Sprintf("%d", uploadResponse.ID))

	// Only clean up if we downloaded from web source and verify the path is what we expect
	common.CleanupDownloadedIcon(d.Get("icon_file_web_source").(string), filePath)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of the Jamf Pro Icon configuration.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resourceID := d.Id()
	iconID, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse icon ID %s: %v", resourceID, err))
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
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

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

// update is responsible for updating the Jamf Pro Icon configuration.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Check if either file source has changed
	if d.HasChange("icon_file_path") || d.HasChange("icon_file_web_source") {
		filePath, err := construct(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to construct icon file path for update: %v", err))
		}

		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			uploadResponse, apiErr := client.UploadIcon(filePath)
			if apiErr != nil {
				return retry.RetryableError(fmt.Errorf("failed to upload icon: %v", apiErr))
			}

			d.SetId(fmt.Sprintf("%d", uploadResponse.ID))
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Icon after retries: %v", err))
		}

		// Only clean up if we downloaded from web source and verify the path is what we expect
		common.CleanupDownloadedIcon(d.Get("icon_file_web_source").(string), filePath)
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for 'deleting' the Jamf Pro Icon configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
