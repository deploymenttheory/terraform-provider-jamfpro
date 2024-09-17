package provider

import "time"

const (
	DefaultContextTimeoutCreate = 65 * time.Second
	DefaultContextTimeoutRead   = 65 * time.Second
	DefaultContextTimeoutUpdate = 65 * time.Second
	DefaultContextTimeoutDelete = 65 * time.Second

	LoadBalancedContextTimeoutCreate = 15 * time.Second
	LoadBalancedContextTimeoutRead   = 10 * time.Second
	LoadBalancedContextTimeoutUpdate = 15 * time.Second
	LoadBalancedContextTimeoutDelete = 10 * time.Second
)

// GetDefaultContextTimeoutCreate returns the appropriate timeout duration for resource creation.
// If load balancer lock is enabled, it returns the LoadBalancedContextTimeoutCreate,
// otherwise it returns the DefaultContextTimeoutCreate.
func GetDefaultContextTimeoutCreate(load_balancer_lock_enabled bool) time.Duration {
	if load_balancer_lock_enabled {
		return LoadBalancedContextTimeoutCreate
	}
	return DefaultContextTimeoutCreate
}

// GetDefaultContextTimeoutRead returns the appropriate timeout duration for resource reading.
// If load balancer lock is enabled, it returns the LoadBalancedContextTimeoutRead,
// otherwise it returns the DefaultContextTimeoutRead.
func GetDefaultContextTimeoutRead(load_balancer_lock_enabled bool) time.Duration {
	if load_balancer_lock_enabled {
		return LoadBalancedContextTimeoutRead
	}
	return DefaultContextTimeoutRead
}

// GetDefaultContextTimeoutUpdate returns the appropriate timeout duration for resource updating.
// If load balancer lock is enabled, it returns the LoadBalancedContextTimeoutUpdate,
// otherwise it returns the DefaultContextTimeoutUpdate.
func GetDefaultContextTimeoutUpdate(load_balancer_lock_enabled bool) time.Duration {
	if load_balancer_lock_enabled {
		return LoadBalancedContextTimeoutUpdate
	}
	return DefaultContextTimeoutUpdate
}

// GetDefaultContextTimeoutDelete returns the appropriate timeout duration for resource deletion.
// If load balancer lock is enabled, it returns the LoadBalancedContextTimeoutDelete,
// otherwise it returns the DefaultContextTimeoutDelete.
func GetDefaultContextTimeoutDelete(load_balancer_lock_enabled bool) time.Duration {
	if load_balancer_lock_enabled {
		return LoadBalancedContextTimeoutDelete
	}
	return DefaultContextTimeoutDelete
}
