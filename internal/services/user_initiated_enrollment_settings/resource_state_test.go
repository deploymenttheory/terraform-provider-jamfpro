package user_initiated_enrollment_settings

import (
	"strconv"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUpdateStatePreservesConfiguredDisabledAdministrator(t *testing.T) {
	for _, parentEnabled := range []bool{true, false} {
		t.Run(strconv.FormatBool(parentEnabled), func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, ResourceJamfProUserInitatedEnrollmentSettings().Schema, map[string]any{
				"user_initiated_enrollment_for_computers": []any{map[string]any{
					"enable_user_initiated_enrollment_for_computers": parentEnabled,
					"ensure_ssh_is_enabled":                          false, "launch_self_service_when_done": parentEnabled,
					"account_driven_device_enrollment": false,
					"managed_local_administrator_account": []any{map[string]any{
						"create_managed_local_administrator_account":                    false,
						"management_account_username":                                   "",
						"hide_managed_local_administrator_account":                      true,
						"allow_ssh_access_for_managed_local_administrator_account_only": false,
					}},
				}},
			})
			prior := data.Get("user_initiated_enrollment_for_computers").(*schema.Set)
			diags := updateState(data, &jamfpro.ResourceEnrollment{
				MacOsEnterpriseEnrollmentEnabled: parentEnabled, LaunchSelfService: parentEnabled,
				CreateManagementAccount: false,
			}, nil, nil)
			require.False(t, diags.HasError(), "%v", diags)
			require.True(t, prior.Equal(data.Get("user_initiated_enrollment_for_computers")), "read removed configured disabled settings")
		})
	}
}

func TestUpdateStateDoesNotInventDisabledComputerSettings(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ResourceJamfProUserInitatedEnrollmentSettings().Schema, map[string]any{})
	diags := updateState(data, &jamfpro.ResourceEnrollment{}, nil, nil)
	require.False(t, diags.HasError(), "%v", diags)
	require.Zero(t, data.Get("user_initiated_enrollment_for_computers").(*schema.Set).Len())
}

func TestUpdateStateAdministratorEnablementRemainsAuthoritative(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(strconv.FormatBool(enabled), func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, ResourceJamfProUserInitatedEnrollmentSettings().Schema, map[string]any{
				"user_initiated_enrollment_for_computers": []any{map[string]any{
					"enable_user_initiated_enrollment_for_computers": true, "ensure_ssh_is_enabled": false, "launch_self_service_when_done": false,
					"managed_local_administrator_account": []any{map[string]any{
						"create_managed_local_administrator_account": true, "management_account_username": "configured",
						"hide_managed_local_administrator_account": false, "allow_ssh_access_for_managed_local_administrator_account_only": false,
					}},
				}},
			})
			diags := updateState(data, &jamfpro.ResourceEnrollment{
				MacOsEnterpriseEnrollmentEnabled: true, CreateManagementAccount: enabled, ManagementUsername: "from-api", HideManagementAccount: true,
			}, nil, nil)
			require.False(t, diags.HasError(), "%v", diags)
			computer := data.Get("user_initiated_enrollment_for_computers").(*schema.Set).List()[0].(map[string]any)
			admin := computer["managed_local_administrator_account"].(*schema.Set).List()[0].(map[string]any)
			require.Equal(t, enabled, admin["create_managed_local_administrator_account"])
			if enabled {
				require.Equal(t, "from-api", admin["management_account_username"])
				require.True(t, admin["hide_managed_local_administrator_account"].(bool))
			}
		})
	}
}
