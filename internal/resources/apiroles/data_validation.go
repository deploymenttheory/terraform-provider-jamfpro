package apiroles

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

//go:embed api_privileges.json
var privilegesFS embed.FS

var validPrivileges []string

func init() {
	// Read the embedded JSON file
	data, err := privilegesFS.ReadFile("api_privileges.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to read api_privileges.json: %v", err))
	}

	// Unmarshal the JSON data
	var privilegesData struct {
		Privileges []string `json:"privileges"`
	}
	if err := json.Unmarshal(data, &privilegesData); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal privileges data: %v", err))
	}

	// Assign the privileges to validPrivileges
	validPrivileges = privilegesData.Privileges
}

// validateResourceApiRolesDataFields checks if a given privilege is in the list of valid privileges
// and groups privileges by category.
func validateResourceApiRolesDataFields(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	categories := make(map[string][]string)
	var nonCrudPrivileges []string
	var MDMCommands []string

	for _, priv := range validPrivileges {
		if v == priv {
			return
		}

		// Split the privilege into operation and category
		parts := strings.SplitN(priv, " ", 2)
		if len(parts) == 2 {
			operation, category := parts[0], parts[1]

			// Group CRUD privileges by category
			if operation == "Create" || operation == "Read" || operation == "Update" || operation == "Delete" {
				categories[category] = append(categories[category], priv)
			} else if operation == "Send" {
				MDMCommands = append(MDMCommands, priv)
			} else {
				nonCrudPrivileges = append(nonCrudPrivileges, priv)
			}
		} else {
			nonCrudPrivileges = append(nonCrudPrivileges, priv)
		}
	}

	var formattedPrivileges strings.Builder

	// Sort categories for consistent ordering
	sortedCategories := make([]string, 0, len(categories))
	for category := range categories {
		sortedCategories = append(sortedCategories, category)
	}
	sort.Strings(sortedCategories)

	for _, category := range sortedCategories {
		privileges := categories[category]
		// Adding a spacer with the category name
		formattedPrivileges.WriteString(fmt.Sprintf("---- Privilege Set: %s ----\n", category))
		formattedPrivileges.WriteString(fmt.Sprintf("    %s\n", strings.Join(privileges, "\n    ")))
		formattedPrivileges.WriteString("---- End ----\n\n")
	}

	if len(MDMCommands) > 0 {
		// Adding a spacer for Send MDM Commands
		formattedPrivileges.WriteString("---- MDM Commands ----\n")
		formattedPrivileges.WriteString(fmt.Sprintf("    %s\n", strings.Join(MDMCommands, "\n    ")))
		formattedPrivileges.WriteString("---- End ----\n\n")
	}

	if len(nonCrudPrivileges) > 0 {
		// Adding a spacer for non-CRUD privileges
		formattedPrivileges.WriteString("---- Other Jamf Pro Operations ----\n")
		formattedPrivileges.WriteString(fmt.Sprintf("    %s\n", strings.Join(nonCrudPrivileges, "\n    ")))
		formattedPrivileges.WriteString("---- End ----\n\n")
	}

	errs = append(errs, fmt.Errorf("%q contains an invalid privilege: %s; must be one of:\n%s", key, v, formattedPrivileges.String()))
	return
}
