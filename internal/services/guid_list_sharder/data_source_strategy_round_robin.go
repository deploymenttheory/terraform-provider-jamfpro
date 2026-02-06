package guid_list_sharder

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// shardByRoundRobin distributes IDs in circular order, guaranteeing equal shard sizes ±1.
// If reservations are provided, reserved IDs are placed first in their designated shards,
// then unreserved IDs are distributed round-robin across all shards.
func shardByRoundRobin(ctx context.Context, ids []string, shardCount int, seed string, reservations *reservationInfo) [][]string {
	if shardCount <= 0 {
		shardCount = 1
	}

	unreservedIDs := ids
	if reservations != nil {
		unreservedIDs = reservations.UnreservedIDs
	}

	shards := make([][]string, shardCount)
	workingIds := prepareIDsForDistribution(unreservedIDs, seed)

	for i, id := range workingIds {
		shardIndex := i % shardCount
		shards[shardIndex] = append(shards[shardIndex], id)
	}

	if reservations != nil {
		for shardName, reservedIDs := range reservations.IDsByShard {
			var shardIndex int
			fmt.Sscanf(shardName, "shard_%d", &shardIndex)
			shards[shardIndex] = append(reservedIDs, shards[shardIndex]...)
		}
	}

	for i := range shards {
		if len(shards[i]) > 0 {
			tflog.Debug(ctx, "Round-Robin: sorting shard", map[string]any{
				"shard": i,
				"size":  len(shards[i]),
			})
		}
		sortIDsNumerically(shards[i])
	}

	return shards
}
