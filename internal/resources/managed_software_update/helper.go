package managed_software_update

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// checkAndEnableManagedSoftwareUpdateFeatureToggle checks the status of the Managed Software Update
// Feature Toggle and enables it if it's not already enabled.
func checkAndEnableManagedSoftwareUpdateFeatureToggle(ctx context.Context, client *jamfpro.Client) error {
	status, err := client.GetManagedSoftwareUpdateFeatureToggle()
	if err != nil {
		return fmt.Errorf("failed to fetch Managed Software Update Feature Toggle status: %v", err)
	}
	fmt.Printf("Fetched Feature Toggle Status: %+v\n", status)

	// If feature toggle is already enabled, no further action is needed
	if status.Toggle {
		return nil
	}

	// Enable the feature toggle
	_, err = client.UpdateManagedSoftwareUpdateFeatureToggle(&jamfpro.ResourceManagedSoftwareUpdateFeatureToggle{Toggle: true})
	if err != nil {
		return fmt.Errorf("failed to enable Managed Software Update Feature Toggle: %v", err)
	}

	const maxRetries = 5
	var retryInterval = 2 * time.Second

	return retry.RetryContext(ctx, time.Duration(maxRetries)*retryInterval, func() *retry.RetryError {
		status, err := client.GetManagedSoftwareUpdateFeatureToggle()
		if err != nil {
			return retry.RetryableError(fmt.Errorf("failed to fetch Managed Software Update Feature Toggle status: %v", err))
		}
		fmt.Printf("Retry Fetched Feature Toggle Status: %+v\n", status)

		if !status.Toggle {
			time.Sleep(retryInterval)
			retryInterval *= 2
			return retry.RetryableError(fmt.Errorf("managed Software Update Feature Toggle is not yet enabled"))
		}

		return nil
	})
}
