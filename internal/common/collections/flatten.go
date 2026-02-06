package collections

import "sort"

// FlattenSortIDs collects non-zero IDs from items, sorts them, and returns a deterministic slice.
func FlattenSortIDs[T any](items []T, idFn func(T) int) []int {
	if len(items) == 0 || idFn == nil {
		return nil
	}

	ids := make([]int, 0, len(items))
	for _, item := range items {
		if id := idFn(item); id != 0 {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return nil
	}

	sort.Ints(ids)
	return ids
}

// FlattenSortStrings collects non-empty strings from items, sorts them, and returns a deterministic slice.
func FlattenSortStrings[T any](items []T, valueFn func(T) string) []string {
	if len(items) == 0 || valueFn == nil {
		return nil
	}

	values := make([]string, 0, len(items))
	for _, item := range items {
		if value := valueFn(item); value != "" {
			values = append(values, value)
		}
	}

	if len(values) == 0 {
		return nil
	}

	sort.Strings(values)
	return values
}
