// advancedusersearches_object.go
package advancedusersearches

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAdvancedUserSearch constructs an advanced user search object for create and update operations.
func construct(d *schema.ResourceData) (*jamfpro.ResourceAdvancedUserSearch, error) {
	resource := &jamfpro.ResourceAdvancedUserSearch{
		Name: d.Get("name").(string),
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
		resource.Criteria.Criterion = criteria
	}

	displayFieldsHcl := d.Get("display_fields").([]interface{})
	if len(displayFieldsHcl) > 0 {
		for _, v := range displayFieldsHcl {
			resource.DisplayFields = append(resource.DisplayFields, jamfpro.DisplayField{Name: v.(string)})
		}
	}

	resource.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return nil, fmt.Errorf("failed to marshal Jamf Pro Advanced User Search '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Advanced User Search XML:\n%s\n", string(resourceXML))

	return resource, nil
}
