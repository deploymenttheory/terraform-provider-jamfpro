package buildings

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.CreateUpdate[jamfpro.ResourceBuilding, jamfpro.ResponseBuildingCreate](
		ctx,
		d,
		meta,
		constructJamfProBuilding,
		meta.(*jamfpro.Client).CreateBuilding,
		readNoCleanup,
	)
}

func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	return common.Read[jamfpro.ResourceBuilding](
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetBuildingByID,
		updateTerraformState,
	)
}

func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.CreateUpdate[jamfpro.ResourceBuilding, jamfpro.ResponseBuildingCreate](
		ctx,
		d,
		meta,
		constructJamfProBuilding,
		meta.(*jamfpro.Client).CreateBuilding,
		readNoCleanup,
	)
}

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteBuildingByID,
	)
}
