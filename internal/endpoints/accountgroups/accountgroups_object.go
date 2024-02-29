// accountgroups_object.go
package accountgroups

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAccountGroup constructs an AccountGroup object from the provided schema data.
func constructJamfProAccountGroup(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceAccountGroup, error) {
	accountGroup := &jamfpro.ResourceAccountGroup{
		Name:         d.Get("name").(string),
		AccessLevel:  d.Get("access_level").(string),
		PrivilegeSet: d.Get("privilege_set").(string),
	}

	// Handle Site
	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		siteData := v.([]interface{})[0].(map[string]interface{})
		accountGroup.Site = jamfpro.SharedResourceSite{
			ID:   siteData["id"].(int),
			Name: siteData["name"].(string),
		}
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

	// Serialize and pretty-print the accountGroup object as XML for logging
	resourceXML, err := xml.MarshalIndent(accountGroup, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Account Group '%s' to XML: %v", accountGroup.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro Account Group XML:\n%s\n", string(resourceXML))

	return accountGroup, nil
}

// constructAccountSubsetPrivileges constructs AccountSubsetPrivileges from schema data.
func constructAccountSubsetPrivileges(d *schema.ResourceData) jamfpro.AccountSubsetPrivileges {
	return jamfpro.AccountSubsetPrivileges{
		JSSObjects:    getStringSliceFromSet(d.Get("jss_objects_privileges").(*schema.Set)),
		JSSSettings:   getStringSliceFromSet(d.Get("jss_settings_privileges").(*schema.Set)),
		JSSActions:    getStringSliceFromSet(d.Get("jss_actions_privileges").(*schema.Set)),
		CasperAdmin:   getStringSliceFromSet(d.Get("casper_admin_privileges").(*schema.Set)),
		CasperRemote:  getStringSliceFromSet(d.Get("casper_remote_privileges").(*schema.Set)),
		CasperImaging: getStringSliceFromSet(d.Get("casper_imaging_privileges").(*schema.Set)),
		Recon:         getStringSliceFromSet(d.Get("recon_privileges").(*schema.Set)),
	}
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
