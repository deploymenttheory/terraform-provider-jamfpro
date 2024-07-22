package apiroles

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create handles the creation of a Jamf Pro API Role.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateJamfApiRole,
		readNoCleanup,
	)
}

// read handles reading a Jamf Pro API Role from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetJamfApiRoleByID,
		updateTerraformState,
	)
}

// resourceJamfProAPIRolesReadWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// resourceJamfProAPIRolesReadNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// resourceJamfProAPIRolesUpdate handles updating a Jamf Pro API Role.
// The function:
// 1. Constructs the updated API role data using the provided Terraform configuration.
// 2. Calls the API to update the role in Jamf Pro.
// 3. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateJamfApiRoleByID,
		readNoCleanup,
	)
}

// resourceJamfProAPIRolesDelete handles the deletion of a Jamf Pro API Role.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteJamfApiRoleByID,
	)
}
