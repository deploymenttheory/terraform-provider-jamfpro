// smartcomputergroup_crud.go
package smart_computer_group

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro Smart Computer Group in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	return sdkv2.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateComputerGroup,
		readNoCleanup,
	)
}

// read is responsible for reading the current state of a Jamf Pro Smart Computer Group from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return sdkv2.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetComputerGroupByID,
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

// update is responsible for updating an existing Jamf Pro Smart Computer Group on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return sdkv2.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateComputerGroupByID,
		readNoCleanup,
	)
}

// delete is responsible for deleting a Jamf Pro Smart Computer Group.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return sdkv2.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteComputerGroupByID,
	)
}
