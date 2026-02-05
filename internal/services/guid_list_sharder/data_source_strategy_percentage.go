package guid_list_sharder

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// shardByPercentage distributes IDs according to specified percentages.
func shardByPercentage(ctx context.Context, ids []string, percentages []int64, seed string) [][]string {
	totalIds := len(ids)
	shardCount := len(percentages)
	shards := make([][]string, shardCount)

	if totalIds == 0 {
		return shards
	}

	workingIds := prepareIDsForDistribution(ids, seed)

	currentIndex := 0
	for i, percentage := range percentages {
		var shardSize int
		if i == shardCount-1 {
			shardSize = totalIds - currentIndex
		} else {
			shardSize = int(float64(totalIds) * float64(percentage) / 100.0)
		}

		if currentIndex+shardSize > totalIds {
			shardSize = totalIds - currentIndex
		}

		shards[i] = workingIds[currentIndex : currentIndex+shardSize]
		currentIndex += shardSize
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

// calculatePercentageTargets calculates target shard sizes for percentage distribution.
func calculatePercentageTargets(totalCount int, percentages []int64) []int {
	shardCount := len(percentages)
	targets := make([]int, shardCount)
	distributed := 0

	for i, pct := range percentages {
		if i == shardCount-1 {
			targets[i] = totalCount - distributed
		} else {
			targets[i] = int(float64(totalCount) * float64(pct) / 100.0)
			distributed += targets[i]
		}
	}
	return targets
}
