package guid_list_sharder

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// shardBySize distributes IDs according to specified absolute sizes.
// Supports -1 in the last position to mean "all remaining IDs".
func shardBySize(ctx context.Context, ids []string, sizes []int64, seed string) [][]string {
	totalIds := len(ids)
	shardCount := len(sizes)
	shards := make([][]string, shardCount)

	if totalIds == 0 {
		return shards
	}

	workingIds := prepareIDsForDistribution(ids, seed)

	currentIndex := 0
	for i, size := range sizes {
		var shardSize int

		if size == -1 {
			shardSize = totalIds - currentIndex
		} else {
			shardSize = int(size)
			if currentIndex+shardSize > totalIds {
				shardSize = totalIds - currentIndex
			}
		}

		if shardSize > 0 && currentIndex < totalIds {
			shards[i] = workingIds[currentIndex : currentIndex+shardSize]
			currentIndex += shardSize
		} else {
			// Initialize empty slice to prevent nil in Terraform state
			shards[i] = []string{}
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

// calculateSizeTargets calculates target shard sizes for size-based distribution.
func calculateSizeTargets(sizes []int64) []int {
	targets := make([]int, len(sizes))
	for i, size := range sizes {
		targets[i] = int(size)
	}
	return targets
}
