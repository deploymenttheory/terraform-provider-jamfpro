// usergroups_object.go
package usergroups

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/constructobject"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProUserGroup constructs a ResourceUserGroup object from the provided schema data.
func constructJamfProUserGroup(d *schema.ResourceData) (*jamfpro.ResourceUserGroup, error) {
	userGroup := &jamfpro.ResourceUserGroup{
		Name:             d.Get("name").(string),
		IsSmart:          d.Get("is_smart").(bool),
		IsNotifyOnChange: d.Get("is_notify_on_change").(bool),
	}

	if v, ok := d.GetOk("site_id"); ok {
		userGroup.Site = constructobject.ConstructSharedResourceSite(v.([]interface{}))
	} else {
		userGroup.Site = constructobject.ConstructSharedResourceSite([]interface{}{})
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

	if v, ok := d.GetOk("users"); ok && len(v.([]interface{})) > 0 {
		usersBlock := v.([]interface{})[0].(map[string]interface{})
		userIDList := usersBlock["id"].([]interface{})
		for _, userID := range userIDList {
			userIDStr, ok := userID.(string)
			if !ok {
				return nil, fmt.Errorf("user ID is not a string as expected: %v", userID)
			}

			intID, err := strconv.Atoi(userIDStr)
			if err != nil {
				return nil, fmt.Errorf("error converting user ID '%s' to integer: %v", userIDStr, err)
			}

			userGroup.Users = append(userGroup.Users, jamfpro.UserGroupSubsetUserItem{
				ID: intID,
			})
		}
	}

	userGroup.UserAdditions = extractUsers(d.Get("user_additions").([]interface{}))
	userGroup.UserDeletions = extractUsers(d.Get("user_deletions").([]interface{}))

	resourceXML, err := xml.MarshalIndent(userGroup, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro User Group  '%s' to XML: %v", userGroup.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro User Group  XML:\n%s\n", string(resourceXML))

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
