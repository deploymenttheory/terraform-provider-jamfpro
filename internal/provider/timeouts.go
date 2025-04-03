package provider

import (
	"time"
)

const (
	DefaultContextTimeout      = 75 * time.Second
	LoadBalancedContextTimeout = 15 * time.Second
)

// GetDefaultContextTimeoutCreate returns the appropriate timeout duration for resource creation.
// If load balancer lock is enabled, it returns the LoadBalancedContextTimeoutCreate,
// otherwise it returns the DefaultContextTimeoutCreate.
func Timeout(load_balancer_lock_enabled bool) time.Duration {
	if load_balancer_lock_enabled {
		return LoadBalancedContextTimeout
	}
	return DefaultContextTimeout
}

// resourceTimeout holds timeouts, used in overrides
type resourceTimeout struct {
	Create, Read, Update, Delete time.Duration
}

// Overrides returns a list of timeout overrides by their resource key
func TimeoutOverrides(lb_lock bool) map[string]resourceTimeout {
	return map[string]resourceTimeout{
		"jamfpro_package": {
			Create: 45 * time.Minute,
			Read:   Timeout(lb_lock),
			Update: 45 * time.Minute,
			Delete: Timeout(lb_lock),
		},
		"jamfpro_smart_computer_group": {
			Create: 75 * time.Second,
			Read:   75 * time.Second,
			Update: 75 * time.Second,
			Delete: 75 * time.Second,
		},
		"jamfpro_static_computer_group": {
			Create: 75 * time.Second,
			Read:   75 * time.Second,
			Update: 75 * time.Second,
			Delete: 75 * time.Second,
		},
		"jamfpro_policy": {
			Create: 30 * time.Second,
		},
	}
}
