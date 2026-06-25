package user

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapUsersListItem maps a list endpoint item (id + name only) into the item model,
// leaving every other field null.
func mapUsersListItem(item jamfpro.UsersListItem) UserItemModel {
	return UserItemModel{
		ID:                  types.StringValue(strconv.Itoa(item.ID)),
		Name:                types.StringValue(item.Name),
		FullName:            types.StringNull(),
		Email:               types.StringNull(),
		EmailAddress:        types.StringNull(),
		PhoneNumber:         types.StringNull(),
		Position:            types.StringNull(),
		EnableCustomPhoto:   types.BoolNull(),
		CustomPhotoURL:      types.StringNull(),
		LDAPServer:          nil,
		ExtensionAttributes: nil,
		Sites:               nil,
		Links:               nil,
	}
}

// mapResourceUser maps a full user resource into the item model, including nested subsets.
func mapResourceUser(resource *jamfpro.ResourceUser) UserItemModel {
	item := UserItemModel{
		ID:                types.StringValue(strconv.Itoa(resource.ID)),
		Name:              types.StringValue(resource.Name),
		FullName:          types.StringValue(resource.FullName),
		Email:             types.StringValue(resource.Email),
		EmailAddress:      types.StringValue(resource.EmailAddress),
		PhoneNumber:       types.StringValue(resource.PhoneNumber),
		Position:          types.StringValue(resource.Position),
		EnableCustomPhoto: types.BoolValue(resource.EnableCustomPhoto),
		CustomPhotoURL:    types.StringValue(resource.CustomPhotoURL),
		LDAPServer: &UserLDAPServerModel{
			ID:   types.StringValue(strconv.Itoa(resource.LDAPServer.ID)),
			Name: types.StringValue(resource.LDAPServer.Name),
		},
		Links: &UserLinksModel{
			Computers:         mapLinkItems(resource.Links.Computers),
			Peripherals:       mapLinkItems(resource.Links.Peripherals),
			MobileDevices:     mapLinkItems(resource.Links.MobileDevices),
			VPPAssignments:    mapLinkItems(resource.Links.VPPAssignments),
			TotalVPPCodeCount: types.Int64Value(int64(resource.Links.TotalVPPCodeCount)),
		},
	}

	item.ExtensionAttributes = make([]UserExtensionAttributeModel, 0, len(resource.ExtensionAttributes.Attributes))
	for _, ea := range resource.ExtensionAttributes.Attributes {
		item.ExtensionAttributes = append(item.ExtensionAttributes, UserExtensionAttributeModel{
			ID:    types.StringValue(strconv.Itoa(ea.ID)),
			Name:  types.StringValue(ea.Name),
			Type:  types.StringValue(ea.Type),
			Value: types.StringValue(ea.Value),
		})
	}

	item.Sites = make([]UserSiteModel, 0, len(resource.Sites))
	for _, site := range resource.Sites {
		item.Sites = append(item.Sites, UserSiteModel{
			ID:   types.StringValue(strconv.Itoa(site.ID)),
			Name: types.StringValue(site.Name),
		})
	}

	return item
}

// mapLinkItems maps a slice of linked object references into item models.
func mapLinkItems(links []jamfpro.UserSubsetLinksListItem) []UserLinkItemModel {
	result := make([]UserLinkItemModel, 0, len(links))
	for _, link := range links {
		result = append(result, UserLinkItemModel{
			ID:   types.StringValue(strconv.Itoa(link.ID)),
			Name: types.StringValue(link.Name),
		})
	}
	return result
}
