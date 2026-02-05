package guid_list_sharder

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// distributeToTargets distributes IDs to shards according to target counts.
// This ensures the final distribution (including reserved IDs) matches the strategy.
func distributeToTargets(ctx context.Context, ids []string, targets []int, seed string) [][]string {
	shardCount := len(targets)
	shards := make([][]string, shardCount)

	for i := range shards {
		shards[i] = make([]string, 0, targets[i])
	}

	idsToDistribute := make([]string, len(ids))
	copy(idsToDistribute, ids)

	if seed != "" {
		rng := createSeededRNG(seed)
		rng.Shuffle(len(idsToDistribute), func(i, j int) {
			idsToDistribute[i], idsToDistribute[j] = idsToDistribute[j], idsToDistribute[i]
		})
	}

	currentShard := 0
	for _, id := range idsToDistribute {
		attempts := 0
		for len(shards[currentShard]) >= targets[currentShard] {
			currentShard = (currentShard + 1) % shardCount
			attempts++
			if attempts >= shardCount {
				tflog.Warn(ctx, "All shards reached targets but IDs remain. Adding to shard 0.")
				shards[0] = append(shards[0], id)
				break
			}
		}

		if attempts < shardCount {
			shards[currentShard] = append(shards[currentShard], id)
			currentShard = (currentShard + 1) % shardCount
		}
	}

	for i := range shards {
		sortIDsNumerically(shards[i])
	}

	return shards
}
