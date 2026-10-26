package user_group

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	crud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2_crud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro User Group in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	return crud.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateUserGroup,
		readNoCleanup,
	)
}

// read is responsible for reading the current state of a Jamf Pro User Group Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	return crud.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetUserGroupByID,
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

// update is responsible for updating an existing Jamf Pro User Group on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Retain prior command state if a request fails; otherwise the SDK may
	// record an unexecuted operation as already submitted.
	d.Partial(true)
	diags := crud.Update(
		ctx,
		d,
		meta,
		construct,
		func(id string, payload *jamfpro.ResourceUserGroup) (*jamfpro.ResponseUserGroupCreateAndUpdate, error) {
			client := meta.(*jamfpro.Client)
			config := d.GetRawConfig()
			if !config.IsNull() && config.GetAttr("assigned_user_ids").IsNull() {
				current, err := client.GetUserGroupByID(id)
				if err != nil {
					return nil, err
				}
				// The SDK always serializes a users element. Preserve current membership
				// when applying commands so an empty element cannot clear the group first.
				preserveOperationMembership(payload, current.Users)
			}
			return client.UpdateUserGroupByID(id, payload)
		},
		readNoCleanup,
	)
	if !diags.HasError() {
		d.Partial(false)
	}
	return diags
}

// delete is responsible for deleting a Jamf Pro User Group.
func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return crud.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteUserGroupByID,
	)
}

func preserveOperationMembership(payload *jamfpro.ResourceUserGroup, members []jamfpro.UserGroupSubsetUserItem) {
	payload.Users = members
	contains := func(user jamfpro.UserGroupSubsetUserItem) bool {
		for _, member := range members {
			if user.ID != 0 && user.ID == member.ID || user.ID == 0 && user.Username == member.Username {
				return true
			}
		}
		return false
	}
	additions := make([]jamfpro.UserGroupSubsetUserItem, 0, len(payload.UserAdditions))
	for _, user := range payload.UserAdditions {
		if !contains(user) {
			additions = append(additions, user)
		}
	}
	deletions := make([]jamfpro.UserGroupSubsetUserItem, 0, len(payload.UserDeletions))
	for _, user := range payload.UserDeletions {
		if contains(user) {
			deletions = append(deletions, user)
		}
	}
	payload.UserAdditions, payload.UserDeletions = additions, deletions
}
