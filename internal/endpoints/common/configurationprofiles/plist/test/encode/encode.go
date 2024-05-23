package main

import (
	"log"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles"
)

func main() {
	// Example usage
	plistData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
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
		<key>PayloadScope</key>
		<string>System</string>
		<key>PayloadIdentifier</key>
		<string>Dafydds-MacBook-Air.5663C461-5EE7-427E-BAED-FA62B1AE88E8</string>
		<key>PayloadType</key>
		<string>Configuration</string>
		<key>PayloadUUID</key>
		<string>5663C461-5EE7-427E-BAED-FA62B1AE88E8</string>
		<key>PayloadVersion</key>
		<integer>1</integer>
	</dict>
	</plist>`)

	decodedData, err := configurationprofiles.DecodePlist(plistData)
	if err != nil {
		log.Fatalf("Failed to decode plist: %v", err)
	}

	fieldsToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"}

	configurationprofiles.RemoveFields(decodedData, fieldsToRemove, "")

	sortedData := configurationprofiles.SortPlistKeys(decodedData)

	log.Printf("Data structure before encoding: %v\n", sortedData)

	encodedPlist, err := configurationprofiles.EncodePlist(sortedData)
	if err != nil {
		log.Fatalf("Failed to encode plist: %v", err)
	}

	log.Printf("Sorted and encoded plist data: %s\n", encodedPlist)
}
