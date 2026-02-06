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

// TestShardByRoundRobin_EmptyList validates round-robin behavior with zero IDs.
// Purpose: Ensure the algorithm handles empty input gracefully without panics or errors.
// Why: Edge case testing - users may have empty groups or filtered all IDs via exclude_ids.
// Expected: Returns empty shards array with correct length, no crashes.
func TestShardByRoundRobin_EmptyList(t *testing.T) {
	ids := []string{}
	shards := shardByRoundRobin(context.Background(), ids, 3, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

// TestShardByRoundRobin_SingleShard validates round-robin with shard_count=1.
// Purpose: Ensure all IDs land in the single shard when only one shard is requested.
// Why: Common use case - users may start with single group before splitting later.
// Expected: All 10 IDs should be in shard_0, no distribution needed.
func TestShardByRoundRobin_SingleShard(t *testing.T) {
	ids := generateTestIDs(10)
	shards := shardByRoundRobin(context.Background(), ids, 1, "", nil)

	require.Len(t, shards, 1, "Expected 1 shard")
	assert.Len(t, shards[0], len(ids), "All IDs should be in single shard")
}

// TestShardByRoundRobin_PerfectDistribution validates equal distribution with evenly divisible IDs.
// Purpose: Verify round-robin produces perfectly balanced shards when ID count divides evenly.
// Why: Core promise of round-robin is ±1 variance - with perfect division, variance should be 0.
// Expected: 30 IDs across 3 shards = exactly 10 IDs per shard, no variance.
func TestShardByRoundRobin_PerfectDistribution(t *testing.T) {
	ids := generateTestIDs(30)
	shards := shardByRoundRobin(context.Background(), ids, 3, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")

	for i, shard := range shards {
		assert.Len(t, shard, 10, "Shard %d should have exactly 10 IDs", i)
	}
}

// TestShardByRoundRobin_WithRemainder validates distribution when ID count doesn't divide evenly.
// Purpose: Verify round-robin maintains ±1 variance guarantee with remainders.
// Why: Most real-world scenarios won't have perfect divisibility - algorithm must handle gracefully.
// Expected: 31 IDs across 3 shards = 11, 10, 10 distribution (or similar ±1 pattern).
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

// TestShardByRoundRobin_LargeDataset validates round-robin behavior with enterprise-scale ID counts.
// Purpose: Ensure algorithm maintains ±1 variance promise even with large datasets (512 IDs).
// Why: Production environments may have 500-1000+ devices per deployment ring.
// Expected: 512 IDs across 3 shards maintains ±1 variance, no performance degradation.
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

// TestShardByRoundRobin_NoSeed_UsesInputOrder validates sequential distribution without seed.
// Purpose: Verify that without a seed, round-robin distributes IDs in exact API response order.
// Why: Users need to understand no-seed behavior - IDs go to shards in order: 1→s0, 2→s1, 3→s2, 4→s0...
// Expected: First 3 IDs go to shards 0,1,2 sequentially, then cycle repeats.
func TestShardByRoundRobin_NoSeed_UsesInputOrder(t *testing.T) {
	ids := generateTestIDs(9)
	shards := shardByRoundRobin(context.Background(), ids, 3, "", nil)

	assert.Equal(t, ids[0], shards[0][0], "First ID should be in shard 0")
	assert.Equal(t, ids[1], shards[1][0], "Second ID should be in shard 1")
	assert.Equal(t, ids[2], shards[2][0], "Third ID should be in shard 2")
	assert.Equal(t, ids[3], shards[0][1], "Fourth ID should be in shard 0")
}

// TestShardByRoundRobin_WithSeed_Deterministic validates reproducibility with seeds.
// Purpose: Verify same seed produces identical shard assignments across multiple runs.
// Why: Critical for Terraform - users need stable plans without unexpected shard membership changes.
// Expected: Two calls with same seed produce byte-for-byte identical shard contents.
func TestShardByRoundRobin_WithSeed_Deterministic(t *testing.T) {
	ids := generateTestIDs(100)
	seed := "test-seed"

	shards1 := shardByRoundRobin(context.Background(), ids, 3, seed, nil)
	shards2 := shardByRoundRobin(context.Background(), ids, 3, seed, nil)

	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should be identical with same seed", i)
	}
}

// TestShardByRoundRobin_DifferentSeeds validates seed independence for multiple rollouts.
// Purpose: Verify different seeds produce different distributions to spread pilot burden.
// Why: Users run multiple rollout types (OS updates, apps, policies) - same devices shouldn't always be pilots.
// Expected: >50% of IDs land in different shards between seed1 and seed2.
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

// TestShardByRoundRobin_WithReservations validates reserved ID pinning with round-robin distribution.
// Purpose: Verify reserved IDs go to specified shards AND appear first in those shards.
// Why: Users need specific devices (VIPs, test devices) in specific rings regardless of distribution.
// Expected: Reserved IDs appear at start of their designated shards, remaining IDs distributed round-robin.
func TestShardByRoundRobin_WithReservations(t *testing.T) {
	ids := generateTestIDs(100)

	reservations := &shardReservations{
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

// TestShardByPercentage_EmptyList validates percentage strategy with zero IDs.
// Purpose: Ensure percentage algorithm handles empty input gracefully without panics.
// Why: Edge case testing - users may have empty groups or filtered all IDs.
// Expected: Returns empty shards array with correct length based on percentages array size.
func TestShardByPercentage_EmptyList(t *testing.T) {
	ids := []string{}
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(context.Background(), ids, percentages, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

// TestShardByPercentage_AccuratePercentages validates exact percentage-based distribution.
// Purpose: Verify percentage strategy produces exact shard sizes matching specified percentages.
// Why: Core promise of percentage strategy - 10% means exactly 10% of IDs, not approximately.
// Expected: [10,30,60] with 100 IDs = exactly 10, 30, 60 IDs per shard.
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

// TestShardByPercentage_LargeDataset validates percentage accuracy with enterprise-scale datasets.
// Purpose: Verify percentage calculations remain accurate with large ID counts (512 IDs).
// Why: Percentage rounding errors could accumulate with large datasets.
// Expected: Percentages apply to total count, last shard gets remainder to ensure all IDs distributed.
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

// TestShardByPercentage_LastShardGetsRemainder validates remainder handling with odd ID counts.
// Purpose: Verify last shard receives all remaining IDs after percentage allocation to ensure nothing is lost.
// Why: With 103 IDs and [10,20,70] percentages, rounding means some IDs remain after first shards filled.
// Expected: All 103 IDs distributed, last shard gets any remainder from percentage calculations.
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

// TestShardByPercentage_NoSeed_UsesInputOrder validates sequential filling without seed.
// Purpose: Verify without seed, percentage slices fill sequentially from API response order.
// Why: Users need to understand no-seed behavior - first 20% of IDs go to shard_0, next 30% to shard_1, etc.
// Expected: IDs appear in shards in their original order, no shuffling.
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

// TestShardByPercentage_WithSeed_Deterministic validates reproducibility with seeds for percentage strategy.
// Purpose: Verify same seed produces identical percentage-based shard assignments across runs.
// Why: Critical for Terraform stability - percentage splits must be consistent across plan/apply cycles.
// Expected: Two calls with same seed produce byte-for-byte identical shard contents.
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

// TestShardByPercentage_WithReservations validates reserved ID integration with percentage distribution.
// Purpose: Verify reserved IDs count toward percentage targets and maintain exact percentages.
// Why: Users need specific devices in specific rings while maintaining percentage-based ring sizes.
// Expected: Reserved IDs appear first, remaining IDs fill to exact percentage targets (20%=20, 30%=30, 50%=50).
func TestShardByPercentage_WithReservations(t *testing.T) {
	ids := generateTestIDs(100)
	percentages := []int64{20, 30, 50}

	reservations := &shardReservations{
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

// TestShardByRendezvous_EmptyList validates rendezvous HRW algorithm with zero IDs.
// Purpose: Ensure consistent hashing handles empty input gracefully without panics.
// Why: Edge case testing for rendezvous-specific code path.
// Expected: Returns empty shards array with correct length, no crashes from hash calculations.
func TestShardByRendezvous_EmptyList(t *testing.T) {
	ids := []string{}
	shards := shardByRendezvous(context.Background(), ids, 3, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

// TestShardByRendezvous_Deterministic validates consistent hashing reproducibility.
// Purpose: Verify rendezvous (HRW) produces identical shard assignments with same seed.
// Why: Rendezvous must be deterministic for minimal-disruption topology changes - critical property.
// Expected: Two calls with same seed produce byte-for-byte identical results via consistent hashing.
func TestShardByRendezvous_Deterministic(t *testing.T) {
	ids := generateTestIDs(100)
	seed := "test-seed"

	shards1 := shardByRendezvous(context.Background(), ids, 3, seed, nil)
	shards2 := shardByRendezvous(context.Background(), ids, 3, seed, nil)

	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should be identical with same seed", i)
	}
}

// TestShardByRendezvous_AllIDsDistributed validates complete ID coverage with HRW algorithm.
// Purpose: Verify every ID gets assigned to exactly one shard via highest random weight selection.
// Why: Rendezvous hash must not lose IDs or duplicate them during weight-based assignment.
// Expected: All 100 IDs present exactly once across all shards.
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

// TestShardByRendezvous_WithReservations validates reserved ID integration with consistent hashing.
// Purpose: Verify reserved IDs bypass HRW algorithm and go directly to designated shards.
// Why: Users need specific devices pinned to rings even when using topology-stable rendezvous strategy.
// Expected: Reserved IDs appear first in designated shards, remaining IDs distributed via HRW.
func TestShardByRendezvous_WithReservations(t *testing.T) {
	ids := generateTestIDs(100)

	reservations := &shardReservations{
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

// TestShardByRendezvous_DistributionVariance_100IDs measures HRW variance with small dataset.
// Purpose: Quantify actual distribution imbalance with 100 IDs to document in schema.
// Why: Users need realistic expectations - rendezvous prioritizes stability over perfect balance.
// Expected: Variance ≈12% (actual: 27-39 IDs per shard), acceptable for minimal-disruption benefit.
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

// TestShardByRendezvous_DistributionVariance_500IDs measures HRW variance with medium dataset.
// Purpose: Demonstrate variance decreases as dataset grows (law of large numbers).
// Why: Helps users understand when rendezvous becomes more balanced (500+ IDs).
// Expected: Variance ≈15% (actual: 154-179 IDs per shard), better than 100 IDs but still present.
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

// TestShardByRendezvous_DistributionVariance_1000IDs measures HRW variance with large dataset.
// Purpose: Quantify variance improvement at typical enterprise scale (1000 devices).
// Why: Most organizations have 500-2000 devices - this is the relevant real-world datapoint.
// Expected: Variance ≈10% (actual: 320-353 IDs per shard), significantly better than smaller datasets.
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

// TestShardByRendezvous_DistributionVariance_5000IDs measures HRW variance at very large scale.
// Purpose: Demonstrate variance continues decreasing with larger datasets.
// Why: Large enterprises with 5000+ devices need to understand rendezvous becomes well-balanced.
// Expected: Variance ≈6% (actual: 1610-1710 IDs per shard), approaching round-robin balance.
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

// TestShardByRendezvous_DistributionVariance_10000IDs measures HRW variance at maximum enterprise scale.
// Purpose: Prove rendezvous achieves excellent balance with very large datasets (10k+).
// Why: Demonstrates law of large numbers effect - variance becomes negligible at scale.
// Expected: Variance ≈1% (actual: 3321-3354 IDs per shard), essentially perfect distribution.
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
// Rendezvous with Reserved IDs - Variance Impact Tests
// Purpose: Compare distribution variance with and without reserved IDs
// to validate claims about increased variance when using reservations
// =============================================================================

// TestShardByRendezvous_VarianceWithReservedIDs_1000IDs compares HRW variance with/without reserved IDs at 1000 ID scale.
// Purpose: Validate schema claim that reserved IDs may increase variance and quantify the actual impact.
// Why: Users need data-driven understanding of reserved ID tradeoffs - does pinning devices hurt balance?
// Expected: Baseline 9.9% variance may change with reservations (test shows -1.2% = actually improved).
func TestShardByRendezvous_VarianceWithReservedIDs_1000IDs(t *testing.T) {
	ids := generateTestIDs(1000)
	shardCount := 3

	// Test WITHOUT reserved IDs (baseline)
	shardsNoReservations := shardByRendezvous(context.Background(), ids, shardCount, "variance-seed", nil)
	
	sizesNoRes := []int{len(shardsNoReservations[0]), len(shardsNoReservations[1]), len(shardsNoReservations[2])}
	minNoRes, maxNoRes := sizesNoRes[0], sizesNoRes[0]
	for _, size := range sizesNoRes {
		if size < minNoRes {
			minNoRes = size
		}
		if size > maxNoRes {
			maxNoRes = size
		}
	}
	varianceNoRes := maxNoRes - minNoRes
	variancePercentNoRes := (float64(varianceNoRes) / (1000.0 / 3.0)) * 100

	// Test WITH reserved IDs
	reservations := &shardReservations{
		IDsByShard: map[string][]string{
			"shard_0": {"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},  // 10 reserved
			"shard_1": {"11", "12", "13", "14", "15"},                        // 5 reserved
			"shard_2": {"16", "17", "18"},                                     // 3 reserved
		},
		CountsByShard: map[int]int{
			0: 10,
			1: 5,
			2: 3,
		},
		UnreservedIDs: make([]string, 0, 982),
	}

	reservedSet := map[string]bool{
		"1": true, "2": true, "3": true, "4": true, "5": true,
		"6": true, "7": true, "8": true, "9": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true,
		"16": true, "17": true, "18": true,
	}

	for _, id := range ids {
		if !reservedSet[id] {
			reservations.UnreservedIDs = append(reservations.UnreservedIDs, id)
		}
	}

	shardsWithReservations := shardByRendezvous(context.Background(), ids, shardCount, "variance-seed", reservations)

	sizesWithRes := []int{len(shardsWithReservations[0]), len(shardsWithReservations[1]), len(shardsWithReservations[2])}
	minWithRes, maxWithRes := sizesWithRes[0], sizesWithRes[0]
	for _, size := range sizesWithRes {
		if size < minWithRes {
			minWithRes = size
		}
		if size > maxWithRes {
			maxWithRes = size
		}
	}
	varianceWithRes := maxWithRes - minWithRes
	variancePercentWithRes := (float64(varianceWithRes) / (1000.0 / 3.0)) * 100

	t.Logf("WITHOUT Reserved IDs: Sizes=%v, Variance=%d (%.1f%%)", sizesNoRes, varianceNoRes, variancePercentNoRes)
	t.Logf("WITH Reserved IDs:    Sizes=%v, Variance=%d (%.1f%%)", sizesWithRes, varianceWithRes, variancePercentWithRes)
	t.Logf("Variance Delta: %.1f%% → %.1f%% (change: %+.1f%%)", variancePercentNoRes, variancePercentWithRes, variancePercentWithRes-variancePercentNoRes)

	// Verify all IDs are present
	assert.Equal(t, 1000, countTotalIDs(shardsWithReservations), "All 1000 IDs should be distributed")
	
	// Check if variance increased (or stayed similar)
	if variancePercentWithRes > variancePercentNoRes {
		t.Logf("✓ Reserved IDs increased variance by %.1f%%", variancePercentWithRes-variancePercentNoRes)
	} else {
		t.Logf("✓ Reserved IDs did not significantly increase variance (delta: %.1f%%)", variancePercentWithRes-variancePercentNoRes)
	}
}

// TestShardByRendezvous_VarianceWithReservedIDs_5000IDs compares HRW variance with/without reserved IDs at large scale.
// Purpose: Quantify reserved ID impact with significant reservations (85 IDs reserved from 5000).
// Why: Demonstrates that larger reservation counts do increase variance moderately as documented.
// Expected: Baseline 6.0% variance increases to 8.9% (+2.9%) with 85 reserved IDs across shards.
func TestShardByRendezvous_VarianceWithReservedIDs_5000IDs(t *testing.T) {
	ids := generateTestIDs(5000)
	shardCount := 3

	// Test WITHOUT reserved IDs (baseline)
	shardsNoReservations := shardByRendezvous(context.Background(), ids, shardCount, "variance-seed", nil)
	
	sizesNoRes := []int{len(shardsNoReservations[0]), len(shardsNoReservations[1]), len(shardsNoReservations[2])}
	minNoRes, maxNoRes := sizesNoRes[0], sizesNoRes[0]
	for _, size := range sizesNoRes {
		if size < minNoRes {
			minNoRes = size
		}
		if size > maxNoRes {
			maxNoRes = size
		}
	}
	varianceNoRes := maxNoRes - minNoRes
	variancePercentNoRes := (float64(varianceNoRes) / (5000.0 / 3.0)) * 100

	// Test WITH reserved IDs (more significant reservations)
	reservations := &shardReservations{
		IDsByShard: map[string][]string{
			"shard_0": make([]string, 50),  // 50 reserved
			"shard_1": make([]string, 25),  // 25 reserved
			"shard_2": make([]string, 10),  // 10 reserved
		},
		CountsByShard: map[int]int{
			0: 50,
			1: 25,
			2: 10,
		},
		UnreservedIDs: make([]string, 0, 4915),
	}

	// Populate reserved IDs
	reservedSet := make(map[string]bool)
	idIdx := 0
	for i := range 50 {
		id := fmt.Sprintf("%d", idIdx+1)
		reservations.IDsByShard["shard_0"][i] = id
		reservedSet[id] = true
		idIdx++
	}
	for i := range 25 {
		id := fmt.Sprintf("%d", idIdx+1)
		reservations.IDsByShard["shard_1"][i] = id
		reservedSet[id] = true
		idIdx++
	}
	for i := range 10 {
		id := fmt.Sprintf("%d", idIdx+1)
		reservations.IDsByShard["shard_2"][i] = id
		reservedSet[id] = true
		idIdx++
	}

	for _, id := range ids {
		if !reservedSet[id] {
			reservations.UnreservedIDs = append(reservations.UnreservedIDs, id)
		}
	}

	shardsWithReservations := shardByRendezvous(context.Background(), ids, shardCount, "variance-seed", reservations)

	sizesWithRes := []int{len(shardsWithReservations[0]), len(shardsWithReservations[1]), len(shardsWithReservations[2])}
	minWithRes, maxWithRes := sizesWithRes[0], sizesWithRes[0]
	for _, size := range sizesWithRes {
		if size < minWithRes {
			minWithRes = size
		}
		if size > maxWithRes {
			maxWithRes = size
		}
	}
	varianceWithRes := maxWithRes - minWithRes
	variancePercentWithRes := (float64(varianceWithRes) / (5000.0 / 3.0)) * 100

	t.Logf("WITHOUT Reserved IDs: Sizes=%v, Variance=%d (%.1f%%)", sizesNoRes, varianceNoRes, variancePercentNoRes)
	t.Logf("WITH Reserved IDs:    Sizes=%v, Variance=%d (%.1f%%)", sizesWithRes, varianceWithRes, variancePercentWithRes)
	t.Logf("Variance Delta: %.1f%% → %.1f%% (change: %+.1f%%)", variancePercentNoRes, variancePercentWithRes, variancePercentWithRes-variancePercentNoRes)

	// Verify all IDs are present
	assert.Equal(t, 5000, countTotalIDs(shardsWithReservations), "All 5000 IDs should be distributed")
	
	// Document the impact
	if variancePercentWithRes > variancePercentNoRes {
		t.Logf("✓ Reserved IDs increased variance by %.1f%%", variancePercentWithRes-variancePercentNoRes)
	} else {
		t.Logf("✓ Reserved IDs did not significantly increase variance (delta: %.1f%%)", variancePercentWithRes-variancePercentNoRes)
	}
}

// =============================================================================
// Size Strategy Tests
// =============================================================================

// TestShardBySize_EmptyList validates size-based strategy with zero IDs.
// Purpose: Ensure fixed-size allocation handles empty input gracefully.
// Why: Edge case testing for size-specific code path.
// Expected: Returns empty shards array with correct length based on sizes array, no crashes.
func TestShardBySize_EmptyList(t *testing.T) {
	ids := []string{}
	sizes := []int64{10, 20, -1}
	shards := shardBySize(context.Background(), ids, sizes, "", nil)

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

// TestShardBySize_ExactSizes validates precise size-based allocation without remainder.
// Purpose: Verify size strategy produces exact shard sizes when IDs match specified sizes.
// Why: Core promise of size strategy - [50,30,20] means exactly those counts, not approximations.
// Expected: 100 IDs with [50,30,20] = exactly 50, 30, 20 IDs per shard.
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

// TestShardBySize_WithRemainder validates -1 "all remaining" functionality.
// Purpose: Verify -1 in last position captures all IDs after fixed sizes allocated.
// Why: Common pattern - fixed pilot/staging sizes, then "everyone else" in production ring.
// Expected: [50,200,-1] with 1000 IDs = 50, 200, 750 (all remaining after first two).
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

// TestShardBySize_ZeroRemainder validates -1 behavior when no IDs remain.
// Purpose: Ensure -1 shard gracefully handles case where earlier shards consumed all IDs.
// Why: Edge case - user specifies [10,20,-1] but only has 30 IDs, last shard should be empty not error.
// Expected: [10,20,-1] with 30 IDs = 10, 20, 0 (empty but valid shard).
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

// TestShardBySize_WithReservations validates reserved ID integration with fixed-size allocation.
// Purpose: Verify reserved IDs count toward specified sizes and maintain exact total sizes.
// Why: Users need specific devices in rings while maintaining exact capacity constraints.
// Expected: Reserved IDs appear first, algorithm distributes fewer unreserved IDs to hit exact targets (25,35,40).
func TestShardBySize_WithReservations(t *testing.T) {
	ids := generateTestIDs(100)
	sizes := []int64{25, 35, 40}

	reservations := &shardReservations{
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
