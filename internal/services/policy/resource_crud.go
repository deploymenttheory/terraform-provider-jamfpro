package policy

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	crud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2_crud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Lock the mutex to ensure only one create can run this function at a time
// var mu sync.Mutex

// create is responsible for creating a new Jamf Pro Policy in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	return crud.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreatePolicy,
		readNoCleanup,
	)
}

// Reads and states
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return crud.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetPolicyByID,
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

// update connstructs, updates and reads
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return crud.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdatePolicyByID,
		readNoCleanup,
	)
}

// Deletes and removes from state
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return crud.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeletePolicyByID,
	)
}
