// usergroups_helpers.go
package usergroups

import "github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

// convertUserItems converts a slice of jamfpro.UserGroupSubsetUserItem structs into a slice of map[string]interface{} for Terraform.
func convertUserItems(userItems []jamfpro.UserGroupSubsetUserItem) []interface{} {
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
