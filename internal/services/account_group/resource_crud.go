package account_group

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro account group in the remote system.
// it follows a non standard pattern to allow for the client to be passed
// in as a parameter to the constructor to perform dynamic lookup for valid
// account group privileges.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d, meta)
	if err != nil {
		//nolint:err113
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Account Group for create: %w", err))
	}

	var createdRole *jamfpro.ResponseAccountGroupCreated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		createdRole, apiErr = client.CreateAccountGroup(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		//nolint:err113
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Account Group after retries: %w", err))
	}

	d.SetId(strconv.Itoa(createdRole.ID))

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of a Jamf Pro Account Group Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetAccountGroupByID,
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

// update is responsible for updating an existing Jamf Pro account group in the remote system.
// it follows a non standard pattern to allow for the client to be passed
// in as a parameter to the constructor to perform dynamic lookup for valid
// account group privileges.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := construct(d, meta)
	if err != nil {
		//nolint:err113
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Account Group for update: %w", err))
	}

	roleID := d.Id()

	var updatedRole *jamfpro.ResourceAccountGroup
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		updatedRole, apiErr = client.UpdateAccountGroupByID(
			roleID,
			resource,
		)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		//nolint:err113
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Account Group after retries: %w", err))
	}

	d.SetId(strconv.Itoa(updatedRole.ID))

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// delete is responsible for deleting a Jamf Pro account group.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteAccountGroupByID,
	)
}
