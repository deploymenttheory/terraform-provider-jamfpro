package guid_list_sharder

import "fmt"

// ApplyExclusions removes excluded IDs from the pool.
func ApplyExclusions(ids []string, excludeIDs []string) []string {
	if len(excludeIDs) == 0 {
		return ids
	}
	excludeSet := make(map[string]bool, len(excludeIDs))
	for _, id := range excludeIDs {
		excludeSet[id] = true
	}
	filtered := make([]string, 0, len(ids))
	for _, id := range ids {
		if !excludeSet[id] {
			filtered = append(filtered, id)
		}
	}
	return filtered
}

// ApplyReservations partitions IDs into reserved and unreserved lists.
//
// reservedMap keys must be shard_N where N is within [0, shardCount).
// Each ID may only appear once across all reserved shards.
func ApplyReservations(ids []string, reservedMap map[string][]string, shardCount int) (*Reservations, error) {
	if shardCount < 1 {
		return nil, fmt.Errorf("%w: shardCount must be >= 1", ErrInvalidShardCount)
	}

	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}

	info := &Reservations{
		IDsByShard:    make(map[string][]string),
		CountsByShard: make(map[int]int),
		UnreservedIDs: ids,
	}
	if len(reservedMap) == 0 {
		return info, nil
	}

	seenIDs := make(map[string]string)

	for shardName, idList := range reservedMap {
		var shardIndex int
		if _, err := fmt.Sscanf(shardName, "shard_%d", &shardIndex); err != nil {
			return nil, fmt.Errorf("%w: %q (must be shard_0, shard_1, etc.)", ErrInvalidShardName, shardName)
		}
		if shardIndex < 0 || shardIndex >= shardCount {
			return nil, fmt.Errorf(
				"%w: %q is out of range for shard_count=%d (valid shard_0..shard_%d)",
				ErrInvalidShardName, shardName, shardCount, shardCount-1,
			)
		}
		filteredList := make([]string, 0, len(idList))
		for _, id := range idList {
			if !idSet[id] {
				continue
			}
			filteredList = append(filteredList, id)
		}

		for _, id := range filteredList {
			if prev, exists := seenIDs[id]; exists {
				return nil, fmt.Errorf(
					"%w: ID %q appears in %q and %q",
					ErrReservedIDInMultipleShards, id, prev, shardName,
				)
			}
			seenIDs[id] = shardName
		}
		info.IDsByShard[shardName] = filteredList
		info.CountsByShard[shardIndex] = len(filteredList)
	}

	if len(seenIDs) > 0 {
		reservedSet := make(map[string]bool, len(seenIDs))
		for id := range seenIDs {
			reservedSet[id] = true
		}
		filtered := make([]string, 0, len(ids))
		for _, id := range ids {
			if !reservedSet[id] {
				filtered = append(filtered, id)
			}
		}
		info.UnreservedIDs = filtered
	}

	return info, nil
}
