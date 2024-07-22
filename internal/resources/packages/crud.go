// packages_crud.go
package packages

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProPackagesCreate is responsible for creating a new Jamf Pro Package in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the package metadata in jamfpro.
// 3. Uploads the package file to the Jamf Pro server.
// 4. Sets the ID of the created package in the Terraform state.
// 5. Perform cleanup of downloaded package if it was from an HTTP(s) source
// 6. Reads the created package to ensure the Terraform state is up-to-date.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, localFilePath, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		var creationResponse *jamfpro.ResponsePackageCreatedAndUpdated

		creationResponse, apiErr = client.CreatePackage(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		log.Printf("[DEBUG] Jamf Pro Package Metadata created: %+v", creationResponse)

		fullFilePath := localFilePath

		_, apiErr = client.UploadPackage(creationResponse.ID, []string{fullFilePath})
		if apiErr != nil {
			log.Printf("[ERROR] Failed to upload package file for package '%s': %v", creationResponse.ID, apiErr)
			return retry.NonRetryableError(apiErr)
		}

		d.SetId(creationResponse.ID)

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create and upload Jamf Pro Package '%s' after retries: %v", resource.PackageName, err))
	}

	if strings.HasPrefix(d.Get("package_file_source").(string), "http") {
		err := os.Remove(localFilePath)
		if err != nil {
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

	return append(diags, updateTerraformState(d, response)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistReadWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// resourceJamfProPackagesUpdate is responsible for updating an existing Jamf Pro Package on the remote system.
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to update the package metadata in jamfpro.
// 3. Uploads the package file to the Jamf Pro server if the file has changed based on the MD5 hash.
// 4. Perform cleanup of downloaded package if it was from an HTTP(s) source
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, localFilePath, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdatePackageByID(resourceID, *resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Package '%s' (ID: %s) after retries: %v", resource.PackageName, resourceID, err))
	}

	// Use the local file path for generating the file hash
	newFileHash, err := generateMD5FileHash(localFilePath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to generate file hash for %s: %v", localFilePath, err))
	}

	oldFileHash, _ := d.Get("md5_file_hash").(string)

	log.Printf("[DEBUG] Comparing MD5 hashes for package update: oldFileHash=%s, newFileHash=%s", oldFileHash, newFileHash)

	if newFileHash != oldFileHash {
		_, err = client.UploadPackage(resourceID, []string{localFilePath})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to upload package file for package '%s': %v", resourceID, err))
		}

		// Update the filename and md5_file_hash in Terraform state to reflect the new file
		// this is done here while jamf JCDS hashes the file and updates the package metadata
		// to ensure that any runs during this window doesnt trigger another file upload.
		d.Set("md5_file_hash", newFileHash)
		d.Set("filename", filepath.Base(localFilePath))
	}

	if strings.HasPrefix(d.Get("package_file_source").(string), "http") {
		err := os.Remove(localFilePath)
		if err != nil {
			log.Printf("[WARN] Failed to remove downloaded package file '%s': %v", localFilePath, err)
		} else {
			log.Printf("[INFO] Successfully removed downloaded package file '%s'", localFilePath)
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// resourceJamfProPackagesDelete is responsible for deleting a Jamf Pro Package.
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
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Script '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
