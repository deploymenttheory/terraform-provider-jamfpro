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

func waitForPackageAvailability(ctx context.Context, client *jamfpro.Client, packageID int, timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err := client.GetPackageByID(packageID)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				// Package not yet available, retryable error
				return retry.RetryableError(fmt.Errorf("package ID %d not available yet", packageID))
			}
			// Non-retryable error
			return retry.NonRetryableError(err)
		}

		// Package found, no need to retry
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
