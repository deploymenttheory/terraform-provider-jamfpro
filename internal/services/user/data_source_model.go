package user

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserDataSourceModel describes the Terraform data source model for Jamf Pro users.
type UserDataSourceModel struct {
	ID       types.String    `tfsdk:"id"`
	UserID   types.String    `tfsdk:"user_id"`
	Name     types.String    `tfsdk:"name"`
	Email    types.String    `tfsdk:"email"`
	ListAll  types.Bool      `tfsdk:"list_all"`
	Items    []UserItemModel `tfsdk:"items"`
	Timeouts timeouts.Value  `tfsdk:"timeouts"`
}

// UserItemModel represents an individual Jamf Pro user.
type UserItemModel struct {
	ID                  types.String                  `tfsdk:"id"`
	Name                types.String                  `tfsdk:"name"`
	FullName            types.String                  `tfsdk:"full_name"`
	Email               types.String                  `tfsdk:"email"`
	EmailAddress        types.String                  `tfsdk:"email_address"`
	PhoneNumber         types.String                  `tfsdk:"phone_number"`
	Position            types.String                  `tfsdk:"position"`
	EnableCustomPhoto   types.Bool                    `tfsdk:"enable_custom_photo_url"`
	CustomPhotoURL      types.String                  `tfsdk:"custom_photo_url"`
	LDAPServer          *UserLDAPServerModel          `tfsdk:"ldap_server"`
	ExtensionAttributes []UserExtensionAttributeModel `tfsdk:"extension_attributes"`
	Sites               []UserSiteModel               `tfsdk:"sites"`
	Links               *UserLinksModel               `tfsdk:"links"`
}

// UserLDAPServerModel represents the LDAP server associated with a user.
type UserLDAPServerModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// UserExtensionAttributeModel represents a user extension attribute.
type UserExtensionAttributeModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

// UserSiteModel represents a site associated with a user.
type UserSiteModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// UserLinksModel represents the objects linked to a user.
type UserLinksModel struct {
	Computers         []UserLinkItemModel `tfsdk:"computers"`
	Peripherals       []UserLinkItemModel `tfsdk:"peripherals"`
	MobileDevices     []UserLinkItemModel `tfsdk:"mobile_devices"`
	VPPAssignments    []UserLinkItemModel `tfsdk:"vpp_assignments"`
	TotalVPPCodeCount types.Int64         `tfsdk:"total_vpp_code_count"`
}

// UserLinkItemModel represents an individual linked object reference.
type UserLinkItemModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
