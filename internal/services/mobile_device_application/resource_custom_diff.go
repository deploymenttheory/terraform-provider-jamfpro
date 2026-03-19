package mobile_device_application

import (
	"context"
	"fmt"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations for Mobile Device Applications.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i any) error {
	if err := sharedschemas.ValidateScopeDirectoryServiceUserGroupNames(ctx, diff, i); err != nil {
		return fmt.Errorf("validating scope directory service user/group names: %w", err)
	}

	return nil
}
