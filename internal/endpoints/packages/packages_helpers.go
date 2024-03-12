// packages_helpers.go
package packages

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"golang.org/x/crypto/sha3"
)

// waitForPackageAvailability repeatedly checks for the availability of a package by its name until it is found or a timeout occurs.
// Instead of checking for a specific package ID, it fetches the list of all packages using GetPackages and searches for the package by name.
// The function uses a retry mechanism to periodically check the packages list. If the package with the specified name is not found, it retries until the package becomes available or a timeout is reached.
// Any errors encountered while fetching the packages list are treated as retryable errors, except for context cancellation which stops the retry loop.
func waitForPackageAvailability(ctx context.Context, client *jamfpro.Client, packageName string, timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		fmt.Printf("Checking for package '%s'...\n", packageName)
		packagesList, err := client.GetPackages()
		if err != nil {
			// Treat errors during retrieval as retryable to keep trying
			return retry.RetryableError(fmt.Errorf("error retrieving packages: %v", err))
		}

		for _, pkg := range packagesList.Package {
			if pkg.Name == packageName {
				fmt.Printf("Package '%s' found with ID %d. Fetching package details...\n", packageName, pkg.ID)

				// Once the package is found by name, fetch its details by ID
				packageDetails, err := client.GetPackageByID(pkg.ID)
				if err != nil {
					// If there's an error fetching the package details, treat it as a retryable error
					return retry.RetryableError(fmt.Errorf("error fetching package '%s' details by ID %d: %v", packageName, pkg.ID, err))
				}

				// Log a simple detail from the packageDetails
				fmt.Printf("Retrieved details for package '%s' with ID %d. Filename: %s\n", packageName, pkg.ID, packageDetails.Filename)

				// Pause for 10 seconds after finding the package
				fmt.Println("Pausing for 30 seconds...")
				time.Sleep(30 * time.Second)
				fmt.Println("Resuming...")

				// If everything is okay with the package details, return success (nil error)
				return nil
			}
		}

		// Package not found yet, log this event and retry
		return retry.RetryableError(fmt.Errorf("package '%s' not found yet, retrying", packageName))
	})
}

// generateFileHash accepts a file path and returns a SHA-3-256 hash of the file's contents.
// It opens the file, creates a new SHA3-256 hash object, copies the file content into the hash object, and computes the SHA3-256 checksum of the file.
// AWS S3 uses the SHA-3-256 hash to verify the integrity of the file during upload and download.
func generateFileHash(filePath string) (string, error) {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// Create a new SHA3-256 hash object
	hash := sha3.New256()

	// Copy the file content into the hash object
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file contents of %s: %v", filePath, err)
	}

	// Compute the SHA3-256 checksum of the file
	hashBytes := hash.Sum(nil)

	// Convert the bytes to a hex string
	hashString := fmt.Sprintf("%x", hashBytes)

	return hashString, nil
}
