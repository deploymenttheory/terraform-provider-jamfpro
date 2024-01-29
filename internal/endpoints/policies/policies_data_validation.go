// policies_data_validation.go
package policies

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateJamfProResourcePolicyDataFields ensures that a Jamf Pro policy meets specific criteria:
// 1. It cannot be set to have an ongoing frequency when the trigger_checkin is true and it applies to all computers.
// 2. The 'offline' field can only be set to true if the 'frequency' is "Ongoing" when 'all_computers' is true.
// 3. The 'retry_event', 'retry_attempts', and 'notify_on_each_failed_retry' fields can only be set if 'frequency' is "Once per computer".
// 4. The 'notify_on_each_failed_retry' field must be false if 'retry_attempts' is -1.
// 5. If 'retry_event' is not "none", then 'retry_attempts' must be between 1 and 10.
// 6. If 'retry_event' is "none", then 'retry_attempts' must be -1 and 'notify_on_each_failed_retry' must be false.
// 7. If 'trigger' is "USER_INITIATED", then 'use_for_self_service' must be true.
// 8. If 'trigger' is "EVENT", then 'use_for_self_service' must be false.
func validateJamfProResourcePolicyDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Validatation scenario 01
	triggerCheckin, triggerCheckinOk := diff.GetOk("general.0.trigger_checkin")
	frequency, frequencyOk := diff.GetOk("general.0.frequency")
	allComputers, allComputersOk := diff.GetOk("scope.0.all_computers")

	if triggerCheckinOk && frequencyOk && allComputersOk {
		if triggerCheckin.(bool) && frequency.(string) == "Ongoing" && allComputers.(bool) {
			return fmt.Errorf("jamf pro policies that update inventory on all computers cannot be set to ongoing frequency at recurring check-in. Please update your terraform configuration")
		}
	}

	// Validatation scenario 02
	offline, offlineOk := diff.GetOk("general.0.offline")
	if offlineOk && offline.(bool) && allComputersOk && allComputers.(bool) {
		if frequencyOk && frequency.(string) != "Ongoing" {
			return fmt.Errorf("jamf pro policies that apply to all computers can only be set to 'offline' if the policy frequency is set to 'Ongoing'. Please update your terraform configuration")
		}
	}

	// Validatation scenario 03
	frequency, frequencyOk = diff.GetOk("general.0.frequency")
	retryEvent, retryEventOk := diff.GetOk("general.0.retry_event")
	retryAttempts, retryAttemptsOk := diff.GetOk("general.0.retry_attempts")
	notifyOnEachFailedRetry, notifyOnEachFailedRetryOk := diff.GetOk("general.0.notify_on_each_failed_retry")

	if frequencyOk && frequency.(string) != "Once per computer" {
		if retryEventOk && retryEvent.(string) != "none" {
			return fmt.Errorf("the 'retry_event' field can only be set if the policy 'frequency' is 'Once per computer'")
		}
		if retryAttemptsOk && retryAttempts.(int) != -1 {
			return fmt.Errorf("the 'retry_attempts' field can only be set if the policy 'frequency' is 'Once per computer'")
		}
		if notifyOnEachFailedRetryOk && notifyOnEachFailedRetry.(bool) {
			return fmt.Errorf("the 'notify_on_each_failed_retry' field can only be set if the policy 'frequency' is 'Once per computer'")
		}
	}

	// Validatation scenario 04
	retryAttempts, retryAttemptsOk = diff.GetOk("general.0.retry_attempts")
	notifyOnEachFailedRetry, notifyOnEachFailedRetryOk = diff.GetOk("general.0.notify_on_each_failed_retry")

	if retryAttemptsOk && notifyOnEachFailedRetryOk {
		if retryAttempts.(int) == -1 && notifyOnEachFailedRetry.(bool) {
			return fmt.Errorf("the 'notify_on_each_failed_retry' field must be false if the 'retry_attempts' field is set to -1. Please update your terraform configuration")
		}
	}

	// Validatation scenario 05
	retryEvent, retryEventOk = diff.GetOk("general.0.retry_event")
	retryAttempts, retryAttemptsOk = diff.GetOk("general.0.retry_attempts")

	if retryEventOk && retryAttemptsOk && retryEvent.(string) != "none" {
		if retryAttemptsVal, ok := retryAttempts.(int); !ok || retryAttemptsVal < 1 || retryAttemptsVal > 10 {
			return fmt.Errorf("when 'retry_event' is not 'none', then 'retry_attempts' must be a value between 1 and 10. Please update your terraform configuration")
		}
	}

	// Validation scenario 06
	retryEvent, retryEventOk = diff.GetOk("general.0.retry_event")
	if retryEventOk && retryEvent.(string) == "none" {
		retryAttempts, retryAttemptsOk := diff.GetOk("general.0.retry_attempts")
		notifyOnEachFailedRetry, notifyOnEachFailedRetryOk := diff.GetOk("general.0.notify_on_each_failed_retry")

		var errMsg string

		if retryAttemptsOk && retryAttempts.(int) != -1 {
			errMsg += "When 'retry_event' is 'none', then 'retry_attempts' must be -1. "
		}
		if notifyOnEachFailedRetryOk && notifyOnEachFailedRetry.(bool) {
			errMsg += "When 'retry_event' is 'none', then 'notify_on_each_failed_retry' must be false. "
		}

		if errMsg != "" {
			return fmt.Errorf(errMsg + "Please update your terraform configuration")
		}
	}

	// Validation scenario 07
	trigger, triggerOk := diff.GetOk("general.0.trigger")
	useForSelfService, useForSelfServiceOk := diff.GetOk("self_service.0.use_for_self_service")

	if triggerOk && trigger.(string) == "EVENT" && useForSelfServiceOk && useForSelfService.(bool) {
		return fmt.Errorf("when 'trigger' is 'EVENT', 'use_for_self_service' must be false. Please update your terraform configuration")
	}

	// Validation scenario 08
	trigger, triggerOk = diff.GetOk("general.0.trigger")
	useForSelfService, useForSelfServiceOk = diff.GetOk("self_service.0.use_for_self_service")

	if triggerOk && trigger.(string) == "USER_INITIATED" && (!useForSelfServiceOk || !useForSelfService.(bool)) {
		return fmt.Errorf("when 'trigger' is 'USER_INITIATED', 'use_for_self_service' must be true. Please update your terraform configuration")
	}

	return nil
}
