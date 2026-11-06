//go:build integration

package incidents_test

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	privileges "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/jamf_privileges"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestLivePrivilegeValidationMatchesCurl(t *testing.T) {
	capturePath := os.Getenv("JAMFPRO_CATALOG_CAPTURE")
	if capturePath == "" {
		t.Skip("set JAMFPRO_CATALOG_CAPTURE to a fresh curl catalog response and supply JAMFPRO_* OAuth variables")
	}
	body, err := os.ReadFile(capturePath)
	require.NoError(t, err)
	var captured jamfpro.ResourceApiRolePrivilegesList
	require.NoError(t, json.Unmarshal(body, &captured))
	p := provider.Provider()
	diags := p.Configure(context.Background(), terraform.NewResourceConfigRaw(map[string]any{
		"jamfpro_instance_fqdn": os.Getenv("JAMFPRO_INSTANCE_FQDN"),
		"auth_method":           "oauth2", "client_id": os.Getenv("JAMFPRO_CLIENT_ID"), "client_secret": os.Getenv("JAMFPRO_CLIENT_SECRET"),
		"hide_sensitive_data": true, "enable_client_sdk_logs": false, "jamfpro_load_balancer_lock": false,
	}))
	require.False(t, diags.HasError(), "%v", diags)
	client, ok := p.Meta().(*jamfpro.Client)
	require.True(t, ok)
	current, err := client.GetJamfAPIPrivileges()
	require.NoError(t, err)
	slices.Sort(captured.Privileges)
	slices.Sort(current.Privileges)
	require.Equal(t, captured.Privileges, current.Privileges, "SDK lookup differs from the fresh curl response")
	require.NoError(t, privileges.ValidateAccountPrivileges(client, jamfpro.AccountSubsetPrivileges{
		JSSObjects: []string{"Read Computers"}, JSSSettings: []string{"Read PKI", "Update PKI", "Read User-Initiated Enrollment"},
	}))
	require.ErrorContains(t, privileges.ValidateAccountPrivileges(client, jamfpro.AccountSubsetPrivileges{
		JSSSettings: []string{"Read TF_TEST_NONEXISTENT"},
	}), "Read TF_TEST_NONEXISTENT")
}
