// policies_object.go
package policies

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProPolicy creates a new ResourcePolicy struct from the given data.
func constructJamfProPolicy(ctx context.Context, data *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {
	// Initialize a new ResourcePolicy struct.
	policy := &jamfpro.ResourcePolicy{
		General:              constructGeneral(data),
		Scope:                constructScope(data),
		SelfService:          constructSelfService(data),
		PackageConfiguration: constructPackageConfiguration(data),
		Scripts:              constructScripts(data),
		Printers:             constructPrinters(data),
		DockItems:            constructDockItems(data),
		AccountMaintenance:   constructAccountMaintenance(data),
		Maintenance:          constructMaintenance(data),
		FilesProcesses:       constructFilesProcesses(data),
		UserInteraction:      constructUserInteraction(data),
		DiskEncryption:       constructDiskEncryption(data),
		Reboot:               constructReboot(data),
	}

	// Marshal the jamf pro policy object into XML for logging
	xmlData, err := xml.MarshalIndent(policy, "", "  ")
	if err != nil {
		// Handle the error if XML marshaling fails
		log.Printf("[ERROR] Error marshaling jamf pro policy object to XML: %s", err)
		return nil, fmt.Errorf("error marshaling jamf pro policy object to XML: %v", err)
	}

	// Log the XML formatted search object
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Constructed jamf pro policy Object:\n%s", string(xmlData)))

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Successfully constructed jamf pro policy with name: %s", policy.General.Name))

	return policy, nil // Return the policy object and no error
}

// constructGeneral creates a General policy subset from the provided data.
func constructGeneral(data *schema.ResourceData) jamfpro.PolicySubsetGeneral {
	generalData := data.Get("general").([]interface{})[0].(map[string]interface{})

	general := jamfpro.PolicySubsetGeneral{
		ID:                         util.GetIntFromInterface(generalData["id"]),
		Name:                       util.GetStringFromInterface(generalData["name"]),
		Enabled:                    util.GetBoolFromInterface(generalData["enabled"]),
		Trigger:                    util.GetStringFromInterface(generalData["trigger"]),
		TriggerCheckin:             util.GetBoolFromInterface(generalData["trigger_checkin"]),
		TriggerEnrollmentComplete:  util.GetBoolFromInterface(generalData["trigger_enrollment_complete"]),
		TriggerLogin:               util.GetBoolFromInterface(generalData["trigger_login"]),
		TriggerLogout:              util.GetBoolFromInterface(generalData["trigger_logout"]),
		TriggerNetworkStateChanged: util.GetBoolFromInterface(generalData["trigger_network_state_changed"]),
		TriggerStartup:             util.GetBoolFromInterface(generalData["trigger_startup"]),
		TriggerOther:               util.GetStringFromInterface(generalData["trigger_other"]),
		Frequency:                  util.GetStringFromInterface(generalData["frequency"]),
		RetryEvent:                 util.GetStringFromInterface(generalData["retry_event"]),
		RetryAttempts:              util.GetIntFromInterface(generalData["retry_attempts"]),
		NotifyOnEachFailedRetry:    util.GetBoolFromInterface(generalData["notify_on_each_failed_retry"]),
		LocationUserOnly:           util.GetBoolFromInterface(generalData["location_user_only"]),
		TargetDrive:                util.GetStringFromInterface(generalData["target_drive"]),
		Offline:                    util.GetBoolFromInterface(generalData["offline"]),
		Category:                   constructPolicyCategory(generalData["category"]),
		DateTimeLimitations:        constructDateTimeLimitations(generalData["date_time_limitations"]),
		NetworkLimitations:         constructNetworkLimitations(generalData["network_limitations"]),
		OverrideDefaultSettings:    constructOverrideDefaultSettings(generalData["override_default_settings"]),
		NetworkRequirements:        util.GetStringFromInterface(generalData["network_requirements"]),
		Site:                       constructSharedResourceSite(generalData["site"]),
	}

	return general
}

// constructPolicyCategory creates a PolicyCategory struct from an interface value.
// This function safely extracts data for PolicyCategory fields from a map
// contained within the interface. It handles nil values and type assertions.
func constructPolicyCategory(val interface{}) jamfpro.PolicyCategory {
	categoryData := util.ConvertToMapFromInterface(val)
	if categoryData == nil {
		return jamfpro.PolicyCategory{}
	}

	return jamfpro.PolicyCategory{
		ID:        util.GetIntFromInterface(categoryData["id"]),
		Name:      util.GetStringFromInterface(categoryData["name"]),
		DisplayIn: util.GetBoolFromInterface(categoryData["display_in"]),
		FeatureIn: util.GetBoolFromInterface(categoryData["feature_in"]),
	}
}

// constructDateTimeLimitations creates a PolicySubsetGeneralDateTimeLimitations struct
// from an interface value. This function extracts and sets various date and time
// limitations related fields for a policy. It processes both direct fields and
// slices (for no_execute_on days) and handles nil values and type assertions.
func constructDateTimeLimitations(val interface{}) jamfpro.PolicySubsetGeneralDateTimeLimitations {
	dateTimeData := util.ConvertToMapFromInterface(val)
	if dateTimeData == nil {
		return jamfpro.PolicySubsetGeneralDateTimeLimitations{}
	}

	return jamfpro.PolicySubsetGeneralDateTimeLimitations{
		ActivationDate:      util.GetStringFromInterface(dateTimeData["activation_date"]),
		ActivationDateEpoch: util.GetIntFromInterface(dateTimeData["activation_date_epoch"]),
		ActivationDateUTC:   util.GetStringFromInterface(dateTimeData["activation_date_utc"]),
		ExpirationDate:      util.GetStringFromInterface(dateTimeData["expiration_date"]),
		ExpirationDateEpoch: util.GetIntFromInterface(dateTimeData["expiration_date_epoch"]),
		ExpirationDateUTC:   util.GetStringFromInterface(dateTimeData["expiration_date_utc"]),
		NoExecuteOn:         constructNoExecuteOn(dateTimeData["no_execute_on"]),
		NoExecuteStart:      util.GetStringFromInterface(dateTimeData["no_execute_start"]),
		NoExecuteEnd:        util.GetStringFromInterface(dateTimeData["no_execute_end"]),
	}
}

// constructNetworkLimitations creates a PolicySubsetGeneralNetworkLimitations struct
// from an interface value. It extracts data related to network limitations
// for a policy, including minimum network connection requirements and any IP address
// applicability. It handles nil values and type assertions.
func constructNetworkLimitations(val interface{}) jamfpro.PolicySubsetGeneralNetworkLimitations {
	networkData := util.ConvertToMapFromInterface(val)
	if networkData == nil {
		return jamfpro.PolicySubsetGeneralNetworkLimitations{}
	}

	return jamfpro.PolicySubsetGeneralNetworkLimitations{
		MinimumNetworkConnection: util.GetStringFromInterface(networkData["minimum_network_connection"]),
		AnyIPAddress:             util.GetBoolFromInterface(networkData["any_ip_address"]),
		NetworkSegments:          util.GetStringFromInterface(networkData["network_segments"]),
	}
}

// constructNoExecuteOn creates a slice of PolicySubsetGeneralDateTimeLimitationsNoExecuteOn structs
// from an interface value. This function is specifically used for extracting the days
// on which policy execution is not to occur ('no_execute_on'). It iterates over a slice
// of strings, each representing a day, and constructs a corresponding struct slice.
// It safely handles type assertions and nil values using utility functions.
func constructNoExecuteOn(val interface{}) []jamfpro.PolicySubsetGeneralDateTimeLimitationsNoExecuteOn {
	noExecuteOnSlice := util.GetStringSliceFromInterface(val)
	var noExecuteOn []jamfpro.PolicySubsetGeneralDateTimeLimitationsNoExecuteOn
	for _, day := range noExecuteOnSlice {
		noExecuteOn = append(noExecuteOn, jamfpro.PolicySubsetGeneralDateTimeLimitationsNoExecuteOn{Day: day})
	}
	return noExecuteOn
}

// constructOverrideDefaultSettings creates a PolicySubsetGeneralOverrideDefaultSettings struct
// from an interface value. This function extracts override settings, including target drive
// and distribution point, and handles nil values and type assertions.
func constructOverrideDefaultSettings(val interface{}) jamfpro.PolicySubsetGeneralOverrideDefaultSettings {
	overrideData := util.ConvertToMapFromInterface(val)
	if overrideData == nil {
		return jamfpro.PolicySubsetGeneralOverrideDefaultSettings{}
	}

	return jamfpro.PolicySubsetGeneralOverrideDefaultSettings{
		TargetDrive:       util.GetStringFromInterface(overrideData["target_drive"]),
		DistributionPoint: util.GetStringFromInterface(overrideData["distribution_point"]),
		ForceAfpSmb:       util.GetBoolFromInterface(overrideData["force_afp_smb"]),
		SUS:               util.GetStringFromInterface(overrideData["sus"]),
	}
}

// constructSharedResourceSite creates a SharedResourceSite struct from an interface value.
// This function is used to extract data for a shared resource site, handling nil values
// and type assertions safely.
func constructSharedResourceSite(val interface{}) jamfpro.SharedResourceSite {
	siteData := util.ConvertToMapFromInterface(val)
	if siteData == nil {
		return jamfpro.SharedResourceSite{}
	}

	return jamfpro.SharedResourceSite{
		ID:   util.GetIntFromInterface(siteData["id"]),
		Name: util.GetStringFromInterface(siteData["name"]),
	}
}

// constructScope creates a PolicySubsetScope struct from the provided data.
func constructScope(data *schema.ResourceData) jamfpro.PolicySubsetScope {
	scopeData := util.ConvertToMapFromInterface(data.Get("scope").([]interface{})[0])

	scope := jamfpro.PolicySubsetScope{
		AllComputers:   util.GetBoolFromInterface(scopeData["all_computers"]),
		AllJSSUsers:    util.GetBoolFromInterface(scopeData["all_jss_users"]),
		Computers:      constructPolicyDataSubsetsComputer(scopeData["computers"]),
		ComputerGroups: constructPolicyDataSubsetsComputerGroup(scopeData["computer_groups"]),
		JSSUsers:       constructPolicyDataSubsetsJSSUser(scopeData["jss_users"]),
		JSSUserGroups:  constructPolicyDataSubsetsJSSUserGroup(scopeData["jss_user_groups"]),
		Buildings:      constructPolicyDataSubsetsBuilding(scopeData["buildings"]),
		Departments:    constructPolicyDataSubsetsDepartment(scopeData["departments"]),
		LimitToUsers:   constructPolicyLimitToUsers(scopeData["limit_to_users"]),
		Limitations:    constructPolicySubsetScopeLimitations(scopeData["limitations"]),
		Exclusions:     constructPolicySubsetScopeExclusions(scopeData["exclusions"]),
	}

	return scope
}

// constructPolicyDataSubsetsComputer creates a slice of PolicyDataSubsetComputer from the provided data.
func constructPolicyDataSubsetsComputer(val interface{}) []jamfpro.PolicyDataSubsetComputer {
	computersInterface, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var computers []jamfpro.PolicyDataSubsetComputer
	for _, ci := range computersInterface {
		computerMap, ok := ci.(map[string]interface{})
		if !ok {
			continue
		}

		computer := jamfpro.PolicyDataSubsetComputer{
			ID:   util.GetIntFromInterface(computerMap["id"]),
			Name: util.GetStringFromInterface(computerMap["name"]),
			UDID: util.GetStringFromInterface(computerMap["udid"]),
		}
		computers = append(computers, computer)
	}
	return computers
}

// constructPolicyDataSubsetsComputerGroup creates a slice of PolicyDataSubsetComputerGroup from the provided data.
func constructPolicyDataSubsetsComputerGroup(val interface{}) []jamfpro.PolicyDataSubsetComputerGroup {
	groupsInterface, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var groups []jamfpro.PolicyDataSubsetComputerGroup
	for _, gi := range groupsInterface {
		groupMap, ok := gi.(map[string]interface{})
		if !ok {
			continue
		}

		group := jamfpro.PolicyDataSubsetComputerGroup{
			ID:   util.GetIntFromInterface(groupMap["id"]),
			Name: util.GetStringFromInterface(groupMap["name"]),
		}
		groups = append(groups, group)
	}
	return groups
}

// constructPolicyDataSubsetsJSSUser creates a slice of PolicyDataSubsetJSSUser from the provided data.
func constructPolicyDataSubsetsJSSUser(val interface{}) []jamfpro.PolicyDataSubsetJSSUser {
	usersInterface, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var users []jamfpro.PolicyDataSubsetJSSUser
	for _, ui := range usersInterface {
		userMap, ok := ui.(map[string]interface{})
		if !ok {
			continue
		}

		user := jamfpro.PolicyDataSubsetJSSUser{
			ID:   util.GetIntFromInterface(userMap["id"]),
			Name: util.GetStringFromInterface(userMap["name"]),
		}
		users = append(users, user)
	}
	return users
}

// constructPolicyDataSubsetsJSSUserGroup creates a slice of PolicyDataSubsetJSSUserGroup from the provided data.
func constructPolicyDataSubsetsJSSUserGroup(val interface{}) []jamfpro.PolicyDataSubsetJSSUserGroup {
	groupsInterface, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var groups []jamfpro.PolicyDataSubsetJSSUserGroup
	for _, gi := range groupsInterface {
		groupMap, ok := gi.(map[string]interface{})
		if !ok {
			continue
		}

		group := jamfpro.PolicyDataSubsetJSSUserGroup{
			ID:   util.GetIntFromInterface(groupMap["id"]),
			Name: util.GetStringFromInterface(groupMap["name"]),
		}
		groups = append(groups, group)
	}
	return groups
}

// constructPolicyDataSubsetsBuilding creates a slice of PolicyDataSubsetBuilding from the provided data.
func constructPolicyDataSubsetsBuilding(val interface{}) []jamfpro.PolicyDataSubsetBuilding {
	buildingsInterface, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var buildings []jamfpro.PolicyDataSubsetBuilding
	for _, bi := range buildingsInterface {
		buildingMap, ok := bi.(map[string]interface{})
		if !ok {
			continue
		}

		building := jamfpro.PolicyDataSubsetBuilding{
			ID:   util.GetIntFromInterface(buildingMap["id"]),
			Name: util.GetStringFromInterface(buildingMap["name"]),
		}
		buildings = append(buildings, building)
	}
	return buildings
}

// constructPolicyDataSubsetsDepartment creates a slice of PolicyDataSubsetDepartment from the provided data.
func constructPolicyDataSubsetsDepartment(val interface{}) []jamfpro.PolicyDataSubsetDepartment {
	departmentsInterface, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var departments []jamfpro.PolicyDataSubsetDepartment
	for _, di := range departmentsInterface {
		departmentMap, ok := di.(map[string]interface{})
		if !ok {
			continue
		}

		department := jamfpro.PolicyDataSubsetDepartment{
			ID:   util.GetIntFromInterface(departmentMap["id"]),
			Name: util.GetStringFromInterface(departmentMap["name"]),
		}
		departments = append(departments, department)
	}
	return departments
}

// constructPolicyLimitToUsers creates a PolicyLimitToUsers struct from the provided data.
func constructPolicyLimitToUsers(val interface{}) jamfpro.PolicyLimitToUsers {
	userGroupsInterface, ok := val.([]interface{})
	if !ok {
		return jamfpro.PolicyLimitToUsers{}
	}

	var userGroups []string
	for _, ug := range userGroupsInterface {
		userGroup, ok := ug.(string)
		if !ok {
			continue
		}
		userGroups = append(userGroups, userGroup)
	}

	return jamfpro.PolicyLimitToUsers{UserGroups: userGroups}
}

// constructPolicySubsetScopeLimitations creates a PolicySubsetScopeLimitations struct from the provided data.
func constructPolicySubsetScopeLimitations(val interface{}) jamfpro.PolicySubsetScopeLimitations {
	limitationsData, ok := val.(map[string]interface{})
	if !ok {
		return jamfpro.PolicySubsetScopeLimitations{}
	}

	var users []jamfpro.PolicyDataSubsetUser
	usersInterface, usersExist := limitationsData["users"].([]interface{})
	if usersExist {
		for _, ui := range usersInterface {
			userMap, ok := ui.(map[string]interface{})
			if !ok {
				continue
			}

			user := jamfpro.PolicyDataSubsetUser{
				ID:   util.GetIntFromInterface(userMap["id"]),
				Name: util.GetStringFromInterface(userMap["name"]),
			}
			users = append(users, user)
		}
	}

	var userGroups []jamfpro.PolicyDataSubsetUserGroup
	userGroupsInterface, userGroupsExist := limitationsData["user_groups"].([]interface{})
	if userGroupsExist {
		for _, ugi := range userGroupsInterface {
			userGroupMap, ok := ugi.(map[string]interface{})
			if !ok {
				continue
			}

			userGroup := jamfpro.PolicyDataSubsetUserGroup{
				ID:   util.GetIntFromInterface(userGroupMap["id"]),
				Name: util.GetStringFromInterface(userGroupMap["name"]),
			}
			userGroups = append(userGroups, userGroup)
		}
	}

	var networkSegments []jamfpro.PolicyDataSubsetNetworkSegment
	networkSegmentsInterface, networkSegmentsExist := limitationsData["network_segments"].([]interface{})
	if networkSegmentsExist {
		for _, nsi := range networkSegmentsInterface {
			networkSegmentMap, ok := nsi.(map[string]interface{})
			if !ok {
				continue
			}

			networkSegment := jamfpro.PolicyDataSubsetNetworkSegment{
				ID:   util.GetIntFromInterface(networkSegmentMap["id"]),
				Name: util.GetStringFromInterface(networkSegmentMap["name"]),
			}
			networkSegments = append(networkSegments, networkSegment)
		}
	}

	var iBeacons []jamfpro.PolicyDataSubsetIBeacon
	iBeaconsInterface, iBeaconsExist := limitationsData["ibeacons"].([]interface{})
	if iBeaconsExist {
		for _, ibi := range iBeaconsInterface {
			iBeaconMap, ok := ibi.(map[string]interface{})
			if !ok {
				continue
			}

			iBeacon := jamfpro.PolicyDataSubsetIBeacon{
				ID:   util.GetIntFromInterface(iBeaconMap["id"]),
				Name: util.GetStringFromInterface(iBeaconMap["name"]),
			}
			iBeacons = append(iBeacons, iBeacon)
		}
	}

	return jamfpro.PolicySubsetScopeLimitations{
		Users:           users,
		UserGroups:      userGroups,
		NetworkSegments: networkSegments,
		IBeacons:        iBeacons,
	}
}

// constructPolicyDataSubsetsUser creates a slice of PolicyDataSubsetUser from the provided data.
func constructPolicyDataSubsetsUser(val interface{}) []jamfpro.PolicyDataSubsetUser {
	usersInterface, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var users []jamfpro.PolicyDataSubsetUser
	for _, ui := range usersInterface {
		userMap, ok := ui.(map[string]interface{})
		if !ok {
			continue
		}

		user := jamfpro.PolicyDataSubsetUser{
			ID:   util.GetIntFromInterface(userMap["id"]),
			Name: util.GetStringFromInterface(userMap["name"]),
		}
		users = append(users, user)
	}
	return users
}

// constructPolicyDataSubsetsUserGroupcreates a slice of PolicyDataSubsetUserGroup from the provided data.
func constructPolicyDataSubsetsUserGroup(val interface{}) []jamfpro.PolicyDataSubsetUserGroup {
	groupsInterface, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var groups []jamfpro.PolicyDataSubsetUserGroup
	for _, gi := range groupsInterface {
		groupMap, ok := gi.(map[string]interface{})
		if !ok {
			continue
		}

		group := jamfpro.PolicyDataSubsetUserGroup{
			ID:   util.GetIntFromInterface(groupMap["id"]),
			Name: util.GetStringFromInterface(groupMap["name"]),
		}
		groups = append(groups, group)
	}
	return groups
}

// constructPolicyDataSubsetsNetworkSegment a slice of PolicyDataSubsetNetworkSegment from the provided data.
func constructPolicyDataSubsetsNetworkSegment(val interface{}) []jamfpro.PolicyDataSubsetNetworkSegment {
	segmentsInterface, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var segments []jamfpro.PolicyDataSubsetNetworkSegment
	for _, si := range segmentsInterface {
		segmentMap, ok := si.(map[string]interface{})
		if !ok {
			continue
		}

		segment := jamfpro.PolicyDataSubsetNetworkSegment{
			ID:   util.GetIntFromInterface(segmentMap["id"]),
			Name: util.GetStringFromInterface(segmentMap["name"]),
			UID:  util.GetStringFromInterface(segmentMap["uid"]),
		}
		segments = append(segments, segment)
	}
	return segments
}

// constructPolicySubsetScopeExclusions creates a PolicySubsetScopeExclusions struct from the provided data.
func constructPolicySubsetScopeExclusions(val interface{}) jamfpro.PolicySubsetScopeExclusions {
	exclusionsData, ok := val.(map[string]interface{})
	if !ok {
		return jamfpro.PolicySubsetScopeExclusions{}
	}

	var computers []jamfpro.PolicyDataSubsetComputer
	computersInterface, computersExist := exclusionsData["computers"].([]interface{})
	if computersExist {
		for _, ci := range computersInterface {
			computerMap, ok := ci.(map[string]interface{})
			if !ok {
				continue
			}

			computer := jamfpro.PolicyDataSubsetComputer{
				ID:   util.GetIntFromInterface(computerMap["id"]),
				Name: util.GetStringFromInterface(computerMap["name"]),
				UDID: util.GetStringFromInterface(computerMap["udid"]),
			}
			computers = append(computers, computer)
		}
	}

	var computerGroups []jamfpro.PolicyDataSubsetComputerGroup
	computerGroupsInterface, computerGroupsExist := exclusionsData["computer_groups"].([]interface{})
	if computerGroupsExist {
		for _, cgi := range computerGroupsInterface {
			computerGroupMap, ok := cgi.(map[string]interface{})
			if !ok {
				continue
			}

			computerGroup := jamfpro.PolicyDataSubsetComputerGroup{
				ID:   util.GetIntFromInterface(computerGroupMap["id"]),
				Name: util.GetStringFromInterface(computerGroupMap["name"]),
			}
			computerGroups = append(computerGroups, computerGroup)
		}
	}

	var users []jamfpro.PolicyDataSubsetUser
	usersInterface, usersExist := exclusionsData["users"].([]interface{})
	if usersExist {
		for _, ui := range usersInterface {
			userMap, ok := ui.(map[string]interface{})
			if !ok {
				continue
			}

			user := jamfpro.PolicyDataSubsetUser{
				ID:   util.GetIntFromInterface(userMap["id"]),
				Name: util.GetStringFromInterface(userMap["name"]),
			}
			users = append(users, user)
		}
	}

	var userGroups []jamfpro.PolicyDataSubsetUserGroup
	userGroupsInterface, userGroupsExist := exclusionsData["user_groups"].([]interface{})
	if userGroupsExist {
		for _, ugi := range userGroupsInterface {
			userGroupMap, ok := ugi.(map[string]interface{})
			if !ok {
				continue
			}

			userGroup := jamfpro.PolicyDataSubsetUserGroup{
				ID:   util.GetIntFromInterface(userGroupMap["id"]),
				Name: util.GetStringFromInterface(userGroupMap["name"]),
			}
			userGroups = append(userGroups, userGroup)
		}
	}

	var buildings []jamfpro.PolicyDataSubsetBuilding
	buildingsInterface, buildingsExist := exclusionsData["buildings"].([]interface{})
	if buildingsExist {
		for _, bi := range buildingsInterface {
			buildingMap, ok := bi.(map[string]interface{})
			if !ok {
				continue
			}

			building := jamfpro.PolicyDataSubsetBuilding{
				ID:   util.GetIntFromInterface(buildingMap["id"]),
				Name: util.GetStringFromInterface(buildingMap["name"]),
			}
			buildings = append(buildings, building)
		}
	}

	var departments []jamfpro.PolicyDataSubsetDepartment
	departmentsInterface, departmentsExist := exclusionsData["departments"].([]interface{})
	if departmentsExist {
		for _, di := range departmentsInterface {
			departmentMap, ok := di.(map[string]interface{})
			if !ok {
				continue
			}

			department := jamfpro.PolicyDataSubsetDepartment{
				ID:   util.GetIntFromInterface(departmentMap["id"]),
				Name: util.GetStringFromInterface(departmentMap["name"]),
			}
			departments = append(departments, department)
		}
	}

	var networkSegments []jamfpro.PolicyDataSubsetNetworkSegment
	networkSegmentsInterface, networkSegmentsExist := exclusionsData["network_segments"].([]interface{})
	if networkSegmentsExist {
		for _, nsi := range networkSegmentsInterface {
			networkSegmentMap, ok := nsi.(map[string]interface{})
			if !ok {
				continue
			}

			networkSegment := jamfpro.PolicyDataSubsetNetworkSegment{
				ID:   util.GetIntFromInterface(networkSegmentMap["id"]),
				Name: util.GetStringFromInterface(networkSegmentMap["name"]),
				UID:  util.GetStringFromInterface(networkSegmentMap["uid"]),
			}
			networkSegments = append(networkSegments, networkSegment)
		}
	}

	var jssUsers []jamfpro.PolicyDataSubsetJSSUser
	jssUsersInterface, jssUsersExist := exclusionsData["jss_users"].([]interface{})
	if jssUsersExist {
		for _, jui := range jssUsersInterface {
			jssUserMap, ok := jui.(map[string]interface{})
			if !ok {
				continue
			}

			jssUser := jamfpro.PolicyDataSubsetJSSUser{
				ID:   util.GetIntFromInterface(jssUserMap["id"]),
				Name: util.GetStringFromInterface(jssUserMap["name"]),
			}
			jssUsers = append(jssUsers, jssUser)
		}
	}

	var jssUserGroups []jamfpro.PolicyDataSubsetJSSUserGroup
	jssUserGroupsInterface, jssUserGroupsExist := exclusionsData["jss_user_groups"].([]interface{})
	if jssUserGroupsExist {
		for _, jugi := range jssUserGroupsInterface {
			jssUserGroupMap, ok := jugi.(map[string]interface{})
			if !ok {
				continue
			}

			jssUserGroup := jamfpro.PolicyDataSubsetJSSUserGroup{
				ID:   util.GetIntFromInterface(jssUserGroupMap["id"]),
				Name: util.GetStringFromInterface(jssUserGroupMap["name"]),
			}
			jssUserGroups = append(jssUserGroups, jssUserGroup)
		}
	}

	var iBeacons []jamfpro.PolicyDataSubsetIBeacon
	iBeaconsInterface, iBeaconsExist := exclusionsData["ibeacons"].([]interface{})
	if iBeaconsExist {
		for _, ibi := range iBeaconsInterface {
			iBeaconMap, ok := ibi.(map[string]interface{})
			if !ok {
				continue
			}

			iBeacon := jamfpro.PolicyDataSubsetIBeacon{
				ID:   util.GetIntFromInterface(iBeaconMap["id"]),
				Name: util.GetStringFromInterface(iBeaconMap["name"]),
			}
			iBeacons = append(iBeacons, iBeacon)
		}
	}

	return jamfpro.PolicySubsetScopeExclusions{
		Computers:       computers,
		ComputerGroups:  computerGroups,
		Users:           users,
		UserGroups:      userGroups,
		Buildings:       buildings,
		Departments:     departments,
		NetworkSegments: networkSegments,
		JSSUsers:        jssUsers,
		JSSUserGroups:   jssUserGroups,
		IBeacons:        iBeacons,
	}
}

// Self service

// constructSelfService creates a PolicySubsetSelfService instance from the provided data.
func constructSelfService(data *schema.ResourceData) jamfpro.PolicySubsetSelfService {
	useForSelfService := util.GetBoolFromInterface(data.Get("use_for_self_service"))
	selfServiceDisplayName := util.GetStringFromInterface(data.Get("self_service_display_name"))
	installButtonText := util.GetStringFromInterface(data.Get("install_button_text"))
	reinstallButtonText := util.GetStringFromInterface(data.Get("re_install_button_text"))
	selfServiceDescription := util.GetStringFromInterface(data.Get("self_service_description"))
	forceUsersToViewDescription := util.GetBoolFromInterface(data.Get("force_users_to_view_description"))

	// Initialize selfServiceIcon as an empty struct
	var selfServiceIcon jamfpro.SharedResourceSelfServiceIcon

	// Construct SelfServiceIcon
	if selfServiceIconInterface, ok := data.Get("self_service_icon").([]interface{}); ok && len(selfServiceIconInterface) > 0 {
		if selfServiceIconData, ok := selfServiceIconInterface[0].(map[string]interface{}); ok {
			selfServiceIcon = jamfpro.SharedResourceSelfServiceIcon{
				ID:       util.GetIntFromInterface(selfServiceIconData["id"]),
				URI:      util.GetStringFromInterface(selfServiceIconData["uri"]),
				Data:     util.GetStringFromInterface(selfServiceIconData["data"]),
				Filename: util.GetStringFromInterface(selfServiceIconData["filename"]),
			}
		}
	}

	featureOnMainPage := util.GetBoolFromInterface(data.Get("feature_on_main_page"))

	// Initialize selfServiceCategories as an empty slice
	var selfServiceCategories []jamfpro.PolicySubsetSelfServiceCategories

	// Construct SelfServiceCategories
	if selfServiceCategoriesData, ok := data.Get("self_service_categories").([]interface{}); ok {
		for _, categoryDataInterface := range selfServiceCategoriesData {
			// Safely assert type of each category data
			if categoryData, ok := categoryDataInterface.(map[string]interface{}); ok {
				category := jamfpro.PolicyCategory{
					Name: util.GetStringFromInterface(categoryData["name"]),
				}
				selfServiceCategories = append(selfServiceCategories, jamfpro.PolicySubsetSelfServiceCategories{Category: category})
			}
		}
	}

	// Construct PolicySubsetSelfService
	selfService := jamfpro.PolicySubsetSelfService{
		UseForSelfService:           useForSelfService,
		SelfServiceDisplayName:      selfServiceDisplayName,
		InstallButtonText:           installButtonText,
		ReinstallButtonText:         reinstallButtonText,
		SelfServiceDescription:      selfServiceDescription,
		ForceUsersToViewDescription: forceUsersToViewDescription,
		SelfServiceIcon:             selfServiceIcon,
		FeatureOnMainPage:           featureOnMainPage,
		SelfServiceCategories:       selfServiceCategories,
	}

	return selfService
}

// Package Configuration

// constructPackageConfiguration creates a PolicySubsetPackageConfiguration instance from the provided data.
func constructPackageConfiguration(data *schema.ResourceData) jamfpro.PolicySubsetPackageConfiguration {
	// Initialize packages as an empty slice
	var packages []jamfpro.PolicySubsetPackageConfigurationPackage

	// Check if packages data exists and is of the correct type
	if packagesData, ok := data.Get("packages").([]interface{}); ok {
		for _, packageDataInterface := range packagesData {
			// Safely assert type of each package data
			if packageItem, ok := packageDataInterface.(map[string]interface{}); ok {
				pkg := jamfpro.PolicySubsetPackageConfigurationPackage{
					ID:                util.GetIntFromInterface(packageItem["id"]),
					Name:              util.GetStringFromInterface(packageItem["name"]),
					Action:            util.GetStringFromInterface(packageItem["action"]),
					FillUserTemplate:  util.GetBoolFromInterface(packageItem["fut"]),
					FillExistingUsers: util.GetBoolFromInterface(packageItem["feu"]),
					UpdateAutorun:     util.GetBoolFromInterface(packageItem["update_autorun"]),
				}
				packages = append(packages, pkg)
			}
		}
	}

	distributionPoint := util.GetStringFromInterface(data.Get("distribution_point"))

	// Construct PolicySubsetPackageConfiguration
	packageConfiguration := jamfpro.PolicySubsetPackageConfiguration{
		Packages:          packages,
		DistributionPoint: distributionPoint,
	}

	return packageConfiguration
}

// Scripts

// constructScripts creates a PolicySubsetScripts instance from the provided data.
func constructScripts(data *schema.ResourceData) jamfpro.PolicySubsetScripts {
	size := util.GetIntFromInterface(data.Get("size"))

	// Initialize scripts as an empty slice
	var scripts []jamfpro.PolicySubsetScript

	// Check if script data exists and is of the correct type
	if scriptData, ok := data.Get("script").([]interface{}); ok {
		for _, scriptItemInterface := range scriptData {
			// Safely assert type of each script data
			if scriptItem, ok := scriptItemInterface.(map[string]interface{}); ok {
				scriptObj := jamfpro.PolicySubsetScript{
					ID:          util.GetStringFromInterface(scriptItem["id"]),
					Name:        util.GetStringFromInterface(scriptItem["name"]),
					Priority:    util.GetStringFromInterface(scriptItem["priority"]),
					Parameter4:  util.GetStringFromInterface(scriptItem["parameter4"]),
					Parameter5:  util.GetStringFromInterface(scriptItem["parameter5"]),
					Parameter6:  util.GetStringFromInterface(scriptItem["parameter6"]),
					Parameter7:  util.GetStringFromInterface(scriptItem["parameter7"]),
					Parameter8:  util.GetStringFromInterface(scriptItem["parameter8"]),
					Parameter9:  util.GetStringFromInterface(scriptItem["parameter9"]),
					Parameter10: util.GetStringFromInterface(scriptItem["parameter10"]),
					Parameter11: util.GetStringFromInterface(scriptItem["parameter11"]),
				}
				scripts = append(scripts, scriptObj)
			}
		}
	}

	// Construct PolicySubsetScripts
	scriptsConfig := jamfpro.PolicySubsetScripts{
		Size:   size,
		Script: scripts,
	}

	return scriptsConfig
}

// Printers

// constructDockItems creates a PolicySubsetDockItems instance from the provided data.
func constructDockItems(data *schema.ResourceData) jamfpro.PolicySubsetDockItems {
	size := util.GetIntFromInterface(data.Get("size"))

	// Initialize dockItems as an empty slice
	var dockItems []jamfpro.PolicySubsetDockItem

	// Check if dockItemData exists and is of the correct type
	if dockItemData, ok := data.Get("dock_item").([]interface{}); ok {
		for _, dockItemInterface := range dockItemData {
			// Safely assert type of each dock item data
			if dockItem, ok := dockItemInterface.(map[string]interface{}); ok {
				dockItemObj := jamfpro.PolicySubsetDockItem{
					ID:     util.GetIntFromInterface(dockItem["id"]),
					Name:   util.GetStringFromInterface(dockItem["name"]),
					Action: util.GetStringFromInterface(dockItem["action"]),
				}
				dockItems = append(dockItems, dockItemObj)
			}
		}
	}

	// Construct PolicySubsetDockItems
	dockItemsConfig := jamfpro.PolicySubsetDockItems{
		Size:     size,
		DockItem: dockItems,
	}

	return dockItemsConfig
}

// Printers

// constructPrinters creates a PolicySubsetPrinters instance from the provided data.
func constructPrinters(data *schema.ResourceData) jamfpro.PolicySubsetPrinters {
	size := util.GetIntFromInterface(data.Get("size"))

	// Initialize printers as an empty slice
	var printers []jamfpro.PolicySubsetPrinter

	// Check if printerData exists and is of the correct type
	if printerData, ok := data.Get("printer").([]interface{}); ok {
		for _, printerInterface := range printerData {
			// Safely assert type of each printer item
			if printer, ok := printerInterface.(map[string]interface{}); ok {
				printerObj := jamfpro.PolicySubsetPrinter{
					ID:          util.GetIntFromInterface(printer["id"]),
					Name:        util.GetStringFromInterface(printer["name"]),
					Action:      util.GetStringFromInterface(printer["action"]),
					MakeDefault: util.GetBoolFromInterface(printer["make_default"]),
				}
				printers = append(printers, printerObj)
			}
		}
	}

	// Construct PolicySubsetPrinters
	printersConfig := jamfpro.PolicySubsetPrinters{
		Size:                 size,
		LeaveExistingDefault: util.GetBoolFromInterface(data.Get("leave_existing_default")),
		Printer:              printers,
	}

	return printersConfig
}

// Account Maintainence

// constructAccountMaintenance creates a PolicySubsetAccountMaintenance instance from the provided data.
func constructAccountMaintenance(data *schema.ResourceData) jamfpro.PolicySubsetAccountMaintenance {
	// Initialize accounts and directory bindings as an empty slice
	var accounts []jamfpro.PolicySubsetAccountMaintenanceAccount
	var directoryBindings []jamfpro.PolicySubsetAccountMaintenanceDirectoryBindings

	if accountData, ok := data.Get("accounts").([]interface{}); ok {
		for _, accountInterface := range accountData {
			// Safely assert type of each account item
			if account, ok := accountInterface.(map[string]interface{}); ok {
				accountObj := jamfpro.PolicySubsetAccountMaintenanceAccount{
					Action:                 util.GetStringFromInterface(account["action"]),
					Username:               util.GetStringFromInterface(account["username"]),
					Realname:               util.GetStringFromInterface(account["realname"]),
					Password:               util.GetStringFromInterface(account["password"]),
					ArchiveHomeDirectory:   util.GetBoolFromInterface(account["archive_home_directory"]),
					ArchiveHomeDirectoryTo: util.GetStringFromInterface(account["archive_home_directory_to"]),
					Home:                   util.GetStringFromInterface(account["home"]),
					Hint:                   util.GetStringFromInterface(account["hint"]),
					Picture:                util.GetStringFromInterface(account["picture"]),
					Admin:                  util.GetBoolFromInterface(account["admin"]),
					FilevaultEnabled:       util.GetBoolFromInterface(account["filevault_enabled"]),
					PasswordSha256:         util.GetStringFromInterface(account["password_sha256"]),
				}
				accounts = append(accounts, accountObj)
			}
		}
	}
	// Safely assert type of each directory binding item
	if directoryBindingsData, ok := data.Get("directory_bindings").([]interface{}); ok {
		for _, bindingInterface := range directoryBindingsData {
			if binding, ok := bindingInterface.(map[string]interface{}); ok {
				bindingObj := jamfpro.PolicySubsetAccountMaintenanceDirectoryBindings{
					ID:   util.GetIntFromInterface(binding["id"]),
					Name: util.GetStringFromInterface(binding["name"]),
				}
				directoryBindings = append(directoryBindings, bindingObj)
			}
		}
	}

	managementAccountData := util.ConvertToMapFromInterface(data.Get("management_account"))
	managementAccountObj := jamfpro.PolicySubsetAccountMaintenanceManagementAccount{
		Action:                util.GetStringFromInterface(managementAccountData["action"]),
		ManagedPassword:       util.GetStringFromInterface(managementAccountData["managed_password"]),
		ManagedPasswordLength: util.GetIntFromInterface(managementAccountData["managed_password_length"]),
	}

	openFirmwareEfiPasswordData := util.ConvertToMapFromInterface(data.Get("open_firmware_efi_password"))
	openFirmwareEfiPasswordObj := jamfpro.PolicySubsetAccountMaintenanceOpenFirmwareEfiPassword{
		OfMode:           util.GetStringFromInterface(openFirmwareEfiPasswordData["of_mode"]),
		OfPassword:       util.GetStringFromInterface(openFirmwareEfiPasswordData["of_password"]),
		OfPasswordSHA256: util.GetStringFromInterface(openFirmwareEfiPasswordData["of_password_sha256"]),
	}

	// Construct PolicySubsetAccountMaintenance
	accountMaintenanceConfig := jamfpro.PolicySubsetAccountMaintenance{
		Accounts:                accounts,
		DirectoryBindings:       directoryBindings,
		ManagementAccount:       managementAccountObj,
		OpenFirmwareEfiPassword: openFirmwareEfiPasswordObj,
	}

	return accountMaintenanceConfig
}

// Maintainence

// constructMaintenance creates a PolicySubsetMaintenance instance from the provided data.
func constructMaintenance(data *schema.ResourceData) jamfpro.PolicySubsetMaintenance {
	maintenanceConfig := jamfpro.PolicySubsetMaintenance{
		Recon:                    util.GetBoolFromInterface(data.Get("recon")),
		ResetName:                util.GetBoolFromInterface(data.Get("reset_name")),
		InstallAllCachedPackages: util.GetBoolFromInterface(data.Get("install_all_cached_packages")),
		Heal:                     util.GetBoolFromInterface(data.Get("heal")),
		Prebindings:              util.GetBoolFromInterface(data.Get("prebindings")),
		Permissions:              util.GetBoolFromInterface(data.Get("permissions")),
		Byhost:                   util.GetBoolFromInterface(data.Get("byhost")),
		SystemCache:              util.GetBoolFromInterface(data.Get("system_cache")),
		UserCache:                util.GetBoolFromInterface(data.Get("user_cache")),
		Verify:                   util.GetBoolFromInterface(data.Get("verify")),
	}
	return maintenanceConfig
}

// Files Processes

// constructFilesProcesses creates a PolicySubsetFilesProcesses instance from the provided data.
func constructFilesProcesses(data *schema.ResourceData) jamfpro.PolicySubsetFilesProcesses {
	filesProcessesConfig := jamfpro.PolicySubsetFilesProcesses{
		SearchByPath:         util.GetStringFromInterface(data.Get("search_by_path")),
		DeleteFile:           util.GetBoolFromInterface(data.Get("delete_file")),
		LocateFile:           util.GetStringFromInterface(data.Get("locate_file")),
		UpdateLocateDatabase: util.GetBoolFromInterface(data.Get("update_locate_database")),
		SpotlightSearch:      util.GetStringFromInterface(data.Get("spotlight_search")),
		SearchForProcess:     util.GetStringFromInterface(data.Get("search_for_process")),
		KillProcess:          util.GetBoolFromInterface(data.Get("kill_process")),
		RunCommand:           util.GetStringFromInterface(data.Get("run_command")),
	}
	return filesProcessesConfig
}

// User Interaction

// constructUserInteraction creates a PolicySubsetUserInteraction instance from the provided data.
func constructUserInteraction(data *schema.ResourceData) jamfpro.PolicySubsetUserInteraction {
	userInteractionConfig := jamfpro.PolicySubsetUserInteraction{
		MessageStart:          util.GetStringFromInterface(data.Get("message_start")),
		AllowUserToDefer:      util.GetBoolFromInterface(data.Get("allow_user_to_defer")),
		AllowDeferralUntilUtc: util.GetStringFromInterface(data.Get("allow_deferral_until_utc")),
		AllowDeferralMinutes:  util.GetIntFromInterface(data.Get("allow_deferral_minutes")),
		MessageFinish:         util.GetStringFromInterface(data.Get("message_finish")),
	}
	return userInteractionConfig
}

// Disk Encryption

// constructDiskEncryption creates a PolicySubsetDiskEncryption instance from the provided data.
func constructDiskEncryption(data *schema.ResourceData) jamfpro.PolicySubsetDiskEncryption {
	diskEncryptionConfig := jamfpro.PolicySubsetDiskEncryption{
		Action:                                 util.GetStringFromInterface(data.Get("action")),
		DiskEncryptionConfigurationID:          util.GetIntFromInterface(data.Get("disk_encryption_configuration_id")),
		AuthRestart:                            util.GetBoolFromInterface(data.Get("auth_restart")),
		RemediateKeyType:                       util.GetStringFromInterface(data.Get("remediate_key_type")),
		RemediateDiskEncryptionConfigurationID: util.GetIntFromInterface(data.Get("remediate_disk_encryption_configuration_id")),
	}
	return diskEncryptionConfig
}

// Reboot

// constructReboot creates a PolicySubsetReboot instance from the provided data.
func constructReboot(data *schema.ResourceData) jamfpro.PolicySubsetReboot {
	rebootConfig := jamfpro.PolicySubsetReboot{
		Message:                     util.GetStringFromInterface(data.Get("message")),
		StartupDisk:                 util.GetStringFromInterface(data.Get("startup_disk")),
		SpecifyStartup:              util.GetStringFromInterface(data.Get("specify_startup")),
		NoUserLoggedIn:              util.GetStringFromInterface(data.Get("no_user_logged_in")),
		UserLoggedIn:                util.GetStringFromInterface(data.Get("user_logged_in")),
		MinutesUntilReboot:          util.GetIntFromInterface(data.Get("minutes_until_reboot")),
		StartRebootTimerImmediately: util.GetBoolFromInterface(data.Get("start_reboot_timer_immediately")),
		FileVault2Reboot:            util.GetBoolFromInterface(data.Get("file_vault_2_reboot")),
	}
	return rebootConfig
}

/*
// constructJamfProPolicy constructs a ResourcePolicy object from the provided schema data.
func constructJamfProPolicy(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {
	// Initialize a new ResourcePolicy struct with all its sub-components.
	policy := &jamfpro.ResourcePolicy{
		General:              jamfpro.PolicySubsetGeneral{},
		Scope:                jamfpro.PolicySubsetScope{},
		SelfService:          jamfpro.PolicySubsetSelfService{},
		PackageConfiguration: jamfpro.PolicySubsetPackageConfiguration{},
		Scripts:              jamfpro.PolicySubsetScripts{},
		Printers:             jamfpro.PolicySubsetPrinters{},
		DockItems:            jamfpro.PolicySubsetDockItems{},
		AccountMaintenance:   jamfpro.PolicySubsetAccountMaintenance{},
		Maintenance:          jamfpro.PolicySubsetMaintenance{},
		FilesProcesses:       jamfpro.PolicySubsetFilesProcesses{},
		UserInteraction:      jamfpro.PolicySubsetUserInteraction{},
		DiskEncryption:       jamfpro.PolicySubsetDiskEncryption{},
		Reboot:               jamfpro.PolicySubsetReboot{},
	}

	// Construct the General section
	if v, ok := d.GetOk("general"); ok {
		generalData := v.([]interface{})[0].(map[string]interface{})
		policy.General = jamfpro.PolicySubsetGeneral{
			Name:                       util.GetStringFromMap(generalData, "name"),
			Enabled:                    util.GetBoolFromMap(generalData, "enabled"),
			Trigger:                    util.GetStringFromMap(generalData, "trigger"),
			TriggerCheckin:             util.GetBoolFromMap(generalData, "trigger_checkin"),
			TriggerEnrollmentComplete:  util.GetBoolFromMap(generalData, "trigger_enrollment_complete"),
			TriggerLogin:               util.GetBoolFromMap(generalData, "trigger_login"),
			TriggerLogout:              util.GetBoolFromMap(generalData, "trigger_logout"),
			TriggerNetworkStateChanged: util.GetBoolFromMap(generalData, "trigger_network_state_changed"),
			TriggerStartup:             util.GetBoolFromMap(generalData, "trigger_startup"),
			TriggerOther:               util.GetStringFromMap(generalData, "trigger_other"),
			Frequency:                  util.GetStringFromMap(generalData, "frequency"),
			RetryEvent:                 util.GetStringFromMap(generalData, "retry_event"),
			RetryAttempts:              util.GetIntFromMap(generalData, "retry_attempts"),
			NotifyOnEachFailedRetry:    util.GetBoolFromMap(generalData, "notify_on_each_failed_retry"),
			LocationUserOnly:           util.GetBoolFromMap(generalData, "location_user_only"),
			TargetDrive:                util.GetStringFromMap(generalData, "target_drive"),
			Offline:                    util.GetBoolFromMap(generalData, "offline"),
			Category: func() jamfpro.PolicyCategory {
				var category jamfpro.PolicyCategory

				if categoryData, ok := generalData["category"].([]interface{}); ok && len(categoryData) > 0 {
					catMap := categoryData[0].(map[string]interface{})
					category.ID = util.GetIntFromMap(catMap, "id")
					category.Name = util.GetStringFromMap(catMap, "name")
				}

				return category
			}(),
			// DateTimeLimitations field
			DateTimeLimitations: func() jamfpro.PolicySubsetGeneralDateTimeLimitations {
				if dtData, ok := generalData["date_time_limitations"].([]interface{}); ok && len(dtData) > 0 {
					dateTimeMap := dtData[0].(map[string]interface{})
					dateTimeLimitations := jamfpro.PolicySubsetGeneralDateTimeLimitations{
						ActivationDate:      util.GetStringFromMap(dateTimeMap, "activation_date"),
						ActivationDateEpoch: util.GetIntFromMap(dateTimeMap, "activation_date_epoch"),
						ActivationDateUTC:   util.GetStringFromMap(dateTimeMap, "activation_date_utc"),
						ExpirationDate:      util.GetStringFromMap(dateTimeMap, "expiration_date"),
						ExpirationDateEpoch: util.GetIntFromMap(dateTimeMap, "expiration_date_epoch"),
						ExpirationDateUTC:   util.GetStringFromMap(dateTimeMap, "expiration_date_utc"),
					}

					// Handling NoExecuteOn field
					if noExecOnData, ok := dateTimeMap["no_execute_on"].([]interface{}); ok {
						var noExecOnDays []jamfpro.PolicySubsetGeneralDateTimeLimitationsNoExecuteOn
						for _, day := range noExecOnData {
							if dayMap, ok := day.(map[string]interface{}); ok {
								noExecOnDays = append(noExecOnDays, jamfpro.PolicySubsetGeneralDateTimeLimitationsNoExecuteOn{
									Day: util.GetStringFromMap(dayMap, "day"),
								})
							}
						}
						dateTimeLimitations.NoExecuteOn = noExecOnDays
					}

					// Handling NoExecuteStart and NoExecuteEnd fields
					dateTimeLimitations.NoExecuteStart = util.GetStringFromMap(dateTimeMap, "no_execute_start")
					dateTimeLimitations.NoExecuteEnd = util.GetStringFromMap(dateTimeMap, "no_execute_end")

					return dateTimeLimitations
				}
				return jamfpro.PolicySubsetGeneralDateTimeLimitations{}
			}(),
			// NetworkLimitations field
			NetworkLimitations: func() jamfpro.PolicySubsetGeneralNetworkLimitations {
				var networkLimitations jamfpro.PolicySubsetGeneralNetworkLimitations

				if networkLimitationsData, ok := generalData["network_limitations"].([]interface{}); ok && len(networkLimitationsData) > 0 {
					netMap := networkLimitationsData[0].(map[string]interface{})

					networkLimitations.MinimumNetworkConnection = util.GetStringFromMap(netMap, "minimum_network_connection")
					networkLimitations.AnyIPAddress = util.GetBoolFromMap(netMap, "any_ip_address")
					networkLimitations.NetworkSegments = util.GetStringFromMap(netMap, "network_segments")
				}

				return networkLimitations
			}(),
			// OverrideDefaultSettings field
			OverrideDefaultSettings: func() jamfpro.PolicySubsetGeneralOverrideDefaultSettings {
				if overrideData, ok := generalData["override_default_settings"].([]interface{}); ok && len(overrideData) > 0 {
					overrideMap := overrideData[0].(map[string]interface{})
					return jamfpro.PolicySubsetGeneralOverrideDefaultSettings{
						TargetDrive:       util.GetStringFromMap(overrideMap, "target_drive"),
						DistributionPoint: util.GetStringFromMap(overrideMap, "distribution_point"),
						ForceAfpSmb:       util.GetBoolFromMap(overrideMap, "force_afp_smb"),
						SUS:               util.GetStringFromMap(overrideMap, "sus"),
						NetbootServer:     util.GetStringFromMap(overrideMap, "netboot_server"),
					}
				}
				return jamfpro.PolicySubsetGeneralOverrideDefaultSettings{}
			}(),
			// NetworkRequirements field
			NetworkRequirements: func() string {
				return util.GetStringFromMap(generalData, "network_requirements")
			}(),
			// Construct the Site fields
			Site: func() jamfpro.SharedResourceSite {
				var site jamfpro.SharedResourceSite

				// Check if values are provided in Terraform
				if siteData, ok := generalData["site"].([]interface{}); ok && len(siteData) > 0 {
					siteMap := siteData[0].(map[string]interface{})

					// Extract values directly from the Terraform data
					site.ID = util.GetIntFromMap(siteMap, "id")
					site.Name = util.GetStringFromMap(siteMap, "name")
				}

				return site
			}(),
		}
	}

	// Construct the Scope section
	if v, ok := d.GetOk("scope"); ok {
		scopeData := v.([]interface{})[0].(map[string]interface{})
		var computers []jamfpro.PolicyDataSubsetComputer
		var computerGroups []jamfpro.PolicyDataSubsetComputerGroup
		var jssUsers []jamfpro.PolicyDataSubsetJSSUser
		var jssUserGroups []jamfpro.PolicyDataSubsetJSSUserGroup
		var buildings []jamfpro.PolicyDataSubsetBuilding
		var departments []jamfpro.PolicyDataSubsetDepartment

		// Construct Computers slice
		if comps, ok := scopeData["computers"].([]interface{}); ok {
			for _, comp := range comps {
				compMap := comp.(map[string]interface{})
				computers = append(computers, jamfpro.PolicyDataSubsetComputer{
					ID:   compMap["id"].(int),
					Name: compMap["name"].(string),
					UDID: compMap["udid"].(string),
				})
			}
		}

		// Construct ComputerGroups slice
		if groups, ok := scopeData["computer_groups"].([]interface{}); ok {
			for _, group := range groups {
				groupMap := group.(map[string]interface{})
				computerGroups = append(computerGroups, jamfpro.PolicyDataSubsetComputerGroup{
					ID:   groupMap["id"].(int),
					Name: groupMap["name"].(string),
				})
			}
		}

		// Construct JSSUsers slice
		if users, ok := scopeData["jss_users"].([]interface{}); ok {
			for _, user := range users {
				userMap := user.(map[string]interface{})
				jssUsers = append(jssUsers, jamfpro.PolicyDataSubsetJSSUser{
					ID:   userMap["id"].(int),
					Name: userMap["name"].(string),
				})
			}
		}

		// Construct JSSUserGroups slice
		if groups, ok := scopeData["jss_user_groups"].([]interface{}); ok {
			for _, group := range groups {
				groupMap := group.(map[string]interface{})
				jssUserGroups = append(jssUserGroups, jamfpro.PolicyDataSubsetJSSUserGroup{
					ID:   groupMap["id"].(int),
					Name: groupMap["name"].(string),
				})
			}
		}

		// Construct Buildings slice
		if bldgs, ok := scopeData["buildings"].([]interface{}); ok {
			for _, bldg := range bldgs {
				bldgMap := bldg.(map[string]interface{})
				buildings = append(buildings, jamfpro.PolicyDataSubsetBuilding{
					ID:   bldgMap["id"].(int),
					Name: bldgMap["name"].(string),
				})
			}
		}

		// Construct Departments slice
		if depts, ok := scopeData["departments"].([]interface{}); ok {
			for _, dept := range depts {
				deptMap := dept.(map[string]interface{})
				departments = append(departments, jamfpro.PolicyDataSubsetDepartment{
					ID:   deptMap["id"].(int),
					Name: deptMap["name"].(string),
				})
			}
		}

		// Construct LimitToUsers field
		var limitToUsers jamfpro.PolicyLimitToUsers
		if luData, ok := scopeData["limit_to_users"].([]interface{}); ok && len(luData) > 0 {
			luMap := luData[0].(map[string]interface{})
			var userGroups []string
			if uGroups, ok := luMap["user_groups"].([]interface{}); ok {
				for _, uGroup := range uGroups {
					userGroups = append(userGroups, uGroup.(string))
				}
			}
			limitToUsers = jamfpro.PolicyLimitToUsers{UserGroups: userGroups}
		}

		// Construct Limitations field
		var limitations jamfpro.PolicySubsetScopeLimitations
		if limitationsData, ok := scopeData["limitations"].([]interface{}); ok && len(limitationsData) > 0 {
			limitationsMap := limitationsData[0].(map[string]interface{})

			// Construct Directory Service/Local Users slice
			var users []jamfpro.PolicyDataSubsetUser
			if directoryServicesUsersData, ok := limitationsMap["users"].([]interface{}); ok {
				for _, user := range directoryServicesUsersData {
					userMap := user.(map[string]interface{})
					users = append(users, jamfpro.PolicyDataSubsetUser{
						ID:   userMap["id"].(int),
						Name: userMap["name"].(string),
					})
				}
			}

			// Construct Directory Service User Groups slice
			var userGroups []jamfpro.PolicyDataSubsetUserGroup
			if userGroupsData, ok := limitationsMap["user_groups"].([]interface{}); ok {
				for _, group := range userGroupsData {
					groupMap := group.(map[string]interface{})
					userGroups = append(userGroups, jamfpro.PolicyDataSubsetUserGroup{
						ID:   groupMap["id"].(int),
						Name: groupMap["name"].(string),
					})
				}
			}

			// Construct NetworkSegments slice
			var networkSegments []jamfpro.PolicyDataSubsetNetworkSegment
			if netSegsData, ok := limitationsMap["network_segments"].([]interface{}); ok {
				for _, seg := range netSegsData {
					segMap := seg.(map[string]interface{})
					networkSegments = append(networkSegments, jamfpro.PolicyDataSubsetNetworkSegment{
						ID:   segMap["id"].(int),
						Name: segMap["name"].(string),
						UID:  segMap["uid"].(string),
					})
				}
			}

			// Construct iBeacons slice
			var iBeacons []jamfpro.PolicyDataSubsetIBeacon
			if beaconsData, ok := limitationsMap["ibeacons"].([]interface{}); ok {
				for _, beacon := range beaconsData {
					beaconMap := beacon.(map[string]interface{})
					iBeacons = append(iBeacons, jamfpro.PolicyDataSubsetIBeacon{
						ID:   beaconMap["id"].(int),
						Name: beaconMap["name"].(string),
					})
				}
			}

			// Assign constructed slices to limitations struct
			limitations = jamfpro.PolicySubsetScopeLimitations{
				Users:           users,
				UserGroups:      userGroups,
				NetworkSegments: networkSegments,
				IBeacons:        iBeacons,
			}
		}

		// Assign Limitations to policy's Scope
		policy.Scope.Limitations = limitations

		// Construct Exclusions field
		var exclusions jamfpro.PolicySubsetScopeExclusions
		if exclusionsData, ok := scopeData["exclusions"].([]interface{}); ok && len(exclusionsData) > 0 {
			exclusionsMap := exclusionsData[0].(map[string]interface{})

			// Construct exclusion Computers slice
			var exclusionComputers []jamfpro.PolicyDataSubsetComputer
			if comps, ok := exclusionsMap["computers"].([]interface{}); ok {
				for _, comp := range comps {
					compMap := comp.(map[string]interface{})
					exclusionComputers = append(exclusionComputers, jamfpro.PolicyDataSubsetComputer{
						ID:   compMap["id"].(int),
						Name: compMap["name"].(string),
						UDID: compMap["udid"].(string),
					})
				}
			}

			// Construct exclusion ComputerGroups slice
			var exclusionComputerGroups []jamfpro.PolicyDataSubsetComputerGroup
			if groups, ok := exclusionsMap["computer_groups"].([]interface{}); ok {
				for _, group := range groups {
					groupMap := group.(map[string]interface{})
					exclusionComputerGroups = append(exclusionComputerGroups, jamfpro.PolicyDataSubsetComputerGroup{
						ID:   groupMap["id"].(int),
						Name: groupMap["name"].(string),
					})
				}
			}

			// Construct exclusion Users slice
			var exclusionUsers []jamfpro.PolicyDataSubsetUser
			if users, ok := exclusionsMap["users"].([]interface{}); ok {
				for _, user := range users {
					userMap := user.(map[string]interface{})
					exclusionUsers = append(exclusionUsers, jamfpro.PolicyDataSubsetUser{
						ID:   userMap["id"].(int),
						Name: userMap["name"].(string),
					})
				}
			}

			// Construct exclusion UserGroups slice
			var exclusionUserGroups []jamfpro.PolicyDataSubsetUserGroup
			if userGroups, ok := exclusionsMap["user_groups"].([]interface{}); ok {
				for _, group := range userGroups {
					groupMap := group.(map[string]interface{})
					exclusionUserGroups = append(exclusionUserGroups, jamfpro.PolicyDataSubsetUserGroup{
						ID:   groupMap["id"].(int),
						Name: groupMap["name"].(string),
					})
				}
			}

			// Construct exclusion Buildings slice
			var exclusionBuildings []jamfpro.PolicyDataSubsetBuilding
			if bldgs, ok := exclusionsMap["buildings"].([]interface{}); ok {
				for _, bldg := range bldgs {
					bldgMap := bldg.(map[string]interface{})
					exclusionBuildings = append(exclusionBuildings, jamfpro.PolicyDataSubsetBuilding{
						ID:   bldgMap["id"].(int),
						Name: bldgMap["name"].(string),
					})
				}
			}

			// Construct exclusion Departments slice
			var exclusionDepartments []jamfpro.PolicyDataSubsetDepartment
			if depts, ok := exclusionsMap["departments"].([]interface{}); ok {
				for _, dept := range depts {
					deptMap := dept.(map[string]interface{})
					exclusionDepartments = append(exclusionDepartments, jamfpro.PolicyDataSubsetDepartment{
						ID:   deptMap["id"].(int),
						Name: deptMap["name"].(string),
					})
				}
			}

			// Construct exclusion NetworkSegments slice
			var exclusionNetworkSegments []jamfpro.PolicyDataSubsetNetworkSegment
			if netSegments, ok := exclusionsMap["network_segments"].([]interface{}); ok {
				for _, segment := range netSegments {
					segmentMap := segment.(map[string]interface{})
					exclusionNetworkSegments = append(exclusionNetworkSegments, jamfpro.PolicyDataSubsetNetworkSegment{
						ID:   segmentMap["id"].(int),
						Name: segmentMap["name"].(string),
						UID:  segmentMap["uid"].(string),
					})
				}
			}

			// Construct exclusion JSSUsers slice
			var exclusionJSSUsers []jamfpro.PolicyDataSubsetJSSUser
			if jssUsers, ok := exclusionsMap["jss_users"].([]interface{}); ok {
				for _, jssUser := range jssUsers {
					jssUserMap := jssUser.(map[string]interface{})
					exclusionJSSUsers = append(exclusionJSSUsers, jamfpro.PolicyDataSubsetJSSUser{
						ID:   jssUserMap["id"].(int),
						Name: jssUserMap["name"].(string),
					})
				}
			}

			// Construct exclusion JSSUserGroups slice
			var exclusionJSSUserGroups []jamfpro.PolicyDataSubsetJSSUserGroup
			if jssUserGroups, ok := exclusionsMap["jss_user_groups"].([]interface{}); ok {
				for _, jssUserGroup := range jssUserGroups {
					jssUserGroupMap := jssUserGroup.(map[string]interface{})
					exclusionJSSUserGroups = append(exclusionJSSUserGroups, jamfpro.PolicyDataSubsetJSSUserGroup{
						ID:   jssUserGroupMap["id"].(int),
						Name: jssUserGroupMap["name"].(string),
					})
				}
			}

			// Construct exclusion iBeacons slice
			var exclusionIBeacons []jamfpro.PolicyDataSubsetIBeacon
			if beacons, ok := exclusionsMap["ibeacons"].([]interface{}); ok {
				for _, beacon := range beacons {
					beaconMap := beacon.(map[string]interface{})
					exclusionIBeacons = append(exclusionIBeacons, jamfpro.PolicyDataSubsetIBeacon{
						ID:   beaconMap["id"].(int),
						Name: beaconMap["name"].(string),
					})
				}
			}

			// Assign constructed slices to exclusions struct
			exclusions = jamfpro.PolicySubsetScopeExclusions{
				Computers:       exclusionComputers,
				ComputerGroups:  exclusionComputerGroups,
				Users:           exclusionUsers,
				UserGroups:      exclusionUserGroups,
				Buildings:       exclusionBuildings,
				Departments:     exclusionDepartments,
				NetworkSegments: exclusionNetworkSegments,
				JSSUsers:        exclusionJSSUsers,
				JSSUserGroups:   exclusionJSSUserGroups,
				IBeacons:        exclusionIBeacons,
			}

		}

		// Assign Exclusions to policy's Scope
		policy.Scope.Exclusions = exclusions

		// Assign constructed fields to the policy's Scope
		policy.Scope = jamfpro.PolicySubsetScope{
			AllComputers:   util.GetBoolFromMap(scopeData, "all_computers"),
			Computers:      computers,
			ComputerGroups: computerGroups,
			JSSUsers:       jssUsers,
			JSSUserGroups:  jssUserGroups,
			Buildings:      buildings,
			Departments:    departments,
			LimitToUsers:   limitToUsers,
			Limitations:    limitations,
			Exclusions:     exclusions,
		}
	}

	// Construct the SelfService section
	if v, ok := d.GetOk("self_service"); ok {
		selfServiceData := v.([]interface{})[0].(map[string]interface{})
		policy.SelfService = jamfpro.PolicySubsetSelfService{
			UseForSelfService:           util.GetBoolFromMap(selfServiceData, "use_for_self_service"),
			SelfServiceDisplayName:      util.GetStringFromMap(selfServiceData, "self_service_display_name"),
			InstallButtonText:           util.GetStringFromMap(selfServiceData, "install_button_text"),
			ReinstallButtonText:         util.GetStringFromMap(selfServiceData, "reinstall_button_text"),
			SelfServiceDescription:      util.GetStringFromMap(selfServiceData, "self_service_description"),
			ForceUsersToViewDescription: util.GetBoolFromMap(selfServiceData, "force_users_to_view_description"),
			FeatureOnMainPage:           util.GetBoolFromMap(selfServiceData, "feature_on_main_page"),
			SelfServiceIcon: func() jamfpro.SharedResourceSelfServiceIcon {
				if iconData, ok := selfServiceData["self_service_icon"].([]interface{}); ok && len(iconData) > 0 {
					iconMap := iconData[0].(map[string]interface{})
					return jamfpro.SharedResourceSelfServiceIcon{
						ID:       util.GetIntFromMap(iconMap, "id"),
						Filename: util.GetStringFromMap(iconMap, "filename"),
						URI:      util.GetStringFromMap(iconMap, "uri"),
					}
				}
				return jamfpro.SharedResourceSelfServiceIcon{}
			}(),
			SelfServiceCategories: func() []jamfpro.PolicySubsetSelfServiceCategories {
				var categories []jamfpro.PolicySubsetSelfServiceCategories
				if catData, ok := selfServiceData["self_service_categories"].([]interface{}); ok {
					for _, cat := range catData {
						catMap := cat.(map[string]interface{})
						if catDetails, ok := catMap["category"].([]interface{}); ok && len(catDetails) > 0 {
							details := catDetails[0].(map[string]interface{})
							category := jamfpro.PolicySubsetSelfServiceCategories{
								Category: jamfpro.PolicyCategory{
									ID:        util.GetIntFromMap(details, "id"),
									Name:      util.GetStringFromMap(details, "name"),
									DisplayIn: util.GetBoolFromMap(details, "display_in"),
									FeatureIn: util.GetBoolFromMap(details, "feature_in"),
								},
							}
							categories = append(categories, category)
						}
					}
				}
				return categories
			}(),
		}
	}
	// Construct the PackageConfiguration section
	if v, ok := d.GetOk("package_configuration"); ok {
		packageConfigData := v.([]interface{})[0].(map[string]interface{})

		packageItems := func() []jamfpro.PolicySubsetPackageConfigurationPackage {
			var items []jamfpro.PolicySubsetPackageConfigurationPackage
			if pkgs, ok := packageConfigData["packages"].([]interface{}); ok {
				for _, pkg := range pkgs {
					pkgMap := util.ConvertToMapFromInterface(pkg)
					if pkgMap == nil {
						continue // Skip if package is nil
					}

					items = append(items, jamfpro.PolicySubsetPackageConfigurationPackage{
						ID:                util.GetIntFromMap(pkgMap, "id"),
						Name:              util.GetStringFromMap(pkgMap, "name"),
						Action:            util.GetStringFromMap(pkgMap, "action"),
						FillUserTemplate:  util.GetBoolFromMap(pkgMap, "fut"),
						FillExistingUsers: util.GetBoolFromMap(pkgMap, "feu"),
						UpdateAutorun:     util.GetBoolFromMap(pkgMap, "update_autorun"),
					})
				}
			}
			return items
		}()
		// Assign the constructed package items to the policy's package configuration
		policy.PackageConfiguration = jamfpro.PolicySubsetPackageConfiguration{
			Packages: packageItems,
		}
	}
	// Construct the Scripts section
	if v, ok := d.GetOk("scripts"); ok {
		scriptsData := v.([]interface{})[0].(map[string]interface{})

		scriptItems := func() []jamfpro.PolicySubsetScript {
			var items []jamfpro.PolicySubsetScript
			if scripts, ok := scriptsData["script"].([]interface{}); ok {
				for _, script := range scripts {
					scriptMap := util.ConvertToMapFromInterface(script)
					if scriptMap == nil {
						continue // Skip if script is nil or not a map
					}

					items = append(items, jamfpro.PolicySubsetScript{
						ID:          util.GetStringFromMap(scriptMap, "id"),
						Name:        util.GetStringFromMap(scriptMap, "name"),
						Priority:    util.GetStringFromMap(scriptMap, "priority"),
						Parameter4:  util.GetStringFromMap(scriptMap, "parameter4"),
						Parameter5:  util.GetStringFromMap(scriptMap, "parameter5"),
						Parameter6:  util.GetStringFromMap(scriptMap, "parameter6"),
						Parameter7:  util.GetStringFromMap(scriptMap, "parameter7"),
						Parameter8:  util.GetStringFromMap(scriptMap, "parameter8"),
						Parameter9:  util.GetStringFromMap(scriptMap, "parameter9"),
						Parameter10: util.GetStringFromMap(scriptMap, "parameter10"),
						Parameter11: util.GetStringFromMap(scriptMap, "parameter11"),
					})
				}
			}
			return items
		}()

		// Assign the constructed script items to the policy's scripts
		policy.Scripts = jamfpro.PolicySubsetScripts{
			Script: scriptItems,
		}
	}
	// Construct the Printers section
	if v, ok := d.GetOk("printers"); ok {
		printersData := v.([]interface{})[0].(map[string]interface{})

		var printerItems []jamfpro.PolicySubsetPrinter
		if printers, ok := printersData["printer"].([]interface{}); ok {
			for _, printer := range printers {
				printerMap := printer.(map[string]interface{})
				printerItems = append(printerItems, jamfpro.PolicySubsetPrinter{
					ID:          util.GetIntFromMap(printerMap, "id"),
					Name:        util.GetStringFromMap(printerMap, "name"),
					Action:      util.GetStringFromMap(printerMap, "action"),
					MakeDefault: util.GetBoolFromMap(printerMap, "make_default"),
				})
			}
		}

		leaveExistingDefault := false
		if val, ok := printersData["leave_existing_default"].(bool); ok {
			leaveExistingDefault = val
		}

		policy.Printers = jamfpro.PolicySubsetPrinters{
			LeaveExistingDefault: leaveExistingDefault,
			Printer:              printerItems,
		}
	}

	// Construct the DockItems section
	if v, ok := d.GetOk("dock_items"); ok {
		dockItemsData := v.([]interface{})[0].(map[string]interface{})

		dockItems := func() []jamfpro.PolicySubsetDockItem {
			var items []jamfpro.PolicySubsetDockItem
			if docks, ok := dockItemsData["dock_item"].([]interface{}); ok {
				for _, dock := range docks {
					dockMap := dock.(map[string]interface{})

					items = append(items, jamfpro.PolicySubsetDockItem{
						ID:     util.GetIntFromMap(dockMap, "id"),
						Name:   util.GetStringFromMap(dockMap, "name"),
						Action: util.GetStringFromMap(dockMap, "action"),
					})
				}
			}
			return items
		}()

		policy.DockItems = jamfpro.PolicySubsetDockItems{
			DockItem: dockItems,
		}
	}

	// Construct the AccountMaintenance section
	if v, ok := d.GetOk("account_maintenance"); ok {
		accountMaintenanceData := v.([]interface{})[0].(map[string]interface{})

		// Construct user accounts
		accounts := func() []jamfpro.PolicySubsetAccountMaintenanceAccount {
			var policyAccountItems []jamfpro.PolicySubsetAccountMaintenanceAccount
			if accs, ok := accountMaintenanceData["accounts"].([]interface{}); ok {
				for _, acc := range accs {
					accMap, ok := acc.(map[string]interface{})
					if !ok || accMap == nil {
						continue // Skip if not a map or if nil
					}

					var account jamfpro.PolicySubsetAccountMaintenanceAccount
					account.Action = util.GetStringFromMap(accMap, "action")
					account.Username = util.GetStringFromMap(accMap, "username")
					account.Realname = util.GetStringFromMap(accMap, "realname")
					account.Password = util.GetStringFromMap(accMap, "password")
					account.ArchiveHomeDirectory = util.GetBoolFromMap(accMap, "archive_home_directory")
					account.ArchiveHomeDirectoryTo = util.GetStringFromMap(accMap, "archive_home_directory_to")
					account.Home = util.GetStringFromMap(accMap, "home")
					account.Hint = util.GetStringFromMap(accMap, "hint")
					account.Picture = util.GetStringFromMap(accMap, "picture")
					account.Admin = util.GetBoolFromMap(accMap, "admin")
					account.FilevaultEnabled = util.GetBoolFromMap(accMap, "filevault_enabled")

					policyAccountItems = append(policyAccountItems, account)
				}
			}
			return policyAccountItems
		}()

		directoryBindings := func() []jamfpro.PolicySubsetAccountMaintenanceDirectoryBindings {
			var directoryBindings []jamfpro.PolicySubsetAccountMaintenanceDirectoryBindings
			if bindingsList, ok := accountMaintenanceData["directory_bindings"].([]interface{}); ok && len(bindingsList) > 0 {
				for _, bindingEntry := range bindingsList {
					bindingData := util.ConvertToMapFromInterface(bindingEntry)
					if bindingData == nil {
						continue // Skip if the map is nil
					}
					if bindings, ok := bindingData["binding"].([]interface{}); ok {
						for _, binding := range bindings {
							bindingMap := util.ConvertToMapFromInterface(binding)
							if bindingMap == nil {
								continue // Skip if the binding map is nil
							}
							directoryBindings = append(directoryBindings, jamfpro.PolicySubsetAccountMaintenanceDirectoryBindings{
								ID:   util.GetIntFromMap(bindingMap, "id"),
								Name: util.GetStringFromMap(bindingMap, "name"),
							})
						}
					}
				}
			}
			return directoryBindings
		}()

		// TODO refactor this section to use default values from schema. recent attempts cause 400 request errors.
		// Action: "doNotChange", is not being correctly passed from the schema despite it's correct config.
		managementAccount := func() jamfpro.PolicySubsetAccountMaintenanceManagementAccount {
			// Initialize with default values
			defaultManagementAccount := jamfpro.PolicySubsetAccountMaintenanceManagementAccount{
				Action:                "doNotChange",
				ManagedPassword:       "",
				ManagedPasswordLength: 0,
			}

			// Check if values are provided in Terraform and override defaults if necessary
			if managementAccountData, ok := accountMaintenanceData["management_account"].(map[string]interface{}); ok {
				defaultManagementAccount.Action = util.GetStringFromMap(managementAccountData, "action")
				defaultManagementAccount.ManagedPassword = util.GetStringFromMap(managementAccountData, "managed_password")
				defaultManagementAccount.ManagedPasswordLength = util.GetIntFromMap(managementAccountData, "managed_password_length")
			}

			return defaultManagementAccount
		}()

		openFirmwareEfiPassword := func() jamfpro.PolicySubsetAccountMaintenanceOpenFirmwareEfiPassword {
			var openFirmwareEfiPassword jamfpro.PolicySubsetAccountMaintenanceOpenFirmwareEfiPassword

			// Check if open firmware EFI password data is provided in Terraform
			if openFirmwareEfiPasswordData, ok := accountMaintenanceData["open_firmware_efi_password"].(map[string]interface{}); ok {
				openFirmwareEfiPasswordDataMap := util.ConvertToMapFromInterface(openFirmwareEfiPasswordData)
				if openFirmwareEfiPasswordDataMap == nil {
					// Skip if the open firmware EFI password data map is nil
					return openFirmwareEfiPassword
				}

				// Extract values from the Terraform data
				openFirmwareEfiPassword.OfMode = util.GetStringFromMap(openFirmwareEfiPasswordDataMap, "of_mode")
				openFirmwareEfiPassword.OfPassword = util.GetStringFromMap(openFirmwareEfiPasswordDataMap, "of_password")
				openFirmwareEfiPassword.OfPasswordSHA256 = util.GetStringFromMap(openFirmwareEfiPasswordDataMap, "of_password_sha256")
			}

			return openFirmwareEfiPassword
		}()

		// Assign all constructed components to AccountMaintenance
		policy.AccountMaintenance = jamfpro.PolicySubsetAccountMaintenance{
			Accounts:                accounts,
			DirectoryBindings:       directoryBindings,
			ManagementAccount:       managementAccount,
			OpenFirmwareEfiPassword: openFirmwareEfiPassword,
		}
	}

	// Construct the Reboot section
	policy.Reboot = func() jamfpro.PolicySubsetReboot {
		var reboot jamfpro.PolicySubsetReboot

		// Check if values are provided in Terraform
		if v, ok := d.GetOk("reboot"); ok {
			rebootData := v.(*schema.Set).List()[0].(map[string]interface{})

			// Extract values from the Terraform data
			reboot.Message = util.GetStringFromMap(rebootData, "message")
			reboot.SpecifyStartup = util.GetStringFromMap(rebootData, "specify_startup")
			reboot.StartupDisk = util.GetStringFromMap(rebootData, "startup_disk")
			reboot.NoUserLoggedIn = util.GetStringFromMap(rebootData, "no_user_logged_in")
			reboot.UserLoggedIn = util.GetStringFromMap(rebootData, "user_logged_in")
			reboot.MinutesUntilReboot = util.GetIntFromMap(rebootData, "minutes_until_reboot")
			reboot.StartRebootTimerImmediately = util.GetBoolFromMap(rebootData, "start_reboot_timer_immediately")
			reboot.FileVault2Reboot = util.GetBoolFromMap(rebootData, "file_vault_2_reboot")
		}

		return reboot
	}()

	// Construct the Maintenance section
	if v, ok := d.GetOk("maintenance"); ok {
		maintenanceData := v.([]interface{})[0].(map[string]interface{})
		policy.Maintenance = jamfpro.PolicySubsetMaintenance{
			Recon:                    util.GetBoolFromMap(maintenanceData, "recon"),
			ResetName:                util.GetBoolFromMap(maintenanceData, "reset_name"),
			InstallAllCachedPackages: util.GetBoolFromMap(maintenanceData, "install_all_cached_packages"),
			Heal:                     util.GetBoolFromMap(maintenanceData, "heal"),
			Prebindings:              util.GetBoolFromMap(maintenanceData, "prebindings"),
			Permissions:              util.GetBoolFromMap(maintenanceData, "permissions"),
			Byhost:                   util.GetBoolFromMap(maintenanceData, "byhost"),
			SystemCache:              util.GetBoolFromMap(maintenanceData, "system_cache"),
			UserCache:                util.GetBoolFromMap(maintenanceData, "user_cache"),
			Verify:                   util.GetBoolFromMap(maintenanceData, "verify"),
		}
	}

	// Construct the FilesProcesses section
	if v, ok := d.GetOk("files_processes"); ok {
		filesProcessesData := v.([]interface{})[0].(map[string]interface{})
		policy.FilesProcesses = jamfpro.PolicySubsetFilesProcesses{
			SearchByPath:         util.GetStringFromMap(filesProcessesData, "search_by_path"),
			DeleteFile:           util.GetBoolFromMap(filesProcessesData, "delete_file"),
			LocateFile:           util.GetStringFromMap(filesProcessesData, "locate_file"),
			UpdateLocateDatabase: util.GetBoolFromMap(filesProcessesData, "update_locate_database"),
			SpotlightSearch:      util.GetStringFromMap(filesProcessesData, "spotlight_search"),
			SearchForProcess:     util.GetStringFromMap(filesProcessesData, "search_for_process"),
			KillProcess:          util.GetBoolFromMap(filesProcessesData, "kill_process"),
			RunCommand:           util.GetStringFromMap(filesProcessesData, "run_command"),
		}
	}

	// Construct the UserInteraction section
	if v, ok := d.GetOk("user_interaction"); ok {
		userInteractionData := v.([]interface{})[0].(map[string]interface{})
		policy.UserInteraction = jamfpro.PolicySubsetUserInteraction{
			MessageStart:          util.GetStringFromMap(userInteractionData, "message_start"),
			AllowUserToDefer:      util.GetBoolFromMap(userInteractionData, "allow_user_to_defer"),
			AllowDeferralUntilUtc: util.GetStringFromMap(userInteractionData, "allow_deferral_until_utc"),
			AllowDeferralMinutes:  util.GetIntFromMap(userInteractionData, "allow_deferral_minutes"),
			MessageFinish:         util.GetStringFromMap(userInteractionData, "message_finish"),
		}
	}

	// Construct the DiskEncryption section
	policy.DiskEncryption = func() jamfpro.PolicySubsetDiskEncryption {
		var diskEncryption jamfpro.PolicySubsetDiskEncryption

		// Check if values are provided in Terraform
		if v, ok := d.GetOk("disk_encryption"); ok && len(v.([]interface{})) > 0 {
			diskEncryptionData := v.([]interface{})[0].(map[string]interface{})

			// Extract values from the Terraform data
			diskEncryption.Action = util.GetStringFromMap(diskEncryptionData, "action")
			diskEncryption.DiskEncryptionConfigurationID = util.GetIntFromMap(diskEncryptionData, "disk_encryption_configuration_id")
			diskEncryption.AuthRestart = util.GetBoolFromMap(diskEncryptionData, "auth_restart")
			diskEncryption.RemediateKeyType = util.GetStringFromMap(diskEncryptionData, "remediate_key_type")
			diskEncryption.RemediateDiskEncryptionConfigurationID = util.GetIntFromMap(diskEncryptionData, "remediate_disk_encryption_configuration_id")
		}

		return diskEncryption
	}()

	// Marshal the jamf pro policy object into XML for logging
	xmlData, err := xml.MarshalIndent(policy, "", "  ")
	if err != nil {
		// Handle the error if XML marshaling fails
		log.Printf("[ERROR] Error marshaling jamf pro policy object to XML: %s", err)
		return nil, fmt.Errorf("error marshaling jamf pro policy object to XML: %v", err)
	}

	// Log the XML formatted search object
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Constructed jamf pro policy Object:\n%s", string(xmlData)))

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Successfully constructed jamf pro policy with name: %s", policy.General.Name))

	return policy, nil
}
*/
