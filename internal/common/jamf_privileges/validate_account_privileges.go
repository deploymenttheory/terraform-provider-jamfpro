package jamf_privileges

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

var errEmptyPrivilegeCatalog = errors.New("cannot validate privileges: the Jamf Pro privilege catalog is empty")

type invalidPrivInfo struct {
	privileges  []string
	suggestions map[string]string
}

type privilegeCatalog interface {
	GetJamfAPIPrivileges() (*jamfpro.ResourceApiRolePrivilegesList, error)
}

// ValidateAccountPrivileges validates names against the connected server's privilege
// catalog. An administrator's assigned privileges are not a complete catalog.
func ValidateAccountPrivileges(client privilegeCatalog, privileges jamfpro.AccountSubsetPrivileges) error {
	catalog, err := client.GetJamfAPIPrivileges()
	if err != nil {
		return fmt.Errorf("failed to fetch the Jamf Pro privilege catalog (requires Read API Roles): %w", err)
	}
	if catalog == nil || len(catalog.Privileges) == 0 {
		return errEmptyPrivilegeCatalog
	}

	invalidPrivileges := make(map[string]invalidPrivInfo)

	validatePrivilegeSet := func(supplied []string, reference []string, category string) {
		for _, priv := range supplied {
			found := false
			for _, refPriv := range reference {
				if priv == refPriv {
					found = true
					break
				}
			}
			if !found {
				info, exists := invalidPrivileges[category]
				if !exists {
					info = invalidPrivInfo{
						privileges:  make([]string, 0),
						suggestions: make(map[string]string),
					}
				}
				info.privileges = append(info.privileges, priv)

				suggestions := FindSimilarPrivileges(priv, reference)
				if len(suggestions) > 0 {
					info.suggestions[priv] = suggestions[0]
				}

				invalidPrivileges[category] = info
			}
		}
	}

	validatePrivilegeSet(privileges.JSSObjects, catalog.Privileges, "JSS Objects")
	validatePrivilegeSet(privileges.JSSSettings, catalog.Privileges, "JSS Settings")
	validatePrivilegeSet(privileges.JSSActions, catalog.Privileges, "JSS Actions")
	validatePrivilegeSet(privileges.Recon, catalog.Privileges, "Recon")
	validatePrivilegeSet(privileges.CasperAdmin, catalog.Privileges, "Casper Admin")
	validatePrivilegeSet(privileges.CasperRemote, catalog.Privileges, "Casper Remote")
	validatePrivilegeSet(privileges.CasperImaging, catalog.Privileges, "Casper Imaging")

	if len(invalidPrivileges) > 0 {
		var msg strings.Builder
		msg.WriteString("Invalid privilege names absent from this Jamf Pro server's privilege catalog:\n")

		for category, info := range invalidPrivileges {
			msg.WriteString(fmt.Sprintf("\nInvalid %s privileges:\n", category))
			for _, priv := range info.privileges {
				msg.WriteString(fmt.Sprintf("- %s\n", priv))
				if suggestion, ok := info.suggestions[priv]; ok {
					msg.WriteString(fmt.Sprintf("  Did you mean: %s?\n", suggestion))
				}
			}
		}
		return fmt.Errorf("%s", msg.String())
	}

	return nil
}
