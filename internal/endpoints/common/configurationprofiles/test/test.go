package main

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles"
)

// TestProcessConfigurationProfile tests the ProcessConfigurationProfile function
func TestProcessConfigurationProfile() {
	originalPlist := `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	<dict>
		<key>PayloadContent</key>
		<array>
			<dict>
				<key>PayloadDisplayName</key>
				<string>Relay #1</string>
				<key>PayloadIdentifier</key>
				<string>com.apple.relay.managed.BC5A86E6-338F-4AAE-81BB-E957661345B5</string>
				<key>PayloadType</key>
				<string>com.apple.relay.managed</string>
				<key>PayloadUUID</key>
				<string>BC5A86E6-338F-4AAE-81BB-E957661345B5</string>
				<key>PayloadVersion</key>
				<integer>1</integer>
				<key>Relays</key>
				<array/>
			</dict>
			<dict>
				<key>OPPrefBiometryAllowed</key>
				<true/>
				<key>OPPreferencesNotifyCompromisedWebsites</key>
				<true/>
				<key>OPPreferencesNotifyOfTOTPCopy</key>
				<true/>
				<key>OPPreferencesNotifyVaultAddedRemoved</key>
				<true/>
				<key>PayloadDisplayName</key>
				<string>1Password 7</string>
				<key>PayloadIdentifier</key>
				<string>com.agilebits.onepassword7.C9F016D2-2813-48E2-9052-6426AB10D470</string>
				<key>PayloadType</key>
				<string>com.agilebits.onepassword7</string>
				<key>PayloadUUID</key>
				<string>C9F016D2-2813-48E2-9052-6426AB10D470</string>
				<key>PayloadVersion</key>
				<integer>1</integer>
			</dict>
		</array>
		<key>PayloadDisplayName</key>
		<string>Untitled</string>
		<key>PayloadIdentifier</key>
		<string>Dafydds-MacBook-Air.5663C461-5EE7-427E-BAED-FA62B1AE88E8</string>
		<key>PayloadType</key>
		<string>Configuration</string>
		<key>PayloadUUID</key>
		<string>5663C461-5EE7-427E-BAED-FA62B1AE88E8</string>
		<key>PayloadVersion</key>
		<integer>1</integer>
	</dict>
	</plist>`

	keysToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization"}
	processedProfile, err := configurationprofiles.ProcessConfigurationProfile(originalPlist, keysToRemove)
	if err != nil {
		log.Fatalf("Error processing configuration profile: %v\n", err)
	}

	fmt.Printf("Processed profile payload:\n%s\n", processedProfile)
}

func main() {
	TestProcessConfigurationProfile()
}
