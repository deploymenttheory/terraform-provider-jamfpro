// packages_crud.go
package packages

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProPackagesCreate is responsible for creating a new Jamf Pro Package in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProPackagesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Construct the package resource from the Terraform schema
	resource, err := constructJamfProPackageCreate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package: %v", err))
	}

	// Step 1: Create the package in Jamf Pro
	var creationResponse *jamfpro.ResponsePackageCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreatePackage(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Package '%s' after retries: %v", resource.PackageName, err))
	}

	// Step 2: Upload the package file
	filePath := d.Get("package_file_path").(string)
	_, err = client.UploadPackage(creationResponse.ID, []string{filePath})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to upload package file for package '%s': %v", creationResponse.ID, err))
	}

	// Set the ID in the Terraform state
	d.SetId(creationResponse.ID)

	return append(diags, ResourceJamfProPackagesRead(ctx, d, meta)...)
}

// ResourceJamfProPackagesRead is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProPackagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return append(diags, diag.FromErr(err)...)
	}

	return append(diags, updateTerraformState(d, response)...)
}

// ResourceJamfProPackagesUpdate is responsible for updating an existing Jamf Pro Package on the remote system.
func ResourceJamfProPackagesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Check if package_file_path has changed
	if d.HasChange("package_file_path") {
		// Step 1: Calculate the new file hash
		filePath := d.Get("package_file_path").(string)
		newFileHash, err := generateMD5FileHash(filePath)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to generate file hash for %s: %v", filePath, err))
		}

		// Step 2: Compare the new file hash with the old one
		oldFileHash, _ := d.GetChange("md5_file_hash")
		if newFileHash != oldFileHash.(string) {
			// The file has changed, upload it
			_, err = client.UploadPackage(resourceID, []string{filePath})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to upload package file for package '%s': %v", resourceID, err))
			}

			// Update the package_uri and md5_file_hash in Terraform state
			d.Set("package_uri", filePath)
			d.Set("md5_file_hash", newFileHash)
			d.Set("filename", filepath.Base(filePath))
		}
	}

	// Construct the updated package resource from the Terraform schema
	resource, err := constructJamfProPackageCreate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package for update: %v", err))
	}

	// Update the package metadata in Jamf Pro
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdatePackageManifestByID(resourceID, *resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Package '%s' (ID: %s) after retries: %v", resource.PackageName, resourceID, err))
	}

	return append(diags, ResourceJamfProPackagesRead(ctx, d, meta)...)
}

// ResourceJamfProPackagesDelete is responsible for deleting a Jamf Pro Package.
func ResourceJamfProPackagesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
