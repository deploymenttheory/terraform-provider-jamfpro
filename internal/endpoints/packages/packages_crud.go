// packages_crud.go
package packages

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProPackagesCreate is responsible for creating a new Jamf Pro Package in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
// ResourceJamfProPackagesCreate is responsible for creating a new Jamf Pro Package in the remote system.
func resourceJamfProPackagesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Construct the package resource from the Terraform schema
	resource, err := constructJamfProPackageCreate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		var creationResponse *jamfpro.ResponsePackageCreatedAndUpdated

		// Step 1: Create the package metadata in Jamf Pro
		creationResponse, apiErr = client.CreatePackage(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		log.Printf("[DEBUG] Jamf Pro Package Metadata created: %+v", creationResponse)

		// Step 2: Upload the package file within the same context
		filePath := d.Get("package_file_path").(string)
		fullFilePath, _ := filepath.Abs(filePath)

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

	return append(diags, resourceJamfProPackagesReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProPackagesRead is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProPackagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
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
func resourceJamfProPackagesReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProPackagesRead(ctx, d, meta, true)
}

// resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProPackagesReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProPackagesRead(ctx, d, meta, false)
}

// resourceJamfProPackagesUpdate is responsible for updating an existing Jamf Pro Package on the remote system.
func resourceJamfProPackagesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProPackageCreate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package for update: %v", err))
	}

	// Step 1: Update the package metadata in Jamf Pro
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

	// Step 2: Upload the package file if it has changed

	filePath := d.Get("package_file_path").(string)
	newFileHash, err := generateMD5FileHash(filePath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to generate file hash for %s: %v", filePath, err))
	}

	oldFileHash, _ := d.Get("md5_file_hash").(string)

	log.Printf("[DEBUG] Comparing MD5 hashes for package update: oldFileHash=%s, newFileHash=%s", oldFileHash, newFileHash)

	if newFileHash != oldFileHash {
		_, err = client.UploadPackage(resourceID, []string{filePath})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to upload package file for package '%s': %v", resourceID, err))
		}

		// Update the filename and md5_file_hash in Terraform state to reflect the new file
		// this is done here while jamf JCDS hashes the file and updates the package metadata
		// to ensure that any runs during this window doesnt trigger another file upload.
		d.Set("md5_file_hash", newFileHash)
		d.Set("filename", filepath.Base(filePath))
	}

	return append(diags, resourceJamfProPackagesReadWithCleanup(ctx, d, meta)...)
}

// resourceJamfProPackagesDelete is responsible for deleting a Jamf Pro Package.
func resourceJamfProPackagesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
