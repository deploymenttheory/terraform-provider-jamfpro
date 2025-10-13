package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Helper functions for getting values with defaults
func getStringValueWithEnvFallback(configValue types.String, envVar string) string {
	if !configValue.IsNull() && configValue.ValueString() != "" {
		return configValue.ValueString()
	}
	return os.Getenv(envVar)
}

func getBoolWithDefault(configValue types.Bool, defaultValue bool) bool {
	if configValue.IsNull() || configValue.IsUnknown() {
		return defaultValue
	}
	return configValue.ValueBool()
}

func getStringWithDefault(configValue types.String, defaultValue string) string {
	if configValue.IsNull() || configValue.IsUnknown() {
		return defaultValue
	}
	return configValue.ValueString()
}

func getInt64WithDefault(configValue types.Int64, defaultValue int64) int64 {
	if configValue.IsNull() || configValue.IsUnknown() {
		return defaultValue
	}
	return configValue.ValueInt64()
}
