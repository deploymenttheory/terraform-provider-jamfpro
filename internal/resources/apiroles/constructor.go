package apiroles

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/jamfprivileges"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct builds an ResourceAPIRole object from the provided schema data
func construct(d *schema.ResourceData, meta interface{}) (*jamfpro.ResourceAPIRole, error) {
	client := meta.(*jamfpro.Client)

	resource := &jamfpro.ResourceAPIRole{
		DisplayName: d.Get("display_name").(string),
	}

	if v, ok := d.GetOk("privileges"); ok {
		privilegesInterface := v.(*schema.Set).List()
		privileges := make([]string, len(privilegesInterface))
		for i, priv := range privilegesInterface {
			if privileges[i], ok = priv.(string); !ok {

				return nil, fmt.Errorf("failed to assert api role privilege to string")
			}
		}

		if err := validateApiRolePrivileges(client, privileges); err != nil {
			return nil, err
		}

		resource.Privileges = privileges
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {

		return nil, fmt.Errorf("failed to marshal Jamf Pro Api Role '%s' to JSON: %v", resource.DisplayName, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Api Role JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// validateApiRolePrivileges validates that all provided privileges exist in Jamf Pro
func validateApiRolePrivileges(client *jamfpro.Client, privileges []string) error {
	versionInfo, err := client.GetJamfProVersion()
	if err != nil {

		return fmt.Errorf("failed to fetch Jamf Pro version: %v", err)
	}

	privilegesList, err := client.GetJamfAPIPrivileges()
	if err != nil {

		return fmt.Errorf("failed to fetch API privileges for validation: %v", err)
	}

	// Create a map of valid privileges for O(1) lookup
	validPrivileges := make(map[string]bool)
	for _, privilege := range privilegesList.Privileges {
		validPrivileges[privilege] = true
	}

	// Collect invalid privileges
	var invalidPrivileges []string
	for _, privilege := range privileges {
		if !validPrivileges[privilege] {
			invalidPrivileges = append(invalidPrivileges, privilege)
		}
	}

	if len(invalidPrivileges) > 0 {
		var msg strings.Builder
		msg.WriteString(fmt.Sprintf("Invalid privileges have been defined when compared to your Jamf Pro version %s:\n", *versionInfo.Version))
		msg.WriteString("\nInvalid privileges:\n")
		for _, p := range invalidPrivileges {
			msg.WriteString(fmt.Sprintf("- %s\n", p))
		}

		msg.WriteString("\nSuggested similar privileges:\n")
		for _, invalid := range invalidPrivileges {
			similars := jamfprivileges.FindSimilarPrivileges(invalid, privilegesList.Privileges)
			if len(similars) > 0 {
				msg.WriteString(fmt.Sprintf("Instead of '%s', did you mean:\n", invalid))
				for _, s := range similars {
					msg.WriteString(fmt.Sprintf("- %s\n", s))
				}
			}
		}

		return fmt.Errorf("%s", msg.String())
	}

	return nil
}
