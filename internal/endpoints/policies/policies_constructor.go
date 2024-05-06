package policies

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"reflect"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func constructPolicy(d *schema.ResourceData) (*jamfpro.ResourcePolicy, error) {
	log.Println(("LOGHERE"))

	// Main Object Definition with primitive values assigned.

	log.Println("LOG-UNPROCESSED FIELDS START")

	out := &jamfpro.ResourcePolicy{
		General: jamfpro.PolicySubsetGeneral{
			// ID computed
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
			// DateTimeLimitations: &jamfpro.PolicySubsetGeneralDateTimeLimitations{ // TODO this later. Need to create a nice way of standardising our HCL date/time inputs.
			// 	ActivationDate:      d.Get("date_time_limitations.0.activation_date").(string),
			// 	ActivationDateEpoch: d.Get("date_time_limitations.0.activation_date_epoch").(int),
			// 	ActivationDateUTC:   d.Get("date_time_limitations.0.activation_date_utc").(string),
			// 	ExpirationDate:      d.Get("date_time_limitations.0.expiration_date").(string),
			// 	ExpirationDateEpoch: d.Get("date_time_limitations.0.expiration_date_epoch").(int),
			// 	ExpirationDateUTC:   d.Get("date_time_limitations.0.expiration_date_utc").(string),
			// 	// no execute on // TODO
			// 	NoExecuteStart: d.Get("no_execute_start").(string),
			// 	NoExecuteEnd:   d.Get("no_execute_end").(string),
			// },

			// Category processed
			// site processed

			NetworkLimitations: &jamfpro.PolicySubsetGeneralNetworkLimitations{
				MinimumNetworkConnection: d.Get("network_limitations.0.minimum_network_connection").(string),
				AnyIPAddress:             d.Get("network_limitations.0.any_ip_address").(bool),
				NetworkSegments:          d.Get("network_limitations.0.network_segments").(string), // TODO is this a string?
			},
		},
		Scope: &jamfpro.PolicySubsetScope{
			AllComputers: d.Get("scope.0.all_computers").(bool),
			AllJSSUsers:  d.Get("scope.0.all_jss_users").(bool),
			// Rest processed
		},
		SelfService: &jamfpro.PolicySubsetSelfService{
			// UseForSelfService:           d.Get("self_service.0.use_for_self_service").(bool),
			// SelfServiceDisplayName:      d.Get("self_service_display_name").(string),
			// InstallButtonText:           d.Get("install_button_text").(string),
			// ReinstallButtonText:         d.Get("reinstall_button_text").(string),
			// SelfServiceDescription:      d.Get("self_service_description").(string),
			// ForceUsersToViewDescription: d.Get("force_users_to_view_description").(bool),
			// TODO self service icon later
			// FeatureOnMainPage: d.Get("feature_on_main_page").(bool),
			// TODO Self service categories later
		},
		// Package Configuration
		// Scripts
		// Printers
		// DockItems
		// Account Maintenance
		// FilesProcesses
		// UserInteraction
		// DiskEncryption
		// Reboot
	}

	// DEBUG
	log.Println("LOG-UNPROCESSED FIELDS END")
	log.Println("STRUCT NOW:")
	json, _ := json.MarshalIndent(out, " ", "    ")
	log.Println(string(json))
	log.Println("LOG-PROCESSED FIELDS START")

	// Processed Fields

	// General

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

	// Scope - Targets
	var err error
	out.Scope.AllComputers = d.Get("scope.0.all_computers").(bool)
	out.Scope.AllJSSUsers = d.Get("scope.0.all_jss_users").(bool)

	// Computers
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetComputer, int]("scope.0.computer_ids", "ID", d, out.Scope.Computers)
	if err != nil {
		return nil, err
	}

	// Computer Groups
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetComputerGroup, int]("scope.0.computer_group_ids", "ID", d, out.Scope.ComputerGroups)
	if err != nil {
		return nil, err
	}

	// JSS Users
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetJSSUser, int]("scope.0.jss_user_ids", "ID", d, out.Scope.JSSUsers)
	if err != nil {
		return nil, err
	}

	// JSS User Groups
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.jss_user_group_ids", "ID", d, out.Scope.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	// Buildings
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetBuilding, int]("scope.0.building_ids", "ID", d, out.Scope.Buildings)
	if err != nil {
		return nil, err
	}

	// Departments
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetDepartment, int]("scope.0.department_ids", "ID", d, out.Scope.Departments)
	if err != nil {
		return nil, err
	}

	// Scope - Limitations

	// Users
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetUser, string]("scope.0.limitations.0.user_names", "Name", d, out.Scope.Limitations.Users)
	if err != nil {
		return nil, err
	}

	// Network Segment
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.limitations.0.network_segment_ids", "ID", d, out.Scope.Limitations.NetworkSegments)
	if err != nil {
		return nil, err
	}

	// IBeacons
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetIBeacon, int]("scope.0.limitations.0.ibeacon_ids", "ID", d, out.Scope.Limitations.IBeacons)
	if err != nil {
		return nil, err
	}

	// User Groups
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetUserGroup, int]("scope.0.limitations.0.user_group_ids", "ID", d, out.Scope.Limitations.UserGroups)
	if err != nil {
		return nil, err
	}

	// Scope - Limitations

	// Computers
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetComputer, int]("scope.0.exclusions.0.computer_ids", "ID", d, out.Scope.Exclusions.Computers)
	if err != nil {
		return nil, err
	}

	// Computer Groups
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetComputerGroup, int]("scope.0.exclusions.0.computer_group_ids", "ID", d, out.Scope.Exclusions.ComputerGroups)
	if err != nil {
		return nil, err
	}

	// Buildings
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetBuilding, int]("scope.0.exclusions.0.building_ids", "ID", d, out.Scope.Exclusions.Buildings)
	if err != nil {
		return nil, err
	}

	// Departments
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetDepartment, int]("scope.0.exclusions.0.department_ids", "ID", d, out.Scope.Exclusions.Departments)
	if err != nil {
		return nil, err
	}

	// Network Segments
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetNetworkSegment, int]("scope.0.exclusions.0.network_segment_ids", "ID", d, out.Scope.Exclusions.NetworkSegments)
	if err != nil {
		return nil, err
	}

	// JSS Users
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetJSSUser, int]("scope.0.exclusions.0.jss_user_ids", "ID", d, out.Scope.Exclusions.JSSUsers)
	if err != nil {
		return nil, err
	}

	// JSS User Groups
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetJSSUserGroup, int]("scope.0.exclusions.0.jss_user_group_ids", "ID", d, out.Scope.Exclusions.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	// IBeacons
	err = GetAttrsListFromHCL[jamfpro.PolicySubsetIBeacon, int]("scope.0.exclusions.0.ibeacon_ids", "ID", d, out.Scope.Exclusions.IBeacons)
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

	// DEBUG
	policyXML, _ := xml.MarshalIndent(out, "", "  ")
	log.Println("LOGEND")
	log.Println(string(policyXML))

	// END
	return out, nil
}

// TODO rename this func and put it somewhere else
func GetAttrsListFromHCL[NestedObjectType any, ListItemPrimitiveType any](path string, target_field string, d *schema.ResourceData, home *[]NestedObjectType) (err error) {
	getAttr, ok := d.GetOk(path)

	if len(getAttr.([]interface{})) == 0 {
		return nil
	}

	if ok {
		outList := make([]NestedObjectType, 0)
		for _, v := range getAttr.([]interface{}) {
			var newObj NestedObjectType
			newObjReflect := reflect.ValueOf(&newObj).Elem()
			idField := newObjReflect.FieldByName(target_field)

			if idField.IsValid() && idField.CanSet() {
				idField.Set(reflect.ValueOf(v.(ListItemPrimitiveType)))
			} else {
				return fmt.Errorf("error cannot set field line 695") // TODO write this error
			}

			outList = append(outList, newObj)

		}

		if len(outList) > 0 {
			*home = outList
		} else {
			log.Println("list is empty")
		}

		return nil
	}
	return fmt.Errorf("no path found/no scoped items at %v", path)
}
