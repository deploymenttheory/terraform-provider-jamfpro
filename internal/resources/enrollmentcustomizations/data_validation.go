package enrollmentcustomizations

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validatePaneCombinations ensures that only valid combinations of pane types are configured
func validatePaneCombinations(d *schema.ResourceData) error {
	hasSSO := len(d.Get("sso_pane").([]interface{})) > 0
	hasLDAP := len(d.Get("ldap_pane").([]interface{})) > 0

	if hasSSO && hasLDAP {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return fmt.Errorf("invalid combination: SSO and LDAP panes cannot be used together")
	}

	return nil
}

// validateHexColor validates a hex color code without the # prefix
func validateHexColor(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	// Check length (6 characters for hex without #)
	if len(v) != 6 {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		errs = append(errs, fmt.Errorf("%q must be exactly 6 characters (without #), got: %d characters", key, len(v)))
		return
	}

	// Check if all characters are valid hex
	for _, c := range v {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			errs = append(errs, fmt.Errorf("%q contains invalid hex character: %c", key, c))
			return
		}
	}

	return
}
