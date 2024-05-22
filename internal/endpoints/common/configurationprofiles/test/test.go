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
				<key>AutoJoin</key>
				<true/>
				<key>CaptiveBypass</key>
				<true/>
				<key>DisableAssociationMACRandomization</key>
				<false/>
				<key>EAPClientConfiguration</key>
				<dict>
					<key>AcceptEAPTypes</key>
					<array>
						<integer>13</integer>
						<integer>17</integer>
						<integer>43</integer>
						<integer>21</integer>
						<integer>25</integer>
						<integer>18</integer>
						<integer>23</integer>
					</array>
					<key>EAPFASTProvisionPAC</key>
					<true/>
					<key>EAPFASTProvisionPACAnonymously</key>
					<true/>
					<key>EAPFASTUsePAC</key>
					<true/>
					<key>TLSAllowTrustExceptions</key>
					<true/>
					<key>TTLSInnerAuthentication</key>
					<string>MSCHAPv2</string>
				</dict>
				<key>EncryptionType</key>
				<string>WPA3</string>
				<key>HIDDEN_NETWORK</key>
				<true/>
				<key>PayloadDescription</key>
				<string/>
				<key>PayloadDisplayName</key>
				<string>WiFi (thing)</string>
				<key>PayloadEnabled</key>
				<true/>
				<key>PayloadType</key>
				<string>com.apple.wifi.managed</string>
				<key>PayloadVersion</key>
				<integer>1</integer>
				<key>ProxyPassword</key>
				<string>thing</string>
				<key>ProxyServer</key>
				<string>server.com:443</string>
				<key>ProxyServerPort</key>
				<integer>0</integer>
				<key>ProxyType</key>
				<string>Manual</string>
				<key>ProxyUsername</key>
				<string>thing</string>
				<key>SSID_STR</key>
				<string>thing</string>
			</dict>
		</array>
		<key>PayloadDescription</key>
		<string/>
		<key>PayloadDisplayName</key>
		<string>mobile-wifi</string>
		<key>PayloadEnabled</key>
		<true/>
		<key>PayloadRemovalDisallowed</key>
		<true/>
		<key>PayloadScope</key>
		<string>System</string>
		<key>PayloadType</key>
		<string>Configuration</string>
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
