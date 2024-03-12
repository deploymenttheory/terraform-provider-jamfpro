package macosconfigurationprofiles

import (
	"fmt"
	"html"

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
			UserRemovable:      d.Get("user_removeable").(bool),
			Level:              d.Get("level").(string),
			UUID:               d.Get("uuid").(string), // TODO not sure if this is needed as it's computed
			// RedeployOnUpdate:   d.Get("redeploy_on_update").(string), // TODO Review this, I don't think it's in the UI
		},
		Scope: jamfpro.MacOSConfigurationProfileSubsetScope{},
		SelfService: jamfpro.MacOSConfigurationProfileSubsetSelfService{
			InstallButtonText:           d.Get("self_service.0.install_button_text").(string),
			SelfServiceDescription:      d.Get("self_service.0.self_service_description").(string),
			ForceUsersToViewDescription: d.Get("self_service.0.force_users_to_view_description").(bool),
			// Self Service Icon - TBA at a later date because jamf is odd
			FeatureOnMainPage: d.Get("self_service.0.feature_on_main_page").(bool),
			// Self Service Categories
			// Notification parsed cos it's stupid and has dupe keys
			NotificationSubject: d.Get("self_service.0.notification_subject").(string),
			NotificationMessage: d.Get("self_service.0.notification_message").(string),
		},
	}

	// Processed Fields

	// Site
	if len(d.Get("site").([]interface{})) != 0 {
		out.General.Site = jamfpro.SharedResourceSite{
			ID:   d.Get("site.0.id").(int),
			Name: d.Get("site.0.name").(string),
		}
	}

	// Category
	if len(d.Get("category").([]interface{})) != 0 {
		out.General.Category = jamfpro.SharedResourceCategory{
			ID:   d.Get("category.0.id").(int),
			Name: d.Get("category.0.name").(string),
		}
	}

	// Payload
	payload, ok := d.GetOk("payload")
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
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetComputer, int]("scope.0.computer_ids", "ID", d, &out.Scope.Computers)
	if err != nil {
		return nil, err
	}

	// Computer Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetComputerGroup, int]("scope.0.computer_group_ids", "ID", d, &out.Scope.ComputerGroups)
	if err != nil {
		return nil, err
	}

	// JSS Users
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetJSSUser, int]("scope.0.jss_user_ids", "ID", d, &out.Scope.JSSUsers)
	if err != nil {
		return nil, err
	}

	// JSS User Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetJSSUserGroup, int]("scope.0.jss_user_group_ids", "ID", d, &out.Scope.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	// Buildings
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetBuilding, int]("scope.0.building_ids", "ID", d, &out.Scope.Buildings)
	if err != nil {
		return nil, err
	}

	// Departments
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetDepartment, int]("scope.0.department_ids", "ID", d, &out.Scope.Departments)
	if err != nil {
		return nil, err
	}

	// Scope - Limitations

	// Users
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetUser, string]("scope.0.limitations.0.user_names", "Name", d, &out.Scope.Limitations.Users)
	if err != nil {
		return nil, err
	}

	// Network Segment
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, int]("scope.0.limitations.0.network_segment_ids", "ID", d, &out.Scope.Limitations.NetworkSegments)
	if err != nil {
		return nil, err
	}

	// IBeacons
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetIBeacon, int]("scope.0.limitations.0.ibeacon_ids", "ID", d, &out.Scope.Limitations.IBeacons)
	if err != nil {
		return nil, err
	}

	// User Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetUserGroup, int]("scope.0.limitations.0.user_group_ids", "ID", d, &out.Scope.Limitations.UserGroups)
	if err != nil {
		return nil, err
	}

	// Scope - Limitations

	// Computers
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetComputer, int]("scope.0.exclusions.0.computer_ids", "ID", d, &out.Scope.Exclusions.Computers)
	if err != nil {
		return nil, err
	}

	// Computer Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetComputerGroup, int]("scope.0.exclusions.0.computer_group_ids", "ID", d, &out.Scope.Exclusions.ComputerGroups)
	if err != nil {
		return nil, err
	}

	// Buildings
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetBuilding, int]("scope.0.exclusions.0.building_ids", "ID", d, &out.Scope.Exclusions.Buildings)
	if err != nil {
		return nil, err
	}

	// Departments
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetDepartment, int]("scope.0.exclusions.0.department_ids", "ID", d, &out.Scope.Exclusions.Departments)
	if err != nil {
		return nil, err
	}

	// Network Segments
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, int]("scope.0.exclusions.0.network_segment_ids", "ID", d, &out.Scope.Exclusions.NetworkSegments)
	if err != nil {
		return nil, err
	}

	// JSS Users
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetJSSUser, int]("scope.0.exclusions.0.jss_user_ids", "ID", d, &out.Scope.Exclusions.JSSUsers)
	if err != nil {
		return nil, err
	}

	// JSS User Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetJSSUserGroup, int]("scope.0.exclusions.0.jss_user_group_ids", "ID", d, &out.Scope.Exclusions.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	// IBeacons
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetIBeacon, int]("scope.0.exclusions.0.ibeacon_ids", "ID", d, &out.Scope.Exclusions.IBeacons)
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

	return &out, nil
}
