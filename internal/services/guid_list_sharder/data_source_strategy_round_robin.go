package guid_list_sharder

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// shardByRoundRobin distributes IDs in circular order, guaranteeing equal shard sizes ±1.
func shardByRoundRobin(ctx context.Context, ids []string, shardCount int, seed string) [][]string {
	if shardCount <= 0 {
		shardCount = 1
	}

	shards := make([][]string, shardCount)
	workingIds := prepareIDsForDistribution(ids, seed)

	for i, id := range workingIds {
		shardIndex := i % shardCount
		shards[shardIndex] = append(shards[shardIndex], id)
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

// calculateRoundRobinTargets calculates target shard sizes for round-robin distribution.
func calculateRoundRobinTargets(totalCount int, shardCount int) []int {
	targets := make([]int, shardCount)
	baseSize := totalCount / shardCount
	remainder := totalCount % shardCount

	for i := 0; i < shardCount; i++ {
		targets[i] = baseSize
		if i < remainder {
			targets[i]++
		}
	}
	return targets
}
