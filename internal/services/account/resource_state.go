// accounts_state.go
package account

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Account information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceAccount) diag.Diagnostics {
	var diags diag.Diagnostics

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

	if resp.Site != nil {
		d.Set("site_id", resp.Site.ID)
	} else {
		d.Set("site_id", -1)
	}

	if err := d.Set("jss_actions_privileges", resp.Privileges.JSSActions); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("jss_objects_privileges", resp.Privileges.JSSObjects); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("jss_settings_privileges", resp.Privileges.JSSSettings); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("recon_privileges", resp.Privileges.Recon); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}
