// advancedmobiledevicesearches_object.go
package advancedmobiledevicesearches

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProAdvancedMobileDeviceSearch constructs a mobile device search object for create and update operations.
func construct(d *schema.ResourceData) (*jamfpro.ResourceAdvancedMobileDeviceSearch, error) {
	siteId := d.Get("site_id").(string)

	resource := &jamfpro.ResourceAdvancedMobileDeviceSearch{
		Name:   d.Get("name").(string),
		SiteId: &siteId,
	}

	if v, ok := d.GetOk("criteria"); ok {
		criteriaList := v.([]interface{})
		criteria := make([]jamfpro.SharedSubsetCriteriaJamfProAPI, len(criteriaList))
		for i, crit := range criteriaList {
			criterionMap := crit.(map[string]interface{})

			criteria[i] = jamfpro.SharedSubsetCriteriaJamfProAPI{
				Name:         criterionMap["name"].(string),
				Priority:     criterionMap["priority"].(int),
				AndOr:        criterionMap["and_or"].(string),
				SearchType:   criterionMap["search_type"].(string),
				Value:        criterionMap["value"].(string),
				OpeningParen: jamfpro.BoolPtr(criterionMap["opening_paren"].(bool)),
				ClosingParen: jamfpro.BoolPtr(criterionMap["closing_paren"].(bool)),
			}
		}
		resource.Criteria = criteria
	}

	if v, ok := d.GetOk("display_fields"); ok {
		displayFieldsSet := v.(*schema.Set)
		displayFields := make([]string, displayFieldsSet.Len())

		for i, field := range displayFieldsSet.List() {
			displayFields[i] = field.(string)
		}
		resource.DisplayFields = displayFields
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return nil, fmt.Errorf("failed to marshal Jamf Pro Advanced Mobile Searches '%s' to JSON: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Advanced Mobile Searches JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
