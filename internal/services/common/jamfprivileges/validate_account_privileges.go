package jamfprivileges

import (
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

type invalidPrivInfo struct {
	privileges  []string
	suggestions map[string]string
}

// ValidateAccountPrivileges validates the privileges of an account against a lookup of all available
// privileges based upon on the first found jamf pro accout with Administrator privilege set.
func ValidateAccountPrivileges(client *jamfpro.Client, privileges jamfpro.AccountSubsetPrivileges) error {
	versionInfo, err := client.GetJamfProVersion()
	if err != nil {
		return fmt.Errorf("failed to fetch Jamf Pro version: %v", err)
	}

	accountsList, err := client.GetAccounts()
	if err != nil {
		return fmt.Errorf("failed to fetch accounts for validation: %v", err)
	}

	var adminAccount *jamfpro.ResourceAccount
	for _, user := range accountsList.Users {
		account, err := client.GetAccountByID(fmt.Sprint(user.ID))
		if err != nil {
			continue
		}
		if account.PrivilegeSet == "Administrator" {
			adminAccount = account
			log.Printf("[WARN] No administrator account found for comparison, privilege validation will be skipped")
			break
		}
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

	validatePrivilegeSet(privileges.JSSObjects, adminAccount.Privileges.JSSObjects, "JSS Objects")
	validatePrivilegeSet(privileges.JSSSettings, adminAccount.Privileges.JSSSettings, "JSS Settings")
	validatePrivilegeSet(privileges.JSSActions, adminAccount.Privileges.JSSActions, "JSS Actions")
	validatePrivilegeSet(privileges.Recon, adminAccount.Privileges.Recon, "Recon")
	validatePrivilegeSet(privileges.CasperAdmin, adminAccount.Privileges.CasperAdmin, "Casper Admin")
	validatePrivilegeSet(privileges.CasperRemote, adminAccount.Privileges.CasperRemote, "Casper Remote")
	validatePrivilegeSet(privileges.CasperImaging, adminAccount.Privileges.CasperImaging, "Casper Imaging")

	if len(invalidPrivileges) > 0 {
		var msg strings.Builder
		msg.WriteString(fmt.Sprintf("Invalid privileges found when compared to Jamf Pro version %s:\n", *versionInfo.Version))

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
