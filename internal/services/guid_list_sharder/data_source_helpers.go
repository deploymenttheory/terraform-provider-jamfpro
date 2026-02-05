package guid_list_sharder

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
	"slices"
	"strconv"
)

// =============================================================================
// Sharding Strategy Implementations
// =============================================================================

// shardByRoundRobin distributes IDs in circular order, guaranteeing equal shard sizes
// Without seed: uses API order (non-deterministic, may change between runs)
// With seed: sorts first for stable input, then shuffles using Fisher-Yates (deterministic, reproducible)
func shardByRoundRobin(ids []string, shardCount int, seed string) [][]string {

	if shardCount <= 0 {
		shardCount = 1
	}

	shards := make([][]string, shardCount)

	// Use deterministic shuffle if seed provided
	workingIds := ids
	if seed != "" {
		// Sort numerically first to ensure consistent input order, then shuffle
		sortedIds := make([]string, len(ids))
		copy(sortedIds, ids)
		slices.SortFunc(sortedIds, func(a, b string) int {
			aInt, _ := strconv.Atoi(a)
			bInt, _ := strconv.Atoi(b)
			return aInt - bInt
		})
		workingIds = shuffleWithSeed(sortedIds, seed)
	}

	// Distribute in round-robin fashion
	for i, id := range workingIds {
		shardIndex := i % shardCount
		shards[shardIndex] = append(shards[shardIndex], id)
	}

	// Sort each shard numerically to match how Jamf API returns them
	// This prevents Terraform from seeing spurious diffs
	for i := range shards {
		slices.SortFunc(shards[i], func(a, b string) int {
			aInt, _ := strconv.Atoi(a)
			bInt, _ := strconv.Atoi(b)
			return aInt - bInt
		})
	}

	return shards
}

// shardByPercentage distributes IDs according to specified percentages
// Without seed: uses API order (non-deterministic, may change between runs)
// With seed: sorts first for stable input, then shuffles using Fisher-Yates (deterministic, reproducible)
func shardByPercentage(ids []string, percentages []int64, seed string) [][]string {
	totalIds := len(ids)
	shardCount := len(percentages)
	shards := make([][]string, shardCount)

	if totalIds == 0 {
		return shards
	}

	// Use deterministic shuffle if seed provided
	workingIds := ids
	if seed != "" {
		// Sort numerically first to ensure consistent input order, then shuffle
		sortedIds := make([]string, len(ids))
		copy(sortedIds, ids)
		slices.SortFunc(sortedIds, func(a, b string) int {
			aInt, _ := strconv.Atoi(a)
			bInt, _ := strconv.Atoi(b)
			return aInt - bInt
		})
		workingIds = shuffleWithSeed(sortedIds, seed)
	}

	// Distribute by percentages
	currentIndex := 0
	for i, percentage := range percentages {
		var shardSize int
		if i == shardCount-1 {
			// Last shard gets all remaining IDs
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

	// Sort each shard numerically to match how Jamf API returns them
	// This prevents Terraform from seeing spurious diffs
	for i := range shards {
		slices.SortFunc(shards[i], func(a, b string) int {
			aInt, _ := strconv.Atoi(a)
			bInt, _ := strconv.Atoi(b)
			return aInt - bInt
		})
	}

	return shards
}

// shardBySize distributes IDs according to specified absolute sizes
// Without seed: uses API order (non-deterministic, may change between runs)
// With seed: sorts first for stable input, then shuffles using Fisher-Yates (deterministic, reproducible)
// Supports -1 in the last position to mean "all remaining IDs"
func shardBySize(ids []string, sizes []int64, seed string) [][]string {
	totalIds := len(ids)
	shardCount := len(sizes)
	shards := make([][]string, shardCount)

	if totalIds == 0 {
		return shards
	}

	// Use deterministic shuffle if seed provided
	workingIds := ids
	if seed != "" {
		// Sort numerically first to ensure consistent input order, then shuffle
		sortedIds := make([]string, len(ids))
		copy(sortedIds, ids)
		slices.SortFunc(sortedIds, func(a, b string) int {
			aInt, _ := strconv.Atoi(a)
			bInt, _ := strconv.Atoi(b)
			return aInt - bInt
		})
		workingIds = shuffleWithSeed(sortedIds, seed)
	}

	// Distribute by sizes
	currentIndex := 0
	for i, size := range sizes {
		var shardSize int

		if size == -1 {
			// -1 means "all remaining IDs"
			shardSize = totalIds - currentIndex
		} else {
			shardSize = int(size)

			// If we don't have enough IDs left, take what's available
			if currentIndex+shardSize > totalIds {
				shardSize = totalIds - currentIndex
			}
		}

		// Always initialize shard, even if empty
		// Why: nil slices become null in Terraform state, breaking HCL expressions like length()
		// Empty slices []string{} become empty sets (length 0) which work correctly in HCL
		if shardSize > 0 && currentIndex < totalIds {
			shards[i] = workingIds[currentIndex : currentIndex+shardSize]
			currentIndex += shardSize
		} else {
			shards[i] = []string{}
		}
	}

	// Sort each shard numerically to match how Jamf API returns them
	// This prevents Terraform from seeing spurious diffs
	for i := range shards {
		slices.SortFunc(shards[i], func(a, b string) int {
			aInt, _ := strconv.Atoi(a)
			bInt, _ := strconv.Atoi(b)
			return aInt - bInt
		})
	}

	return shards
}

// =============================================================================
// Seeding and Shuffle Helpers
// =============================================================================

// createSeededRNG creates a deterministic random number generator from a seed string
// Uses SHA-256 to convert string seed into int64 for reproducible randomization
func createSeededRNG(seed string) *rand.Rand {
	hash := sha256.Sum256([]byte(seed))
	seedValue := int64(binary.BigEndian.Uint64(hash[:8]))
	return rand.New(rand.NewSource(seedValue))
}

// shuffle performs Fisher-Yates shuffle on a copy of the input slice using provided RNG
// Returns shuffled copy without mutating the original slice
func shuffle(ids []string, rng *rand.Rand) []string {
	// Create a copy to avoid mutating original slice
	shuffled := make([]string, len(ids))
	copy(shuffled, ids)

	// Fisher-Yates shuffle algorithm
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

// shuffleWithSeed combines seeding and shuffling for convenience
// Used by round-robin and percentage strategies when seed is provided for reproducible randomization
func shuffleWithSeed(ids []string, seed string) []string {
	rng := createSeededRNG(seed)
	return shuffle(ids, rng)
}

// shardByRendezvous distributes IDs using Highest Random Weight (HRW) algorithm
// Each ID computes a score for every shard and is assigned to the shard with the highest score
// This provides superior stability when shard counts change - only ~1/n IDs move when adding a shard
// Always deterministic (reproducible across runs) - seed affects which shard wins for each ID
func shardByRendezvous(ids []string, shardCount int, seed string) [][]string {
	if shardCount <= 0 {
		shardCount = 1
	}

	shards := make([][]string, shardCount)

	// Initialize all shards as empty slices to prevent nil
	// Why: nil slices become null in Terraform state, breaking HCL expressions like length()
	// Empty slices []string{} become empty sets (length 0) which work correctly in HCL
	for i := 0; i < shardCount; i++ {
		shards[i] = []string{}
	}

	// For each ID, compute weight for every shard and assign to highest
	for _, id := range ids {
		highestWeight := uint64(0)
		selectedShard := 0

		// Evaluate this ID against all shards
		for shardIdx := 0; shardIdx < shardCount; shardIdx++ {
			// Combine ID + shard identifier + seed for deterministic weight
			// Format: "id:shard_N:seed" ensures each ID-shard pair gets unique hash
			input := fmt.Sprintf("%s:shard_%d:%s", id, shardIdx, seed)
			hash := sha256.Sum256([]byte(input))

			// Use first 8 bytes of hash as weight (uint64 for large range)
			weight := binary.BigEndian.Uint64(hash[:8])

			// Track shard with highest weight for this ID
			if weight > highestWeight {
				highestWeight = weight
				selectedShard = shardIdx
			}
		}

		shards[selectedShard] = append(shards[selectedShard], id)
	}

	// Sort each shard numerically to match how Jamf API returns them
	// This prevents Terraform from seeing spurious diffs
	for i := range shards {
		slices.SortFunc(shards[i], func(a, b string) int {
			aInt, _ := strconv.Atoi(a)
			bInt, _ := strconv.Atoi(b)
			return aInt - bInt
		})
	}

	return shards
}
