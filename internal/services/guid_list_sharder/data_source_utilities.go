package guid_list_sharder

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"slices"
	"strconv"
)

// createSeededRNG creates a deterministic random number generator from a seed string.
func createSeededRNG(seed string) *rand.Rand {
	hash := sha256.Sum256([]byte(seed))
	seedValue := int64(binary.BigEndian.Uint64(hash[:8]))
	return rand.New(rand.NewSource(seedValue))
}

// shuffleIDs performs Fisher-Yates shuffle using the provided seed.
// Returns a shuffled copy without mutating the original slice.
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

// sortIDsNumerically sorts a slice of string IDs by their numeric value.
// Modifies the slice in place.
func sortIDsNumerically(ids []string) {
	slices.SortFunc(ids, func(a, b string) int {
		aInt, _ := strconv.Atoi(a)
		bInt, _ := strconv.Atoi(b)
		return aInt - bInt
	})
}

// prepareIDsForDistribution sorts IDs numerically, then shuffles if seed provided.
// Returns a new slice ready for distribution.
func prepareIDsForDistribution(ids []string, seed string) []string {
	if seed == "" {
		return ids
	}

	sorted := make([]string, len(ids))
	copy(sorted, ids)
	sortIDsNumerically(sorted)
	return shuffleIDs(sorted, seed)
}
