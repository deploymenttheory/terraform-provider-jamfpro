package provider

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	errEmptyVersionString  = errors.New("empty version string")
	errMissingMinorSegment = errors.New("missing minor version segment")
	errInvalidMajorSegment = errors.New("invalid major version segment")
	errInvalidMinorSegment = errors.New("invalid minor version segment")
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

// versionSupportsRequirement compares the current version against the minimum required version.
func versionSupportsRequirement(current, minimum string) (bool, error) {
	currentMajor, currentMinor, err := extractMajorMinor(current)
	if err != nil {
		return false, err
	}

	minimumMajor, minimumMinor, err := extractMajorMinor(minimum)
	if err != nil {
		return false, err
	}

	switch {
	case currentMajor > minimumMajor:
		return true, nil
	case currentMajor < minimumMajor:
		return false, nil
	default:
		return currentMinor >= minimumMinor, nil
	}
}

// extractMajorMinor extracts the major and minor version numbers from a version string.
func extractMajorMinor(input string) (int, int, error) {
	trimmed := strings.TrimSpace(input)
	if idx := strings.Index(trimmed, "-"); idx != -1 {
		trimmed = trimmed[:idx]
	}

	if trimmed == "" {
		return 0, 0, errEmptyVersionString
	}

	parts := strings.Split(trimmed, ".")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("version %q must include major and minor segments: %w", input, errMissingMinorSegment)
	}

	major, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid major segment in %q: %w", input, errInvalidMajorSegment)
	}

	minor, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minor segment in %q: %w", input, errInvalidMinorSegment)
	}

	return major, minor, nil
}
