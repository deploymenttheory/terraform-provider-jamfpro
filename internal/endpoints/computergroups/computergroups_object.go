// computergroup_object.go
package computergroups

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProComputerGroup constructs a ResourceComputerGroup object from the provided schema data.
func constructJamfProComputerGroup(d *schema.ResourceData) (*jamfpro.ResourceComputerGroup, error) {
	group := &jamfpro.ResourceComputerGroup{
		Name:    d.Get("name").(string),
		IsSmart: d.Get("is_smart").(bool),
	}

	// Handle nested "site" field
	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		siteData := v.([]interface{})[0].(map[string]interface{})
		group.Site = jamfpro.SharedResourceSite{
			ID:   siteData["id"].(int),
			Name: siteData["name"].(string),
		}
	}

	// Handle "criteria" field
	if v, ok := d.GetOk("criteria"); ok {
		criteria := constructGroupCriteria(v.([]interface{}))
		group.Criteria = jamfpro.SharedContainerCriteria{
			Criterion: criteria,
		}
	}

	// Handle "computers" field
	if v, ok := d.GetOk("computers"); ok && !group.IsSmart {
		group.Computers = constructGroupComputers(v.([]interface{}))
	}

	// Serialize and pretty-print the Computer Group object as XML for logging
	resourceXML, err := xml.MarshalIndent(group, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Computer Group '%s' to XML: %v", group.Name, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Computer Group XML:\n%s\n", string(resourceXML))

	return group, nil
}

// Helper function to construct group criteria from schema data
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

// Helper function to construct group computers from schema data
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
