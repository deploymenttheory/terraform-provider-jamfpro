// accountgroups_object.go
package accountgroups

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAccountGroup constructs an AccountGroup object from the provided schema data.
func constructJamfProAccountGroup(d *schema.ResourceData) (*jamfpro.ResourceAccountGroup, error) {
	accountGroup := &jamfpro.ResourceAccountGroup{
		Name:         d.Get("name").(string),
		AccessLevel:  d.Get("access_level").(string),
		PrivilegeSet: d.Get("privilege_set").(string),
	}

	accountGroup.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))
	accountGroup.Privileges = constructAccountSubsetPrivileges(d)

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

	if v, ok := d.GetOk("identity_server"); ok && len(v.([]interface{})) > 0 {
		identityServerData := v.([]interface{})[0].(map[string]interface{})
		accountGroup.LDAPServer = jamfpro.AccountGroupSubsetLDAPServer{
			ID: identityServerData["id"].(int),
		}
	}

	resourceXML, err := xml.MarshalIndent(accountGroup, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Account Group '%s' to XML: %v", accountGroup.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Account Group XML:\n%s\n", string(resourceXML))

	return accountGroup, nil
}

// Helper functions for nested structures

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
// TODO does this need to be moved out?
func getStringSliceFromSet(set *schema.Set) []string {
	list := set.List()
	slice := make([]string, len(list))
	for i, item := range list {
		slice[i] = item.(string) // Direct assertion to string, assuming all items are strings.
	}
	return slice
}
