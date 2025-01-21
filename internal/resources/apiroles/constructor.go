package apiroles

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct builds an ResourceAPIRole object from the provided schema data.
func construct(d *schema.ResourceData, meta interface{}) (*jamfpro.ResourceAPIRole, error) {
	client := meta.(*jamfpro.Client)

	resource := &jamfpro.ResourceAPIRole{
		DisplayName: d.Get("display_name").(string),
	}

	if v, ok := d.GetOk("privileges"); ok {
		privilegesInterface := v.(*schema.Set).List()
		privileges := make([]string, len(privilegesInterface))
		for i, priv := range privilegesInterface {
			var ok bool
			privileges[i], ok = priv.(string)
			if !ok {
				return nil, fmt.Errorf("failed to assert privilege to string")
			}
		}

		if err := validatePrivileges(client, privileges); err != nil {
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

// validatePrivileges validates that all provided privileges exist in Jamf Pro
func validatePrivileges(client *jamfpro.Client, privileges []string) error {
	versionInfo, err := client.GetJamfProVersion()
	if err != nil {
		return fmt.Errorf("failed to fetch Jamf Pro version: %v", err)
	}

	apiRoles, err := client.GetJamfAPIRoles("")
	if err != nil {
		return fmt.Errorf("failed to fetch API roles for privilege validation: %v", err)
	}

	validPrivileges := make(map[string]bool)
	var allPrivileges []string
	for _, role := range apiRoles.Results {
		for _, privilege := range role.Privileges {
			if !validPrivileges[privilege] {
				validPrivileges[privilege] = true
				allPrivileges = append(allPrivileges, privilege)
			}
		}
	}

	sort.Strings(allPrivileges)

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

		msg.WriteString("\nAvailable privileges that might be similar:\n")
		for _, invalid := range invalidPrivileges {
			similars := findSimilarPrivileges(invalid, allPrivileges)
			if len(similars) > 0 {
				msg.WriteString(fmt.Sprintf("Instead of '%s', did you mean:\n", invalid))
				for _, s := range similars {
					msg.WriteString(fmt.Sprintf("- %s\n", s))
				}
			}
		}

		return fmt.Errorf(msg.String())
	}

	return nil
}

// findSimilarPrivileges finds privileges that might be similar to the invalid one
func findSimilarPrivileges(invalid string, validPrivileges []string) []string {
	var similar []string

	// Split the invalid privilege into action and resource
	parts := strings.SplitN(invalid, " ", 2)
	if len(parts) != 2 {
		return similar
	}

	action := strings.ToLower(parts[0])   // e.g., "Create", "Read", "Delete"
	resource := strings.ToLower(parts[1]) // e.g., "Jamf Content Distribution Server Files"

	// Get the main resource terms by splitting and removing common words
	resourceTerms := strings.Fields(resource)
	var significantTerms []string
	commonWords := map[string]bool{
		"a": true, "an": true, "and": true, "the": true, "in": true, "on": true, "at": true,
		"to": true, "for": true, "of": true, "with": true, "by": true,
	}

	for _, term := range resourceTerms {
		term = strings.ToLower(term)
		if !commonWords[term] && len(term) > 2 {
			significantTerms = append(significantTerms, term)
		}
	}

	// Score each valid privilege based on matching criteria
	type scoredPrivilege struct {
		privilege string
		score     int
	}
	var scored []scoredPrivilege

	for _, valid := range validPrivileges {
		validLower := strings.ToLower(valid)
		score := 0

		// Match the action (Create/Read/Delete)
		if strings.HasPrefix(validLower, action) {
			score += 5
		}

		// Match significant terms from the resource
		for _, term := range significantTerms {
			if strings.Contains(validLower, term) {
				// Give higher score for matches in the same position
				if strings.Index(validLower, term) > strings.Index(strings.ToLower(invalid), term)-2 &&
					strings.Index(validLower, term) < strings.Index(strings.ToLower(invalid), term)+2 {
					score += 3
				} else {
					score += 1
				}
			}
		}

		// Boost score if the privilege deals with similar system components
		systemComponents := []string{"jamf", "server", "files", "distribution"}
		for _, comp := range systemComponents {
			if strings.Contains(validLower, comp) && strings.Contains(strings.ToLower(invalid), comp) {
				score += 2
			}
		}

		if score > 0 {
			scored = append(scored, scoredPrivilege{valid, score})
		}
	}

	// Sort by score
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Take top 3 matches with a minimum score
	for i := 0; i < len(scored) && len(similar) < 3; i++ {
		if scored[i].score >= 3 { // Only include if it's a reasonably good match
			similar = append(similar, scored[i].privilege)
		}
	}

	return similar
}
