package main

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/configurationprofiles/plist"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers"
)

// plistPreJamfProUpload is the plist data for the configuration profile prior to being uploaded into Jamf Pro.
// plistPostJamfProUpload is the plist data for the configuration profile after being uploaded into Jamf Pro.

// TestProcessAndCompareConfigurationProfiles processes and compares two configuration profiles.
func TestProcessAndCompareConfigurationProfiles() {
	plistPreJamfProUpload := `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
		<dict>
			<key>PayloadContent</key>
			<array>
				<dict>
					<key>AllowCacheDelete</key>
					<true/>
					<key>AllowPersonalCaching</key>
					<true/>
					<key>AllowSharedCaching</key>
					<true/>
					<key>AutoActivation</key>
					<true/>
					<key>AutoEnableTetheredCaching</key>
					<true/>
					<key>CacheLimit</key>
					<integer>1000</integer>
					<key>DataPath</key>
					<string>/Library/Application Support/Apple/AssetCache/Data</string>
					<key>DisplayAlerts</key>
					<true/>
					<key>KeepAwake</key>
					<true/>
					<key>ListenRangesOnly</key>
					<false/>
					<key>LocalSubnetsOnly</key>
					<true/>
					<key>LogClientIdentity</key>
					<true/>
					<key>ParentSelectionPolicy</key>
					<string>round-robin</string>
					<key>PayloadDescription</key>
					<string/>
					<key>PayloadDisplayName</key>
					<string>Content Caching</string>
					<key>PayloadEnabled</key>
					<true/>
					<key>PayloadIdentifier</key>
					<string>CC8EF85E-7C43-44BD-9736-7D1748FF1D9F</string>
					<key>PayloadOrganization</key>
					<string>Lloyds Bank</string>
					<key>PayloadType</key>
					<string>com.apple.AssetCache.managed</string>
					<key>PayloadUUID</key>
					<string>CC8EF85E-7C43-44BD-9736-7D1748FF1D9F</string>
					<key>PayloadVersion</key>
					<integer>1</integer>
					<key>PeerLocalSubnetsOnly</key>
					<true/>
					<key>Port</key>
					<integer>443</integer>
				</dict>
			</array>
			<key>PayloadDescription</key>
			<string/>
			<key>PayloadDisplayName</key>
			<string>content-caching</string>
			<key>PayloadEnabled</key>
			<true/>
			<key>PayloadIdentifier</key>
			<string>8F5784BA-1A97-4F54-AB90-FCA5FA051194</string>
			<key>PayloadOrganization</key>
			<string>Lloyds Bank</string>
			<key>PayloadRemovalDisallowed</key>
			<true/>
			<key>PayloadScope</key>
			<string>System</string>
			<key>PayloadType</key>
			<string>Configuration</string>
			<key>PayloadUUID</key>
			<string>8F5784BA-1A97-4F54-AB90-FCA5FA051194</string>
			<key>PayloadVersion</key>
			<integer>1</integer>
		</dict>
	</plist>
	`

	plistPostJamfProUpload := `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
		<dict>
			<key>PayloadContent</key>
			<array>
				<dict>
					<key>AllowCacheDelete</key>
					<true/>
					<key>AllowPersonalCaching</key>
					<true/>
					<key>AllowSharedCaching</key>
					<true/>
					<key>AutoActivation</key>
					<true/>
					<key>AutoEnableTetheredCaching</key>
					<true/>
					<key>CacheLimit</key>
					<integer>1000</integer>
					<key>DataPath</key>
					<string>/Library/Application Support/Apple/AssetCache/Data</string>
					<key>DisplayAlerts</key>
					<true/>
					<key>KeepAwake</key>
					<true/>
					<key>ListenRangesOnly</key>
					<false/>
					<key>LocalSubnetsOnly</key>
					<true/>
					<key>LogClientIdentity</key>
					<true/>
					<key>ParentSelectionPolicy</key>
					<string>round-robin</string>
					<key>PayloadDescription</key>
					<string/>
					<key>PayloadDisplayName</key>
					<string>Content Caching</string>
					<key>PayloadEnabled</key>
					<true/>
					<key>PayloadIdentifier</key>
					<string>CC8EF85E-7C43-44BD-9736-7D1748FF1D9F</string>
					<key>PayloadOrganization</key>
					<string>Lloyds Bank</string>
					<key>PayloadType</key>
					<string>com.apple.AssetCache.managed</string>
					<key>PayloadUUID</key>
					<string>CC8EF85E-7C43-44BD-9736-7D1748FF1D9F</string>
					<key>PayloadVersion</key>
					<integer>1</integer>
					<key>PeerLocalSubnetsOnly</key>
					<true/>
					<key>Port</key>
					<integer>443</integer>
				</dict>
			</array>
			<key>PayloadDescription</key>
			<string/>
			<key>PayloadDisplayName</key>
			<string>content-caching</string>
			<key>PayloadEnabled</key>
			<true/>
			<key>PayloadIdentifier</key>
			<string>8F5784BA-1A97-4F54-AB90-FCA5FA051194</string>
			<key>PayloadOrganization</key>
			<string>Lloyds Bank</string>
			<key>PayloadRemovalDisallowed</key>
			<true/>
			<key>PayloadScope</key>
			<string>System</string>
			<key>PayloadType</key>
			<string>Configuration</string>
			<key>PayloadUUID</key>
			<string>8F5784BA-1A97-4F54-AB90-FCA5FA051194</string>
			<key>PayloadVersion</key>
			<integer>1</integer>
		</dict>
	</plist>`

	keysToRemove := []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"}

	log.Println("Processing the first configuration profile...")
	processedProfile1, err := plist.ProcessConfigurationProfileForDiffSuppression(plistPreJamfProUpload, keysToRemove)
	if err != nil {
		log.Fatalf("Error processing the first configuration profile: %v\n", err)
	}
	log.Printf("Processed first profile payload:\n%s\n", processedProfile1)

	log.Println("Processing the second configuration profile...")
	processedProfile2, err := plist.ProcessConfigurationProfileForDiffSuppression(plistPostJamfProUpload, keysToRemove)
	if err != nil {
		log.Fatalf("Error processing the second configuration profile: %v\n", err)
	}
	log.Printf("Processed second profile payload:\n%s\n", processedProfile2)

	log.Println("Hashing the first processed profile payload...")
	hash1 := helpers.HashString(processedProfile1)
	log.Printf("Hash of the first profile payload: %s\n", hash1)

	log.Println("Hashing the second processed profile payload...")
	hash2 := helpers.HashString(processedProfile2)
	log.Printf("Hash of the second profile payload: %s\n", hash2)

	log.Println("Comparing the hashes...")
	if hash1 == hash2 {
		log.Println("The hashes are identical.")
	} else {
		log.Println("The hashes are different.")
	}

	fmt.Printf("Processed profile 1 payload (pre jamf pro upload):\n%s\n", processedProfile1)
	fmt.Printf("Processed profile 2 payload (post jamf pro upload):\n%s\n", processedProfile2)
	fmt.Printf("Hash of profile 1: %s\n", hash1)
	fmt.Printf("Hash of profile 2: %s\n", hash2)
	if hash1 == hash2 {
		fmt.Println("The hashes are identical.")
	} else {
		fmt.Println("The hashes are different.")
	}
}

func main() {
	TestProcessAndCompareConfigurationProfiles()
}
