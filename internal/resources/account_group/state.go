// accountgroups_state.go
package account_group

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Account Groupinformation from the Jamf Pro API.
func updateState(d *schema.ResourceData, response *jamfpro.ResourceAccountGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("name", response.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("access_level", response.AccessLevel); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("privilege_set", response.PrivilegeSet); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if response.LDAPServer.ID != 0 {
		d.Set("identity_server_id", response.LDAPServer.ID)
	}

	d.Set("site_id", response.Site.ID)

	if err := d.Set("jss_actions_privileges", response.Privileges.JSSActions); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("jss_objects_privileges", response.Privileges.JSSObjects); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("jss_settings_privileges", response.Privileges.JSSSettings); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("casper_admin_privileges", response.Privileges.CasperAdmin); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if len(response.Members) > 0 {
		var member_ids []int
		for _, v := range response.Members {
			member_ids = append(member_ids, v.ID)
		}

		if err := d.Set("member_ids", member_ids); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

	} else {
		d.Set("member_ids", nil)
	}

	return diags
}
