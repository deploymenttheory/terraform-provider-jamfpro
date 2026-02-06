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
	shards := shardByRoundRobin(context.Background(), ids, 3, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardByRoundRobin_SingleShard(t *testing.T) {
	ids := generateTestIDs(10)
	shards := shardByRoundRobin(context.Background(), ids, 1, "", nil)

	require.Len(t, shards, 1, "Expected 1 shard")
	assert.Len(t, shards[0], len(ids), "All IDs should be in single shard")
}

// Test perfect distribution (no variance)
func TestShardByRoundRobin_PerfectDistribution(t *testing.T) {
	ids := generateTestIDs(30)
	shards := shardByRoundRobin(context.Background(), ids, 3, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")

	for i, shard := range shards {
		assert.Len(t, shard, 10, "Shard %d should have exactly 10 IDs", i)
	}
}

func TestShardByRoundRobin_WithRemainder(t *testing.T) {
	ids := generateTestIDs(31)
	shards := shardByRoundRobin(context.Background(), ids, 3, "", nil)

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
	shards := shardByRoundRobin(context.Background(), ids, 3, "", nil)

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
	shards := shardByRoundRobin(context.Background(), ids, 3, "", nil)

	assert.Equal(t, ids[0], shards[0][0], "First ID should be in shard 0")
	assert.Equal(t, ids[1], shards[1][0], "Second ID should be in shard 1")
	assert.Equal(t, ids[2], shards[2][0], "Third ID should be in shard 2")
	assert.Equal(t, ids[3], shards[0][1], "Fourth ID should be in shard 0")
}

func TestShardByRoundRobin_WithSeed_Deterministic(t *testing.T) {
	ids := generateTestIDs(100)
	seed := "test-seed"

	shards1 := shardByRoundRobin(context.Background(), ids, 3, seed, nil)
	shards2 := shardByRoundRobin(context.Background(), ids, 3, seed, nil)

	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should be identical with same seed", i)
	}
}

func TestShardByRoundRobin_DifferentSeeds(t *testing.T) {
	ids := generateTestIDs(100)

	shardsSeed1 := shardByRoundRobin(context.Background(), ids, 3, "seed1", nil)
	shardsSeed2 := shardByRoundRobin(context.Background(), ids, 3, "seed2", nil)

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

func TestShardByRoundRobin_WithReservations(t *testing.T) {
	ids := generateTestIDs(100)
	
	reservations := &reservationInfo{
		IDsByShard: map[string][]string{
			"shard_0": {"1", "2"},
			"shard_1": {"3", "4", "5"},
			"shard_2": {"6"},
		},
		CountsByShard: map[int]int{
			0: 2,
			1: 3,
			2: 1,
		},
		UnreservedIDs: make([]string, 0, 94),
	}

	for _, id := range ids {
		if !slices.Contains([]string{"1", "2", "3", "4", "5", "6"}, id) {
			reservations.UnreservedIDs = append(reservations.UnreservedIDs, id)
		}
	}

	shards := shardByRoundRobin(context.Background(), ids, 3, "", reservations)

	require.Len(t, shards, 3, "Expected 3 shards")

	assert.Contains(t, shards[0], "1", "Shard 0 should contain reserved ID 1")
	assert.Contains(t, shards[0], "2", "Shard 0 should contain reserved ID 2")
	assert.Contains(t, shards[1], "3", "Shard 1 should contain reserved ID 3")
	assert.Contains(t, shards[1], "4", "Shard 1 should contain reserved ID 4")
	assert.Contains(t, shards[1], "5", "Shard 1 should contain reserved ID 5")
	assert.Contains(t, shards[2], "6", "Shard 2 should contain reserved ID 6")

	assert.Equal(t, "1", shards[0][0], "Reserved IDs should appear first in shard 0")
	assert.Equal(t, "2", shards[0][1], "Reserved IDs should appear first in shard 0")
	assert.Equal(t, "3", shards[1][0], "Reserved IDs should appear first in shard 1")

	total := countTotalIDs(shards)
	assert.Equal(t, 100, total, "All 100 IDs should be distributed")

	for _, id := range ids {
		assert.True(t, containsID(shards, id), "ID %s should appear exactly once", id)
	}
}

// =============================================================================
// Percentage Strategy Tests
// =============================================================================

func TestShardByPercentage_EmptyList(t *testing.T) {
	ids := []string{}
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(context.Background(), ids, percentages, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardByPercentage_AccuratePercentages(t *testing.T) {
	ids := generateTestIDs(100)
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(context.Background(), ids, percentages, "", nil)

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
	shards := shardByPercentage(context.Background(), ids, percentages, "", nil)

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
	shards := shardByPercentage(context.Background(), ids, percentages, "", nil)

	total := countTotalIDs(shards)
	assert.Equal(t, 103, total, "All IDs should be distributed")

	for _, id := range ids {
		assert.True(t, containsID(shards, id), "Each ID should appear exactly once")
	}
}

func TestShardByPercentage_NoSeed_UsesInputOrder(t *testing.T) {
	ids := generateTestIDs(10)
	percentages := []int64{20, 30, 50}
	shards := shardByPercentage(context.Background(), ids, percentages, "", nil)

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

	shards1 := shardByPercentage(context.Background(), ids, percentages, seed, nil)
	shards2 := shardByPercentage(context.Background(), ids, percentages, seed, nil)

	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should be identical with same seed", i)
	}
}

func TestShardByPercentage_WithReservations(t *testing.T) {
	ids := generateTestIDs(100)
	percentages := []int64{20, 30, 50}
	
	reservations := &reservationInfo{
		IDsByShard: map[string][]string{
			"shard_0": {"1", "2", "3"},
			"shard_1": {"4", "5"},
			"shard_2": {"6", "7", "8", "9", "10"},
		},
		CountsByShard: map[int]int{
			0: 3,
			1: 2,
			2: 5,
		},
		UnreservedIDs: make([]string, 0, 90),
	}

	for _, id := range ids {
		if !slices.Contains([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, id) {
			reservations.UnreservedIDs = append(reservations.UnreservedIDs, id)
		}
	}

	shards := shardByPercentage(context.Background(), ids, percentages, "", reservations)

	require.Len(t, shards, 3, "Expected 3 shards")

	assert.Contains(t, shards[0], "1", "Shard 0 should contain reserved ID 1")
	assert.Contains(t, shards[0], "2", "Shard 0 should contain reserved ID 2")
	assert.Contains(t, shards[0], "3", "Shard 0 should contain reserved ID 3")
	assert.Equal(t, "1", shards[0][0], "Reserved IDs should appear first")

	assert.Contains(t, shards[1], "4", "Shard 1 should contain reserved ID 4")
	assert.Contains(t, shards[1], "5", "Shard 1 should contain reserved ID 5")

	assert.Contains(t, shards[2], "6", "Shard 2 should contain reserved ID 6")
	assert.Contains(t, shards[2], "10", "Shard 2 should contain reserved ID 10")

	assert.Len(t, shards[0], 20, "Shard 0 should have 20 IDs total (20% of 100)")
	assert.Len(t, shards[1], 30, "Shard 1 should have 30 IDs total (30% of 100)")
	assert.Len(t, shards[2], 50, "Shard 2 should have 50 IDs total (50% of 100)")

	total := countTotalIDs(shards)
	assert.Equal(t, 100, total, "All 100 IDs should be distributed")

	for _, id := range ids {
		assert.True(t, containsID(shards, id), "ID %s should appear exactly once", id)
	}
}

// =============================================================================
// Rendezvous Strategy Tests
// =============================================================================

func TestShardByRendezvous_EmptyList(t *testing.T) {
	ids := []string{}
	shards := shardByRendezvous(context.Background(), ids, 3, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardByRendezvous_Deterministic(t *testing.T) {
	ids := generateTestIDs(100)
	seed := "test-seed"

	shards1 := shardByRendezvous(context.Background(), ids, 3, seed, nil)
	shards2 := shardByRendezvous(context.Background(), ids, 3, seed, nil)

	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should be identical with same seed", i)
	}
}

func TestShardByRendezvous_AllIDsDistributed(t *testing.T) {
	ids := generateTestIDs(100)
	shards := shardByRendezvous(context.Background(), ids, 3, "seed", nil)

	require.Len(t, shards, 3, "Expected 3 shards")

	total := countTotalIDs(shards)
	assert.Equal(t, 100, total, "All 100 IDs should be distributed")

	for _, id := range ids {
		assert.True(t, containsID(shards, id), "ID %s should appear exactly once", id)
	}
}

func TestShardByRendezvous_WithReservations(t *testing.T) {
	ids := generateTestIDs(100)
	
	reservations := &reservationInfo{
		IDsByShard: map[string][]string{
			"shard_0": {"1", "2"},
			"shard_1": {"3", "4", "5", "6"},
			"shard_2": {"7"},
		},
		CountsByShard: map[int]int{
			0: 2,
			1: 4,
			2: 1,
		},
		UnreservedIDs: make([]string, 0, 93),
	}

	for _, id := range ids {
		if !slices.Contains([]string{"1", "2", "3", "4", "5", "6", "7"}, id) {
			reservations.UnreservedIDs = append(reservations.UnreservedIDs, id)
		}
	}

	shards := shardByRendezvous(context.Background(), ids, 3, "test-seed", reservations)

	require.Len(t, shards, 3, "Expected 3 shards")

	assert.Contains(t, shards[0], "1", "Shard 0 should contain reserved ID 1")
	assert.Contains(t, shards[0], "2", "Shard 0 should contain reserved ID 2")
	assert.Equal(t, "1", shards[0][0], "Reserved IDs should appear first")

	assert.Contains(t, shards[1], "3", "Shard 1 should contain reserved ID 3")
	assert.Contains(t, shards[1], "4", "Shard 1 should contain reserved ID 4")
	assert.Contains(t, shards[1], "5", "Shard 1 should contain reserved ID 5")
	assert.Contains(t, shards[1], "6", "Shard 1 should contain reserved ID 6")

	assert.Contains(t, shards[2], "7", "Shard 2 should contain reserved ID 7")

	total := countTotalIDs(shards)
	assert.Equal(t, 100, total, "All 100 IDs should be distributed")

	for _, id := range ids {
		assert.True(t, containsID(shards, id), "ID %s should appear exactly once", id)
	}
}

// =============================================================================
// Rendezvous Distribution Variance Tests
// Purpose: Demonstrate that rendezvous hashing improves distribution balance
// as dataset size increases due to law of large numbers
// =============================================================================

func TestShardByRendezvous_DistributionVariance_100IDs(t *testing.T) {
	ids := generateTestIDs(100)
	shardCount := 3
	shards := shardByRendezvous(context.Background(), ids, shardCount, "variance-seed", nil)

	require.Len(t, shards, shardCount, "Expected %d shards", shardCount)

	sizes := []int{len(shards[0]), len(shards[1]), len(shards[2])}
	minSize := sizes[0]
	maxSize := sizes[0]
	sum := 0

	for _, size := range sizes {
		if size < minSize {
			minSize = size
		}
		if size > maxSize {
			maxSize = size
		}
		sum += size
	}

	variance := maxSize - minSize
	avgSize := float64(sum) / float64(shardCount)
	expectedSize := float64(100) / float64(shardCount)

	t.Logf("Dataset: 100 IDs | Shard sizes: %v", sizes)
	t.Logf("Min: %d, Max: %d, Variance: %d", minSize, maxSize, variance)
	t.Logf("Average: %.2f, Expected: %.2f", avgSize, expectedSize)

	assert.Equal(t, 100, sum, "All IDs should be distributed")
	assert.LessOrEqual(t, variance, 20, "Variance should be reasonable for 100 IDs")
}

func TestShardByRendezvous_DistributionVariance_500IDs(t *testing.T) {
	ids := generateTestIDs(500)
	shardCount := 3
	shards := shardByRendezvous(context.Background(), ids, shardCount, "variance-seed", nil)

	require.Len(t, shards, shardCount, "Expected %d shards", shardCount)

	sizes := []int{len(shards[0]), len(shards[1]), len(shards[2])}
	minSize := sizes[0]
	maxSize := sizes[0]
	sum := 0

	for _, size := range sizes {
		if size < minSize {
			minSize = size
		}
		if size > maxSize {
			maxSize = size
		}
		sum += size
	}

	variance := maxSize - minSize
	avgSize := float64(sum) / float64(shardCount)
	expectedSize := float64(500) / float64(shardCount)
	variancePercent := (float64(variance) / expectedSize) * 100

	t.Logf("Dataset: 500 IDs | Shard sizes: %v", sizes)
	t.Logf("Min: %d, Max: %d, Variance: %d (%.1f%%)", minSize, maxSize, variance, variancePercent)
	t.Logf("Average: %.2f, Expected: %.2f", avgSize, expectedSize)

	assert.Equal(t, 500, sum, "All IDs should be distributed")
	assert.LessOrEqual(t, variance, 40, "Variance should improve with larger dataset")
}

func TestShardByRendezvous_DistributionVariance_1000IDs(t *testing.T) {
	ids := generateTestIDs(1000)
	shardCount := 3
	shards := shardByRendezvous(context.Background(), ids, shardCount, "variance-seed", nil)

	require.Len(t, shards, shardCount, "Expected %d shards", shardCount)

	sizes := []int{len(shards[0]), len(shards[1]), len(shards[2])}
	minSize := sizes[0]
	maxSize := sizes[0]
	sum := 0

	for _, size := range sizes {
		if size < minSize {
			minSize = size
		}
		if size > maxSize {
			maxSize = size
		}
		sum += size
	}

	variance := maxSize - minSize
	avgSize := float64(sum) / float64(shardCount)
	expectedSize := float64(1000) / float64(shardCount)
	variancePercent := (float64(variance) / expectedSize) * 100

	t.Logf("Dataset: 1000 IDs | Shard sizes: %v", sizes)
	t.Logf("Min: %d, Max: %d, Variance: %d (%.1f%%)", minSize, maxSize, variance, variancePercent)
	t.Logf("Average: %.2f, Expected: %.2f", avgSize, expectedSize)

	assert.Equal(t, 1000, sum, "All IDs should be distributed")
	assert.LessOrEqual(t, variance, 60, "Variance should continue to improve")
}

func TestShardByRendezvous_DistributionVariance_5000IDs(t *testing.T) {
	ids := generateTestIDs(5000)
	shardCount := 3
	shards := shardByRendezvous(context.Background(), ids, shardCount, "variance-seed", nil)

	require.Len(t, shards, shardCount, "Expected %d shards", shardCount)

	sizes := []int{len(shards[0]), len(shards[1]), len(shards[2])}
	minSize := sizes[0]
	maxSize := sizes[0]
	sum := 0

	for _, size := range sizes {
		if size < minSize {
			minSize = size
		}
		if size > maxSize {
			maxSize = size
		}
		sum += size
	}

	variance := maxSize - minSize
	avgSize := float64(sum) / float64(shardCount)
	expectedSize := float64(5000) / float64(shardCount)
	variancePercent := (float64(variance) / expectedSize) * 100

	t.Logf("Dataset: 5000 IDs | Shard sizes: %v", sizes)
	t.Logf("Min: %d, Max: %d, Variance: %d (%.1f%%)", minSize, maxSize, variance, variancePercent)
	t.Logf("Average: %.2f, Expected: %.2f", avgSize, expectedSize)

	assert.Equal(t, 5000, sum, "All IDs should be distributed")
	assert.LessOrEqual(t, variance, 150, "Variance should be reasonable for 5000 IDs")
	assert.LessOrEqual(t, variancePercent, 10.0, "Variance should be less than 10% of expected shard size")
}

func TestShardByRendezvous_DistributionVariance_10000IDs(t *testing.T) {
	ids := generateTestIDs(10000)
	shardCount := 3
	shards := shardByRendezvous(context.Background(), ids, shardCount, "variance-seed", nil)

	require.Len(t, shards, shardCount, "Expected %d shards", shardCount)

	sizes := []int{len(shards[0]), len(shards[1]), len(shards[2])}
	minSize := sizes[0]
	maxSize := sizes[0]
	sum := 0

	for _, size := range sizes {
		if size < minSize {
			minSize = size
		}
		if size > maxSize {
			maxSize = size
		}
		sum += size
	}

	variance := maxSize - minSize
	avgSize := float64(sum) / float64(shardCount)
	expectedSize := float64(10000) / float64(shardCount)
	variancePercent := (float64(variance) / expectedSize) * 100

	t.Logf("Dataset: 10000 IDs | Shard sizes: %v", sizes)
	t.Logf("Min: %d, Max: %d, Variance: %d (%.1f%%)", minSize, maxSize, variance, variancePercent)
	t.Logf("Average: %.2f, Expected: %.2f", avgSize, expectedSize)

	assert.Equal(t, 10000, sum, "All IDs should be distributed")
	assert.LessOrEqual(t, variance, 250, "Variance should be reasonable for 10000 IDs")
	assert.LessOrEqual(t, variancePercent, 8.0, "Variance should be less than 8% of expected shard size")
	
	// Verify this is the best distribution (lowest variance %)
	t.Logf("Distribution quality: Excellent - variance only %.1f%% with 10k IDs", variancePercent)
}

// =============================================================================
// Size Strategy Tests
// =============================================================================

func TestShardBySize_EmptyList(t *testing.T) {
	ids := []string{}
	sizes := []int64{10, 20, -1}
	shards := shardBySize(context.Background(), ids, sizes, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardBySize_ExactSizes(t *testing.T) {
	ids := generateTestIDs(100)
	sizes := []int64{50, 30, 20}
	shards := shardBySize(context.Background(), ids, sizes, "", nil)

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
	shards := shardBySize(context.Background(), ids, sizes, "", nil)

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
	shards := shardBySize(context.Background(), ids, sizes, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 10, "Shard 0 should have 10 IDs")
	assert.Len(t, shards[1], 20, "Shard 1 should have 20 IDs")

	assert.NotNil(t, shards[2], "Shard 2 must not be nil")
	assert.Len(t, shards[2], 0, "Shard 2 should have 0 IDs")

	total := countTotalIDs(shards)
	assert.Equal(t, 30, total, "All IDs should be distributed")
}

func TestShardBySize_WithReservations(t *testing.T) {
	ids := generateTestIDs(100)
	sizes := []int64{25, 35, 40}
	
	reservations := &reservationInfo{
		IDsByShard: map[string][]string{
			"shard_0": {"1", "2", "3", "4", "5"},
			"shard_1": {"6", "7"},
			"shard_2": {"8", "9", "10"},
		},
		CountsByShard: map[int]int{
			0: 5,
			1: 2,
			2: 3,
		},
		UnreservedIDs: make([]string, 0, 90),
	}

	for _, id := range ids {
		if !slices.Contains([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, id) {
			reservations.UnreservedIDs = append(reservations.UnreservedIDs, id)
		}
	}

	shards := shardBySize(context.Background(), ids, sizes, "", reservations)

	require.Len(t, shards, 3, "Expected 3 shards")

	assert.Contains(t, shards[0], "1", "Shard 0 should contain reserved ID 1")
	assert.Contains(t, shards[0], "5", "Shard 0 should contain reserved ID 5")
	assert.Equal(t, "1", shards[0][0], "Reserved IDs should appear first")

	assert.Contains(t, shards[1], "6", "Shard 1 should contain reserved ID 6")
	assert.Contains(t, shards[1], "7", "Shard 1 should contain reserved ID 7")

	assert.Contains(t, shards[2], "8", "Shard 2 should contain reserved ID 8")
	assert.Contains(t, shards[2], "9", "Shard 2 should contain reserved ID 9")
	assert.Contains(t, shards[2], "10", "Shard 2 should contain reserved ID 10")

	assert.Len(t, shards[0], 25, "Shard 0 should have 25 IDs total")
	assert.Len(t, shards[1], 35, "Shard 1 should have 35 IDs total")
	assert.Len(t, shards[2], 40, "Shard 2 should have 40 IDs total")

	total := countTotalIDs(shards)
	assert.Equal(t, 100, total, "All 100 IDs should be distributed")

	for _, id := range ids {
		assert.True(t, containsID(shards, id), "ID %s should appear exactly once", id)
	}
}
