// accountgroups_object.go
package accountgroups

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAccountGroup constructs an AccountGroup object from the provided schema data.
func constructJamfProAccountGroup(d *schema.ResourceData) (*jamfpro.ResourceAccountGroup, error) {
	accountGroup := &jamfpro.ResourceAccountGroup{
		Name:         d.Get("name").(string),
		AccessLevel:  d.Get("access_level").(string),
		PrivilegeSet: d.Get("privilege_set").(string),
	}

	// Handle Site
	if v, ok := d.GetOk("site"); ok {
		accountGroup.Site = constructSharedResourceSite(v.([]interface{}))
	} else {
		// Set default values if 'site' data is not provided
		accountGroup.Site = constructSharedResourceSite([]interface{}{})
	}

	// Handle Privileges
	accountGroup.Privileges = constructAccountSubsetPrivileges(d)

	// Handle Members
	if v, ok := d.GetOk("members"); ok {
		memberList := v.([]interface{})
		accountGroup.Members = make(jamfpro.AccountGroupSubsetMembers, len(memberList))
		for i, member := range memberList {
			memberData := member.(map[string]interface{})
			accountGroup.Members[i].User = jamfpro.MemberUser{
				ID:   memberData["id"].(int),
				Name: memberData["name"].(string),
			}
		}
	}

	// Handle Identity Server (LDAP Server). Fields are used for both LDAP and IdP configuration
	if v, ok := d.GetOk("identity_server"); ok && len(v.([]interface{})) > 0 {
		identityServerData := v.([]interface{})[0].(map[string]interface{})
		accountGroup.LDAPServer = jamfpro.AccountGroupSubsetLDAPServer{
			ID: identityServerData["id"].(int),
		}
	}

	// Serialize and pretty-print the accountGroup object as XML for logging
	resourceXML, err := xml.MarshalIndent(accountGroup, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Account Group '%s' to XML: %v", accountGroup.Name, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Account Group XML:\n%s\n", string(resourceXML))

	return accountGroup, nil
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
