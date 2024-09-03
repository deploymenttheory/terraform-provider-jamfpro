package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {
	client, err := jamfpro.BuildClientWithEnv()
	if err != nil {
		log.Fatalf("Failed to initialize Jamf Pro client: %v", err)
	}

	appInstallers, err := client.GetJamfAPIPrivileges()
	if err != nil {
		log.Fatalf("Error fetching Jamf API Privileges list: %v", err)
	}

	apiPrivilegesJSON, err := json.MarshalIndent(appInstallers, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling Jamf API Privileges list data: %v", err)
	}

	err = os.WriteFile("api_privileges.json", apiPrivilegesJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("App Catalog data written to api_privileges.json")
}
