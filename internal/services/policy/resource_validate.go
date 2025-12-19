package policy

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateDateTime validates the input string is in the format 'YYYY-MM-DD HH:mm:ss'
func validateDateTime(v any, k string) (warns []string, errs []error) {
	value := v.(string)
	if _, err := time.Parse("2006-01-02 15:04:05", value); err != nil {
		errs = append(errs, fmt.Errorf("%q must be in the format 'YYYY-MM-DD HH:mm:ss', got: %s", k, value))
	}
	return
}

// validateDateTimeUTC validates the input string is in the format 'YYYY-MM-DDThh:mm:ss.sss+0000'
func validateDateTimeUTC(v any, k string) (warns []string, errs []error) {
	value := v.(string)
	if _, err := time.Parse("2006-01-02T15:04:05.000-0700", value); err != nil {
		errs = append(errs, fmt.Errorf("%q must be in the format 'YYYY-MM-DDThh:mm:ss.sss+0000', got: %s", k, value))
	}
	return
}

// validateEpochMillis validates the input integer is a positive number
func validateEpochMillis(v any, k string) (warns []string, errs []error) {
	value := v.(int)
	if value < 0 {
		errs = append(errs, fmt.Errorf("%q must be a positive integer, got: %d", k, value))
	}
	return
}

// validateDayOfWeek validates the input string is a valid day of the week
func validate12HourTime(v any, k string) (warns []string, errs []error) {
	value := v.(string)
	pattern := regexp.MustCompile(`^(1[0-2]|0?[1-9]):[0-5][0-9] (AM|PM)$`)
	if !pattern.MatchString(value) {
		errs = append(errs, fmt.Errorf("%q must be in 12-hour format (h:mm AM/PM), got: %s", k, value))
	}
	return
}

// validateSelfServiceConfig validates that when use_for_self_service is false,
// no other self service attributes should be configured with non-default values
func validateSelfServiceConfig(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if selfService, ok := d.GetOk("self_service"); ok {
		selfServiceList := selfService.([]interface{})
		if len(selfServiceList) > 0 {
			selfServiceMap := selfServiceList[0].(map[string]interface{})
			useForSelfService := selfServiceMap["use_for_self_service"].(bool)

			// If self service is disabled, check that no other attributes are set to non-default values
			if !useForSelfService {
				defaults := map[string]interface{}{
					"self_service_display_name":       "",
					"install_button_text":             "Install",
					"reinstall_button_text":           "Reinstall",
					"self_service_description":        "",
					"force_users_to_view_description": false,
					"self_service_icon_id":            0,
					"feature_on_main_page":            false,
					"notification":                    false,
					"notification_type":               "Self Service",
					"notification_subject":            "",
					"notification_message":            "",
				}

				var nonDefaultAttrs []string
				for key, defaultValue := range defaults {
					if actualValue, exists := selfServiceMap[key]; exists && actualValue != defaultValue {
						nonDefaultAttrs = append(nonDefaultAttrs, key)
					}
				}

				// Check for self_service_category
				if categories, exists := selfServiceMap["self_service_category"]; exists {
					if categoriesList, ok := categories.([]interface{}); ok && len(categoriesList) > 0 {
						nonDefaultAttrs = append(nonDefaultAttrs, "self_service_category")
					}
				}

				if len(nonDefaultAttrs) > 0 {
					return fmt.Errorf("when use_for_self_service is false, the following self_service attributes should not be configured: %v. Either set use_for_self_service to true or remove these attributes", nonDefaultAttrs)
				}
			}
		}
	}
	return nil
}
