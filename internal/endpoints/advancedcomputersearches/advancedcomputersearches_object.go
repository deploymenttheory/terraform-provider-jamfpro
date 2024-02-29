// advancedcomputersearches_resource.go
package advancedcomputersearches

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAdvancedComputerSearch constructs an advanced computer search object for create and update operations.
func constructJamfProAdvancedComputerSearch(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceAdvancedComputerSearch, error) {
	search := &jamfpro.ResourceAdvancedComputerSearch{
		Name:   d.Get("name").(string),
		ViewAs: d.Get("view_as").(string),
		Sort1:  d.Get("sort1").(string),
		Sort2:  d.Get("sort2").(string),
		Sort3:  d.Get("sort3").(string),
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
		return nil, fmt.Errorf("failed to marshal Jamf Pro Advanced Computer Search '%s' to XML: %v", search.Name, err)
	}
	fmt.Printf("Constructed Advanced Computer Search XML:\n%s\n", string(xmlData))

	return search, nil
}
