// accounts_object.go
package accounts

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAccount constructs an Account object from the provided schema data.
func constructJamfProAccount(d *schema.ResourceData) (*jamfpro.ResourceAccount, error) {
	var resource *jamfpro.ResourceAccount

	resource = &jamfpro.ResourceAccount{
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
		resource.LdapServer = jamfpro.AccountSubsetLdapServer{
			ID: ldapServerData["id"].(int),
		}
	}

	resource.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))
	resource.Privileges = constructAccountSubsetPrivileges(d)

	if v, ok := d.GetOk("groups"); ok {
		groupsSet := v.(*schema.Set)
		for _, groupItem := range groupsSet.List() {
			groupData := groupItem.(map[string]interface{})
			group := jamfpro.AccountsListSubsetGroups{
				ID:   groupData["id"].(int),
				Name: groupData["name"].(string),
			}

			resource.Groups = append(resource.Groups, group)
		}
	}

	return resource, nil
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
		slice[i] = item.(string)
	}
	return slice
}
