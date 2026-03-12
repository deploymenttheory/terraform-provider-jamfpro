package guid_list_sharder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShard_RoundRobin(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}

	shards, err := Shard(ids, StrategyRoundRobin, 2, nil, nil, "", nil)
	require.NoError(t, err)
	require.Len(t, shards, 2)

	// round-robin: 1,3,5 and 2,4 (sorted per shard)
	require.Equal(t, []string{"1", "3", "5"}, shards[0])
	require.Equal(t, []string{"2", "4"}, shards[1])
}

func TestShard_Percentage(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	shards, err := Shard(ids, StrategyPercentage, 0, []int{20, 30, 50}, nil, "seed", nil)
	require.NoError(t, err)
	require.Len(t, shards, 3)
	require.Len(t, shards[0], 2)
	require.Len(t, shards[1], 3)
	require.Len(t, shards[2], 5)
}

func TestShard_Size_WithRemainder(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5", "6"}
	shards, err := Shard(ids, StrategySize, 0, nil, []int{2, -1}, "seed", nil)
	require.NoError(t, err)
	require.Len(t, shards, 2)
	require.Len(t, shards[0], 2)
	require.Len(t, shards[1], 4)
}

func TestShard_Rendezvous_IsDeterministic(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5", "6"}

	a, err := Shard(ids, StrategyRendezvous, 3, nil, nil, "seed", nil)
	require.NoError(t, err)
	b, err := Shard(ids, StrategyRendezvous, 3, nil, nil, "seed", nil)
	require.NoError(t, err)

	require.Equal(t, a, b)
}

func TestReservations_PinToShard0_Size(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}
	reserved := map[string][]string{"shard_0": {"5"}}

	res, err := ApplyReservations(ids, reserved, 2)
	require.NoError(t, err)

	shards, err := Shard(ids, StrategySize, 0, nil, []int{2, -1}, "seed", res)
	require.NoError(t, err)
	require.Contains(t, shards[0], "5")
}

func TestMergeUniqueStrings(t *testing.T) {
	a := []string{"1", "2"}
	b := []string{"2", "3"}
	out := mergeUniqueStrings(a, b)
	require.ElementsMatch(t, []string{"1", "2", "3"}, out)
}
