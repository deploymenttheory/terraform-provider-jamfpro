package jamfprivileges

import (
	"log"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// FindSimilarPrivileges tries to suggest resource names similar to the supplied “invalid privilege”
// It uses fuzzy string matching across *all* validPrivileges from Jamf Pro and is used by
// Accout, AccountGroup and Api role resources to suggest similar privileges when an invalid one is detected.
func FindSimilarPrivileges(invalid string, validPrivileges []string) []string {
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
