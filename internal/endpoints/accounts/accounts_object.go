// accounts_resource.go
package accounts

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAccount constructs an Account object from the provided schema data.
func constructJamfProAccount(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceAccount, error) {
	account := &jamfpro.ResourceAccount{
		Name:                d.Get("name").(string),
		DirectoryUser:       d.Get("directory_user").(bool),
		FullName:            d.Get("full_name").(string),
		Email:               d.Get("email").(string),
		EmailAddress:        d.Get("email_address").(string),
		Enabled:             d.Get("enabled").(string),
		ForcePasswordChange: d.Get("force_password_change").(bool),
		AccessLevel:         d.Get("access_level").(string),
		Password:            d.Get("password").(string),
		PrivilegeSet:        d.Get("privilege_set").(string),
	}

	// Handle LdapServer
	if v, ok := d.GetOk("ldap_server"); ok && len(v.([]interface{})) > 0 {
		ldapServerData := v.([]interface{})[0].(map[string]interface{})
		account.LdapServer = jamfpro.AccountSubsetLdapServer{
			ID:   ldapServerData["id"].(int),
			Name: ldapServerData["name"].(string),
		}
	}

	// Handle Site
	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		siteData := v.([]interface{})[0].(map[string]interface{})
		account.Site = jamfpro.SharedResourceSite{
			ID:   siteData["id"].(int),
			Name: siteData["name"].(string),
		}
	}

	// Handle Groups
	if v, ok := d.GetOk("groups"); ok {
		groupsSet := v.(*schema.Set)
		for _, groupItem := range groupsSet.List() {
			groupData := groupItem.(map[string]interface{})
			group := jamfpro.AccountsListSubsetGroups{
				Name:       groupData["name"].(string),
				Privileges: constructGroupPrivileges(groupData),
			}

			// Handle Site for Group
			if site, ok := groupData["site"].([]interface{}); ok && len(site) > 0 {
				siteMap := site[0].(map[string]interface{})
				group.Site = jamfpro.SharedResourceSite{
					ID:   siteMap["id"].(int),
					Name: siteMap["name"].(string),
				}
			}

			account.Groups = append(account.Groups, group)
		}
	}

	// Serialize and pretty-print the account object as XML for logging
	resourceXML, err := xml.MarshalIndent(account, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Account '%s' to XML: %v", account.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro Account XML:\n%s\n", string(resourceXML))

	return account, nil
}

// constructGroupPrivileges constructs a Privileges object from group data.
func constructGroupPrivileges(groupData map[string]interface{}) jamfpro.AccountSubsetPrivileges {
	return jamfpro.AccountSubsetPrivileges{
		JSSObjects:    getStringSliceFromInterface(groupData["jss_objects_privileges"]),
		JSSSettings:   getStringSliceFromInterface(groupData["jss_settings_privileges"]),
		JSSActions:    getStringSliceFromInterface(groupData["jss_actions_privileges"]),
		CasperAdmin:   getStringSliceFromInterface(groupData["casper_admin_privileges"]),
		CasperRemote:  getStringSliceFromInterface(groupData["casper_remote_privileges"]),
		CasperImaging: getStringSliceFromInterface(groupData["casper_imaging_privileges"]),
		Recon:         getStringSliceFromInterface(groupData["recon_privileges"]),
	}
}

// getStringSliceFromInterface helps in converting an interface{} to a slice of strings.
func getStringSliceFromInterface(i interface{}) []string {
	var slice []string
	if i != nil {
		for _, item := range i.([]interface{}) {
			slice = append(slice, item.(string))
		}
	}
	return slice
}
