package mobiledeviceapplications

import (
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// Helper function to check if there are any exclusions
func hasExclusions(exclusions jamfpro.MobileDeviceApplicationSubsetExclusion) bool {
	return len(exclusions.MobileDevices) > 0 ||
		len(exclusions.Buildings) > 0 ||
		len(exclusions.Departments) > 0 ||
		len(exclusions.MobileDeviceGroups) > 0 ||
		len(exclusions.JSSUsers) > 0 ||
		len(exclusions.JSSUserGroups) > 0
}

// normalizeWhitespace removes leading/trailing whitespace and normalizes newlines
func normalizeWhitespace(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}
