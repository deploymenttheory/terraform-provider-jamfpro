// advancedusersearches_object.go
package advancedusersearches

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAdvancedUserSearch constructs an advanced user search object for create and update operations.
func constructJamfProAdvancedUserSearch(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceAdvancedUserSearch, error) {
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

	// Handle 'site' field
	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		siteData := v.([]interface{})[0].(map[string]interface{})
		search.Site = jamfpro.SharedResourceSite{
			ID:   siteData["id"].(int),
			Name: siteData["name"].(string),
		}
	}

	// Serialize and pretty-print the search object as XML for logging
	xmlData, err := xml.MarshalIndent(search, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Advanced User Search '%s' to XML: %v", search.Name, err)
	}
	fmt.Printf("[DEBUG] Constructed Advanced User Search Object:\n%s\n", string(xmlData))

	return search, nil
}
