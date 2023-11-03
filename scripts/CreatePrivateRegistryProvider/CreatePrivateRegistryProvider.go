package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	tfe "github.com/hashicorp/go-tfe"
)

func main() {
	ctx := context.Background()

	// Configure the Terraform Enterprise client
	config := &tfe.Config{
		Address: os.Getenv("TFE_ADDRESS"), // e.g., "https://app.terraform.io"
		Token:   os.Getenv("TFE_TOKEN"),   // Your Terraform Cloud API token
	}

	// Initialize the Terraform Enterprise client
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating TFE client: %v", err)
	}

	// Create a new private registry provider
	providerOptions := tfe.RegistryProviderCreateOptions{
		Type:         "registry-providers", // This is set by the SDK and is not user-defined
		Name:         "terraform-provider-jamfpro",
		Namespace:    "deploymenttheory",  // Same as your organization name for private providers
		RegistryName: tfe.PrivateRegistry, // This is for a private provider
	}
	newProvider, err := client.RegistryProviders.Create(ctx, "deploymenttheory", providerOptions)
	if err != nil {
		log.Fatalf("Error creating provider: %v", err)
	}

	// Create a new provider version
	versionOptions := tfe.RegistryProviderVersionCreateOptions{
		Version: "1.0.0", // Specify the version of your provider
	}
	newVersion, err := client.RegistryProviderVersions.Create(ctx, tfe.RegistryProviderID{
		OrganizationName: "deploymenttheory",
		RegistryName:     tfe.PrivateRegistry,
		Namespace:        newProvider.Namespace,
		Name:             newProvider.Name,
	}, versionOptions)
	if err != nil {
		log.Fatalf("Error creating provider version: %v", err)
	}

	// Assign a platform to the provider version
	platformOptions := tfe.RegistryProviderVersionPlatformCreateOptions{
		Os:   "linux",
		Arch: "amd64",
	}
	newPlatform, err := client.RegistryProviderVersions.CreatePlatform(ctx, newVersion.ID, platformOptions)
	if err != nil {
		log.Fatalf("Error assigning platform to provider version: %v", err)
	}

	// Upload the provider binary
	// The upload URL is provided by the newVersion object.
	if newVersion.UploadLink == "" {
		log.Fatal("No upload link was provided for the provider version")
	}

	// The path to your provider binary
	filePath := "/path/to/provider/binary"
	err = uploadProviderBinary(ctx, newVersion.UploadLink, filePath)
	if err != nil {
		log.Fatalf("Error uploading provider binary: %v", err)
	}

	log.Printf("Provider %s version %s for platform %s/%s uploaded successfully.",
		newProvider.Name, newVersion.Version, newPlatform.Os, newPlatform.Arch)
}

// uploadProviderBinary uploads the provider binary to the given URL.
func uploadProviderBinary(ctx context.Context, uploadURL string, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new request
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, file)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		return errors.New("upload failed with status code: " + resp.Status)
	}

	return nil
}
