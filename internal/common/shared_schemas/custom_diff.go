package sharedschemas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	errJamfClientNotConfigured           = errors.New("jamf client is not configured for LDAP validation")
	errDirectoryServiceUserGroupNotFound = errors.New("directory service user group name was not found in Jamf Pro")
)

// ValidateScopeDirectoryServiceUserGroupNames validates user group names when the scope
// resource is nested at the root level (e.g. scope.0.limitations...).
func ValidateScopeDirectoryServiceUserGroupNames(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return fmt.Errorf("%w", errJamfClientNotConfigured)
	}

	type entry struct {
		description string
		paths       []string
	}

	entries := []entry{
		{
			description: "scope limitations directory service user group names",
			paths: []string{
				"scope.0.limitations.0.directory_service_usergroup_names",
			},
		},
		{
			description: "scope exclusions directory service user group names",
			paths: []string{
				"scope.0.exclusions.0.directory_service_usergroup_names",
			},
		},
	}

	checked := make(map[string]bool)

	for _, entry := range entries {
		var (
			value any
			ok    bool
		)
		for _, path := range entry.paths {
			value, ok = d.GetOk(path)
			if ok {
				break
			}
		}
		if !ok {
			continue
		}

		set, ok := value.(*schema.Set)
		if !ok || set.Len() == 0 {
			continue
		}

		for _, raw := range set.List() {
			trimmed := strings.TrimSpace(raw.(string))
			if trimmed == "" {
				continue
			}

			exists, cached := checked[trimmed]
			if !cached {
				resp, err := client.GetLdapGroupsV1(trimmed)
				if err != nil {
					return fmt.Errorf("failed to validate directory service user group %q: %w", trimmed, err)
				}

				exists = false
				for _, group := range resp.Results {
					if group.Name == trimmed {
						exists = true
						break
					}
				}
				checked[trimmed] = exists
			}

			if !exists {
				return fmt.Errorf("%w: %q defined in %s", errDirectoryServiceUserGroupNotFound, trimmed, entry.description)
			}
		}
	}

	return nil
}
