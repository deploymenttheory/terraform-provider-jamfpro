package policies

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func constructPolicy(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {
	log.Println(("LOGHERE-CONSTRUCT"))

	// Main Object Definition with primitive values assigned.

	log.Println("LOG-UNPROCESSED FIELDS START")

	// Main obj
	out := &jamfpro.ResourcePolicy{}

	// General

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

	// Category
	log.Println("CATEGORY")

	suppliedCategory := d.Get("category").([]interface{})
	if len(suppliedCategory) > 0 {
		// construct category if provided
		outCat := &jamfpro.SharedResourceCategory{
			ID: suppliedCategory[0].(map[string]interface{})["id"].(int),
		}
		out.General.Category = outCat
	} else {
		// if no category, supply empty cat to remove it.
		out.General.Category = &jamfpro.SharedResourceCategory{
			ID: 0,
		}
	}

	// Site
	log.Println("SITE")

	suppliedSite := d.Get("site").([]interface{})
	if len(suppliedSite) > 0 {
		// If site provided, construct
		outSite := &jamfpro.SharedResourceSite{
			ID: suppliedSite[0].(map[string]interface{})["id"].(int),
		}
		out.General.Site = outSite
	} else {
		// If no site, construct no site obj. We have to do this for the site to be removed.
		out.General.Site = &jamfpro.SharedResourceSite{
			ID: 0,
		}
	}

	// Scope

	log.Println("SCOPE")

	out.Scope = &jamfpro.PolicySubsetScope{
		Computers:      &[]jamfpro.PolicySubsetComputer{},
		ComputerGroups: &[]jamfpro.PolicySubsetComputerGroup{},
		JSSUsers:       &[]jamfpro.PolicySubsetJSSUser{},
		JSSUserGroups:  &[]jamfpro.PolicySubsetJSSUserGroup{},
		Buildings:      &[]jamfpro.PolicySubsetBuilding{},
		Departments:    &[]jamfpro.PolicySubsetDepartment{},
	}

	// Scope - Targets
	var err error
	out.Scope.AllComputers = d.Get("scope.0.all_computers").(bool)
	out.Scope.AllJSSUsers = d.Get("scope.0.all_jss_users").(bool)

	log.Println("CONSTRUCT-FLAG-1")

	// Computers
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputer, int]("scope.0.computer_ids", "ID", d, out.Scope.Computers)
	if err != nil {
		return nil, err
	}

	log.Printf("%+v", out.Scope.Computers)
	log.Println("CONSTRUCT-FLAG-2")

	// Computer Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputerGroup, int]("scope.0.computer_group_ids", "ID", d, out.Scope.ComputerGroups)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-3")

	// JSS Users
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUser, int]("scope.0.jss_user_ids", "ID", d, out.Scope.JSSUsers)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-4")

	// JSS User Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.jss_user_group_ids", "ID", d, out.Scope.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-5")

	// Buildings
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetBuilding, int]("scope.0.building_ids", "ID", d, out.Scope.Buildings)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-6")

	// Departments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetDepartment, int]("scope.0.department_ids", "ID", d, out.Scope.Departments)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-7")

	// Scope - Limitations

	log.Println("CONSTRUCT-FLAG-8")

	out.Scope.Limitations = &jamfpro.PolicySubsetScopeLimitations{
		Users:           &[]jamfpro.PolicySubsetUser{},
		UserGroups:      &[]jamfpro.PolicySubsetUserGroup{},
		NetworkSegments: &[]jamfpro.PolicySubsetNetworkSegment{},
		IBeacons:        &[]jamfpro.PolicySubsetIBeacon{},
	}

	// Network Segments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.limitations.0.network_segment_ids", "ID", d, out.Scope.Limitations.NetworkSegments)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-9")

	// IBeacons
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetIBeacon, int]("scope.0.limitations.0.ibeacon_ids", "ID", d, out.Scope.Limitations.IBeacons)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-10")

	// User Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetUserGroup, int]("scope.0.limitations.0.user_group_ids", "ID", d, out.Scope.Limitations.UserGroups)
	if err != nil {
		return nil, err
	}

	// TODO Users

	log.Println("CONSTRUCT-FLAG-11")

	// Scope - Exclusions

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
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-12")

	// Computer Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetComputerGroup, int]("scope.0.exclusions.0.computer_group_ids", "ID", d, out.Scope.Exclusions.ComputerGroups)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-13")

	// Buildings
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetBuilding, int]("scope.0.exclusions.0.building_ids", "ID", d, out.Scope.Exclusions.Buildings)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-14")

	// Departments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetDepartment, int]("scope.0.exclusions.0.department_ids", "ID", d, out.Scope.Exclusions.Departments)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-15")

	// Network Segments
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.exclusions.0.network_segment_ids", "ID", d, out.Scope.Exclusions.NetworkSegments)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-16")

	// JSS Users
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUser, int]("scope.0.exclusions.0.jss_user_ids", "ID", d, out.Scope.Exclusions.JSSUsers)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-17")

	// JSS User Groups
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.exclusions.0.jss_user_group_ids", "ID", d, out.Scope.Exclusions.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-18")

	// IBeacons
	err = GetAttrsListFromHCLForPointers[jamfpro.PolicySubsetIBeacon, int]("scope.0.exclusions.0.ibeacon_ids", "ID", d, out.Scope.Exclusions.IBeacons)
	if err != nil {
		return nil, err
	}

	log.Println("CONSTRUCT-FLAG-19")

	// TODO make this better, it works for now
	if out.Scope.AllComputers && (out.Scope.Computers != nil ||
		out.Scope.ComputerGroups != nil ||
		out.Scope.Departments != nil ||
		out.Scope.Buildings != nil) {
		return nil, fmt.Errorf("invalid combination - all computers with scoped endpoints")
	}

	log.Println("CONSTRUCT-FLAG-20")

	// Self Service
	// out.SelfService = &jamfpro.PolicySubsetSelfService{
	// 	// UseForSelfService:           d.Get("self_service.0.use_for_self_service").(bool),
	// 	// SelfServiceDisplayName:      d.Get("self_service_display_name").(string),
	// 	// InstallButtonText:           d.Get("install_button_text").(string),
	// 	// ReinstallButtonText:         d.Get("reinstall_button_text").(string),
	// 	// SelfServiceDescription:      d.Get("self_service_description").(string),
	// 	// ForceUsersToViewDescription: d.Get("force_users_to_view_description").(bool),
	// 	// TODO self service icon later
	// 	// FeatureOnMainPage: d.Get("feature_on_main_page").(bool),
	// 	// TODO Self service categories later
	// }

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
	log.Println("LOGEND")
	log.Println(string(policyXML))

	// END

	return out, nil
}
