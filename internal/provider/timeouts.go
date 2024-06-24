package provider

import "time"

const (
	DefaultContextTimeoutCreate = 65 * time.Second
	DefaultContextTimeoutRead   = 65 * time.Second
	DefaultContextTimeoutUpdate = 65 * time.Second
	DefaultContextTimeoutDelete = 65 * time.Second

	LoadBalancedContextTimeoutCreate = 15 * time.Second
	LoadBalancedContextTimeoutRead   = 5 * time.Second
	LoadBalancedContextTimeoutUpdate = 15 * time.Second
	LoadBalancedContextTimeoutDelete = 10 * time.Second
)

// TODO func comment
func GetDefaultContextTimeoutCreate(load_balancer_lock_enabled bool) time.Duration {
	if load_balancer_lock_enabled {
		return LoadBalancedContextTimeoutCreate
	}
	return DefaultContextTimeoutCreate
}

// TODO func comment
func GetDefaultContextTimeoutRead(load_balancer_lock_enabled bool) time.Duration {
	if load_balancer_lock_enabled {
		return LoadBalancedContextTimeoutRead
	}
	return DefaultContextTimeoutRead
}

// TODO func comment
func GetDefaultContextTimeoutUpdate(load_balancer_lock_enabled bool) time.Duration {
	if load_balancer_lock_enabled {
		return LoadBalancedContextTimeoutUpdate
	}
	return DefaultContextTimeoutUpdate
}

// TODO func comment
func GetDefaultContextTimeoutDelete(load_balancer_lock_enabled bool) time.Duration {
	if load_balancer_lock_enabled {
		return LoadBalancedContextTimeoutDelete
	}
	return DefaultContextTimeoutDelete
}
