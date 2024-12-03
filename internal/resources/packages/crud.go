// packages_crud.go
package packages

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro Package in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the package metadata in jamfpro.
// 3. Uploads the package file to the Jamf Pro server.
// 4. Verifes upload hash
// 5. Sets the ID of the created package in the Terraform state.
// 6. Perform cleanup of downloaded package if it was from an HTTP(s) source
// 7. Reads the created package to ensure the Terraform state is up-to-date.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, localFilePath, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package: %v", err))
	}

	// Set the filename in the resource based on the local file path
	resource.FileName = filepath.Base(localFilePath)

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		initialHash, err := jamfpro.CalculateSHA3_512(localFilePath)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to calculate SHA3-512: %v", err))
		}

		creationResponse, err := client.CreatePackage(*resource)
		if err != nil {
			return retry.RetryableError(fmt.Errorf("failed to create package metadata in Jamf Pro: %v", err))
		}

		log.Printf("[INFO] Jamf Pro package metadata created successfully with package ID: %s", creationResponse.ID)

		_, err = client.UploadPackage(creationResponse.ID, []string{localFilePath})
		if err != nil {
			log.Printf("[ERROR] Failed to upload package file '%s': %v", resource.FileName, err)
			return retry.NonRetryableError(fmt.Errorf("failed to upload package file: %v", err))
		}

		log.Printf("[INFO] Package %s file uploaded successfully", resource.FileName)

		// Wait for JCDS to process the upload
		time.Sleep(3 * time.Second)

		uploadedPackage, err := client.GetPackageByID(creationResponse.ID)
		if err != nil {
			return retry.RetryableError(fmt.Errorf("failed to verify uploaded package: %v", err))
		}

		if uploadedPackage.HashType != "SHA3_512" || uploadedPackage.HashValue != initialHash {
			return retry.NonRetryableError(fmt.Errorf("hash verification failed: initial=%s, uploaded=%s (type: %s)",
				initialHash, uploadedPackage.HashValue, uploadedPackage.HashType))
		}

		log.Printf("[INFO] Package %s SHA3-512 verification successful with hash: %s", resource.FileName, uploadedPackage.HashValue)
		d.SetId(creationResponse.ID)

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create and upload Jamf Pro Package '%s': %v", resource.PackageName, err))
	}

	// Clean up downloaded file if it was from an HTTP source
	if strings.HasPrefix(d.Get("package_file_source").(string), "http") {
		if err := os.Remove(localFilePath); err != nil {
			log.Printf("[WARN] Failed to remove downloaded package file '%s': %v", localFilePath, err)
		} else {
			log.Printf("[INFO] Successfully removed downloaded package file '%s'", localFilePath)
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
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

// update is responsible for updating an existing Jamf Pro Package on the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. calculates the SHA3-512 hash of the local file.
// 3. Calls the API to update the package metadata in jamfpro.
// 4. Uploads the package file to the Jamf Pro server.
// 5. Verifes upload hash
// 6. Reads the updated package to ensure the Terraform state is up-to-date.
// 7. Perform cleanup of downloaded package if it was from an HTTP(s) source
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, localFilePath, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package for update: %v", err))
	}

	newFileHash, err := jamfpro.CalculateSHA3_512(localFilePath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to calculate SHA3-512 hash for %s: %v", localFilePath, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdatePackageByID(resourceID, *resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		currentPackage, apiErr := client.GetPackageByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(fmt.Errorf("failed to get current package state: %v", apiErr))
		}

		if newFileHash != currentPackage.HashValue {
			_, apiErr = client.UploadPackage(resourceID, []string{localFilePath})
			if apiErr != nil {
				return retry.RetryableError(fmt.Errorf("failed to upload package file: %v", apiErr))
			}

			time.Sleep(5 * time.Second)

			uploadedPackage, apiErr := client.GetPackageByID(resourceID)
			if apiErr != nil {
				return retry.RetryableError(fmt.Errorf("failed to verify uploaded package: %v", apiErr))
			}

			if uploadedPackage.HashValue != newFileHash {
				return retry.NonRetryableError(fmt.Errorf("hash verification failed after upload: expected=%s, got=%s",
					newFileHash, uploadedPackage.HashValue))
			}

			log.Printf("[INFO] Package %s SHA3-512 verification successful with hash: %s",
				filepath.Base(localFilePath), uploadedPackage.HashValue)
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Package '%s' (ID: %s): %v", resource.PackageName, resourceID, err))
	}

	if strings.HasPrefix(d.Get("package_file_source").(string), "http") {
		if err := os.Remove(localFilePath); err != nil {
			log.Printf("[WARN] Failed to remove downloaded package file '%s': %v", localFilePath, err)
		} else {
			log.Printf("[INFO] Successfully removed downloaded package file '%s'", localFilePath)
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for deleting a Jamf Pro Package.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeletePackageByID(resourceID)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteScriptByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Package '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
