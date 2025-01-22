package apiroles

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lithammer/fuzzysearch/fuzzy"
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
			similars := findSimilarPrivileges(invalid, privilegesList.Privileges)
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

// findSimilarPrivileges tries to suggest resource names similar to “invalidPrivilege”
// using fuzzy matching across *all* validPrivileges from Jamf Pro.
func findSimilarPrivileges(invalid string, validPrivileges []string) []string {
	// 1) Parse the invalid string into [action, resource].
	parts := strings.SplitN(invalid, " ", 2)
	var invalidAction, invalidResource string

	// If we detect exactly two parts (e.g. "Create Something"),
	// we treat the first word as the action.
	if len(parts) == 2 {
		invalidAction = strings.ToLower(parts[0])
		invalidResource = parts[1]
	} else {
		// If there's only one token, we can't reliably parse out the verb.
		// So either treat the entire string as the resource,
		// or skip action-based filtering altogether.
		invalidResource = invalid
	}

	type candidate struct {
		priv string
		dist int
	}
	candidates := make([]candidate, 0, len(validPrivileges))

	for _, vp := range validPrivileges {
		// Split each valid privilege into [action, resource].
		vparts := strings.SplitN(vp, " ", 2)
		if len(vparts) == 2 {
			vpAction := strings.ToLower(vparts[0])
			vpResource := vparts[1]

			// **Only** consider suggestions whose action matches the user’s typed action.
			if vpAction == invalidAction {
				dist := fuzzy.LevenshteinDistance(
					strings.ToLower(invalidResource),
					strings.ToLower(vpResource),
				)
				candidates = append(candidates, candidate{priv: vp, dist: dist})
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].dist < candidates[j].dist
	})

	// Keep top suggestion
	maxSuggestions := 1
	if len(candidates) > maxSuggestions {
		candidates = candidates[:maxSuggestions]
	}

	suggestions := make([]string, 0, len(candidates))
	for _, c := range candidates {
		suggestions = append(suggestions, c.priv)
	}

	log.Printf("[DEBUG] findSimilarPrivileges: invalid='%s', action='%s', resource='%s', suggestions=%v",
		invalid, invalidAction, invalidResource, suggestions)

	return suggestions
}
