// usergroups_state.go
package usergroups

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest UserGroup information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceUserGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", resource.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_notify_on_change", resource.IsNotifyOnChange); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	site := []interface{}{}
	if resource.Site.ID != -1 {
		site = append(site, map[string]interface{}{
			"id": resource.Site.ID,
		})
	}
	if len(site) > 0 {
		if err := d.Set("site", site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	criteria := make([]interface{}, len(resource.Criteria))
	for i, criterion := range resource.Criteria {
		criteria[i] = map[string]interface{}{
			"name":          criterion.Name,
			"priority":      criterion.Priority,
			"and_or":        criterion.AndOr,
			"search_type":   criterion.SearchType,
			"value":         criterion.Value,
			"opening_paren": criterion.OpeningParen,
			"closing_paren": criterion.ClosingParen,
		}
	}
	d.Set("criteria", criteria)

	if !resource.IsSmart {
		var userIDStrList []string
		for _, user := range resource.Users {
			userIDStrList = append(userIDStrList, strconv.Itoa(user.ID))
		}

		if err := d.Set("users", []interface{}{
			map[string]interface{}{
				"id": userIDStrList,
			},
		}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if err := d.Set("user_additions", setUserItem(resource.UserAdditions)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("user_deletions", setUserItem(resource.UserDeletions)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags

}

// setUserItem converts a slice of jamfpro.UserGroupSubsetUserItem structs into a slice of map[string]interface{} for Terraform.
func setUserItem(userItems []jamfpro.UserGroupSubsetUserItem) []interface{} {
	var tfUserItems []interface{}

	for _, userItem := range userItems {
		tfUserItem := make(map[string]interface{})
		tfUserItem["id"] = userItem.ID
		tfUserItem["username"] = userItem.Username
		tfUserItem["full_name"] = userItem.FullName
		tfUserItem["phone_number"] = userItem.PhoneNumber
		tfUserItem["email_address"] = userItem.EmailAddress

		tfUserItems = append(tfUserItems, tfUserItem)
	}

	return tfUserItems
}
