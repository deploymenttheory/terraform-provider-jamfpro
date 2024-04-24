// webhooks_data_validation.go
package webhooks

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mainCustomDiffFunc orchestrates all custom diff validations.
func mainCustomDiffFunc(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	// Validate based on 'authentication_type' attribute.
	if err := validateAuthenticationRequirements(ctx, diff, i); err != nil {
		return err
	}

	// Validate based on 'event' attribute and the requirement for 'smart_group_id'.
	if err := validateSmartGroupIDRequirement(ctx, diff, i); err != nil {
		return err
	}

	return nil
}

// validateAuthenticationRequirements checks the conditions related to the 'authentication_type' attribute.
func validateAuthenticationRequirements(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	authType, ok := diff.GetOk("authentication_type")
	if !ok || authType.(string) != "Basic Authentication" {
		// If not using Basic Authentication or attribute not set, no need to validate further.
		return nil
	}

	// Check for existence of 'username' and 'password' when 'Basic Authentication' is used.
	username, usernameOk := diff.GetOk("username")
	password, passwordOk := diff.GetOk("password")

	if !usernameOk || username == "" {
		return fmt.Errorf("when 'authentication_type' is set to 'Basic Authentication', 'username' must be provided")
	}
	if !passwordOk || password == "" {
		return fmt.Errorf("when 'authentication_type' is set to 'Basic Authentication', 'password' must be provided")
	}

	return nil
}

// validateSmartGroupIDRequirement checks if the specified events require a smart_group_id and validates its presence.
func validateSmartGroupIDRequirement(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	event, ok := diff.GetOk("event")
	if !ok {
		return nil // Event not set, no further validation needed
	}

	// List of events that require a smart_group_id
	requiredEvents := []string{
		"SmartGroupComputerMembershipChange",
		"SmartGroupMobileDeviceMembershipChange",
		"SmartGroupUserMembershipChange",
	}

	// Check if the current event is in the list of required events
	for _, reqEvent := range requiredEvents {
		if event.(string) == reqEvent {
			smartGroupID, smartGroupIDOk := diff.GetOk("smart_group_id")
			if !smartGroupIDOk || smartGroupID == 0 {
				return fmt.Errorf("when 'event' is set to '%s', 'smart_group_id' must be provided and must be a valid non-zero integer", event)
			}
			break
		}
	}

	return nil
}
