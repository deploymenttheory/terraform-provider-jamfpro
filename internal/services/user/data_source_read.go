package user

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const defaultReadTimeout = 30 * time.Second

// Read fetches Jamf Pro user data and maps it into the Terraform state.
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := data.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	var items []UserItemModel
	var selector string

	switch {
	case !data.UserID.IsNull() && data.UserID.ValueString() != "":
		selector = "user_id-" + data.UserID.ValueString()
		userID := data.UserID.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Reading Jamf Pro user by ID: %s", userID))
		resource, err := d.client.GetUserByID(userID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Jamf Pro User",
				fmt.Sprintf("Failed to read user with ID '%s': %v", userID, err),
			)
			return
		}
		items = []UserItemModel{mapResourceUser(resource)}

	case !data.Name.IsNull() && data.Name.ValueString() != "":
		selector = "name-" + data.Name.ValueString()
		name := data.Name.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Reading Jamf Pro user by name: %s", name))
		resource, err := d.client.GetUserByName(name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Jamf Pro User",
				fmt.Sprintf("Failed to read user with name '%s': %v", name, err),
			)
			return
		}
		items = []UserItemModel{mapResourceUser(resource)}

	case !data.Email.IsNull() && data.Email.ValueString() != "":
		selector = "email-" + data.Email.ValueString()
		email := data.Email.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Reading Jamf Pro user by email: %s", email))
		resource, err := d.client.GetUserByEmail(email)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Jamf Pro User",
				fmt.Sprintf("Failed to read user with email '%s': %v", email, err),
			)
			return
		}
		items = []UserItemModel{mapResourceUser(resource)}

	case !data.ListAll.IsNull() && data.ListAll.ValueBool():
		selector = "list_all"
		tflog.Debug(ctx, "Listing all Jamf Pro users")
		usersList, err := d.client.GetUsers()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Listing Jamf Pro Users",
				fmt.Sprintf("Failed to list users: %v", err),
			)
			return
		}
		items = make([]UserItemModel, 0, len(usersList.Users))
		for _, u := range usersList.Users {
			items = append(items, mapUsersListItem(u))
		}

	default:
		resp.Diagnostics.AddError(
			"Missing Lookup Attribute",
			"One of 'user_id', 'name', 'email', or 'list_all' must be provided.",
		)
		return
	}

	data.Items = items
	data.ID = types.StringValue("jamfpro_user-" + selector)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
