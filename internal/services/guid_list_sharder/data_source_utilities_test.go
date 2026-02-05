package guid_list_sharder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSeededRNG_Deterministic(t *testing.T) {
	seed := "test-seed"

	rng1 := createSeededRNG(seed)
	rng2 := createSeededRNG(seed)

	for i := 0; i < 10; i++ {
		val1 := rng1.Intn(1000)
		val2 := rng2.Intn(1000)
		assert.Equal(t, val1, val2, "Same seed should produce identical random sequences")
	}
}

func TestCreateSeededRNG_DifferentSeeds(t *testing.T) {
	rng1 := createSeededRNG("seed1")
	rng2 := createSeededRNG("seed2")

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

func TestShuffleIDs_Deterministic(t *testing.T) {
	ids := generateTestIDs(20)

	shuffled1 := shuffleIDs(ids, "test-seed")
	shuffled2 := shuffleIDs(ids, "test-seed")

	assert.Equal(t, shuffled1, shuffled2, "Same seed should produce identical shuffle order")
}

func TestSortIDsNumerically(t *testing.T) {
	ids := []string{"10", "2", "100", "1", "50"}
	sortIDsNumerically(ids)

	expected := []string{"1", "2", "10", "50", "100"}
	assert.Equal(t, expected, ids, "IDs should be sorted numerically")
}

func TestPrepareIDsForDistribution_WithSeed(t *testing.T) {
	ids := generateTestIDs(10)

	prepared1 := prepareIDsForDistribution(ids, "test-seed")
	prepared2 := prepareIDsForDistribution(ids, "test-seed")

	assert.Equal(t, prepared1, prepared2, "Same seed should produce identical preparation")
	assert.NotEqual(t, ids, prepared1, "Prepared IDs should be shuffled")
}

func TestPrepareIDsForDistribution_NoSeed(t *testing.T) {
	ids := generateTestIDs(10)

	prepared := prepareIDsForDistribution(ids, "")

	assert.Equal(t, ids, prepared, "Without seed, should return original IDs")
}
