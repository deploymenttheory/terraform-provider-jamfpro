package smartmobiledevicegroups

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProSmartMobileGroup constructs a ResourceMobileDeviceGroup object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceMobileDeviceGroup, error) {
	resource := &jamfpro.ResourceMobileDeviceGroup{
		Name:    d.Get("name").(string),
		IsSmart: true,
	}
	if (d.Get("site_id").(int)) == 0 || (d.Get("site_id").(int)) == -1 {
		resource.Site = jamfpro.SharedResourceSite{ID: -1, Name: ""}

	} else {
		resource.Site = jamfpro.SharedResourceSite{ID: (d.Get("site_id").(int))}
	}

	if v, ok := d.GetOk("criteria"); ok {
		resource.Criteria = constructMobileGroupSubsetContainerCriteria(v.([]interface{}))
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {

		return nil, fmt.Errorf("failed to marshal Jamf Pro Mobile Group '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Mobile Device Group XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// constructMobileGroupSubsetContainerCriteria constructs a constructMobileGroupSubsetContainerCriteria object from the provided schema data.
func constructMobileGroupSubsetContainerCriteria(criteriaList []interface{}) jamfpro.SharedContainerCriteria {
	criteria := jamfpro.SharedContainerCriteria{
		Size:      len(criteriaList),
		Criterion: []jamfpro.SharedSubsetCriteria{},
	}

	for _, item := range criteriaList {
		criterionData := item.(map[string]interface{})
		criterion := jamfpro.SharedSubsetCriteria{
			Name:         criterionData["name"].(string),
			Priority:     criterionData["priority"].(int),
			AndOr:        criterionData["and_or"].(string),
			SearchType:   criterionData["search_type"].(string),
			Value:        criterionData["value"].(string),
			OpeningParen: criterionData["opening_paren"].(bool),
			ClosingParen: criterionData["closing_paren"].(bool),
		}
		criteria.Criterion = append(criteria.Criterion, criterion)
	}

	return criteria
}
