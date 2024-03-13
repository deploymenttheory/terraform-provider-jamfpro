// usergroups_object.go
package usergroups

import (
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProUserGroup constructs a ResourceUserGroup object from the provided schema data.
func constructJamfProUserGroup(d *schema.ResourceData) (*jamfpro.ResourceUserGroup, error) {
	userGroup := &jamfpro.ResourceUserGroup{
		Name:             d.Get("name").(string),
		IsSmart:          d.Get("is_smart").(bool),
		IsNotifyOnChange: d.Get("is_notify_on_change").(bool),
	}

	if v, ok := d.GetOk("id"); ok {
		userGroup.ID = v.(int)
	}

	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		site := v.([]interface{})[0].(map[string]interface{})
		userGroup.Site = jamfpro.SharedResourceSite{
			ID:   site["id"].(int),
			Name: site["name"].(string),
		}
	}

	criteria := d.Get("criteria").([]interface{})
	for _, criterion := range criteria {
		c := criterion.(map[string]interface{})
		userGroup.Criteria = append(userGroup.Criteria, jamfpro.SharedSubsetCriteria{
			Name:         c["name"].(string),
			Priority:     c["priority"].(int),
			AndOr:        c["and_or"].(string),
			SearchType:   c["search_type"].(string),
			Value:        c["value"].(string),
			OpeningParen: c["opening_paren"].(bool),
			ClosingParen: c["closing_paren"].(bool),
		})
	}

	userGroup.Users = extractUsers(d.Get("users").([]interface{}))
	userGroup.UserAdditions = extractUsers(d.Get("user_additions").([]interface{}))
	userGroup.UserDeletions = extractUsers(d.Get("user_deletions").([]interface{}))

	// Debugging - Serialize and pretty-print the user group object
	resourceXML, err := xml.MarshalIndent(userGroup, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user group '%s' to XML: %v", userGroup.Name, err)
	}
	fmt.Printf("Constructed User Group XML:\n%s\n", string(resourceXML))

	return userGroup, nil
}

// extractUsers converts a slice of interface{} that represents user data into a slice
// of jamfpro.UserGroupSubsetUserItem. It iterates over each user in the interface
// slice, extracts the relevant fields, and constructs a UserGroupSubsetUserItem for
// each user. The resulting slice of UserGroupSubsetUserItem is suitable for use in
// constructing a jamfpro.ResourceUserGroup object.
func extractUsers(usersInterface []interface{}) []jamfpro.UserGroupSubsetUserItem {
	var users []jamfpro.UserGroupSubsetUserItem
	for _, user := range usersInterface {
		u := user.(map[string]interface{})
		userItem := jamfpro.UserGroupSubsetUserItem{
			ID:           u["id"].(int),
			Username:     u["username"].(string),
			FullName:     u["full_name"].(string),
			PhoneNumber:  u["phone_number"].(string),
			EmailAddress: u["email_address"].(string),
		}
		users = append(users, userItem)
	}
	return users
}
