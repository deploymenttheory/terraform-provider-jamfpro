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
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	var state GuidListSharderDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sourceIDs []string
	switch state.SourceType.ValueString() {
	case "computer_inventory":
		sourceIDs = d.listAllComputers(ctx, resp)
	case "mobile_device_inventory":
		sourceIDs = d.listAllMobileDevices(ctx, resp)
	case "computer_group_membership":
		sourceIDs = d.listAllComputerGroupMembers(ctx, resp, state.GroupId.ValueString())
	case "mobile_device_group_membership":
		sourceIDs = d.listAllMobileDeviceGroupMembers(ctx, resp, state.GroupId.ValueString())
	case "user_accounts":
		sourceIDs = d.listAllUsers(ctx, resp)
	default:
		resp.Diagnostics.AddError(
			"Invalid Source Type",
			fmt.Sprintf("Unknown source_type: %s", state.SourceType.ValueString()),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d IDs for source_type '%s'", len(sourceIDs), state.SourceType.ValueString()))

	filteredIDs, totalIDCount := d.applyExclusions(ctx, sourceIDs, &state)

	tflog.Debug(ctx, fmt.Sprintf("After exclusions: %d IDs remain", len(filteredIDs)))

	reservations := d.applyReservations(ctx, resp, filteredIDs, &state)
	if resp.Diagnostics.HasError() {
		return
	}

	shards := d.applyShardingStrategy(ctx, filteredIDs, totalIDCount, &state, reservations)

	if err := setStateToTerraform(ctx, &state, shards); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Set Computed State",
			fmt.Sprintf("Error setting computed state attributes: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", DataSourceName))
}

func (d *guidListSharderDataSource) listAllComputers(ctx context.Context, resp *datasource.ReadResponse) []string {
	params := url.Values{}
	params.Set("section", "GENERAL")

	computers, err := d.client.GetComputersInventory(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Computer Inventory",
			fmt.Sprintf("Failed to retrieve computer inventory: %s", err.Error()),
		)
		return nil
	}

	var ids []string
	managedCount, unmanagedCount := 0, 0

	for _, computer := range computers.Results {
		if computer.General.RemoteManagement.Managed {
			ids = append(ids, computer.ID)
			managedCount++
		} else {
			unmanagedCount++
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Computer inventory: %d managed, %d unmanaged (excluded)", managedCount, unmanagedCount))
	return ids
}

func (d *guidListSharderDataSource) listAllMobileDevices(ctx context.Context, resp *datasource.ReadResponse) []string {
	mobileDevices, err := d.client.GetMobileDevices()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Mobile Devices",
			fmt.Sprintf("Failed to retrieve mobile devices: %s", err.Error()),
		)
		return nil
	}

	var ids []string
	managedCount, unmanagedCount := 0, 0

	for _, device := range mobileDevices.MobileDevices {
		if device.Managed {
			ids = append(ids, strconv.Itoa(device.ID))
			managedCount++
		} else {
			unmanagedCount++
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mobile devices: %d managed, %d unmanaged (excluded)", managedCount, unmanagedCount))
	return ids
}

func (d *guidListSharderDataSource) listAllComputerGroupMembers(ctx context.Context, resp *datasource.ReadResponse, groupID string) []string {
	if groupID == "" {
		resp.Diagnostics.AddError(
			"Missing Group ID",
			"group_id is required when source_type is 'computer_group_membership'",
		)
		return nil
	}

	group, err := d.client.GetComputerGroupByID(groupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Computer Group",
			fmt.Sprintf("Failed to retrieve computer group with ID %s: %s", groupID, err.Error()),
		)
		return nil
	}

	var ids []string
	if group.Computers != nil {
		for _, computer := range *group.Computers {
			ids = append(ids, strconv.Itoa(computer.ID))
		}
	}

	return ids
}

func (d *guidListSharderDataSource) listAllMobileDeviceGroupMembers(ctx context.Context, resp *datasource.ReadResponse, groupID string) []string {
	if groupID == "" {
		resp.Diagnostics.AddError(
			"Missing Group ID",
			"group_id is required when source_type is 'mobile_device_group_membership'",
		)
		return nil
	}

	group, err := d.client.GetMobileDeviceGroupByID(groupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Mobile Device Group",
			fmt.Sprintf("Failed to retrieve mobile device group with ID %s: %s", groupID, err.Error()),
		)
		return nil
	}

	var ids []string
	if group.MobileDevices != nil {
		for _, device := range *group.MobileDevices {
			ids = append(ids, strconv.Itoa(device.ID))
		}
	}

	return ids
}

func (d *guidListSharderDataSource) listAllUsers(ctx context.Context, resp *datasource.ReadResponse) []string {
	users, err := d.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Users",
			fmt.Sprintf("Failed to retrieve users: %s", err.Error()),
		)
		return nil
	}

	var ids []string
	for _, user := range users.Users {
		ids = append(ids, strconv.Itoa(user.ID))
	}

	return ids
}

func (d *guidListSharderDataSource) applyExclusions(ctx context.Context, ids []string, state *GuidListSharderDataSourceModel) (filteredIDs []string, totalCount int) {
	if state.ExcludeIds.IsNull() || len(state.ExcludeIds.Elements()) == 0 {
		return ids, len(ids)
	}

	var excludeIDs []string
	state.ExcludeIds.ElementsAs(ctx, &excludeIDs, false)

	excludeMap := make(map[string]bool, len(excludeIDs))
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}

	filtered := make([]string, 0, len(ids))
	for _, id := range ids {
		if !excludeMap[id] {
			filtered = append(filtered, id)
		}
	}

	if excludedCount := len(ids) - len(filtered); excludedCount > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Excluded %d IDs from sharding", excludedCount))
	}

	return filtered, len(filtered)
}

func (d *guidListSharderDataSource) applyReservations(ctx context.Context, resp *datasource.ReadResponse, ids []string, state *GuidListSharderDataSourceModel) *reservationInfo {
	info := &reservationInfo{
		IDsByShard:    make(map[string][]string),
		CountsByShard: make(map[int]int),
		UnreservedIDs: ids,
	}

	if state.ReservedIds.IsNull() || len(state.ReservedIds.Elements()) == 0 {
		return info
	}

	var reservedMap map[string][]string
	var excludeIDs []string
	if !state.ExcludeIds.IsNull() {
		state.ExcludeIds.ElementsAs(ctx, &excludeIDs, false)
	}
	excludeMap := make(map[string]bool, len(excludeIDs))
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}

	shardCount := d.getShardCount(ctx, state)

	diags := state.ReservedIds.ElementsAs(ctx, &reservedMap, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return nil
	}

	seenIDs := make(map[string]string)

	for shardName, idList := range reservedMap {
		var shardIndex int
		if _, err := fmt.Sscanf(shardName, "shard_%d", &shardIndex); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Reserved ID Shard Name",
				fmt.Sprintf("Invalid shard name '%s' in reserved_ids. Must be in format 'shard_0', 'shard_1', etc.", shardName),
			)
			return nil
		}

		if shardIndex < 0 || shardIndex >= shardCount {
			resp.Diagnostics.AddError(
				"Invalid Reserved ID Shard Index",
				fmt.Sprintf("Shard name '%s' in reserved_ids is out of range. With shard_count=%d, valid shards are shard_0 to shard_%d.", shardName, shardCount, shardCount-1),
			)
			return nil
		}

		for _, id := range idList {
			if excludeMap[id] {
				resp.Diagnostics.AddError(
					"Reserved ID Conflict",
					fmt.Sprintf("ID '%s' appears in both exclude_ids and reserved_ids. Exclusion takes precedence - please remove it from reserved_ids.", id),
				)
				return nil
			}

			if prevShard, exists := seenIDs[id]; exists {
				resp.Diagnostics.AddError(
					"Duplicate Reserved ID",
					fmt.Sprintf("ID '%s' appears in multiple shards in reserved_ids: '%s' and '%s'. Each ID can only be assigned to one shard.", id, prevShard, shardName),
				)
				return nil
			}
			seenIDs[id] = shardName
		}

		info.IDsByShard[shardName] = idList
		info.CountsByShard[shardIndex] = len(idList)
	}

	if len(seenIDs) > 0 {
		reservedSet := make(map[string]bool, len(seenIDs))
		for id := range seenIDs {
			reservedSet[id] = true
		}

		filtered := make([]string, 0, len(ids))
		for _, id := range ids {
			if !reservedSet[id] {
				filtered = append(filtered, id)
			}
		}

		info.UnreservedIDs = filtered
		tflog.Debug(ctx, fmt.Sprintf("Reserved %d IDs for specific shards, %d remain for distribution", len(seenIDs), len(filtered)))
	}

	return info
}

func (d *guidListSharderDataSource) applyShardingStrategy(ctx context.Context, unreservedIDs []string, totalIDCount int, state *GuidListSharderDataSourceModel, reservations *reservationInfo) [][]string {
	strategy := state.Strategy.ValueString()
	seed := state.Seed.ValueString()
	shardCount := d.getShardCount(ctx, state)

	if strategy == "rendezvous" {
		result := shardByRendezvous(ctx, unreservedIDs, shardCount, seed)
		d.appendReservedIds(ctx, result, reservations, nil)
		return result
	}

	var targetSizes []int
	switch strategy {
	case "round-robin":
		targetSizes = calculateRoundRobinTargets(totalIDCount, shardCount)
	case "percentage":
		var percentages []int64
		state.ShardPercentages.ElementsAs(ctx, &percentages, false)
		targetSizes = calculatePercentageTargets(totalIDCount, percentages)
	case "size":
		var sizes []int64
		state.ShardSizes.ElementsAs(ctx, &sizes, false)
		targetSizes = calculateSizeTargets(sizes)
	}

	distributionTargets := make([]int, shardCount)
	for i := range shardCount {
		distributionTargets[i] = targetSizes[i] - reservations.CountsByShard[i]
	}

	tflog.Debug(ctx, fmt.Sprintf("Targets: %v, Reserved: %v, Distribution: %v", targetSizes, reservations.CountsByShard, distributionTargets))

	result := distributeIdsToShards(ctx, unreservedIDs, distributionTargets, seed)
	d.appendReservedIds(ctx, result, reservations, targetSizes)
	return result
}

func (d *guidListSharderDataSource) getShardCount(ctx context.Context, state *GuidListSharderDataSourceModel) int {
	if !state.ShardPercentages.IsNull() {
		return len(state.ShardPercentages.Elements())
	}
	if !state.ShardSizes.IsNull() {
		return len(state.ShardSizes.Elements())
	}
	return int(state.ShardCount.ValueInt64())
}

func (d *guidListSharderDataSource) appendReservedIds(ctx context.Context, shards [][]string, reservations *reservationInfo, targetSizes []int) {
	if len(reservations.IDsByShard) == 0 {
		return
	}

	for shardName, reservedIDs := range reservations.IDsByShard {
		var shardIndex int
		fmt.Sscanf(shardName, "shard_%d", &shardIndex)

		shards[shardIndex] = append(reservedIDs, shards[shardIndex]...)

		finalSize := len(shards[shardIndex])
		if targetSizes != nil {
			tflog.Debug(ctx, fmt.Sprintf("Added %d reserved to %s (final: %d, target: %d)", len(reservedIDs), shardName, finalSize, targetSizes[shardIndex]))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Added %d reserved to %s (final: %d)", len(reservedIDs), shardName, finalSize))
		}
	}
}
