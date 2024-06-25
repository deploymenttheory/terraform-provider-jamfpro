// accountgroups_state.go
package accountgroups

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/utilities"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Account Groupinformation from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, response *jamfpro.ResourceAccountGroup) diag.Diagnostics {
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
		d.Set("identity_server", response.LDAPServer.ID)
	}

	d.Set("site_id", response.Site.ID)

	privilegeAttributes := map[string][]string{
		"jss_objects_privileges":  response.Privileges.JSSObjects,
		"jss_settings_privileges": response.Privileges.JSSSettings,
		"jss_actions_privileges":  response.Privileges.JSSActions,
		"casper_admin_privileges": response.Privileges.CasperAdmin,
	}

	for attrName, privileges := range privilegeAttributes {
		if err := d.Set(attrName, schema.NewSet(schema.HashString, utilities.ConvertToStringInterface(privileges))); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
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

	// members := make([]interface{}, 0)
	// for _, memberStruct := range response.Members {
	// 	member := memberStruct.User
	// 	memberMap := map[string]interface{}{
	// 		"id":   member.ID,
	// 		"name": member.Name,
	// 	}
	// 	members = append(members, memberMap)
	// }
	// if err := d.Set("members", members); err != nil {
	// 	diags = append(diags, diag.FromErr(err)...)
	// }

	return diags
}
