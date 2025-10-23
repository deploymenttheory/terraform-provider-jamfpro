// packages_crud.go
package packages

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/errors"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/files"
	crud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2_crud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// A separate, shorter timeout for the metadata call which does not have a large payload
const PackagesMetaTimeout time.Duration = 10 * time.Minute

// create handles the creation of a Jamf Pro package resource:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calculates initial SHA3-512 hash of the package file.
// 3. Calculates the MD5 hash of the package file.
// 4. Calls the API to create the package metadata in Jamf Pro.
// 5. Uploads the package file to the Jamf Pro server.
// 6. Verifies the uploaded package hash matches the initial hash.
// 7. If verification fails, deletes the package from Jamf Pro.
// 8. If verification succeeds, sets the package ID in Terraform state.
// 8. Performs cleanup of downloaded package if it was from an HTTP(s) source.
// 9. Reads the created package to ensure the Terraform state is up-to-date.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	var packageID string

	resource, localFilePath, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package: %v", err))
	}

	resource.FileName = filepath.Base(localFilePath)

	initialHash, err := jamfpro.CalculateSHA3_512(localFilePath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to calculate SHA3-512: %v", err))
	}

	md5Hash, err := jamfpro.CalculateMD5(localFilePath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to calculate MD5: %v", err))
	}
	resource.MD5 = md5Hash

	// Meta
	err = retry.RetryContext(ctx, PackagesMetaTimeout, func() *retry.RetryError {
		creationResponse, err := client.CreatePackage(*resource)

		if err != nil {
			return retry.RetryableError(fmt.Errorf("failed to create package metadata in Jamf Pro: %v", err))
		}

		packageID = creationResponse.ID

		log.Printf("[INFO] Jamf Pro package metadata created successfully with package ID: %s", packageID)

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to make the metadata, exiting: %v", err))
	}

	// Package
	client.HTTP.ModifyHttpTimeout(d.Timeout(schema.TimeoutCreate))
	defer client.HTTP.ResetTimeout()

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, err = client.UploadPackage(packageID, []string{localFilePath})

		if err != nil {
			log.Printf("[ERROR] Failed to upload package file '%s': %v", resource.FileName, err)
			return retry.NonRetryableError(fmt.Errorf("failed to upload package file: %v", err))
		}

		log.Printf("[INFO] Package %s file uploaded successfully", resource.FileName)

		return nil
	})

	if err != nil {
		// Cleans up the metadata so the next run doesn't hit an error trying to remake it, duplicate names are not allowed
		cleanupErr := client.DeletePackageByID(packageID)

		if cleanupErr != nil {
			return diag.FromErr(fmt.Errorf("failed to upload package: %v and failed to delete metadata: %v", err, cleanupErr))
		}
		return diag.FromErr(fmt.Errorf("failed to upload Jamf Pro Package '%s': %v", resource.PackageName, err))
	}

	if err := verifyPackageUpload(ctx, client, packageID, resource.FileName, initialHash,
		d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to verify Jamf Pro Package '%s': %v", resource.PackageName, err))
	}

	d.SetId(packageID)

	files.CleanupDownloadedPackage(d.Get("package_file_source").(string), localFilePath)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	var response *jamfpro.ResourcePackage

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetPackageByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, errors.HandleResourceNotFoundError(err, d, cleanup)...)
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

// update handles the updating of a Jamf Pro package resource:
//  1. Constructs the updated attribute data from the Terraform configuration.
//  2. Calculates SHA3-512 hash of the new package file.
//  3. Calculates the MD5 hash of the package file.
//  4. Updates the package metadata in Jamf Pro.
//  5. If the file hash differs from current:
//     a. Uploads the new package file.
//     b. Verifies the uploaded package hash matches.
//     c. If verification fails, attempts to revert metadata changes.
//  6. Performs cleanup of downloaded package if it was from an HTTP(s) source.
//  7. Reads the updated package to ensure the Terraform state is up-to-date.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Check if this is a file-related update or metadata-only update
	fileChanged := d.HasChange("package_file_source")

	resource, localFilePath, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package for update: %v", err))
	}

	// Handle metadata update
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, err := client.UpdatePackageByID(resourceID, *resource)
		if err != nil {
			return retry.RetryableError(fmt.Errorf("failed to update package metadata: %v", err))
		}
		log.Printf("[INFO] Package metadata updated successfully")
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Package metadata '%s' (ID: %s): %v",
			resource.PackageName, resourceID, err))
	}

	// Only handle file operations if the package file has changed
	if fileChanged {
		newFileHash, err := jamfpro.CalculateSHA3_512(localFilePath)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to calculate SHA3-512 hash for %s: %v", localFilePath, err))
		}

		md5Hash, err := jamfpro.CalculateMD5(localFilePath)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to calculate MD5: %v", err))
		}
		resource.MD5 = md5Hash

		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			_, err := client.UploadPackage(resourceID, []string{localFilePath})
			if err != nil {
				return retry.RetryableError(fmt.Errorf("failed to upload package file: %v", err))
			}

			log.Printf("[INFO] Package file uploaded successfully")
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to upload new package file: %v", err))
		}

		if err := verifyPackageUpload(ctx, client, resourceID, resource.FileName, newFileHash,
			d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(fmt.Errorf("failed to verify updated package file: %v", err))
		}

		files.CleanupDownloadedPackage(d.Get("package_file_source").(string), localFilePath)
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for deleting a Jamf Pro Package.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return crud.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeletePackageByID,
	)
}
