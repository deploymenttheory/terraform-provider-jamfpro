// packages_helpers.go
package packages

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// verifyPackageUpload attempts to verify a package upload by comparing SHA3-512 hashes.
// The function will:
// 1. Continuously retry while the server calculates the hash (empty hash value)
// 2. Once a hash value is present, perform a one-time verification against the expected hash
// 3. If verification fails or encounters an error, attempt to delete the package
// The function returns an error if verification fails, cleanup fails, or any other error occurs.
func verifyPackageUpload(ctx context.Context, client *jamfpro.Client, packageID string, fileName string, expectedHash string, timeout time.Duration) error {
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		uploadedPackage, err := client.GetPackageByID(packageID)
		if err != nil {
			return retry.RetryableError(fmt.Errorf("failed to verify uploaded package: %v", err))
		}

		if uploadedPackage.HashValue == "" {
			return retry.RetryableError(fmt.Errorf("waiting for package hash calculation to complete"))
		}

		if uploadedPackage.HashType != "SHA3_512" || uploadedPackage.HashValue != expectedHash {
			return retry.NonRetryableError(fmt.Errorf("package hash verification failed: expected=%s, got=%s (type: %s)",
				expectedHash, uploadedPackage.HashValue, uploadedPackage.HashType))
		}

		log.Printf("[INFO] Package %s SHA3-512 verification successful with hash: %s",
			fileName, uploadedPackage.HashValue)
		return nil
	})

	if err != nil {
		cleanupErr := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			if packageID == "" {
				return nil
			}

			if err := client.DeletePackageByID(packageID); err != nil {
				return retry.RetryableError(fmt.Errorf("failed to clean up package %s: %v", packageID, err))
			}

			log.Printf("[INFO] Successfully cleaned up package %s after verification failure", packageID)
			return nil
		})

		if cleanupErr != nil {
			log.Printf("[WARN] Failed to clean up package %s after verification failure: %v", packageID, cleanupErr)
		}

		return fmt.Errorf("failed to verify package upload: %v", err)
	}

	return nil
}
