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
// 1. Continuously retry while the server hash is empty
// 2. Once a hash value is present, verify it against the expected hash
// 3. If deleteOnFailure is true and verification fails, attempt to delete the package
//
// If refreshInventory is true, each retry iteration will call RefreshCloudDistributionPointInventoryV1
// to trigger JCDS hash recomputation. This is needed for updates where an existing file is replaced,
// but not for creates where JCDS processes the file automatically.
func verifyPackageUpload(ctx context.Context, client *jamfpro.Client, packageID string, fileName string, expectedHash string, timeout time.Duration, deleteOnFailure bool, refreshInventory bool) error {
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		if refreshInventory {
			time.Sleep(5 * time.Second)
			if err := client.RefreshCloudDistributionPointInventoryV1(fileName); err != nil {
				log.Printf("[WARN] Failed to refresh JCDS inventory for %s during verification: %v", fileName, err)
			}
		}

		uploadedPackage, err := client.GetPackageByID(packageID)
		if err != nil {
			return retry.RetryableError(fmt.Errorf("failed to verify uploaded package: %v", err))
		}

		if uploadedPackage.HashValue == "" {
			return retry.RetryableError(fmt.Errorf("waiting for package hash calculation to complete"))
		}

		if uploadedPackage.HashType != "SHA3_512" || uploadedPackage.HashValue != expectedHash {
			return retry.RetryableError(fmt.Errorf("package hash verification pending: expected=%s, got=%s (type: %s)",
				expectedHash, uploadedPackage.HashValue, uploadedPackage.HashType))
		}

		log.Printf("[INFO] Package %s SHA3-512 verification successful with hash: %s",
			fileName, uploadedPackage.HashValue)
		return nil
	})

	if err != nil {
		if deleteOnFailure {
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
		}

		return fmt.Errorf("failed to verify package upload: %v", err)
	}

	return nil
}
