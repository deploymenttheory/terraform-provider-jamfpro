package guid_list_sharder

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// shardByPercentage distributes IDs according to specified percentages.
// Target shard sizes are calculated based on percentages of total IDs.
// If reservations are provided, reserved counts are subtracted from targets,
// and reserved IDs are placed first in their designated shards before distribution.
func shardByPercentage(ctx context.Context, ids []string, percentages []int64, seed string, reservations *reservationInfo) [][]string {
	unreservedIDs := ids
	totalIds := len(ids)
	
	if reservations != nil {
		unreservedIDs = reservations.UnreservedIDs
	}

	shardCount := len(percentages)
	shards := make([][]string, shardCount)

	if len(unreservedIDs) == 0 {
		return shards
	}

	workingIds := prepareIDsForDistribution(unreservedIDs, seed)

	currentIndex := 0
	for i, percentage := range percentages {
		var shardSize int
		if i == shardCount-1 {
			shardSize = len(unreservedIDs) - currentIndex
		} else {
			shardSize = int(float64(totalIds) * float64(percentage) / 100.0)
			if reservations != nil {
				shardSize -= reservations.CountsByShard[i]
			}
		}

		if currentIndex+shardSize > len(unreservedIDs) {
			shardSize = len(unreservedIDs) - currentIndex
		}

		if shardSize > 0 {
			shards[i] = workingIds[currentIndex : currentIndex+shardSize]
			currentIndex += shardSize
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
			tflog.Debug(ctx, "Percentage: sorting shard", map[string]any{
				"shard": i,
				"size":  len(shards[i]),
			})
		}
		sortIDsNumerically(shards[i])
	}

	return shards
}
