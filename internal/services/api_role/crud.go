package api_role

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create handles the creation of a Jamf Pro API Role.
// it follows a non standard pattern to allow for the client to be passed
// in as a parameter to the constructor to perform dynamic lookup for valid
// privileges.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro API Role for create: %v", err))
	}

	var createdRole *jamfpro.ResourceAPIRole
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		createdRole, apiErr = client.CreateJamfApiRole(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro API Role after retries: %v", err))
	}

	d.SetId(createdRole.ID)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read handles reading a Jamf Pro API Role from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return sdkv2.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetJamfApiRoleByID,
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

// update handles the update of a Jamf Pro API Role.
// it follows a non standard pattern to allow for the client to be passed
// in as a parameter to the constructor to perform dynamic lookup for valid
// privileges.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro API Role for update: %v", err))
	}

	roleID := d.Id()

	var updatedRole *jamfpro.ResourceAPIRole
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		updatedRole, apiErr = client.UpdateJamfApiRoleByID(
			roleID,
			resource,
		)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro API Role after retries: %v", err))
	}

	d.SetId(updatedRole.ID)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete handles the deletion of a Jamf Pro API Role.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return sdkv2.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteJamfApiRoleByID,
	)
}
