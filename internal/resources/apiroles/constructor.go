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

	validPrivileges := make(map[string]bool)
	allPrivileges := privilegesList.Privileges
	for _, privilege := range allPrivileges {
		validPrivileges[privilege] = true
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
//  1. Splits privileges into action (e.g Create/Read/Delete) and resource parts
//  2. Matches only privileges with the same action type
//  3. Calculates similarity scores based on sequence matches and term matches
//  4. Returns only privileges with at least 70% similarity score
func findSimilarPrivileges(invalid string, validPrivileges []string) []string {
	similar := make([]string, 0, len(validPrivileges))
	scored := make([]scoredPrivilege, 0, len(validPrivileges))

	action, resource, ok := splitPrivilegeIntoActionAndResource(invalid)
	if !ok {
		return similar
	}

	invalidTerms := strings.Fields(resource)
	totalTerms := len(invalidTerms)

	for _, valid := range validPrivileges {
		validAction, validResource, ok := splitPrivilegeIntoActionAndResource(valid)
		if !ok || validAction != action {
			continue
		}

		validTerms := strings.Fields(validResource)
		seqMatches, longestSequence := countConsecutiveTermMatches(invalidTerms, validTerms)
		termMatches := calculateTermMatches(invalidTerms, validTerms)

		score := float64(seqMatches+termMatches) / float64(totalTerms)
		if longestSequence >= 2 {
			score += float64(longestSequence) * 0.1
		}

		if score >= 0.7 {
			scored = append(scored, scoredPrivilege{valid, score})
		}
	}

	if len(scored) > 0 {
		sort.Slice(scored, func(i, j int) bool {
			return scored[i].score > scored[j].score
		})

		similar = make([]string, len(scored))
		for i, s := range scored {
			similar[i] = s.privilege
		}
	}

	return similar
}

// splitPrivilegeIntoActionAndResource splits a privilege into action and resource components
func splitPrivilegeIntoActionAndResource(privilege string) (action, resource string, ok bool) {
	parts := strings.SplitN(privilege, " ", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return strings.ToLower(parts[0]), strings.ToLower(parts[1]), true
}

// countConsecutiveTermMatches evaluates matching terms in sequence
func countConsecutiveTermMatches(invalidTerms, validTerms []string) (matchedTerms, longestSequence int) {
	currentSequence := 0
	for i := 0; i < len(invalidTerms) && i < len(validTerms); i++ {
		if termSimilarity(invalidTerms[i], validTerms[i]) >= 0.85 {
			currentSequence++
			matchedTerms++
			continue
		}

		if currentSequence > longestSequence {
			longestSequence = currentSequence
		}
		currentSequence = 0
	}

	if currentSequence > longestSequence {
		longestSequence = currentSequence
	}

	return matchedTerms, longestSequence
}

// calculateTermMatches counts matching terms in any position with optimized matching
// Sorts terms by length for better matching (try shorter terms first)
func calculateTermMatches(invalidTerms, validTerms []string) int {
	if len(invalidTerms) == 0 || len(validTerms) == 0 {
		return 0
	}

	matchedValid := make(map[int]bool, len(validTerms))
	matchedTerms := 0

	sort.Slice(invalidTerms, func(i, j int) bool {
		return len(invalidTerms[i]) < len(invalidTerms[j])
	})

	for _, invalidTerm := range invalidTerms {
		for j, validTerm := range validTerms {
			if !matchedValid[j] && termSimilarity(invalidTerm, validTerm) >= 0.85 {
				matchedValid[j] = true
				matchedTerms++
				break
			}
		}
	}

	return matchedTerms
}

// termSimilarity calculates a normalized similarity score between two strings using
// their Longest Common Subsequence (LCS). It returns a float64 between 0 and 1, where:
//   - 1.0 indicates identical strings
//   - 0.0 indicates completely different strings
//   - Values between 0 and 1 indicate the degree of similarity
//
// The score is calculated by finding the length of the LCS and dividing it by the
// length of the longer string. For example, comparing "keystore" and "keystores":
//   - LCS length = 8 ("keystore")
//   - Max length = 9 ("keystores")
//   - Score = 8/9 ≈ 0.89
//
// The comparison is case-insensitive and returns 0 if either string is empty.
func termSimilarity(term1, term2 string) float64 {
	if term1 == term2 {
		return 1.0
	}

	s1 := strings.ToLower(term1)
	s2 := strings.ToLower(term2)

	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}

	m := len(s1)
	n := len(s2)
	curr := make([]int, n+1)
	prev := make([]int, n+1)

	for i := 1; i <= m; i++ {
		prev, curr = curr, prev

		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				curr[j] = prev[j-1] + 1
			} else {
				curr[j] = max(curr[j-1], prev[j])
			}
		}
	}

	matchLength := float64(curr[n])
	maxLength := float64(max(len(s1), len(s2)))
	return matchLength / maxLength
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
