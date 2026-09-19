package user_group

import (
	"context"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOperationReadAndReapply(t *testing.T) {
	r := ResourceJamfProUserGroups()
	input := map[string]any{}
	for k, v := range setFixture() {
		if k != "assigned_user_ids" {
			input[k] = v
		}
	}
	data := schema.TestResourceDataRaw(t, r.Schema, input)
	data.SetId("42")
	payload, err := construct(data)
	require.NoError(t, err)
	require.Len(t, payload.UserAdditions, 2)
	require.Len(t, payload.UserDeletions, 2)
	response := &jamfpro.ResourceUserGroup{ID: 42, Name: "set-test", Users: []jamfpro.UserGroupSubsetUserItem{{ID: 4}, {ID: 5}}}
	response.Site = sharedschemas.ConstructSharedResourceSite(-1)
	require.Empty(t, updateState(data, response))
	require.Equal(t, 2, data.Get("user_additions").(*schema.Set).Len())
	require.Equal(t, 2, data.Get("user_deletions").(*schema.Set).Len())
	require.Equal(t, 2, data.Get("assigned_user_ids").(*schema.Set).Len())
	reverseCollections(input)
	diff, err := r.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(input), nil)
	require.NoError(t, err)
	require.True(t, diff == nil || diff.Empty())
	// A subsequent metadata-only update does not replay historical commands.
	restored := r.Data(data.State())
	payload, err = construct(restored)
	require.NoError(t, err)
	require.Empty(t, payload.UserAdditions)
	require.Empty(t, payload.UserDeletions)
}

func TestOperationIdentityExcludesComputedFields(t *testing.T) {
	r := ResourceJamfProUserGroups()
	input := setFixture()
	data := schema.TestResourceDataRaw(t, r.Schema, input)
	before := data.Get("user_additions").(*schema.Set)
	values := before.List()
	values[0].(map[string]any)["full_name"] = "Resolved name"
	values[0].(map[string]any)["email_address"] = "test@example.invalid"
	require.NoError(t, data.Set("user_additions", values))
	require.True(t, before.Equal(data.Get("user_additions").(*schema.Set)))
	users, err := extractUsers([]any{map[string]any{"id": "12", "username": "", "full_name": "", "phone_number": "", "email_address": ""}})
	require.NoError(t, err)
	require.Equal(t, 12, users[0].ID)
	_, err = extractUsers([]any{map[string]any{"id": "invalid"}})
	require.Error(t, err)
}

func TestOperationUpdatesPreserveMembershipAndFilterAlreadyAppliedCommands(t *testing.T) {
	payload := &jamfpro.ResourceUserGroup{
		UserAdditions: []jamfpro.UserGroupSubsetUserItem{{ID: 1}, {Username: "bravo"}, {ID: 3}},
		UserDeletions: []jamfpro.UserGroupSubsetUserItem{{ID: 1}, {ID: 9}},
	}
	members := []jamfpro.UserGroupSubsetUserItem{{ID: 1}, {ID: 2, Username: "bravo"}}
	preserveOperationMembership(payload, members)
	require.Equal(t, members, payload.Users)
	require.Equal(t, []jamfpro.UserGroupSubsetUserItem{{ID: 3}}, payload.UserAdditions)
	require.Equal(t, []jamfpro.UserGroupSubsetUserItem{{ID: 1}}, payload.UserDeletions)
	preserveOperationMembership(payload, []jamfpro.UserGroupSubsetUserItem{{ID: 2}, {ID: 3}})
	require.Empty(t, payload.UserAdditions)
	require.Empty(t, payload.UserDeletions)
}

func TestFailedOperationDoesNotBecomeSubmittedState(t *testing.T) {
	r := ResourceJamfProUserGroups()
	input := map[string]any{"name": "operations", "is_smart": false, "user_additions": []any{map[string]any{"id": "1"}}}
	data := schema.TestResourceDataRaw(t, r.Schema, input)
	data.SetId("42")
	prior := data.State()
	input["user_additions"] = []any{map[string]any{"id": "invalid"}}
	diff, err := r.Diff(context.Background(), prior, terraform.NewResourceConfigRaw(input), nil)
	require.NoError(t, err)
	state, diags := r.Apply(context.Background(), prior, diff, nil)
	require.True(t, diags.HasError())
	require.NotNil(t, state)
	require.Equal(t, prior.Attributes, state.Attributes)
}
