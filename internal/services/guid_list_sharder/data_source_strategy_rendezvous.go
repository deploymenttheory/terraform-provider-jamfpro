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
// If reservations are provided, reserved IDs are placed first in their designated shards,
// then unreserved IDs are distributed using the rendezvous hashing algorithm.
func shardByRendezvous(ctx context.Context, ids []string, shardCount int, seed string, reservations *reservationInfo) [][]string {
	if shardCount <= 0 {
		shardCount = 1
	}

	unreservedIDs := ids
	if reservations != nil {
		unreservedIDs = reservations.UnreservedIDs
	}

	shards := make([][]string, shardCount)
	for i := range shardCount {
		shards[i] = []string{}
	}

	for _, id := range unreservedIDs {
		highestWeight := uint64(0)
		selectedShard := 0

		for shardIdx := range shardCount {
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

	if reservations != nil {
		for shardName, reservedIDs := range reservations.IDsByShard {
			var shardIndex int
			fmt.Sscanf(shardName, "shard_%d", &shardIndex)
			shards[shardIndex] = append(reservedIDs, shards[shardIndex]...)
		}
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
