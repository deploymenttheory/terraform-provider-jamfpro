package allowed_file_extension

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Lock the mutex to ensure only one create can run this function at a time
// var mu sync.Mutex

// create is responsible for creating a new AllowedFileExtension in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	return common.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateAllowedFileExtension,
		readNoCleanup,
	)
}

// read is responsible for reading the current state of an Allowed File Extension Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetAllowedFileExtensionByID,
		updateState,
	)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update handles the update operation for an AllowedFileExtension resource in Terraform.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	deleteDiags := delete(ctx, d, meta)
	if deleteDiags.HasError() {
		return deleteDiags
	}

	return create(ctx, d, meta)
}

// delete is responsible for deleting an Allowed File Extension in Jamf Pro.
// This function will delete the resource based on its ID from the Terraform state.
// If the resource cannot be found by ID, it will attempt to delete by the 'extension' attribute.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteAllowedFileExtensionByID,
	)
}
