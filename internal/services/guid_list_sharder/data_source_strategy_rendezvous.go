package guid_list_sharder

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// shardByRendezvous distributes IDs using Highest Random Weight (HRW) algorithm.
// Provides superior stability when shard counts change - only ~1/n IDs move when adding a shard.
func shardByRendezvous(ctx context.Context, ids []string, shardCount int, seed string) [][]string {
	if shardCount <= 0 {
		shardCount = 1
	}

	shards := make([][]string, shardCount)

	// Initialize all shards as empty slices to prevent nil in Terraform state
	for i := 0; i < shardCount; i++ {
		shards[i] = []string{}
	}

	for _, id := range ids {
		highestWeight := uint64(0)
		selectedShard := 0

		for shardIdx := 0; shardIdx < shardCount; shardIdx++ {
			input := fmt.Sprintf("%s:shard_%d:%s", id, shardIdx, seed)
			hash := sha256.Sum256([]byte(input))
			weight := binary.BigEndian.Uint64(hash[:8])

			if weight > highestWeight {
				highestWeight = weight
				selectedShard = shardIdx
			}
		}

		shards[selectedShard] = append(shards[selectedShard], id)
	}

	for i := range shards {
		if len(shards[i]) > 0 {
			tflog.Debug(ctx, "Rendezvous: sorting shard", map[string]any{
				"shard": i,
				"size":  len(shards[i]),
			})
		}
		sortIDsNumerically(shards[i])
	}

	return shards
}
