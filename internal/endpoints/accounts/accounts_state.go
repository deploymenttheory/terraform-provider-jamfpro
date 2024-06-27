// accounts_state.go
package accounts

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/utilities"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Account information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceAccount) diag.Diagnostics {
	var diags diag.Diagnostics

	log.Println("LOGHERE-IN")
	jsonData, _ := json.MarshalIndent(resp, " ", "    ")
	log.Println(string(jsonData))

	// Update Terraform state with the resource information
	if err := d.Set("id", strconv.Itoa(resp.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("directory_user", resp.DirectoryUser); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("full_name", resp.FullName); err != nil {
		diags = append(diags, diag.FromErr(err)...)

	}
	if err := d.Set("email", resp.Email); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("enabled", resp.Enabled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if resp.LdapServer.ID != 0 {
		d.Set("identity_server_id", resp.LdapServer.ID)
	}

	d.Set("force_password_change", resp.ForcePasswordChange)
	d.Set("access_level", resp.AccessLevel)
	d.Set("privilege_set", resp.PrivilegeSet)
	log.Println("LOGHERE")
	// log.Printf("%+v", resp)
	if resp.Site != nil {
		d.Set("site_id", resp.Site.ID)
	} else {
		d.Set("site_id", -1)
	}

	// TODO review this.
	privilegeAttributes := map[string][]string{
		"jss_objects_privileges":  resp.Privileges.JSSObjects,
		"jss_settings_privileges": resp.Privileges.JSSSettings,
		"jss_actions_privileges":  resp.Privileges.JSSActions,
		// "casper_admin_privileges":   resp.Privileges.CasperAdmin,
		// "casper_remote_privileges":  resp.Privileges.CasperRemote,
		// "casper_imaging_privileges": resp.Privileges.CasperImaging,
		"recon_privileges": resp.Privileges.Recon,
	}

	log.Println("LOGHERE")
	log.Printf("%+v", resp.Privileges.CasperAdmin)
	log.Println(privilegeAttributes)

	for attrName, privileges := range privilegeAttributes {
		if err := d.Set(attrName, schema.NewSet(schema.HashString, utilities.ConvertToStringInterface(privileges))); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
