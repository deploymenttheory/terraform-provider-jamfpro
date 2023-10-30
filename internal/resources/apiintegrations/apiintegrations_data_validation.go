// apiintegrations_data_validation.go
package apiintegrations

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// the authorization_scopes attribute always contains at least one value.
func validateResourceAPIIntegrationsDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Extract the authorization_scopes list
	scopes, ok := diff.GetOk("authorization_scopes")
	if !ok || len(scopes.([]interface{})) == 0 {
		return fmt.Errorf("authorization_scopes must include at least one authorization scope")
	}

	return nil
}
