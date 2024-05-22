// computergroup_object.go
package computergroups

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/constructobject"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProComputerGroup constructs a ResourceComputerGroup object from the provided schema data.
func constructJamfProComputerGroup(d *schema.ResourceData) (*jamfpro.ResourceComputerGroup, error) {
	group := &jamfpro.ResourceComputerGroup{
		Name:    d.Get("name").(string),
		IsSmart: d.Get("is_smart").(bool),
	}

	// Handle Site
	if v, ok := d.GetOk("site"); ok {
		group.Site = constructobject.ConstructSharedResourceSite(v.([]interface{}))
	} else {
		// Set default values if 'site' data is not provided
		group.Site = constructobject.ConstructSharedResourceSite([]interface{}{})
	}

	// Handle "criteria" field only if it has entries
	if v, ok := d.GetOk("criteria"); ok && len(v.([]interface{})) > 0 {
		criteria := constructGroupCriteria(v.([]interface{}))
		group.Criteria = jamfpro.SharedContainerCriteria{
			Criterion: criteria,
		}
	}

	// Handle "computers" field

	if !group.IsSmart {
		computers, ok := d.GetOk("computers")

		if len(computers.([]interface{})) > 0 && ok {
			group.Computers = constructGroupComputers(computers.([]interface{}))

		} else if !ok {
			return nil, fmt.Errorf("failed to get computers")
	} else {
		group.Computers = nil
	}


	log.Printf("%+v", group)

	// Serialize and pretty-print the Computer Group object as XML for logging
	resourceXML, err := xml.MarshalIndent(group, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Group '%s' to XML: %v", group.Name, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Computer Group XML:\n%s\n", string(resourceXML))

	return group, nil
}

// Helper functions for nested structures

// constructGroupCriteria constructs a slice of SharedSubsetCriteria from the provided schema data.
func constructGroupCriteria(criteriaData []interface{}) []jamfpro.SharedSubsetCriteria {
	var criteria []jamfpro.SharedSubsetCriteria
	for _, crit := range criteriaData {
		criterionMap := crit.(map[string]interface{})
		criteria = append(criteria, jamfpro.SharedSubsetCriteria{
			Name:         criterionMap["name"].(string),
			Priority:     criterionMap["priority"].(int),
			AndOr:        criterionMap["and_or"].(string),
			SearchType:   criterionMap["search_type"].(string),
			Value:        criterionMap["value"].(string),
			OpeningParen: criterionMap["opening_paren"].(bool),
			ClosingParen: criterionMap["closing_paren"].(bool),
		})
	}

	return criteria
}

// constructGroupComputers constructs a slice of ComputerGroupSubsetComputer from the provided schema data.
func constructGroupComputers(computersData []interface{}) []jamfpro.ComputerGroupSubsetComputer {
	var computers []jamfpro.ComputerGroupSubsetComputer
	for _, comp := range computersData {
		computerMap := comp.(map[string]interface{})
		computers = append(computers, jamfpro.ComputerGroupSubsetComputer{
			ID:            computerMap["id"].(int),
			Name:          computerMap["name"].(string),
			SerialNumber:  computerMap["serial_number"].(string),
			MacAddress:    computerMap["mac_address"].(string),
			AltMacAddress: computerMap["alt_mac_address"].(string),
		})
	}

	return computers
}
