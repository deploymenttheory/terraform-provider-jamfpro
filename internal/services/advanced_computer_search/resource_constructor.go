// advancedcomputersearches_resource.go
package advanced_computer_search

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sharedschemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAdvancedComputerSearch constructs an advanced computer search object for create and update operations.
func construct(d *schema.ResourceData) (*jamfpro.ResourceAdvancedComputerSearch, error) {
	resource := &jamfpro.ResourceAdvancedComputerSearch{
		Name:   d.Get("name").(string),
		ViewAs: d.Get("view_as").(string),
		Sort1:  d.Get("sort1").(string),
		Sort2:  d.Get("sort2").(string),
		Sort3:  d.Get("sort3").(string),
	}

	if v, ok := d.GetOk("criteria"); ok {
		criteriaList := v.([]any)
		criteria := make([]jamfpro.SharedSubsetCriteria, len(criteriaList))
		for i, crit := range criteriaList {
			criterionMap := crit.(map[string]any)
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
		resource.Criteria.Criterion = criteria
	}

	displayFieldsHcl := d.Get("display_fields").([]any)
	if len(displayFieldsHcl) > 0 {
		for _, v := range displayFieldsHcl {
			resource.DisplayFields = append(resource.DisplayFields, jamfpro.DisplayField{Name: v.(string)})
		}
	}

	resource.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Advanced Computer Search '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Advanced Computer Search XML:\n%s\n", string(resourceXML))

	return resource, nil
}
