package provider

import (
	"fmt"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/go-version"
)

// MinimumSupportedJamfProVersion is the minimum Jamf Pro version supported by this provider
const MinimumSupportedJamfProVersion = "11.22.0"

// CheckJamfProVersion validates that the Jamf Pro instance version meets minimum requirements
// Returns a warning message if the version is below minimum supported, or error if version cannot be determined
func CheckJamfProVersion(client *jamfpro.Client) (warning string, err error) {

	versionResp, err := client.GetJamfProVersion()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve Jamf Pro version: %w", err)
	}

	if versionResp == nil || versionResp.Version == nil {
		return "", fmt.Errorf("received empty version response from Jamf Pro")
	}

	// Parse the version string to extract the semantic version
	// Format is like "11.22.1-t1762179835791", we only care about "11.22.1"
	fullVersion := *versionResp.Version
	versionParts := strings.Split(fullVersion, "-")
	if len(versionParts) == 0 {
		return "", fmt.Errorf("invalid version format received: %s", fullVersion)
	}

	semanticVersion := versionParts[0]

	currentVersion, err := version.NewVersion(semanticVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse Jamf Pro version '%s': %w", semanticVersion, err)
	}

	minVersion, err := version.NewVersion(MinimumSupportedJamfProVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse minimum supported version '%s': %w", MinimumSupportedJamfProVersion, err)
	}

	if currentVersion.LessThan(minVersion) {
		return fmt.Sprintf(
			"The detected Jamf Pro version is (%s), which is below the minimum supported version of this provider (%s). "+
				"Some provider features may not function correctly with your version of Jamf Pro and could cause issues with state management. "+
				"You are strongly advised to use a compatible version of the provider for your specific Jamf Pro instance version. "+
				"Proceed at your own risk.",
			semanticVersion,
			MinimumSupportedJamfProVersion,
		), nil
	}

	return "", nil
}
