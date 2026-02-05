package guid_list_sharder

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (d *guidListSharderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GuidListSharderDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceType := state.SourceType.ValueString()
	strategy := state.Strategy.ValueString()

	var groupId string
	if !state.GroupId.IsNull() {
		groupId = state.GroupId.ValueString()
	}

	var seed string
	if !state.Seed.IsNull() {
		seed = state.Seed.ValueString()
	}

	var shardCount int
	var percentages []int64
	var sizes []int64

	if !state.ShardPercentages.IsNull() {
		diags = state.ShardPercentages.ElementsAs(ctx, &percentages, false)

		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		shardCount = len(percentages)

	} else if !state.ShardSizes.IsNull() {
		diags = state.ShardSizes.ElementsAs(ctx, &sizes, false)

		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		shardCount = len(sizes)

	} else {
		shardCount = int(state.ShardCount.ValueInt64())
	}

	// get source IDs based upon context
	var ids []string

	switch sourceType {
	case "computer_inventory":
		ids = d.listAllComputers(ctx, resp)
	case "mobile_device_inventory":
		ids = d.listAllMobileDevices(ctx, resp)
	case "computer_group_membership":
		ids = d.listAllComputerGroupMembers(ctx, resp, groupId)
	case "mobile_device_group_membership":
		ids = d.listAllMobileDeviceGroupMembers(ctx, resp, groupId)
	case "user_accounts":
		ids = d.listAllUsers(ctx, resp)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d IDs for source_type '%s'", len(ids), sourceType))

	// Apply the sharding strategy based upon context
	var shards [][]string

	switch strategy {
	case "round-robin":
		shards = shardByRoundRobin(ctx, ids, shardCount, seed)
	case "percentage":
		shards = shardByPercentage(ctx, ids, percentages, seed)
	case "size":
		shards = shardBySize(ctx, ids, sizes, seed)
	case "rendezvous":
		shards = shardByRendezvous(ctx, ids, shardCount, seed)
	}

	if err := setStateToTerraform(ctx, &state, shards, sourceType, shardCount, strategy); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Set Computed State",
			fmt.Sprintf("Error setting computed state attributes: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", DataSourceName))
}

// listAllComputers retrieves all managed computer IDs from Jamf Pro computer inventory
// Filters out unmanaged computers as they cannot be added to static groups
func (d *guidListSharderDataSource) listAllComputers(ctx context.Context, resp *datasource.ReadResponse) []string {
	var ids []string

	params := url.Values{}
	params.Set("section", "GENERAL") // Only need basic info for IDs and managed status

	computers, err := d.client.GetComputersInventory(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Computer Inventory",
			fmt.Sprintf("Failed to retrieve computer inventory: %s", err.Error()),
		)
		return nil
	}

	managedCount := 0
	unmanagedCount := 0

	for _, computer := range computers.Results {
		if computer.General.RemoteManagement.Managed {
			ids = append(ids, computer.ID)
			managedCount++
		} else {
			unmanagedCount++
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Computer inventory filtered: %d managed, %d unmanaged (excluded)", managedCount, unmanagedCount))

	return ids
}

// listAllMobileDevices retrieves all managed mobile device IDs from Jamf Pro
// Filters out unmanaged devices as they cannot be added to static groups
func (d *guidListSharderDataSource) listAllMobileDevices(ctx context.Context, resp *datasource.ReadResponse) []string {
	var ids []string

	mobileDevices, err := d.client.GetMobileDevices()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Mobile Devices",
			fmt.Sprintf("Failed to retrieve mobile devices: %s", err.Error()),
		)
		return nil
	}

	managedCount := 0
	unmanagedCount := 0

	for _, device := range mobileDevices.MobileDevices {
		if device.Managed {
			ids = append(ids, strconv.Itoa(device.ID))
			managedCount++
		} else {
			unmanagedCount++
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mobile devices filtered: %d managed, %d unmanaged (excluded)", managedCount, unmanagedCount))

	return ids
}

// listAllComputerGroupMembers retrieves all computer IDs from a specific computer group
func (d *guidListSharderDataSource) listAllComputerGroupMembers(ctx context.Context, resp *datasource.ReadResponse, groupId string) []string {
	var ids []string

	if groupId == "" {
		resp.Diagnostics.AddError(
			"Missing Group ID",
			"group_id is required when source_type is 'computer_group_membership'",
		)
		return nil
	}

	group, err := d.client.GetComputerGroupByID(groupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Computer Group",
			fmt.Sprintf("Failed to retrieve computer group with ID %s: %s", groupId, err.Error()),
		)
		return nil
	}

	if group.Computers != nil {
		for _, computer := range *group.Computers {
			ids = append(ids, strconv.Itoa(computer.ID))
		}
	}

	return ids
}

// listAllMobileDeviceGroupMembers retrieves all mobile device IDs from a specific mobile device group
func (d *guidListSharderDataSource) listAllMobileDeviceGroupMembers(ctx context.Context, resp *datasource.ReadResponse, groupId string) []string {
	var ids []string

	if groupId == "" {
		resp.Diagnostics.AddError(
			"Missing Group ID",
			"group_id is required when source_type is 'mobile_device_group_membership'",
		)
		return nil
	}

	group, err := d.client.GetMobileDeviceGroupByID(groupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Mobile Device Group",
			fmt.Sprintf("Failed to retrieve mobile device group with ID %s: %s", groupId, err.Error()),
		)
		return nil
	}

	if group.MobileDevices != nil {
		for _, device := range *group.MobileDevices {
			ids = append(ids, strconv.Itoa(device.ID))
		}
	}

	return ids
}

// listAllUsers retrieves all user IDs from Jamf Pro
func (d *guidListSharderDataSource) listAllUsers(ctx context.Context, resp *datasource.ReadResponse) []string {
	var ids []string

	users, err := d.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Users",
			fmt.Sprintf("Failed to retrieve users: %s", err.Error()),
		)
		return nil
	}

	for _, user := range users.Users {
		ids = append(ids, strconv.Itoa(user.ID))
	}

	return ids
}
