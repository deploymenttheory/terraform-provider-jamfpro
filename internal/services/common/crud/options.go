package crud

import "time"

// ReadWithRetryOptions configures the retry behavior for reading resource state
type ReadWithRetryOptions struct {
	// MaxRetries is the maximum number of retry attempts (default: 30)
	MaxRetries int
	// InitialRetryInterval is the initial time to wait between retries (default: 2 seconds)
	InitialRetryInterval time.Duration
	// MaxRetryInterval is the maximum time to wait between retries (default: 30 seconds)
	MaxRetryInterval time.Duration
	// BackoffMultiplier is the multiplier for exponential backoff (default: 1.5)
	BackoffMultiplier float64
	// Operation is the name of the operation for logging (e.g., "Create", "Update")
	Operation string
	// ResourceTypeName is the optional resource type name for logging
	ResourceTypeName string
}

// DefaultReadWithRetryOptions returns sensible default options for most use cases
func DefaultReadWithRetryOptions() ReadWithRetryOptions {
	return ReadWithRetryOptions{
		MaxRetries:           30,
		InitialRetryInterval: 2 * time.Second,
		MaxRetryInterval:     30 * time.Second,
		BackoffMultiplier:    1.5,
		Operation:            "Operation",
	}
}
