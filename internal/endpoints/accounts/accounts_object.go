// accounts_object.go
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

	// Handle Site
	if v, ok := d.GetOk("site"); ok {
		account.Site = constructSharedResourceSite(v.([]interface{}))
	} else {
		// Set default values if 'site' data is not provided
		account.Site = constructSharedResourceSite([]interface{}{})
	}

	// Handle Privileges
	account.Privileges = constructAccountSubsetPrivileges(d)

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

// Helper functions for nested structures

// constructSharedResourceSite constructs a SharedResourceSite object from the provided schema data,
// setting default values if none are presented.
func constructSharedResourceSite(data []interface{}) jamfpro.SharedResourceSite {
	// Check if 'site' data is provided and non-empty
	if len(data) > 0 && data[0] != nil {
		site := data[0].(map[string]interface{})

		// Return the 'site' object with data from the schema
		return jamfpro.SharedResourceSite{
			ID:   site["id"].(int),
			Name: site["name"].(string),
		}
	}

	// Return default 'site' values if no data is provided or it is empty
	return jamfpro.SharedResourceSite{
		ID:   -1,     // Default ID
		Name: "None", // Default name
	}
}

// constructAccountSubsetPrivileges constructs AccountSubsetPrivileges from schema data.
func constructAccountSubsetPrivileges(d *schema.ResourceData) jamfpro.AccountSubsetPrivileges {
	privileges := jamfpro.AccountSubsetPrivileges{}

	if v, ok := d.GetOk("jss_objects_privileges"); ok {
		privileges.JSSObjects = getStringSliceFromSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("jss_settings_privileges"); ok {
		privileges.JSSSettings = getStringSliceFromSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("jss_actions_privileges"); ok {
		privileges.JSSActions = getStringSliceFromSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("casper_admin_privileges"); ok {
		privileges.CasperAdmin = getStringSliceFromSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("casper_remote_privileges"); ok {
		privileges.CasperRemote = getStringSliceFromSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("casper_imaging_privileges"); ok {
		privileges.CasperImaging = getStringSliceFromSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("recon_privileges"); ok {
		privileges.Recon = getStringSliceFromSet(v.(*schema.Set))
	}

	return privileges
}

// getStringSliceFromSet converts a *schema.Set to a slice of strings.
func getStringSliceFromSet(set *schema.Set) []string {
	list := set.List()
	slice := make([]string, len(list))
	for i, item := range list {
		slice[i] = item.(string) // Direct assertion to string, assuming all items are strings.
	}
	return slice
}
