package policies

// TODO remove all the log.print's. Debug use only
// TODO handle all toxic combinations

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Primary

func constructPolicy(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {
	log.Println(("LOGHERE-CONSTRUCT"))
	var err error

	// Main obj
	out := &jamfpro.ResourcePolicy{}

	// General
	err = constructGeneral(d, out)
	if err != nil {
		return nil, err
	}

	// Scope
	err = constructScope(d, out)
	if err != nil {
		return nil, err
	}

	// Self Service
	err = constructSelfService(d, out)
	if err != nil {
		return nil, err
	}

	// Package Configuration
	// Scripts
	// Printers
	// DockItems
	// Account Maintenance
	// FilesProcesses
	// UserInteraction
	// DiskEncryption
	// Reboot

	// DEBUG
	log.Println("XMLOUT")
	policyXML, _ := xml.MarshalIndent(out, "", "  ")
	log.Println(string(policyXML))

	// END
	log.Println("LOGEND")
	return out, nil
}

// Child funcs

func constructGeneral(d *schema.ResourceData, out *jamfpro.ResourcePolicy) error {
	log.Println("GENERAL")

	// Primitive fields
	out.General = jamfpro.PolicySubsetGeneral{
		Name:                       d.Get("name").(string),
		Enabled:                    d.Get("enabled").(bool),
		TriggerCheckin:             d.Get("trigger_checkin").(bool),
		TriggerEnrollmentComplete:  d.Get("trigger_enrollment_complete").(bool),
		TriggerLogin:               d.Get("trigger_login").(bool),
		TriggerNetworkStateChanged: d.Get("trigger_network_state_changed").(bool),
		TriggerStartup:             d.Get("trigger_startup").(bool),
		TriggerOther:               d.Get("trigger_other").(string),
		Frequency:                  d.Get("frequency").(string),
		RetryEvent:                 d.Get("retry_event").(string),
		RetryAttempts:              d.Get("retry_attempts").(int),
		NotifyOnEachFailedRetry:    d.Get("notify_on_each_failed_retry").(bool),
		TargetDrive:                d.Get("target_drive").(string),
		Offline:                    d.Get("offline").(bool),
	}

	// TODO Do we need these set or can we just set the default to nil?
	// Category
	log.Println("CATEGORY")

	suppliedCategory := d.Get("category").([]interface{})
	if len(suppliedCategory) > 0 {
		outCat := &jamfpro.SharedResourceCategory{
			ID: suppliedCategory[0].(map[string]interface{})["id"].(int),
		}
		out.General.Category = outCat
	} else {
		out.General.Category = &jamfpro.SharedResourceCategory{
			ID: 0,
		}
	}

	// Site
	log.Println("SITE")

	suppliedSite := d.Get("site").([]interface{})
	if len(suppliedSite) > 0 {
		outSite := &jamfpro.SharedResourceSite{
			ID: suppliedSite[0].(map[string]interface{})["id"].(int),
		}
		out.General.Site = outSite
	} else {
		out.General.Site = &jamfpro.SharedResourceSite{
			ID: 0,
		}
	}

	return nil
}

func constructScope(d *schema.ResourceData, out *jamfpro.ResourcePolicy) error {
	log.Println("SCOPE")
	var err error

	if len(d.Get("scope").([]interface{})) == 0 {
		return nil
	}

	// Targets

	out.Scope = &jamfpro.PolicySubsetScope{
		Computers:      &[]jamfpro.PolicySubsetComputer{},
		ComputerGroups: &[]jamfpro.PolicySubsetComputerGroup{},
		JSSUsers:       &[]jamfpro.PolicySubsetJSSUser{},
		JSSUserGroups:  &[]jamfpro.PolicySubsetJSSUserGroup{},
		Buildings:      &[]jamfpro.PolicySubsetBuilding{},
		Departments:    &[]jamfpro.PolicySubsetDepartment{},
	}

	// Bools
	out.Scope.AllComputers = d.Get("scope.0.all_computers").(bool)
	out.Scope.AllJSSUsers = d.Get("scope.0.all_jss_users").(bool)

	// Computers
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputer, int]("scope.0.computer_ids", "ID", d, out.Scope.Computers)
	if err != nil {
		return err
	}

	// Computer Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputerGroup, int]("scope.0.computer_group_ids", "ID", d, out.Scope.ComputerGroups)
	if err != nil {
		return err
	}

	// JSS Users
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUser, int]("scope.0.jss_user_ids", "ID", d, out.Scope.JSSUsers)
	if err != nil {
		return err
	}

	// JSS User Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.jss_user_group_ids", "ID", d, out.Scope.JSSUserGroups)
	if err != nil {
		return err
	}

	// Buildings
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetBuilding, int]("scope.0.building_ids", "ID", d, out.Scope.Buildings)
	if err != nil {
		return err
	}

	// Departments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetDepartment, int]("scope.0.department_ids", "ID", d, out.Scope.Departments)
	if err != nil {
		return err
	}

	// Limitations

	out.Scope.Limitations = &jamfpro.PolicySubsetScopeLimitations{
		Users:           &[]jamfpro.PolicySubsetUser{},
		UserGroups:      &[]jamfpro.PolicySubsetUserGroup{},
		NetworkSegments: &[]jamfpro.PolicySubsetNetworkSegment{},
		IBeacons:        &[]jamfpro.PolicySubsetIBeacon{},
	}

	// Network Segments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.limitations.0.network_segment_ids", "ID", d, out.Scope.Limitations.NetworkSegments)
	if err != nil {
		return err
	}

	// IBeacons
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetIBeacon, int]("scope.0.limitations.0.ibeacon_ids", "ID", d, out.Scope.Limitations.IBeacons)
	if err != nil {
		return err
	}

	// User Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetUserGroup, int]("scope.0.limitations.0.user_group_ids", "ID", d, out.Scope.Limitations.UserGroups)
	if err != nil {
		return err
	}

	// TODO User Limitations

	// Exclusions

	// TODO I don't really want this here but it won't work without it. I think it's defeating the purpose of the struct layout slightly.
	out.Scope.Exclusions = &jamfpro.PolicySubsetScopeExclusions{
		Computers:       &[]jamfpro.PolicySubsetComputer{},
		ComputerGroups:  &[]jamfpro.PolicySubsetComputerGroup{},
		Users:           &[]jamfpro.PolicySubsetUser{},
		UserGroups:      &[]jamfpro.PolicySubsetUserGroup{},
		Buildings:       &[]jamfpro.PolicySubsetBuilding{},
		Departments:     &[]jamfpro.PolicySubsetDepartment{},
		NetworkSegments: &[]jamfpro.PolicySubsetNetworkSegment{},
		JSSUsers:        &[]jamfpro.PolicySubsetJSSUser{},
		JSSUserGroups:   &[]jamfpro.PolicySubsetJSSUserGroup{},
		IBeacons:        &[]jamfpro.PolicySubsetIBeacon{},
	}

	// Computers
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputer, int]("scope.0.exclusions.0.computer_ids", "ID", d, out.Scope.Exclusions.Computers)
	if err != nil {
		return err
	}

	// Computer Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputerGroup, int]("scope.0.exclusions.0.computer_group_ids", "ID", d, out.Scope.Exclusions.ComputerGroups)
	if err != nil {
		return err
	}

	// Buildings
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetBuilding, int]("scope.0.exclusions.0.building_ids", "ID", d, out.Scope.Exclusions.Buildings)
	if err != nil {
		return err
	}

	// Departments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetDepartment, int]("scope.0.exclusions.0.department_ids", "ID", d, out.Scope.Exclusions.Departments)
	if err != nil {
		return err
	}

	// Network Segments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.exclusions.0.network_segment_ids", "ID", d, out.Scope.Exclusions.NetworkSegments)
	if err != nil {
		return err
	}

	// JSS Users
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUser, int]("scope.0.exclusions.0.jss_user_ids", "ID", d, out.Scope.Exclusions.JSSUsers)
	if err != nil {
		return err
	}

	// JSS User Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.exclusions.0.jss_user_group_ids", "ID", d, out.Scope.Exclusions.JSSUserGroups)
	if err != nil {
		return err
	}

	// IBeacons
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetIBeacon, int]("scope.0.exclusions.0.ibeacon_ids", "ID", d, out.Scope.Exclusions.IBeacons)
	if err != nil {
		return err
	}

	// TODO make this better, it works for now
	if out.Scope.AllComputers && (out.Scope.Computers != nil ||
		out.Scope.ComputerGroups != nil ||
		out.Scope.Departments != nil ||
		out.Scope.Buildings != nil) {
		return fmt.Errorf("invalid combination - all computers with scoped endpoints")
	}

	return nil
}

func constructSelfService(d *schema.ResourceData, out *jamfpro.ResourcePolicy) error {

	if len(d.Get("self_service").([]interface{})) > 0 {
		out.SelfService = &jamfpro.PolicySubsetSelfService{
			UseForSelfService:           d.Get("self_service.0.use_for_self_service").(bool),
			SelfServiceDisplayName:      d.Get("self_service.0.self_service_display_name").(string),
			InstallButtonText:           d.Get("self_service.0.install_button_text").(string),
			ReinstallButtonText:         d.Get("self_service.0.reinstall_button_text").(string),
			SelfServiceDescription:      d.Get("self_service.0.self_service_description").(string),
			ForceUsersToViewDescription: d.Get("self_service.0.force_users_to_view_description").(bool),
			// TODO self service icon later
			FeatureOnMainPage: d.Get("self_service.0.feature_on_main_page").(bool),
			// TODO Self service categories later
		}
	}

	return nil
}
