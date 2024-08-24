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

	appInstallers, err := client.GetJamfAppCatalogAppInstallerTitles("")
	if err != nil {
		log.Fatalf("Error fetching Jamf App Catalog App Installer list: %v", err)
	}

	appInstallerJSON, err := json.MarshalIndent(appInstallers, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling Jamf App Catalog App Installer list data: %v", err)
	}

	err = os.WriteFile("app_catalog_app_installer_titles.json", appInstallerJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("App Catalog data written to app_catalog_app_installer_titles.json")
}
