// accountgroups_object.go
package account_group

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/jamfprivileges"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a jamf pro Account group object from the provided schema data.
func construct(d *schema.ResourceData, meta interface{}) (*jamfpro.ResourceAccountGroup, error) {
	client := meta.(*jamfpro.Client)

	privileges := constructAccountSubsetPrivileges(d)

	if err := jamfprivileges.ValidateAccountPrivileges(client, privileges); err != nil {
		return nil, err
	}

	resource := &jamfpro.ResourceAccountGroup{
		Name:         d.Get("name").(string),
		AccessLevel:  d.Get("access_level").(string),
		PrivilegeSet: d.Get("privilege_set").(string),
	}

	resource.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))
	resource.Privileges = constructAccountSubsetPrivileges(d)

	members_ids := d.Get("member_ids").([]interface{})
	if len(members_ids) > 0 {
		for _, v := range members_ids {
			resource.Members = append(resource.Members, jamfpro.MemberUser{ID: v.(int)})
		}
	}

	resource.LDAPServer = jamfpro.SharedResourceLdapServer{
		ID: d.Get("identity_server_id").(int),
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Account Group '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Account Group XML:\n%s\n", string(resourceXML))

	return resource, nil
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
		slice[i] = item.(string)
	}
	return slice
}
