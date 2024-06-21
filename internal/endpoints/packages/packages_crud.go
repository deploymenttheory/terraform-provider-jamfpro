package packages

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProPackagesRead is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProPackagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetPackageByID(resourceIDInt)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProPackagesUpdate is responsible for updating an existing Jamf Pro Package on the remote system.
func ResourceJamfProPackagesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	packageID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting package ID '%s' to integer: %v", d.Id(), err))
	}

	if d.HasChange("package_file_path") {
		filePath := d.Get("package_file_path").(string)
		newFileHash, err := generateMD5FileHash(filePath)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to generate file hash for %s: %v", filePath, err))
		}

		oldFileHash, _ := d.GetChange("md5_file_hash")
		if newFileHash != oldFileHash.(string) {
			fileUploadResponse, err := client.CreateJCDS2PackageV2(filePath)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to upload file to JCDS 2.0 with file path '%s': %v", filePath, err))
			}

			d.Set("package_uri", fileUploadResponse.URI)
			d.Set("md5_file_hash", newFileHash)
			d.Set("filename", filepath.Base(filePath))
		}
	}

	packageResource, err := constructJamfProPackageCreate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct package for update: %v", err))
	}

	_, err = client.UpdatePackageByID(packageID, packageResource)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update package with ID %d: %v", packageID, err))
	}

	return append(diags, ResourceJamfProPackagesRead(ctx, d, meta)...)
}

// ResourceJamfProPackagesDelete is responsible for deleting a Jamf Pro Package.
func ResourceJamfProPackagesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeletePackageByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeletePackageByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Package '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
