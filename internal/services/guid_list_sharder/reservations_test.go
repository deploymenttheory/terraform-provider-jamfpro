package guid_list_sharder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyReservations_RemovesReservedFromUnreserved(t *testing.T) {
	ids := []string{"1", "2", "3"}
	reserved := map[string][]string{"shard_0": {"2"}}

	res, err := ApplyReservations(ids, reserved, 2)
	require.NoError(t, err)
	require.Equal(t, []string{"1", "3"}, res.UnreservedIDs)
	require.Equal(t, 1, res.CountsByShard[0])
}

func TestApplyReservations_DuplicateIDFails(t *testing.T) {
	ids := []string{"1", "2", "3"}
	reserved := map[string][]string{
		"shard_0": {"2"},
		"shard_1": {"2"},
	}

	_, err := ApplyReservations(ids, reserved, 2)
	require.Error(t, err)
}

func TestApplyReservations_IgnoresIDsNotInSource(t *testing.T) {
	ids := []string{"1", "2", "3"}
	reserved := map[string][]string{"shard_0": {"999"}}

	res, err := ApplyReservations(ids, reserved, 2)
	require.NoError(t, err)
	require.Empty(t, res.IDsByShard["shard_0"])
}
