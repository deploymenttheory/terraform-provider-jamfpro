package appinstallers

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// checkJamfAppCatalogAppInstallerTermsAndConditions checks and accepts the terms and conditions for
// the Jamf App Catalog App Installer if it not enabled. this is required per account trying to interact
// with the Jamf App Catalog App Installer.
func checkJamfAppCatalogAppInstallerTermsAndConditions(ctx context.Context, client *jamfpro.Client) error {
	status, err := client.GetJamfAppCatalogAppInstallerTermsAndConditionsStatus()
	if err != nil {

		return fmt.Errorf("failed to fetch Jamf Pro App Installer terms and conditions status: %v", err)
	}
	fmt.Printf("Fetched Terms and Conditions Status: %+v\n", status)

	// If terms and conditions are already accepted, no further action is needed
	if status.Accepted {
		return nil
	}

	_, err = client.AcceptJamfAppCatalogAppInstallerTermsAndConditions()
	if err != nil {

		return fmt.Errorf("failed to accept Jamf Pro App Installer terms and conditions: %v", err)
	}

	const maxRetries = 5
	var retryInterval = 2 * time.Second

	return retry.RetryContext(ctx, time.Duration(maxRetries)*retryInterval, func() *retry.RetryError {
		status, err := client.GetJamfAppCatalogAppInstallerTermsAndConditionsStatus()
		if err != nil {

			return retry.RetryableError(fmt.Errorf("failed to fetch Jamf Pro App Installer terms and conditions status: %v", err))
		}
		fmt.Printf("Retry Fetched Terms and Conditions Status: %+v\n", status)

		if !status.Accepted {
			time.Sleep(retryInterval)
			retryInterval *= 2

			return retry.RetryableError(fmt.Errorf("terms and conditions are not yet accepted"))
		}

		return nil
	})
}
