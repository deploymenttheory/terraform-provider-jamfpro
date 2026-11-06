package jamf_privileges

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/stretchr/testify/require"
)

var errTestCatalogAccess = errors.New("catalog access denied")

type fakePrivilegeCatalog struct {
	response *jamfpro.ResourceApiRolePrivilegesList
	err      error
}

func (f fakePrivilegeCatalog) GetJamfAPIPrivileges() (*jamfpro.ResourceApiRolePrivilegesList, error) {
	return f.response, f.err
}

func capturedCatalog(t *testing.T) *jamfpro.ResourceApiRolePrivilegesList {
	t.Helper()
	body, err := os.ReadFile("testdata/api_role_privileges_11_33_3.json")
	require.NoError(t, err)
	var catalog jamfpro.ResourceApiRolePrivilegesList
	require.NoError(t, json.Unmarshal(body, &catalog))
	return &catalog
}

func TestValidateAccountPrivilegesUsesLiveCatalogNames(t *testing.T) {
	// These user-facing values were present in the live catalog and matched
	// the corresponding search and Classic account/group responses.
	err := ValidateAccountPrivileges(fakePrivilegeCatalog{response: capturedCatalog(t)}, jamfpro.AccountSubsetPrivileges{
		JSSObjects:  []string{"Read Computers"},
		JSSSettings: []string{"Read PKI", "Update PKI", "Read User-Initiated Enrollment", "Update User-Initiated Enrollment"},
	})
	require.NoError(t, err)
}

func TestValidateAccountPrivilegesRejectsNamesAbsentFromLiveCatalog(t *testing.T) {
	for _, name := range []string{"PKI", "Update Allowed File Extension", "Read TF_TEST_NONEXISTENT"} {
		t.Run(name, func(t *testing.T) {
			err := ValidateAccountPrivileges(fakePrivilegeCatalog{response: capturedCatalog(t)}, jamfpro.AccountSubsetPrivileges{JSSSettings: []string{name}})
			require.ErrorContains(t, err, name)
			require.ErrorContains(t, err, "JSS Settings")
		})
	}
}

func TestValidateAccountPrivilegesCatalogFailure(t *testing.T) {
	err := ValidateAccountPrivileges(fakePrivilegeCatalog{err: errTestCatalogAccess}, jamfpro.AccountSubsetPrivileges{})
	require.ErrorIs(t, err, errTestCatalogAccess)
	for _, response := range []*jamfpro.ResourceApiRolePrivilegesList{nil, {}} {
		err := ValidateAccountPrivileges(fakePrivilegeCatalog{response: response}, jamfpro.AccountSubsetPrivileges{})
		require.ErrorIs(t, err, errEmptyPrivilegeCatalog)
	}
}
