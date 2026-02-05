package guid_list_sharder

import (
	"context"
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
		ids[i] = fmt.Sprintf("%d", i+1)
	}
	return ids
}

// countTotalIDs counts total IDs across all shards
func countTotalIDs(shards [][]string) int {
	total := 0
	for _, shard := range shards {
		total += len(shard)
	}
	return total
}

// containsID checks if an ID exists in any shard
func containsID(shards [][]string, id string) bool {
	for _, shard := range shards {
		if slices.Contains(shard, id) {
			return true
		}
	}
	return false
}

// findIDShard finds which shard contains a specific ID
func findIDShard(shards [][]string, id string) int {
	for i, shard := range shards {
		if slices.Contains(shard, id) {
			return i
		}
	}
	return -1
}

// abs returns the absolute value
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// =============================================================================
// Round-Robin Strategy Tests
// =============================================================================

func TestShardByRoundRobin_EmptyList(t *testing.T) {
	ids := []string{}
	shards := shardByRoundRobin(context.Background(), ids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardByRoundRobin_SingleShard(t *testing.T) {
	ids := generateTestIDs(10)
	shards := shardByRoundRobin(context.Background(), ids, 1, "")

	require.Len(t, shards, 1, "Expected 1 shard")
	assert.Len(t, shards[0], len(ids), "All IDs should be in single shard")
}

// Test perfect distribution (no variance)
func TestShardByRoundRobin_PerfectDistribution(t *testing.T) {
	ids := generateTestIDs(30)
	shards := shardByRoundRobin(context.Background(), ids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	for i, shard := range shards {
		assert.Len(t, shard, 10, "Shard %d should have exactly 10 IDs", i)
	}
}

func TestShardByRoundRobin_WithRemainder(t *testing.T) {
	ids := generateTestIDs(31)
	shards := shardByRoundRobin(context.Background(), ids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}

	for i := 0; i < len(counts)-1; i++ {
		diff := counts[i] - counts[i+1]
		assert.LessOrEqual(t, abs(diff), 1, "Adjacent shards should differ by at most 1")
	}

	total := countTotalIDs(shards)
	assert.Equal(t, 31, total, "Total should be 31")
}

func TestShardByRoundRobin_LargeDataset(t *testing.T) {
	ids := generateTestIDs(512)
	shards := shardByRoundRobin(context.Background(), ids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}

	for i := 0; i < len(counts)-1; i++ {
		diff := counts[i] - counts[i+1]
		assert.LessOrEqual(t, abs(diff), 1, "Shards should differ by at most 1")
	}

	total := countTotalIDs(shards)
	assert.Equal(t, 512, total, "Total should be 512")
}

func TestShardByRoundRobin_NoSeed_UsesInputOrder(t *testing.T) {
	ids := generateTestIDs(9)
	shards := shardByRoundRobin(context.Background(), ids, 3, "")

	assert.Equal(t, ids[0], shards[0][0], "First ID should be in shard 0")
	assert.Equal(t, ids[1], shards[1][0], "Second ID should be in shard 1")
	assert.Equal(t, ids[2], shards[2][0], "Third ID should be in shard 2")
	assert.Equal(t, ids[3], shards[0][1], "Fourth ID should be in shard 0")
}

func TestShardByRoundRobin_WithSeed_Deterministic(t *testing.T) {
	ids := generateTestIDs(100)
	seed := "test-seed"

	shards1 := shardByRoundRobin(context.Background(), ids, 3, seed)
	shards2 := shardByRoundRobin(context.Background(), ids, 3, seed)

	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should be identical with same seed", i)
	}
}

func TestShardByRoundRobin_DifferentSeeds(t *testing.T) {
	ids := generateTestIDs(100)

	shardsSeed1 := shardByRoundRobin(context.Background(), ids, 3, "seed1")
	shardsSeed2 := shardByRoundRobin(context.Background(), ids, 3, "seed2")

	differentBetweenSeeds := 0
	for _, id := range ids {
		seed1Shard := findIDShard(shardsSeed1, id)
		seed2Shard := findIDShard(shardsSeed2, id)

		if seed1Shard != seed2Shard {
			differentBetweenSeeds++
		}
	}

	assert.Greater(t, differentBetweenSeeds, 50, "Different seeds should produce different distributions")
}

// =============================================================================
// Percentage Strategy Tests
// =============================================================================

func TestShardByPercentage_EmptyList(t *testing.T) {
	ids := []string{}
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(context.Background(), ids, percentages, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardByPercentage_AccuratePercentages(t *testing.T) {
	ids := generateTestIDs(100)
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(context.Background(), ids, percentages, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	assert.Len(t, shards[0], 10, "Shard 0 should have 10 IDs")
	assert.Len(t, shards[1], 30, "Shard 1 should have 30 IDs")
	assert.Len(t, shards[2], 60, "Shard 2 should have 60 IDs")

	total := countTotalIDs(shards)
	assert.Equal(t, 100, total, "Total should be 100")
}

func TestShardByPercentage_LargeDataset(t *testing.T) {
	ids := generateTestIDs(512)
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(context.Background(), ids, percentages, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}

	assert.InDelta(t, 51, counts[0], 2, "Shard 0 should have ~10%")
	assert.InDelta(t, 154, counts[1], 2, "Shard 1 should have ~30%")

	expectedRemainder := 512 - counts[0] - counts[1]
	assert.Equal(t, expectedRemainder, counts[2], "Last shard gets remainder")

	total := countTotalIDs(shards)
	assert.Equal(t, 512, total, "Total should be 512")
}

func TestShardByPercentage_LastShardGetsRemainder(t *testing.T) {
	ids := generateTestIDs(103)
	percentages := []int64{10, 20, 70}
	shards := shardByPercentage(context.Background(), ids, percentages, "")

	total := countTotalIDs(shards)
	assert.Equal(t, 103, total, "All IDs should be distributed")

	for _, id := range ids {
		assert.True(t, containsID(shards, id), "Each ID should appear exactly once")
	}
}

func TestShardByPercentage_NoSeed_UsesInputOrder(t *testing.T) {
	ids := generateTestIDs(10)
	percentages := []int64{20, 30, 50}
	shards := shardByPercentage(context.Background(), ids, percentages, "")

	assert.Contains(t, shards[0], ids[0])
	assert.Contains(t, shards[0], ids[1])
	assert.Len(t, shards[0], 2)

	assert.Contains(t, shards[1], ids[2])
	assert.Len(t, shards[1], 3)

	assert.Len(t, shards[2], 5)
}

func TestShardByPercentage_WithSeed_Deterministic(t *testing.T) {
	ids := generateTestIDs(100)
	percentages := []int64{10, 30, 60}
	seed := "test-seed"

	shards1 := shardByPercentage(context.Background(), ids, percentages, seed)
	shards2 := shardByPercentage(context.Background(), ids, percentages, seed)

	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should be identical with same seed", i)
	}
}

// =============================================================================
// Size Strategy Tests
// =============================================================================

func TestShardBySize_EmptyList(t *testing.T) {
	ids := []string{}
	sizes := []int64{10, 20, -1}
	shards := shardBySize(context.Background(), ids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardBySize_ExactSizes(t *testing.T) {
	ids := generateTestIDs(100)
	sizes := []int64{50, 30, 20}
	shards := shardBySize(context.Background(), ids, sizes, "")

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
	shards := shardBySize(context.Background(), ids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 50, "Shard 0 should have 50 IDs")
	assert.Len(t, shards[1], 200, "Shard 1 should have 200 IDs")
	assert.Len(t, shards[2], 750, "Shard 2 should have 750 IDs (remainder)")

	total := countTotalIDs(shards)
	assert.Equal(t, 1000, total, "All 1000 IDs should be distributed")
}

func TestShardBySize_ZeroRemainder(t *testing.T) {
	ids := generateTestIDs(30)
	sizes := []int64{10, 20, -1}
	shards := shardBySize(context.Background(), ids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 10, "Shard 0 should have 10 IDs")
	assert.Len(t, shards[1], 20, "Shard 1 should have 20 IDs")

	assert.NotNil(t, shards[2], "Shard 2 must not be nil")
	assert.Len(t, shards[2], 0, "Shard 2 should have 0 IDs")

	total := countTotalIDs(shards)
	assert.Equal(t, 30, total, "All IDs should be distributed")
}
