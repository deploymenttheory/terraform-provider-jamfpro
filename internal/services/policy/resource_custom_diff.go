package policy

import (
	"context"
	"fmt"
	"strings"

	sharedschemas "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/shared_schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations for Policies.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i any) error {
	if err := sharedschemas.ValidateScopeDirectoryServiceUserGroupNames(ctx, diff, i); err != nil {
		return fmt.Errorf("validating scope directory service user/group names: %w", err)
	}

	return nil
}

func suppressInactiveSelfServiceNotificationDiff(k, old, new string, d *schema.ResourceData) bool {
	notificationKey := siblingSelfServiceFieldKey(k, "notification")
	if notificationKey == "" {
		return false
	}

	notificationEnabled, ok := d.Get(notificationKey).(bool)
	if !ok {
		return false
	}

	return !notificationEnabled
}

func siblingSelfServiceFieldKey(k, field string) string {
	fieldSeparator := strings.LastIndex(k, ".")
	if fieldSeparator == -1 {
		return ""
	}

	return k[:fieldSeparator+1] + field
}
