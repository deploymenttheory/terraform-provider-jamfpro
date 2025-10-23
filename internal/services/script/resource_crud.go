package script

import (
	"context"
	"sync"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro Script in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	return common.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateScript,
		readNoCleanup,
	)
}

// read is responsible for reading the current state of a Jamf Pro Script Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetScriptByID,
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

// update is responsible for updating an existing Jamf Pro Script on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return common.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateScriptByID,
		readNoCleanup,
	)
}

var mu sync.Mutex

// delete is responsible for deleting a Jamf Pro Script.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	mu.Lock()
	defer mu.Unlock()
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteScriptByID,
	)
}
