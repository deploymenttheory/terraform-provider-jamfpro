package guid_list_sharder

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
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

	distributionIDs := sortAndShuffleIfSeed(unreservedIDs, seed)

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
			shards[i] = distributionIDs[currentIndex : currentIndex+shardSize]
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
	distributionIDs := sortAndShuffleIfSeed(unreservedIDs, seed)

	for i, id := range distributionIDs {
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

	distributionIDs := sortAndShuffleIfSeed(unreservedIDs, seed)

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
			shards[i] = distributionIDs[currentIndex : currentIndex+shardSize]
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
