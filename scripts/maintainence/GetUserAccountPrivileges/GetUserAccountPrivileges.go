package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {
	client, err := jamfpro.BuildClientWithEnv()
	if err != nil {
		log.Fatalf("Failed to initialize Jamf Pro client: %v", err)
	}

	versionInfo, err := client.GetJamfProVersion()
	if err != nil {
		log.Fatalf("Error getting Jamf Pro version: %v", err)
	}

	if versionInfo.Version == nil {
		log.Fatalf("Received nil version from Jamf Pro")
	}

	versionDir := *versionInfo.Version

	err = os.MkdirAll(versionDir, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	accountDetail := &jamfpro.ResourceAccount{
		Name:          "export_available_privileges",
		DirectoryUser: false,
		FullName:      "Export Available Privileges",
		Email:         "export_privileges@example.com",
		EmailAddress:  "export_privileges@example.com",
		Enabled:       "Enabled",
		AccessLevel:   "Full Access",
		PrivilegeSet:  "Administrator",
		Password:      "SecurePassword123!",
	}

	createdAccount, err := client.CreateAccount(accountDetail)
	if err != nil {
		log.Fatalf("Error creating account: %v", err)
	}

	fmt.Printf("Created Full Access account with ID: %d\n", createdAccount.ID)

	fullAccountDetails, err := client.GetAccountByID(fmt.Sprintf("%d", createdAccount.ID))
	if err != nil {
		log.Fatalf("Error retrieving account details: %v", err)
	}

	accountXML, err := xml.MarshalIndent(fullAccountDetails, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling account details: %v", err)
	}
	fmt.Printf("Full Account Details:\n%s\n", string(accountXML))

	if fullAccountDetails.Privileges.JSSObjects == nil &&
		fullAccountDetails.Privileges.JSSSettings == nil &&
		fullAccountDetails.Privileges.JSSActions == nil {
		log.Println("Warning: All privilege fields are nil")
	}

	exportPrivilegesToJSON(fullAccountDetails.Privileges, versionDir)

	err = client.DeleteAccountByID(fmt.Sprintf("%d", createdAccount.ID))
	if err != nil {
		log.Fatalf("Error deleting account: %v", err)
	}

	fmt.Println("Full Access account successfully deleted. Export completed.")
}

func exportPrivilegesToJSON(privileges jamfpro.AccountSubsetPrivileges, versionDir string) {
	exportToJSONFile(privileges.JSSObjects, filepath.Join(versionDir, "jss_objects_privileges.json"))
	exportToJSONFile(privileges.JSSSettings, filepath.Join(versionDir, "jss_settings_privileges.json"))
	exportToJSONFile(privileges.JSSActions, filepath.Join(versionDir, "jss_actions_privileges.json"))
}

func exportToJSONFile(privileges []string, filename string) {
	if privileges == nil {
		log.Printf("Warning: Privileges are nil for file %s\n", filename)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(privileges)
	if err != nil {
		log.Fatalf("Error encoding privileges to JSON for %s: %v", filename, err)
	}

	fmt.Printf("Exported available Jamf Pro privileges to %s\n", filename)
}
