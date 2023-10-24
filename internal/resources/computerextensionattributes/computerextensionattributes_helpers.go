// computerextensionattributes_helpers.go
package computerextensionattributes

import (
	"math/rand"
	"time"
)

const (
	maxRetries = 10
	maxDelay   = 30 * time.Second // Max delay of 30 seconds
)

// exponentialBackoffWithJitter computes the wait duration for the current Resource Read retry attempt.
// The function uses exponential backoff with jitter to introduce randomness. This helps distribute the
// request load evenly and reduces the chances of "thundering herd" issues. The formula starts with a
// delay of 3 seconds (3<<attempt) and introduces a random jitter.
func exponentialBackoffWithJitter(attempt int) time.Duration {
	expBackoff := time.Duration(3<<attempt) * time.Second
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
