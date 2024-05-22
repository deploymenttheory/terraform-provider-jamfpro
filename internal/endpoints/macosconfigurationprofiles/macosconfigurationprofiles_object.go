package macosconfigurationprofiles

import (
	"encoding/xml"
	"fmt"
	"html"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProMacOSConfigurationProfile constructs a ResourceMacOSConfigurationProfile object from the provided schema data.
func constructJamfProMacOSConfigurationProfile(d *schema.ResourceData) (*jamfpro.ResourceMacOSConfigurationProfile, error) {
	// Main obj with fields which do not require processing
	out := jamfpro.ResourceMacOSConfigurationProfile{
		General: jamfpro.MacOSConfigurationProfileSubsetGeneral{
			Name:               d.Get("name").(string),
			Description:        d.Get("description").(string),
			DistributionMethod: d.Get("distribution_method").(string),
			UserRemovable:      d.Get("user_removable").(bool),
			Level:              d.Get("level").(string),
			UUID:               d.Get("uuid").(string), // TODO not sure if this is needed as it's computed
			RedeployOnUpdate:   d.Get("redeploy_on_update").(string),
		},
		Scope: jamfpro.MacOSConfigurationProfileSubsetScope{},
		SelfService: jamfpro.MacOSConfigurationProfileSubsetSelfService{
			InstallButtonText:           d.Get("self_service.0.install_button_text").(string),
			SelfServiceDescription:      d.Get("self_service.0.self_service_description").(string),
			ForceUsersToViewDescription: d.Get("self_service.0.force_users_to_view_description").(bool),
			// Self Service Icon - // TODO at a later date because jamf is odd
			FeatureOnMainPage: d.Get("self_service.0.feature_on_main_page").(bool),
			// Self Service Categories
			// Notification parsed cos it's stupid and has dupe keys
			NotificationSubject: d.Get("self_service.0.notification_subject").(string),
			NotificationMessage: d.Get("self_service.0.notification_message").(string),
		},
	}

	// Processed Fields

	// Site
	suppliedSite := d.Get("site").([]interface{})
	if len(suppliedSite) > 0 {
		// If site provided, construct
		outSite := jamfpro.SharedResourceSite{
			ID: suppliedSite[0].(map[string]interface{})["id"].(int),
		}
		out.General.Site = outSite
	} else {
		// If no site, construct no site obj. We have to do this for the site to be removed.
		out.General.Site = jamfpro.SharedResourceSite{
			ID: 0,
		}
	}

	// Category
	suppliedCategory := d.Get("category").([]interface{})
	if len(suppliedCategory) > 0 {
		// construct category if provided
		outCat := jamfpro.SharedResourceCategory{
			ID: suppliedCategory[0].(map[string]interface{})["id"].(int),
		}
		out.General.Category = outCat
	} else {
		// if no category, supply empty cat to remove it.
		out.General.Category = jamfpro.SharedResourceCategory{
			ID: 0,
		}
	}

	// Payload
	payload, ok := d.GetOk("payloads")
	if ok {
		payload = html.EscapeString(payload.(string))
		out.General.Payloads = payload.(string)
	} else {
		return nil, fmt.Errorf("an error occurred setting the payload")
	}

	// Scope
	var err error

	// Scope - Targets
	out.Scope.AllComputers = d.Get("scope.0.all_computers").(bool)
	out.Scope.AllJSSUsers = d.Get("scope.0.all_jss_users").(bool)

	// Computers
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetComputer, int]("scope.0.computer_ids", "ID", d, &out.Scope.Computers)
	if err != nil {
		return nil, err
	}

	// Computer Groups
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.computer_group_ids", "ID", d, &out.Scope.ComputerGroups)
	if err != nil {
		return nil, err
	}

	// JSS Users
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.jss_user_ids", "ID", d, &out.Scope.JSSUsers)
	if err != nil {
		return nil, err
	}

	// JSS User Groups
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.jss_user_group_ids", "ID", d, &out.Scope.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	// Buildings
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.building_ids", "ID", d, &out.Scope.Buildings)
	if err != nil {
		return nil, err
	}

	// Departments
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.department_ids", "ID", d, &out.Scope.Departments)
	if err != nil {
		return nil, err
	}

	// Scope - Limitations

	// Network Segment
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, int]("scope.0.limitations.0.network_segment_ids", "ID", d, &out.Scope.Limitations.NetworkSegments)
	if err != nil {
		return nil, err
	}

	// directory service Users
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, string]("scope.0.limitations.0.directory_service_or_local_usernames", "Name", d, &out.Scope.Limitations.Users)
	if err != nil {
		return nil, err
	}

	// directory service User Groups
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.limitations.0.directory_service_usergroup_ids", "ID", d, &out.Scope.Limitations.UserGroups)
	if err != nil {
		return nil, err
	}

	// IBeacons
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.limitations.0.ibeacon_ids", "ID", d, &out.Scope.Limitations.IBeacons)
	if err != nil {
		return nil, err
	}

	// Scope - Exclusions

	// Computers
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetComputer, int]("scope.0.exclusions.0.computer_ids", "ID", d, &out.Scope.Exclusions.Computers)
	if err != nil {
		return nil, err
	}

	// Computer Groups
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.exclusions.0.computer_group_ids", "ID", d, &out.Scope.Exclusions.ComputerGroups)
	if err != nil {
		return nil, err
	}

	// JSS Users
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.exclusions.0.jss_user_ids", "ID", d, &out.Scope.Exclusions.JSSUsers)
	if err != nil {
		return nil, err
	}

	// JSS User Groups
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.exclusions.0.jss_user_group_ids", "ID", d, &out.Scope.Exclusions.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	// Buildings
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.exclusions.0.building_ids", "ID", d, &out.Scope.Exclusions.Buildings)
	if err != nil {
		return nil, err
	}

	// Departments
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.exclusions.0.department_ids", "ID", d, &out.Scope.Exclusions.Departments)
	if err != nil {
		return nil, err
	}

	// Network Segments
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, int]("scope.0.exclusions.0.network_segment_ids", "ID", d, &out.Scope.Exclusions.NetworkSegments)
	if err != nil {
		return nil, err
	}

	// directory service Users
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, string]("scope.0.exclusions.0.directory_service_or_local_usernames", "Name", d, &out.Scope.Limitations.Users)
	if err != nil {
		return nil, err
	}

	// directory service User Groups
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.exclusions.0.directory_service_usergroup_ids", "ID", d, &out.Scope.Limitations.UserGroups)
	if err != nil {
		return nil, err
	}

	// IBeacons
	err = ExtractNestedObjectsFromSchema[jamfpro.MacOSConfigurationProfileSubsetScopeEntity, int]("scope.0.exclusions.0.ibeacon_ids", "ID", d, &out.Scope.Exclusions.IBeacons)
	if err != nil {
		return nil, err
	}

	// TODO make this better, it works for now
	if out.Scope.AllComputers && (out.Scope.Computers != nil ||
		out.Scope.ComputerGroups != nil ||
		out.Scope.Departments != nil ||
		out.Scope.Buildings != nil) {
		return nil, fmt.Errorf("invalid combination - all computers with scoped endpoints")
	}

	// Self Service
	// TODO move this to a helper or omit whole key. Logic bad.
	value, ok := d.GetOk("self_service.0.self_service_categories")
	if ok {
		listOfVals := value.([]interface{})
		for _, v := range listOfVals {
			mapOfVals := v.(map[string]interface{})
			catId := mapOfVals["id"]
			displayIn := mapOfVals["display_in"]
			featureIn := mapOfVals["feature_in"]
			name := mapOfVals["name"]
			out.SelfService.SelfServiceCategories = append(out.SelfService.SelfServiceCategories, jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory{
				Name:      name.(string),
				ID:        catId.(int),
				DisplayIn: displayIn.(bool),
				FeatureIn: featureIn.(bool),
			})
		}
	}

	// Serialize and pretty-print the macOS Configuration Profile object as XML for logging
	resourceXML, err := xml.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro macOS Configuration Profile '%s' to XML: %v", out.General.Name, err)
	}

	// Use log.Printf instead of fmt.Printf for logging within the Terraform provider context
	log.Printf("[DEBUG] Constructed Jamf Pro macOS Configuration Profile XML:\n%s\n", string(resourceXML))

	return &out, nil
}
