package guid_list_sharder

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// shardBySize distributes IDs according to specified absolute sizes.
// Supports -1 in the last position to mean "all remaining IDs".
// If reservations are provided, reserved counts are subtracted from targets,
// and reserved IDs are placed first in their designated shards before distribution.
func shardBySize(ctx context.Context, ids []string, sizes []int64, seed string, reservations *reservationInfo) [][]string {
	unreservedIDs := ids
	if reservations != nil {
		unreservedIDs = reservations.UnreservedIDs
	}

	shardCount := len(sizes)
	shards := make([][]string, shardCount)

	if len(unreservedIDs) == 0 {
		return shards
	}

	workingIds := prepareIDsForDistribution(unreservedIDs, seed)

	currentIndex := 0
	for i, size := range sizes {
		var shardSize int

		if size == -1 {
			shardSize = len(unreservedIDs) - currentIndex
		} else {
			shardSize = int(size)
			if reservations != nil {
				shardSize -= reservations.CountsByShard[i]
			}
			if currentIndex+shardSize > len(unreservedIDs) {
				shardSize = len(unreservedIDs) - currentIndex
			}
		}

		if shardSize > 0 && currentIndex < len(unreservedIDs) {
			shards[i] = workingIds[currentIndex : currentIndex+shardSize]
			currentIndex += shardSize
		} else {
			shards[i] = []string{}
		}
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
			tflog.Debug(ctx, "Size: sorting shard", map[string]any{
				"shard": i,
				"size":  len(shards[i]),
			})
		}
		sortIDsNumerically(shards[i])
	}

	return shards
}
