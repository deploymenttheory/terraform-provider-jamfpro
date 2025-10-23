package file_share_distribution_point

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	JamfProResourceDistributionPoint = "Distribution Point"
)

// create is responsible for creating a new file share
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return common.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateDistributionPoint,
		readNoCleanup,
	)
}

// read is responsible for reading the current state of a file share distribution point
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetDistributionPointByID,
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

// update is responsible for updating an existing Jamf Pro File Share Distribution Point on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return common.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateDistributionPointByID,
		readNoCleanup,
	)
}

// delete is responsible for deleting a Jamf Pro file share distribution point from the remote system.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteDistributionPointByID,
	)
}
