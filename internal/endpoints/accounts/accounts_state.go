// accounts_state.go
package accounts

import (
	"log"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/utilities"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Account information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceAccount) diag.Diagnostics {

	var diags diag.Diagnostics

	// Update Terraform state with the resource information
	if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("directory_user", resource.DirectoryUser); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("full_name", resource.FullName); err != nil {
		diags = append(diags, diag.FromErr(err)...)

	}
	if err := d.Set("email", resource.Email); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("enabled", resource.Enabled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if resource.LdapServer.ID != 0 || resource.LdapServer.Name != "" {
		ldapServer := make(map[string]interface{})
		ldapServer["id"] = resource.LdapServer.ID
		d.Set("identity_server", []interface{}{ldapServer})
	} else {
		d.Set("identity_server", []interface{}{})
	}

	d.Set("force_password_change", resource.ForcePasswordChange)
	d.Set("access_level", resource.AccessLevel)
	d.Set("privilege_set", resource.PrivilegeSet)

	if resource.Site.ID != 0 || resource.Site.Name != "" {
		site := make(map[string]interface{})
		site["id"] = resource.Site.ID
		site["name"] = resource.Site.Name
		d.Set("site", []interface{}{site})
	} else {
		d.Set("site", []interface{}{})
	}

	groups := make([]interface{}, len(resource.Groups))
	for i, group := range resource.Groups {
		groupMap := make(map[string]interface{})
		groupMap["name"] = group.Name
		groupMap["id"] = group.ID

		groups[i] = groupMap
	}

	if err := d.Set("groups", groups); err != nil {
		return diag.FromErr(err)
	}

	privilegeAttributes := map[string][]string{
		"jss_objects_privileges":  resource.Privileges.JSSObjects,
		"jss_settings_privileges": resource.Privileges.JSSSettings,
		"jss_actions_privileges":  resource.Privileges.JSSActions,
		// "casper_admin_privileges":   resource.Privileges.CasperAdmin,
		// "casper_remote_privileges":  resource.Privileges.CasperRemote,
		// "casper_imaging_privileges": resource.Privileges.CasperImaging,
		"recon_privileges": resource.Privileges.Recon,
	}

	log.Println("LOGHERE")
	log.Printf("%+v", resource.Privileges.CasperAdmin)
	log.Println(privilegeAttributes)

	for attrName, privileges := range privilegeAttributes {
		if err := d.Set(attrName, schema.NewSet(schema.HashString, utilities.ConvertToStringInterface(privileges))); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
