package building

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create creates and states a jamfpro building
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// mu.Lock()
	// defer mu.Unlock()
	return sdkv2.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateBuilding,
		readNoCleanup,
	)
}

// read reads and states a jamfpro building
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return sdkv2.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetBuildingByID,
		updateState,
	)
}

// readWithCleanup reads a resources and states with cleanup
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads a resource without cleanup
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update updates a jamfpro building
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return sdkv2.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateBuildingByID,
		readNoCleanup,
	)
}

func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return sdkv2.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteBuildingByID,
	)
}
