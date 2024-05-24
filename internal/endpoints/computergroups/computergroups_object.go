// computergroup_object.go
package computergroups

import (
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/sharedschemas/sharedschema_constructors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProComputerGroup constructs a ResourceComputerGroup object from the provided schema data.
func constructJamfProComputerGroup(d *schema.ResourceData) (*jamfpro.ResourceComputerGroup, error) {
	log.Println("LOGHERE")
	out := &jamfpro.ResourceComputerGroup{
		Name:      d.Get("name").(string),
		IsSmart:   d.Get("is_smart").(bool),
		Criteria:  &jamfpro.ComputerGroupSubsetContainerCriteria{},
		Computers: &jamfpro.ComputerGroupSubsetComputersContainer{},
	}
	log.Println("FLAG-1")

	site, err := sharedschema_constructors.GetSite(d)
	if err != nil {
		return nil, err
	}
	out.Site = site

	log.Println("FLAG-2")

	constructGroupCriteria(d, out.Criteria)
	log.Printf("%+v", out.Criteria.Criterion)

	log.Println("FLAG-3")

	return out, nil

}

// Helper functions for nested structures

// constructGroupCriteria constructs a slice of SharedSubsetCriteria from the provided schema data.
func constructGroupCriteria(d *schema.ResourceData, home *jamfpro.ComputerGroupSubsetContainerCriteria) {
	log.Println("FLAG-2.1")
	criteria := d.Get("criteria")
	log.Println("FLAG-2.2")
	if criteria == nil {
		log.Println("FLAG-2.3")
		home = &jamfpro.ComputerGroupSubsetContainerCriteria{
			Size: 0,
		}
		return
	}
	log.Println("FLAG-2.4")
	home.Criterion = &[]jamfpro.SharedSubsetCriteria{}
	for _, crit := range criteria.([]interface{}) {
		log.Println("FLAG-2.5")
		*home.Criterion = append(*home.Criterion, jamfpro.SharedSubsetCriteria{
			Name:         crit.(map[string]interface{})["name"].(string),
			Priority:     crit.(map[string]interface{})["priority"].(int),
			AndOr:        crit.(map[string]interface{})["and_or"].(string),
			SearchType:   crit.(map[string]interface{})["search_type"].(string),
			Value:        crit.(map[string]interface{})["value"].(string),
			OpeningParen: crit.(map[string]interface{})["opening_paren"].(bool),
			ClosingParen: crit.(map[string]interface{})["closing_paren"].(bool),
		})
		log.Println("FLAG-2.6")
	}
}

// constructGroupComputers constructs a slice of ComputerGroupSubsetComputer from the provided schema data.
func constructGroupComputers(computersData []interface{}) *[]jamfpro.ComputerGroupSubsetComputer {
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

	return &computers
}
