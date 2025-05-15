// apiintegrations_data_validation.go
package apiintegrations

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateResourceAPIIntegrationsDataFields ensures that the authorization_scopes attribute always contains at least one value.
func validateResourceAPIIntegrationsDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	scopes, ok := diff.GetOk("authorization_scopes")
	if !ok {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return fmt.Errorf("authorization_scopes must be provided")
	}

	scopesSet, ok := scopes.(*schema.Set)
	if !ok {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return fmt.Errorf("authorization_scopes is not a valid set")
	}

	scopesList := scopesSet.List()
	if len(scopesList) == 0 {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return fmt.Errorf("authorization_scopes must include at least one authorization scope")
	}

	return nil
}
