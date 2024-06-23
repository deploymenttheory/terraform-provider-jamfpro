// advancedmobiledevicesearches_object.go
package advancedmobiledevicesearches

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/constructobject"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAdvancedMobileDeviceSearch constructs a mobile device search object for create and update operations.
func constructJamfProAdvancedMobileDeviceSearch(d *schema.ResourceData) (*jamfpro.ResourceAdvancedMobileDeviceSearch, error) {
	search := &jamfpro.ResourceAdvancedMobileDeviceSearch{
		Name:   d.Get("name").(string),
		ViewAs: d.Get("view_as").(string),
		Sort1:  d.Get("sort1").(string),
		Sort2:  d.Get("sort2").(string),
		Sort3:  d.Get("sort3").(string),
	}

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

	if v, ok := d.GetOk("site_id"); ok {
		search.Site = constructobject.ConstructSharedResourceSite(v.([]interface{}))
	} else {

		search.Site = constructobject.ConstructSharedResourceSite([]interface{}{})
	}

	resourceXML, err := xml.MarshalIndent(search, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Advanced Mobile Device Search '%s' to XML: %v", search.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Advanced Mobile Device Search XML:\n%s\n", string(resourceXML))

	return search, nil
}
