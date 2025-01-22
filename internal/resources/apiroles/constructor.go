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

// scoredPrivilege is a helper struct for sorting privileges by similarity score
type scoredPrivilege struct {
	privilege string
	score     float64
}

// construct builds an ResourceAPIRole object from the provided schema data.
// It performs dynamic validation of the privileges against the Jamf Pro server.
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

		return fmt.Errorf("%s", msg.String())
	}

	return nil
}

// findSimilarPrivileges finds privileges similar to an invalid one by analyzing term matches
// and similarity scores. It performs the following operations:
//  1. Splits privileges into action (Create/Read/Delete) and resource parts
//  2. Matches only privileges with the same action type
//  3. Calculates similarity scores based on:
//     - Exact term matches in sequence
//     - Term matches in any position
//     - Plural/singular variations using simple 's' suffix matching
//     - Sequential term matches (2 or more terms in sequence)
//  4. Applies scoring bonuses for:
//     - Long sequences of matching terms (+0.1 per term)
//     - Exact matches after normalizing plurals (+0.2)
//  5. Returns only privileges with at least 70% similarity score
//
// The similarity score is calculated as (matched terms / total terms) plus any bonuses.
// Returns a slice of similar privileges sorted by similarity score.
func findSimilarPrivileges(invalid string, validPrivileges []string) []string {
	var similar []string

	parts := strings.SplitN(invalid, " ", 2)
	if len(parts) != 2 {
		return similar
	}

	action := strings.ToLower(parts[0])
	resource := strings.ToLower(parts[1])

	var scored []scoredPrivilege

	for _, valid := range validPrivileges {
		validLower := strings.ToLower(valid)
		validParts := strings.SplitN(validLower, " ", 2)
		if len(validParts) != 2 {
			continue
		}
		validAction := validParts[0]
		validResource := validParts[1]

		if validAction != action {
			continue
		}

		invalidTerms := strings.Fields(resource)
		validTerms := strings.Fields(validResource)

		matchedTerms := 0
		totalTerms := len(invalidTerms)

		longestSequence := 0
		currentSequence := 0
		for i := 0; i < len(invalidTerms) && i < len(validTerms); i++ {
			if invalidTerms[i] == validTerms[i] {
				currentSequence++
				matchedTerms++
			} else {
				invalidTerm := strings.TrimSuffix(invalidTerms[i], "s")
				validTerm := strings.TrimSuffix(validTerms[i], "s")
				if invalidTerm == validTerm {
					currentSequence++
					matchedTerms++
				} else {
					if currentSequence > longestSequence {
						longestSequence = currentSequence
					}
					currentSequence = 0
				}
			}
		}
		if currentSequence > longestSequence {
			longestSequence = currentSequence
		}

		for _, invalidTerm := range invalidTerms {
			for _, validTerm := range validTerms {
				if invalidTerm == validTerm {
					matchedTerms++
					break
				} else {
					invalidSingular := strings.TrimSuffix(invalidTerm, "s")
					validSingular := strings.TrimSuffix(validTerm, "s")
					if invalidSingular == validSingular {
						matchedTerms++
						break
					}
				}
			}
		}

		similarityScore := float64(matchedTerms) / float64(totalTerms)

		if longestSequence >= 2 {
			similarityScore += float64(longestSequence) * 0.1
		}

		if strings.TrimSuffix(resource, "s") == strings.TrimSuffix(validResource, "s") {
			similarityScore += 0.2
		}

		if similarityScore >= 0.7 {
			scored = append(scored, scoredPrivilege{valid, similarityScore})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	for _, s := range scored {
		similar = append(similar, s.privilege)
	}

	return similar
}
