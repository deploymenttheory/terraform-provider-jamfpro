package guid_list_sharder

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Test Data Helpers
// =============================================================================

// generateTestIDs creates a predictable set of IDs for testing
func generateTestIDs(count int) []string {
	ids := make([]string, count)
	for i := range count {
		// Generate predictable IDs: 1, 2, 3, ..., 99, 100
		ids[i] = fmt.Sprintf("%d", i+1)
	}
	return ids
}

// Helper to count total IDs across all shards
func countTotalIDs(shards [][]string) int {
	total := 0
	for _, shard := range shards {
		total += len(shard)
	}
	return total
}

// Helper to check if an ID exists in any shard
func containsID(shards [][]string, id string) bool {
	for _, shard := range shards {
		if slices.Contains(shard, id) {
			return true
		}
	}
	return false
}

// Helper to find which shard contains a specific ID
func findIDShard(shards [][]string, id string) int {
	for i, shard := range shards {
		if slices.Contains(shard, id) {
			return i
		}
	}
	return -1 // Not found
}

// =============================================================================
// createSeededRNG Tests
// =============================================================================

func TestCreateSeededRNG_Deterministic(t *testing.T) {
	seed := "test-seed"

	rng1 := createSeededRNG(seed)
	rng2 := createSeededRNG(seed)

	// Generate some random numbers to verify they're identical
	for i := 0; i < 10; i++ {
		val1 := rng1.Intn(1000)
		val2 := rng2.Intn(1000)
		assert.Equal(t, val1, val2, "Same seed should produce identical random sequences")
	}
}

func TestCreateSeededRNG_DifferentSeeds(t *testing.T) {
	rng1 := createSeededRNG("seed1")
	rng2 := createSeededRNG("seed2")

	// Generate some random numbers - they should be different
	differentCount := 0
	for i := 0; i < 10; i++ {
		val1 := rng1.Intn(1000)
		val2 := rng2.Intn(1000)
		if val1 != val2 {
			differentCount++
		}
	}

	assert.Greater(t, differentCount, 0, "Different seeds should produce different random sequences")
}

// =============================================================================
// shuffle Tests
// =============================================================================

func TestShuffle_EmptyList(t *testing.T) {
	ids := []string{}
	rng := createSeededRNG("test-seed")
	shuffled := shuffle(ids, rng)

	assert.Empty(t, shuffled, "Expected empty list")
}

func TestShuffle_SingleItem(t *testing.T) {
	ids := []string{"1"}
	rng := createSeededRNG("test-seed")
	shuffled := shuffle(ids, rng)

	require.Len(t, shuffled, 1, "Expected 1 item")
	assert.Equal(t, ids[0], shuffled[0], "Single item should remain unchanged")
}

func TestShuffle_MultipleItems(t *testing.T) {
	ids := generateTestIDs(10)
	rng := createSeededRNG("test-seed")
	shuffled := shuffle(ids, rng)

	assert.Len(t, shuffled, len(ids), "Shuffled list should have same length as input")

	// Verify all original items are present (no loss or duplication)
	for _, id := range ids {
		assert.Contains(t, shuffled, id, "Original ID should be in shuffled list")
	}
}

func TestShuffle_DoesNotMutateOriginal(t *testing.T) {
	original := generateTestIDs(10)
	originalCopy := make([]string, len(original))
	copy(originalCopy, original)

	rng := createSeededRNG("test-seed")
	_ = shuffle(original, rng)

	assert.Equal(t, originalCopy, original, "Original slice should not be mutated")
}

func TestShuffle_Deterministic(t *testing.T) {
	ids := generateTestIDs(20)

	rng1 := createSeededRNG("test-seed")
	shuffled1 := shuffle(ids, rng1)

	rng2 := createSeededRNG("test-seed")
	shuffled2 := shuffle(ids, rng2)

	require.Len(t, shuffled1, len(shuffled2), "Shuffled lists should have same length")
	assert.Equal(t, shuffled1, shuffled2, "Same RNG seed should produce identical shuffle order")
}

// =============================================================================
// shardByRoundRobin Tests - Perfect Distribution Verification
// =============================================================================

func TestShardByRoundRobin_EmptyList(t *testing.T) {
	ids := []string{}
	shards := shardByRoundRobin(ids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardByRoundRobin_SingleShard(t *testing.T) {
	ids := generateTestIDs(10)
	shards := shardByRoundRobin(ids, 1, "")

	require.Len(t, shards, 1, "Expected 1 shard")
	assert.Len(t, shards[0], len(ids), "All IDs should be in single shard")
}

// Test perfect distribution (no variance)
func TestShardByRoundRobin_PerfectDistribution(t *testing.T) {
	// Test with exactly divisible count
	ids := generateTestIDs(30)
	shards := shardByRoundRobin(ids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	// Each shard should have exactly 10 IDs
	for i, shard := range shards {
		assert.Len(t, shard, 10, "Shard %d should have exactly 10 IDs", i)
	}
}

// Test perfect distribution with remainder
func TestShardByRoundRobin_PerfectDistribution_WithRemainder(t *testing.T) {
	// 31 IDs / 3 shards = 10, 10, 11
	ids := generateTestIDs(31)
	shards := shardByRoundRobin(ids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}

	// Shards should be within ±1 of each other
	for i := 0; i < len(counts)-1; i++ {
		diff := counts[i] - counts[i+1]
		assert.LessOrEqual(t, abs(diff), 1, "Adjacent shards should differ by at most 1")
	}

	// Total should be 31
	total := countTotalIDs(shards)
	assert.Equal(t, 31, total, "Total should be 31")
}

// Test realistic perfect distribution (512 computers, 3 shards)
func TestShardByRoundRobin_RealisticDistribution_512Computers(t *testing.T) {
	ids := generateTestIDs(512)
	shards := shardByRoundRobin(ids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	// 512 / 3 = 170 remainder 2, so distribution should be: 171, 171, 170
	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}

	t.Logf("Shard counts: %v", counts)

	// All shards should be within 1 of each other
	for i := 0; i < len(counts)-1; i++ {
		diff := counts[i] - counts[i+1]
		assert.LessOrEqual(t, abs(diff), 1, "Shards should differ by at most 1")
	}

	total := countTotalIDs(shards)
	assert.Equal(t, 512, total, "Total should be 512")
}

// Test that WITHOUT seed, order is based on input (API order)
func TestShardByRoundRobin_NoSeed_UsesInputOrder(t *testing.T) {
	ids := generateTestIDs(9)
	shards := shardByRoundRobin(ids, 3, "")

	// Without seed, round-robin uses input order
	// ID 1 → shard 0, ID 2 → shard 1, ID 3 → shard 2, ID 4 → shard 0, ...

	assert.Equal(t, ids[0], shards[0][0], "First ID should be in shard 0, position 0")
	assert.Equal(t, ids[1], shards[1][0], "Second ID should be in shard 1, position 0")
	assert.Equal(t, ids[2], shards[2][0], "Third ID should be in shard 2, position 0")
	assert.Equal(t, ids[3], shards[0][1], "Fourth ID should be in shard 0, position 1")
}

// Test that WITH seed, order is shuffled first, then round-robin
func TestShardByRoundRobin_WithSeed_Deterministic(t *testing.T) {
	ids := generateTestIDs(100)
	seed := "test-seed"

	shards1 := shardByRoundRobin(ids, 3, seed)
	shards2 := shardByRoundRobin(ids, 3, seed)

	// Verify each shard has identical contents with same seed
	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should have identical IDs and order with same seed", i)
	}
}

// Test that different seeds produce different distributions
func TestShardByRoundRobin_DifferentSeeds_DifferentDistributions(t *testing.T) {
	ids := generateTestIDs(100)

	shardsNoSeed := shardByRoundRobin(ids, 3, "")
	shardsSeed1 := shardByRoundRobin(ids, 3, "seed1")
	shardsSeed2 := shardByRoundRobin(ids, 3, "seed2")

	// Count how many IDs are in different shards
	differentFromNoSeed := 0
	differentBetweenSeeds := 0

	for _, id := range ids {
		noSeedShard := findIDShard(shardsNoSeed, id)
		seed1Shard := findIDShard(shardsSeed1, id)
		seed2Shard := findIDShard(shardsSeed2, id)

		if noSeedShard != seed1Shard {
			differentFromNoSeed++
		}
		if seed1Shard != seed2Shard {
			differentBetweenSeeds++
		}
	}

	assert.Greater(t, differentFromNoSeed, 50, "At least 50%% of IDs should be in different shards (no seed vs seed)")
	assert.Greater(t, differentBetweenSeeds, 50, "At least 50%% of IDs should be in different shards (seed1 vs seed2)")
}

// =============================================================================
// shardByPercentage Tests - Precise Percentage Verification
// =============================================================================

func TestShardByPercentage_EmptyList(t *testing.T) {
	ids := []string{}
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(ids, percentages, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

// Test that percentages are accurately applied
func TestShardByPercentage_AccuratePercentages_100Computers(t *testing.T) {
	ids := generateTestIDs(100)
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(ids, percentages, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	// With 100 IDs: 10% = 10, 30% = 30, 60% = 60
	assert.Len(t, shards[0], 10, "Shard 0 should have 10 IDs (10%%)")
	assert.Len(t, shards[1], 30, "Shard 1 should have 30 IDs (30%%)")
	assert.Len(t, shards[2], 60, "Shard 2 should have 60 IDs (60%%)")

	total := countTotalIDs(shards)
	assert.Equal(t, 100, total, "Total should be 100")
}

// Test realistic percentages (512 computers, 10/30/60 split)
func TestShardByPercentage_RealisticPercentages_512Computers(t *testing.T) {
	ids := generateTestIDs(512)
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(ids, percentages, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	// Expected: 10% = 51.2 → 51, 30% = 153.6 → 153, 60% = remainder (308)
	// Last shard gets all remaining

	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}
	t.Logf("Shard counts: %v", counts)

	// Shard 0: ~10% = ~51
	assert.InDelta(t, 51, counts[0], 2, "Shard 0 should have ~51 IDs (10%%)")

	// Shard 1: ~30% = ~153
	assert.InDelta(t, 154, counts[1], 2, "Shard 1 should have ~154 IDs (30%%)")

	// Shard 2: Gets remainder
	expectedRemainder := 512 - counts[0] - counts[1]
	assert.Equal(t, expectedRemainder, counts[2], "Shard 2 should get all remaining IDs")

	total := countTotalIDs(shards)
	assert.Equal(t, 512, total, "Total should be 512")
}

// Test that last shard gets all remaining IDs (no loss)
func TestShardByPercentage_LastShardGetsRemainder(t *testing.T) {
	ids := generateTestIDs(103) // Odd number to ensure remainder
	percentages := []int64{10, 20, 70}
	shards := shardByPercentage(ids, percentages, "")

	// Verify all IDs are accounted for
	total := countTotalIDs(shards)
	assert.Equal(t, 103, total, "All 103 IDs should be distributed")

	// Verify each ID appears exactly once
	for _, id := range ids {
		assert.True(t, containsID(shards, id), "ID should be in a shard")
	}
}

// Test that WITHOUT seed, order is based on input
func TestShardByPercentage_NoSeed_UsesInputOrder(t *testing.T) {
	ids := generateTestIDs(10)
	percentages := []int64{20, 30, 50}
	shards := shardByPercentage(ids, percentages, "")

	// Without seed: first 20% (2 IDs) go to shard 0, next 30% (3 IDs) to shard 1, rest to shard 2

	// Shard 0 should have first 2 IDs
	assert.Contains(t, shards[0], ids[0], "Shard 0 should contain first ID")
	assert.Contains(t, shards[0], ids[1], "Shard 0 should contain second ID")
	assert.Len(t, shards[0], 2, "Shard 0 should have 2 IDs")

	// Shard 1 should have next 3 IDs
	assert.Contains(t, shards[1], ids[2], "Shard 1 should contain third ID")
	assert.Len(t, shards[1], 3, "Shard 1 should have 3 IDs")

	// Shard 2 should have remaining 5 IDs
	assert.Len(t, shards[2], 5, "Shard 2 should have 5 IDs")
}

// Test that WITH seed is deterministic
func TestShardByPercentage_WithSeed_Deterministic(t *testing.T) {
	ids := generateTestIDs(100)
	percentages := []int64{10, 30, 60}
	seed := "test-seed"

	shards1 := shardByPercentage(ids, percentages, seed)
	shards2 := shardByPercentage(ids, percentages, seed)

	// Verify each shard has identical contents with same seed
	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should have identical IDs and order with same seed", i)
	}
}

// =============================================================================
// shardBySize Tests
// =============================================================================

func TestShardBySize_EmptyList(t *testing.T) {
	ids := []string{}
	sizes := []int64{10, 20, -1}
	shards := shardBySize(ids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardBySize_ExactSizes(t *testing.T) {
	ids := generateTestIDs(100)
	sizes := []int64{50, 30, 20}
	shards := shardBySize(ids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 50, "Shard 0 should have 50 IDs")
	assert.Len(t, shards[1], 30, "Shard 1 should have 30 IDs")
	assert.Len(t, shards[2], 20, "Shard 2 should have 20 IDs")

	total := countTotalIDs(shards)
	assert.Equal(t, 100, total, "All 100 IDs should be distributed")
}

func TestShardBySize_WithRemainder(t *testing.T) {
	ids := generateTestIDs(1000)
	sizes := []int64{50, 200, -1}
	shards := shardBySize(ids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 50, "Shard 0 should have 50 IDs")
	assert.Len(t, shards[1], 200, "Shard 1 should have 200 IDs")
	assert.Len(t, shards[2], 750, "Shard 2 should have 750 IDs (remainder)")

	total := countTotalIDs(shards)
	assert.Equal(t, 1000, total, "All 1000 IDs should be distributed")
}

func TestShardBySize_ZeroRemainder(t *testing.T) {
	// Edge case: sizes sum to exactly total IDs, leaving 0 for -1 shard
	ids := generateTestIDs(30)
	sizes := []int64{10, 20, -1} // 10+20=30, so 0 remaining
	shards := shardBySize(ids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 10, "Shard 0 should have 10 IDs")
	assert.Len(t, shards[1], 20, "Shard 1 should have 20 IDs")

	// Critical: must be empty slice []string{}, NOT nil
	// nil becomes null in Terraform state, breaking HCL expressions like length()
	assert.NotNil(t, shards[2], "Shard 2 must not be nil (would become null in Terraform)")
	assert.Len(t, shards[2], 0, "Shard 2 should have 0 IDs (exact match, no remainder)")

	total := countTotalIDs(shards)
	assert.Equal(t, 30, total, "All 30 IDs should be distributed")
}

// Helper function for absolute value
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
