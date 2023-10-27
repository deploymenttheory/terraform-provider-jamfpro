// provider_helpers.go
package providerhelpers

import (
	"math/rand"
	"time"
)

const (
	maxReadResourceRetries = 10
	maxDelay               = 120 * time.Second // Max delay of 30 seconds
)

// exponentialBackoffWithJitter computes the wait duration for the current Resource Read retry attempt.
// The function uses exponential backoff with jitter to introduce randomness. This helps distribute the
// request load evenly and reduces the chances of "thundering herd" issues. The formula starts with a
// delay of 1 seconds (1<<attempt) and introduces a random jitter.
func exponentialBackoffWithJitter(attempt int) time.Duration {
	expBackoff := time.Duration(1<<attempt) * time.Second
	jitter := time.Duration(rand.Int63n(int64(expBackoff)))
	return time.Duration(minDuration(expBackoff+jitter, maxDelay))
}

// minDuration returns the smaller of the two provided time durations.
func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
