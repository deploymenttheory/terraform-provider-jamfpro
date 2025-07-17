package user_group

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest UserGroup information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceUserGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", strconv.Itoa(resp.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_smart", resp.IsSmart); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("is_notify_on_change", resp.IsNotifyOnChange); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	d.Set("site_id", resp.Site.ID)

	criteria := make([]interface{}, len(resp.Criteria))
	for i, criterion := range resp.Criteria {
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

	if !resp.IsSmart {
		var userIDStrList []int
		if len(resp.Users) > 0 {
			for _, v := range resp.Users {
				userIDStrList = append(userIDStrList, v.ID)
			}
		}

		err := d.Set("assigned_user_ids", userIDStrList)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if err := d.Set("user_additions", setUserItem(resp.UserAdditions)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("user_deletions", setUserItem(resp.UserDeletions)); err != nil {
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
