package user_group

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProUserGroup constructs a ResourceUserGroup object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceUserGroup, error) {
	resource := &jamfpro.ResourceUserGroup{
		Name:             d.Get("name").(string),
		IsSmart:          d.Get("is_smart").(bool),
		IsNotifyOnChange: d.Get("is_notify_on_change").(bool),
	}

	resource.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))

	criteria := d.Get("criteria").([]interface{})
	for _, criterion := range criteria {
		c := criterion.(map[string]interface{})
		resource.Criteria = append(resource.Criteria, jamfpro.SharedSubsetCriteria{
			Name:         c["name"].(string),
			Priority:     c["priority"].(int),
			AndOr:        c["and_or"].(string),
			SearchType:   c["search_type"].(string),
			Value:        c["value"].(string),
			OpeningParen: c["opening_paren"].(bool),
			ClosingParen: c["closing_paren"].(bool),
		})
	}

	if !resource.IsSmart {
		assignedUsers := d.Get("assigned_user_ids").([]interface{})
		if len(assignedUsers) > 0 {
			for _, v := range assignedUsers {
				resource.Users = append(resource.Users, jamfpro.UserGroupSubsetUserItem{
					ID: v.(int),
				})
			}
		}
	}

	resource.UserAdditions = extractUsers(d.Get("user_additions").([]interface{}))
	resource.UserDeletions = extractUsers(d.Get("user_deletions").([]interface{}))

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro User Group  '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro User Group  XML:\n%s\n", string(resourceXML))

	return resource, nil
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
