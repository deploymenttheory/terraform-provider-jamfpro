// packages_helpers.go
package packages

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"golang.org/x/crypto/sha3"
)

// waitForPackageAvailability repeatedly checks for the availability of a package by its ID until it is found or a timeout occurs.
// It uses a retry mechanism to periodically make the GetPackageByID call. If the package is not found (404 error), it retries until the package becomes available.
// Any other error will immediately stop the retry loop and return the error.
// Console printouts are included for status updates and debugging.
func waitForPackageAvailability(ctx context.Context, client *jamfpro.Client, packageID int, timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		fmt.Printf("Checking availability for package ID %d...\n", packageID)
		_, err := client.GetPackageByID(packageID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				// Package not yet available, log this event and retry
				fmt.Printf("Package ID %d is not available yet, retrying...\n", packageID)
				return retry.RetryableError(fmt.Errorf("package ID %d not available yet", packageID))
			}
			// Log the non-retryable error and return it
			fmt.Printf("Encountered a non-retryable error while checking package ID %d: %s\n", packageID, err)
			return retry.NonRetryableError(err)
		}

		// Package found, log success and return nil to stop retrying
		fmt.Println("Package found: ID", packageID)
		return nil
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
