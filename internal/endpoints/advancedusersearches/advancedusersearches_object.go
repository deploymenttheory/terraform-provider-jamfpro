// advancedusersearches_object.go
package advancedusersearches

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAdvancedUserSearch constructs an advanced user search object for create and update operations.
func constructJamfProAdvancedUserSearch(d *schema.ResourceData) (*jamfpro.ResourceAdvancedUserSearch, error) {
	search := &jamfpro.ResourceAdvancedUserSearch{
		Name: d.Get("name").(string),
	}

	// Handle 'criteria' field
	if v, ok := d.GetOk("criteria"); ok {
		criteriaList := v.([]interface{})
		criteria := make([]jamfpro.SharedSubsetCriteria, len(criteriaList))
		for i, crit := range criteriaList {
			criterionMap := crit.(map[string]interface{})
			criteria[i] = jamfpro.SharedSubsetCriteria{
				Name:         criterionMap["name"].(string),
				Priority:     criterionMap["priority"].(int),
				AndOr:        criterionMap["and_or"].(string),
				SearchType:   criterionMap["search_type"].(string),
				Value:        criterionMap["value"].(string),
				OpeningParen: criterionMap["opening_paren"].(bool),
				ClosingParen: criterionMap["closing_paren"].(bool),
			}
		}
		search.Criteria.Criterion = criteria
	}

	// Handle 'display_fields' field
	if v, ok := d.GetOk("display_fields"); ok {
		displayFieldsSet := v.(*schema.Set).List()
		displayFields := make([]jamfpro.SharedAdvancedSearchSubsetDisplayField, len(displayFieldsSet))
		for i, field := range displayFieldsSet {
			fieldMap := field.(map[string]interface{})
			displayFields[i] = jamfpro.SharedAdvancedSearchSubsetDisplayField{
				Name: fieldMap["name"].(string),
			}
		}
		search.DisplayFields = []jamfpro.SharedAdvancedSearchContainerDisplayField{{DisplayField: displayFields}}
	}

	// Handle Site
	if v, ok := d.GetOk("site"); ok {
		search.Site = constructSharedResourceSite(v.([]interface{}))
	} else {
		// Set default values if 'site' data is not provided
		search.Site = constructSharedResourceSite([]interface{}{})
	}

	// Serialize and pretty-print the Advanced User Search object as XML for logging
	resourceXML, err := xml.MarshalIndent(search, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Advanced User Search '%s' to XML: %v", search.Name, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Advanced User Search XML:\n%s\n", string(resourceXML))

	return search, nil
}

// Helper functions for nested structures

// constructSharedResourceSite constructs a SharedResourceSite object from the provided schema data,
// setting default values if none are presented.
func constructSharedResourceSite(data []interface{}) jamfpro.SharedResourceSite {
	// Check if 'site' data is provided and non-empty
	if len(data) > 0 && data[0] != nil {
		site := data[0].(map[string]interface{})

		// Return the 'site' object with data from the schema
		return jamfpro.SharedResourceSite{
			ID:   site["id"].(int),
			Name: site["name"].(string),
		}
	}

	// Return default 'site' values if no data is provided or it is empty
	return jamfpro.SharedResourceSite{
		ID:   -1,     // Default ID
		Name: "None", // Default name
	}
}
