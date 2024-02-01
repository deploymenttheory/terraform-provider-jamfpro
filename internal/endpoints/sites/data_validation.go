// sites_data_validation.go
package sites

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateJamfProResourceSiteDataFields(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Use GetOk to safely check if the name attribute is set and to get its value.
	proposedName, nameIsSet := diff.GetOk("name")
	if !nameIsSet {
		// The name attribute is not set in the configuration, so no validation is needed.
		return nil
	}

	// Convert proposedName to string safely.
	proposedNameStr, ok := proposedName.(string)
	if !ok {
		return fmt.Errorf("expected name to be a string, got: %T", proposedName)
	}

	// Perform duplicate checks only for new resources or if the name has changed.
	if diff.Id() == "" || diff.HasChange("name") {
		apiclient, ok := v.(*client.APIClient)
		if !ok {
			return fmt.Errorf("error asserting meta as *client.APIClient")
		}
		conn := apiclient.Conn

		// Use the API client to check if a site with the proposed name already exists.
		existingSite, err := conn.GetSiteByName(proposedNameStr)
		if err != nil {
			// Optionally, handle not found errors differently if your API client distinguishes them.
			return fmt.Errorf("error checking if site exists: %s", err)
		}

		// If an existing site is returned and has a valid ID, it means a duplicate name exists.
		if existingSite != nil && existingSite.ID > 0 {
			return fmt.Errorf("a site with the name '%s' already exists", proposedNameStr)
		}
	}

	return nil
}
