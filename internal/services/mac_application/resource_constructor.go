package mac_application

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/constructors"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	errMarshalMacAppXML = errors.New("failed to marshal Mac Application to XML")
)

// construct constructs a ResourceMacApplications object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceMacApplications, error) {
	resource := &jamfpro.ResourceMacApplications{
		General: jamfpro.MacApplicationsSubsetGeneral{
			Name:           d.Get("name").(string),
			Version:        d.Get("version").(string),
			BundleID:       d.Get("bundle_id").(string),
			URL:            d.Get("url").(string),
			IsFree:         jamfpro.BoolPtr(d.Get("is_free").(bool)),
			DeploymentType: d.Get("deployment_type").(string),
		},
	}

	resource.General.Site = sharedschemas.ConstructSharedResourceSite(d.Get("site_id").(int))
	resource.General.Category = sharedschemas.ConstructSharedResourceCategory(d.Get("category_id").(int))

	if _, ok := d.GetOk("scope"); ok {
		resource.Scope = constructMacApplicationScope(d)
	} else {
		log.Printf("[DEBUG] construct: No scope block found or it's empty.")
	}

	if v, ok := d.GetOk("self_service"); ok && len(v.([]any)) > 0 {
		selfServiceMap := v.([]any)[0].(map[string]any)
		selfService := jamfpro.MacAppSubsetSelfService{
			InstallButtonText:           selfServiceMap["install_button_text"].(string),
			SelfServiceDescription:      selfServiceMap["self_service_description"].(string),
			ForceUsersToViewDescription: jamfpro.BoolPtr(selfServiceMap["force_users_to_view_description"].(bool)),
			FeatureOnMainPage:           jamfpro.BoolPtr(selfServiceMap["feature_on_main_page"].(bool)),
			Notification:                selfServiceMap["notification"].(string),
			NotificationSubject:         selfServiceMap["notification_subject"].(string),
			NotificationMessage:         selfServiceMap["notification_message"].(string),
		}

		if categories, ok := selfServiceMap["self_service_category"].([]any); ok {
			var selfServiceCategories []jamfpro.MacAppSubsetSelfServiceCategories
			for _, cat := range categories {
				category := cat.(map[string]any)
				selfServiceCategories = append(selfServiceCategories, jamfpro.MacAppSubsetSelfServiceCategories{
					ID:        category["id"].(int),
					Name:      category["name"].(string),
					DisplayIn: jamfpro.BoolPtr(category["display_in"].(bool)),
					FeatureIn: jamfpro.BoolPtr(category["feature_in"].(bool)),
				})
			}
			selfService.SelfServiceCategories = selfServiceCategories
		}

		if icon, ok := selfServiceMap["self_service_icon"].([]any); ok && len(icon) > 0 {
			iconMap := icon[0].(map[string]any)
			selfService.SelfServiceIcon = jamfpro.SharedResourceSelfServiceIcon{
				ID:   iconMap["id"].(int),
				Data: iconMap["data"].(string),
				URI:  iconMap["uri"].(string),
			}
		}

		resource.SelfService = selfService
	}

	if v, ok := d.GetOk("vpp"); ok && len(v.([]any)) > 0 {
		vppMap := v.([]any)[0].(map[string]any)
		resource.VPP = jamfpro.MacAppSubsetVPP{
			AssignVPPDeviceBasedLicenses: jamfpro.BoolPtr(vppMap["assign_vpp_device_based_licenses"].(bool)),
			VPPAdminAccountID:            vppMap["vpp_admin_account_id"].(int),
		}
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		log.Printf("[ERROR] Failed to marshal Mac Application '%s' to XML: %v", resource.General.Name, err)
		return nil, fmt.Errorf("%w", errMarshalMacAppXML)
	}

	log.Printf("[DEBUG] Constructed Mac Application XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// constructMacApplicationScope constructs the scope from the provided schema data.
func constructMacApplicationScope(d *schema.ResourceData) jamfpro.MacApplicationsSubsetScope {
	scope := jamfpro.MacApplicationsSubsetScope{}
	scopeData := d.Get("scope").([]any)[0].(map[string]any)

	scope.AllComputers = jamfpro.BoolPtr(scopeData["all_computers"].(bool))
	scope.AllJSSUsers = jamfpro.BoolPtr(scopeData["all_jss_users"].(bool))

	var buildings []jamfpro.MacAppSubsetScopeBuilding
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeBuilding, int](
		"scope.0.building_ids", "ID", d, &buildings); err == nil {
		scope.Buildings = buildings
	}

	var departments []jamfpro.MacAppSubsetScopeDepartment
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeDepartment, int](
		"scope.0.department_ids", "ID", d, &departments); err == nil {
		scope.Departments = departments
	}

	var computers []jamfpro.MacAppSubsetScopeComputer
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeComputer, int](
		"scope.0.computer_ids", "ID", d, &computers); err == nil {
		scope.Computers = computers
	}

	var computerGroups []jamfpro.MacAppSubsetScopeComputerGroup
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeComputerGroup, int](
		"scope.0.computer_group_ids", "ID", d, &computerGroups); err == nil {
		scope.ComputerGroups = computerGroups
	}

	var jssUsers []jamfpro.MacAppSubsetScopeUser
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeUser, int](
		"scope.0.jss_user_ids", "ID", d, &jssUsers); err == nil {
		scope.JSSUsers = jssUsers
	}

	var jssUserGroups []jamfpro.MacAppSubsetScopeUserGroup
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeUserGroup, int](
		"scope.0.jss_user_group_ids", "ID", d, &jssUserGroups); err == nil {
		scope.JSSUserGroups = jssUserGroups
	}

	if _, ok := d.GetOk("scope.0.limitations"); ok {
		scope.Limitations = constructLimitations(d)
	}

	if _, ok := d.GetOk("scope.0.exclusions"); ok {
		scope.Exclusions = constructExclusions(d)
	}

	return scope
}

// constructLimitations constructs the limitations from the provided schema data.
func constructLimitations(d *schema.ResourceData) jamfpro.MacAppScopeLimitations {
	limitations := jamfpro.MacAppScopeLimitations{}

	var users []jamfpro.MacAppSubsetScopeUser
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeUser, string](
		"scope.0.limitations.0.users", "Name", d, &users); err == nil {
		limitations.Users = users
	}

	var userGroups []jamfpro.MacAppSubsetScopeUserGroup
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeUserGroup, int](
		"scope.0.limitations.0.user_groups", "ID", d, &userGroups); err == nil {
		limitations.UserGroups = userGroups
	}

	var networkSegments []jamfpro.MacAppSubsetScopeNetworkSegment
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeNetworkSegment, int](
		"scope.0.limitations.0.network_segments", "ID", d, &networkSegments); err == nil {
		limitations.NetworkSegments = networkSegments
	}

	return limitations
}

// constructExclusions constructs the exclusions from the provided schema data.
func constructExclusions(d *schema.ResourceData) jamfpro.MacAppScopeExclusions {
	exclusions := jamfpro.MacAppScopeExclusions{}

	var buildings []jamfpro.MacAppSubsetScopeBuilding
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeBuilding, int](
		"scope.0.exclusions.0.building_ids", "ID", d, &buildings); err == nil {
		exclusions.Buildings = buildings
	}

	var departments []jamfpro.MacAppSubsetScopeDepartment
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeDepartment, int](
		"scope.0.exclusions.0.department_ids", "ID", d, &departments); err == nil {
		exclusions.Departments = departments
	}

	var users []jamfpro.MacAppSubsetScopeUser
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeUser, string](
		"scope.0.exclusions.0.users", "Name", d, &users); err == nil {
		exclusions.Users = users
	}

	var userGroups []jamfpro.MacAppSubsetScopeUserGroup
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeUserGroup, int](
		"scope.0.exclusions.0.user_groups", "ID", d, &userGroups); err == nil {
		exclusions.UserGroups = userGroups
	}

	var networkSegments []jamfpro.MacAppSubsetScopeNetworkSegment
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeNetworkSegment, int](
		"scope.0.exclusions.0.network_segments", "ID", d, &networkSegments); err == nil {
		exclusions.NetworkSegments = networkSegments
	}

	var computers []jamfpro.MacAppSubsetScopeComputer
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeComputer, int](
		"scope.0.exclusions.0.computer_ids", "ID", d, &computers); err == nil {
		exclusions.Computers = computers
	}

	var computerGroups []jamfpro.MacAppSubsetScopeComputerGroup
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeComputerGroup, int](
		"scope.0.exclusions.0.computer_group_ids", "ID", d, &computerGroups); err == nil {
		exclusions.ComputerGroups = computerGroups
	}

	var jssUsers []jamfpro.MacAppSubsetScopeUser
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeUser, int](
		"scope.0.exclusions.0.jss_user_ids", "ID", d, &jssUsers); err == nil {
		exclusions.JSSUsers = jssUsers
	}

	var jssUserGroups []jamfpro.MacAppSubsetScopeUserGroup
	if err := constructors.MapSetToStructs[jamfpro.MacAppSubsetScopeUserGroup, int](
		"scope.0.exclusions.0.jss_user_group_ids", "ID", d, &jssUserGroups); err == nil {
		exclusions.JSSUserGroups = jssUserGroups
	}

	return exclusions
}
