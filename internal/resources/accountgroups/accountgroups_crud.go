package accountgroups

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProAccountGroupCreate is responsible for creating a new Jamf Pro Script in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Create[jamfpro.ResourceAccountGroup, jamfpro.ResponseAccountGroupCreated](
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateAccountGroup,
		readNoCleanup,
	)
}

// resourceJamfProAccountGroupRead is responsible for reading the current state of a Jamf Pro Account Group Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	return common.Read[jamfpro.ResourceAccountGroup](
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetAccountGroupByID,
		updateTerraformState,
	)
}

// resourceJamfProAccountGroupReadWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// resourceJamfProAccountGroupReadNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// resourceJamfProAccountGroupUpdate is responsible for updating an existing Jamf Pro Account Group on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Update[jamfpro.ResourceAccountGroup, jamfpro.ResourceAccountGroup](
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateAccountGroupByID,
		readNoCleanup,
	)
}

// resourceJamfProAccountGroupDelete is responsible for deleting a Jamf Pro account group.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteAccountGroupByID,
	)
}
