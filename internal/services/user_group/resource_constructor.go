package user_group

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
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

	criteria := d.Get("criteria").([]any)
	for _, criterion := range criteria {
		c := criterion.(map[string]any)
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
		assignedUsers := d.Get("assigned_user_ids").(*schema.Set).List()
		if len(assignedUsers) > 0 {
			for _, v := range assignedUsers {
				resource.Users = append(resource.Users, jamfpro.UserGroupSubsetUserItem{
					ID: v.(int),
				})
			}
		}
	}

	var err error
	resource.UserAdditions, err = extractUsers(newUserOperations(d, "user_additions"))
	if err != nil {
		return resource, err
	}
	resource.UserDeletions, err = extractUsers(newUserOperations(d, "user_deletions"))
	if err != nil {
		return resource, err
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro User Group  '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro User Group  XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// extractUsers converts a slice of any that represents user data into a slice
// of jamfpro.UserGroupSubsetUserItem. It iterates over each user in the interface
// slice, extracts the relevant fields, and constructs a UserGroupSubsetUserItem for
// each user. The resulting slice of UserGroupSubsetUserItem is suitable for use in
// constructing a jamfpro.ResourceUserGroup object.
func extractUsers(usersInterface []any) ([]jamfpro.UserGroupSubsetUserItem, error) {
	var users []jamfpro.UserGroupSubsetUserItem
	for _, user := range usersInterface {
		u := user.(map[string]any)
		id := 0
		if raw := u["id"].(string); raw != "" {
			var err error
			id, err = strconv.Atoi(raw)
			if err != nil || id <= 0 {
				return nil, fmt.Errorf("user ID must be a positive integer: %q", raw)
			}
		}
		if id == 0 && u["username"].(string) == "" {
			return nil, fmt.Errorf("user operation requires id or username")
		}
		userItem := jamfpro.UserGroupSubsetUserItem{
			ID:           id,
			Username:     u["username"].(string),
			FullName:     u["full_name"].(string),
			PhoneNumber:  u["phone_number"].(string),
			EmailAddress: u["email_address"].(string),
		}
		users = append(users, userItem)
	}
	return users, nil
}

// Operation blocks record submitted commands. Removing a block does not invert it.
func newUserOperations(d *schema.ResourceData, key string) []any {
	if !d.HasChange(key) {
		return nil
	}
	old, next := d.GetChange(key)
	return next.(*schema.Set).Difference(old.(*schema.Set)).List()
}
