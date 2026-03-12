package guid_list_sharder

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
)

type Strategy string

const (
	StrategyRoundRobin Strategy = "round-robin"
	StrategyPercentage Strategy = "percentage"
	StrategySize       Strategy = "size"
	StrategyRendezvous Strategy = "rendezvous"
)

type Reservations struct {
	IDsByShard    map[string][]string
	CountsByShard map[int]int
	UnreservedIDs []string
}

// Shard splits ids into shard slices using the selected strategy.
//
// Notes:
// - ids are stringified numeric IDs (Jamf uses integer IDs).
// - if Seed is empty, round-robin/percentage/size preserve input order.
// - rendezvous is always deterministic for a given (id, shardCount, seed).
func Shard(ids []string, strategy Strategy, shardCount int, shardPercentages []int, shardSizes []int, seed string, reservations *Reservations) ([][]string, error) {
	switch strategy {
	case StrategyRoundRobin:
		return shardByRoundRobin(ids, shardCount, seed, reservations), nil
	case StrategyRendezvous:
		return shardByRendezvous(ids, shardCount, seed, reservations), nil
	case StrategyPercentage:
		return shardByPercentage(ids, shardPercentages, seed, reservations), nil
	case StrategySize:
		return shardBySize(ids, shardSizes, seed, reservations), nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownStrategy, strategy)
	}
}

func shardByRoundRobin(ids []string, shardCount int, seed string, reservations *Reservations) [][]string {
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
		shards[i%shardCount] = append(shards[i%shardCount], id)
	}

	applyReservationsToShards(shards, reservations)
	for i := range shards {
		sortIDsNumerically(shards[i])
	}

	return shards
}

func shardByPercentage(ids []string, percentages []int, seed string, reservations *Reservations) [][]string {
	unreservedIDs := ids
	totalIDs := len(ids)
	if reservations != nil {
		unreservedIDs = reservations.UnreservedIDs
	}

	shardCount := len(percentages)
	shards := make([][]string, shardCount)
	if len(unreservedIDs) == 0 {
		applyReservationsToShards(shards, reservations)
		return shards
	}

	distributionIDs := sortAndShuffleIfSeed(unreservedIDs, seed)

	currentIndex := 0
	for i, percentage := range percentages {
		var shardSize int
		if i == shardCount-1 {
			shardSize = len(unreservedIDs) - currentIndex
		} else {
			shardSize = int(float64(totalIDs) * float64(percentage) / 100.0)
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

	applyReservationsToShards(shards, reservations)
	for i := range shards {
		sortIDsNumerically(shards[i])
	}

	return shards
}

func shardBySize(ids []string, sizes []int, seed string, reservations *Reservations) [][]string {
	unreservedIDs := ids
	if reservations != nil {
		unreservedIDs = reservations.UnreservedIDs
	}

	shardCount := len(sizes)
	shards := make([][]string, shardCount)
	if len(unreservedIDs) == 0 {
		applyReservationsToShards(shards, reservations)
		return shards
	}

	distributionIDs := sortAndShuffleIfSeed(unreservedIDs, seed)

	currentIndex := 0
	for i, size := range sizes {
		var shardSize int
		if size == -1 {
			shardSize = len(unreservedIDs) - currentIndex
		} else {
			shardSize = size
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

	applyReservationsToShards(shards, reservations)
	for i := range shards {
		sortIDsNumerically(shards[i])
	}

	return shards
}

func shardByRendezvous(ids []string, shardCount int, seed string, reservations *Reservations) [][]string {
	if shardCount <= 0 {
		shardCount = 1
	}

	unreservedIDs := ids
	if reservations != nil {
		unreservedIDs = reservations.UnreservedIDs
	}

	shards := make([][]string, shardCount)

	for _, id := range unreservedIDs {
		highestWeight := uint64(0)
		selectedShard := 0

		for shardIdx := range shards {
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

	applyReservationsToShards(shards, reservations)
	for i := range shards {
		sortIDsNumerically(shards[i])
	}

	return shards
}

func applyReservationsToShards(shards [][]string, reservations *Reservations) {
	if reservations == nil {
		return
	}
	for shardName, reservedIDs := range reservations.IDsByShard {
		var idx int
		if _, err := fmt.Sscanf(shardName, "shard_%d", &idx); err != nil {
			continue
		}
		if idx < 0 || idx >= len(shards) {
			continue
		}
		shards[idx] = append(reservedIDs, shards[idx]...)
	}
}

func sortAndShuffleIfSeed(ids []string, seed string) []string {
	if seed == "" {
		return ids
	}

	sorted := make([]string, len(ids))
	copy(sorted, ids)
	sortIDsNumerically(sorted)
	return shuffleIDs(sorted, seed)
}

func shuffleIDs(ids []string, seed string) []string {
	rng := createSeededRNG(seed)
	shuffled := make([]string, len(ids))
	copy(shuffled, ids)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

func createSeededRNG(seed string) *rand.Rand {
	hash := sha256.Sum256([]byte(seed))
	seedValue := int64(binary.BigEndian.Uint64(hash[:8]) & math.MaxInt64)
	// #nosec G404 -- deterministic shuffling for stable rollouts (not security sensitive)
	return rand.New(rand.NewSource(seedValue))
}

func sortIDsNumerically(ids []string) {
	slices.SortFunc(ids, func(a, b string) int {
		aInt, _ := strconv.Atoi(a)
		bInt, _ := strconv.Atoi(b)
		return aInt - bInt
	})
}
