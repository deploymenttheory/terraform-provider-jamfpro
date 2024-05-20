// restrictedsoftware_object.go
package restrictedsoftware

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/constructobject"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProRestrictedSoftware constructs a RestrictedSoftware object from the provided schema data.
func constructJamfProRestrictedSoftware(d *schema.ResourceData) (*jamfpro.ResourceRestrictedSoftware, error) {
	restrictedSoftware := &jamfpro.ResourceRestrictedSoftware{
		General: jamfpro.RestrictedSoftwareSubsetGeneral{
			Name:                  d.Get("name").(string),
			ProcessName:           d.Get("process_name").(string),
			MatchExactProcessName: d.Get("match_exact_process_name").(bool),
			SendNotification:      d.Get("send_notification").(bool),
			KillProcess:           d.Get("kill_process").(bool),
			DeleteExecutable:      d.Get("delete_executable").(bool),
			DisplayMessage:        d.Get("display_message").(string),
		},
	}

	// Handle Site
	if v, ok := d.GetOk("site"); ok {
		restrictedSoftware.General.Site = constructobject.ConstructSharedResourceSite(v.([]interface{}))
	} else {
		// Set default values if 'site' data is not provided
		restrictedSoftware.General.Site = constructobject.ConstructSharedResourceSite([]interface{}{})
	}

	// Handle Scope
	if v, ok := d.GetOk("scope"); ok {
		scopeData := v.([]interface{})[0].(map[string]interface{})
		scope := jamfpro.RestrictedSoftwareSubsetScope{
			AllComputers: scopeData["all_computers"].(bool),
		}

		if computerIDs, ok := scopeData["computer_ids"]; ok {
			scope.Computers = constructScopeEntitiesFromIds(computerIDs.([]interface{}))
		}
		if computerGroupIDs, ok := scopeData["computer_group_ids"]; ok {
			scope.ComputerGroups = constructScopeEntitiesFromIds(computerGroupIDs.([]interface{}))
		}
		if buildingIDs, ok := scopeData["building_ids"]; ok {
			scope.Buildings = constructScopeEntitiesFromIds(buildingIDs.([]interface{}))
		}
		if departmentIDs, ok := scopeData["department_ids"]; ok {
			scope.Departments = constructScopeEntitiesFromIds(departmentIDs.([]interface{}))
		}

		// Handle Exclusions
		if exclusions, ok := scopeData["exclusions"]; ok && len(exclusions.([]interface{})) > 0 {
			exclusionData := exclusions.([]interface{})[0].(map[string]interface{})
			scope.Exclusions = jamfpro.RestrictedSoftwareSubsetScopeExclusions{}

			if computerIDs, ok := exclusionData["computer_ids"]; ok {
				scope.Exclusions.Computers = constructScopeEntitiesFromIds(computerIDs.([]interface{}))
			}
			if computerGroupIDs, ok := exclusionData["computer_group_ids"]; ok {
				scope.Exclusions.ComputerGroups = constructScopeEntitiesFromIds(computerGroupIDs.([]interface{}))
			}
			if buildingIDs, ok := exclusionData["building_ids"]; ok {
				scope.Exclusions.Buildings = constructScopeEntitiesFromIds(buildingIDs.([]interface{}))
			}
			if departmentIDs, ok := exclusionData["department_ids"]; ok {
				scope.Exclusions.Departments = constructScopeEntitiesFromIds(departmentIDs.([]interface{}))
			}
			if userNames, ok := exclusionData["jss_user_names"]; ok {
				scope.Exclusions.Users = constructScopeEntitiesFromIdsFromNames(userNames.([]interface{}))
			}
		}

		restrictedSoftware.Scope = scope
	}

	// Serialize and pretty-print the restrictedSoftware object as XML for logging
	resourceXML, err := xml.MarshalIndent(restrictedSoftware, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Restricted Software '%s' to XML: %v", restrictedSoftware.General.Name, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro Restricted Software XML:\n%s\n", string(resourceXML))

	return restrictedSoftware, nil
}

// Helper functions for nested structures

// constructScopeEntitiesFromIds constructs a slice of RestrictedSoftwareSubsetScopeEntity from a list of IDs.
func constructScopeEntitiesFromIds(ids []interface{}) []jamfpro.RestrictedSoftwareSubsetScopeEntity {
	scopeEntities := make([]jamfpro.RestrictedSoftwareSubsetScopeEntity, len(ids))
	for i, id := range ids {
		scopeEntities[i] = jamfpro.RestrictedSoftwareSubsetScopeEntity{
			ID: id.(int),
		}
	}
	return scopeEntities
}

// constructScopeEntitiesFromIdsFromNames constructs a slice of RestrictedSoftwareSubsetScopeEntity from a list of names.
func constructScopeEntitiesFromIdsFromNames(names []interface{}) []jamfpro.RestrictedSoftwareSubsetScopeEntity {
	scopeEntities := make([]jamfpro.RestrictedSoftwareSubsetScopeEntity, len(names))
	for i, name := range names {
		scopeEntities[i] = jamfpro.RestrictedSoftwareSubsetScopeEntity{
			Name: name.(string),
		}
	}
	return scopeEntities
}
