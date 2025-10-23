// self_service_branding_ios_crud.go
package self_service_branding_ios

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	crud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2_crud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create creates a new self-service branding configuration using the shared common helper
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return crud.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateSelfServiceBrandingIOS,
		readNoCleanup,
	)
}

// read reads the resource state using the shared common helper
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return crud.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetSelfServiceBrandingIOSByID,
		updateState,
	)
}

func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, true)
}
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update updates the resource using the shared common helper
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return crud.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateSelfServiceBrandingIOSByID,
		readNoCleanup,
	)
}

// delete deletes the resource using the shared common helper
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return crud.Delete(ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteSelfServiceBrandingIOSByID,
	)
}
