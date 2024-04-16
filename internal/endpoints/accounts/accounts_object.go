// File: accounts_object.go
package accounts

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAccount constructs an Account object from the provided schema data.
func constructJamfProAccount(d *schema.ResourceData) (*jamfpro.ResourceAccount, error) {
	account := &jamfpro.ResourceAccount{
		Name:                d.Get("name").(string),
		DirectoryUser:       d.Get("directory_user").(bool),
		FullName:            d.Get("full_name").(string),
		Email:               d.Get("email").(string),
		Enabled:             d.Get("enabled").(string),
		ForcePasswordChange: d.Get("force_password_change").(bool),
		AccessLevel:         d.Get("access_level").(string),
		Password:            d.Get("password").(string),
		PrivilegeSet:        d.Get("privilege_set").(string),
	}

	if v, ok := d.GetOk("identity_server"); ok && len(v.([]interface{})) > 0 {
		ldapServerData := v.([]interface{})[0].(map[string]interface{})
		account.LdapServer = jamfpro.AccountSubsetLdapServer{
			ID: ldapServerData["id"].(int),
		}
	}

	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		siteData := v.([]interface{})[0].(map[string]interface{})
		account.Site = jamfpro.SharedResourceSite{
			ID:   siteData["id"].(int),
			Name: siteData["name"].(string),
		}
	}

	account.Privileges = jamfpro.AccountSubsetPrivileges{
		JSSObjects:    getStringSliceFromInterface(d.Get("jss_objects_privileges")),
		JSSSettings:   getStringSliceFromInterface(d.Get("jss_settings_privileges")),
		JSSActions:    getStringSliceFromInterface(d.Get("jss_actions_privileges")),
		Recon:         getStringSliceFromInterface(d.Get("recon_privileges")),
		CasperAdmin:   getStringSliceFromInterface(d.Get("casper_admin_privileges")),
		CasperRemote:  getStringSliceFromInterface(d.Get("casper_remote_privileges")),
		CasperImaging: getStringSliceFromInterface(d.Get("casper_imaging_privileges")),
	}

	if v, ok := d.GetOk("groups"); ok {
		groupsSet := v.(*schema.Set)
		for _, groupItem := range groupsSet.List() {
			groupData := groupItem.(map[string]interface{})
			group := jamfpro.AccountsListSubsetGroups{
				ID:   groupData["id"].(int),
				Name: groupData["name"].(string),
			}

			account.Groups = append(account.Groups, group)
		}
	}

	resourceXML, err := xml.MarshalIndent(account, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Account '%s' to XML: %v", account.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Account XML:\n%s\n", string(resourceXML))

	return account, nil
}

func getStringSliceFromInterface(i interface{}) []string {
	var slice []string
	if i != nil {
		for _, item := range i.([]interface{}) {
			slice = append(slice, item.(string))
		}
	}
	return slice
}
