package user_group

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro User Group in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return common.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateUserGroup,
		readNoCleanup,
	)
}

// read is responsible for reading the current state of a Jamf Pro User Group Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetUserGroupByID,
		updateState,
	)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating an existing Jamf Pro User Group on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateUserGroupByID,
		readNoCleanup,
	)
}

// delete is responsible for deleting a Jamf Pro User Group.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteUserGroupByID,
	)
}
