package ldap_server

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a ResourceLDAPServers object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceLDAPServers, error) {
	resource := &jamfpro.ResourceLDAPServers{
		Connection:       constructConnection(d),
		MappingsForUsers: constructMappings(d),
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro LDAP Server '%s' to XML: %v", resource.Connection.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro LDAP Server XML:\n%s\n", string(resourceXML))

	return resource, nil
}

// Helper functions for nested structures

func constructConnection(d *schema.ResourceData) jamfpro.LDAPServerSubsetConnection {
	connection := jamfpro.LDAPServerSubsetConnection{
		Name:               d.Get("name").(string),
		Hostname:           d.Get("hostname").(string),
		ServerType:         d.Get("server_type").(string),
		Port:               d.Get("port").(int),
		UseSSL:             d.Get("use_ssl").(bool),
		AuthenticationType: d.Get("authentication_type").(string),
		OpenCloseTimeout:   d.Get("open_close_timeout").(int),
		SearchTimeout:      d.Get("search_timeout").(int),
		ReferralResponse:   d.Get("referral_response").(string),
		UseWildcards:       d.Get("use_wildcards").(bool),
	}

	// Handle account credentials
	if v, ok := d.GetOk("account"); ok && len(v.([]interface{})) > 0 {
		accountMap := v.([]interface{})[0].(map[string]interface{})
		connection.Account = jamfpro.LDAPServerSubsetConnectionAccount{
			DistinguishedUsername: accountMap["distinguished_username"].(string),
			Password:              accountMap["password"].(string),
		}
	}

	return connection
}

func constructMappings(d *schema.ResourceData) jamfpro.LDAPServerContainerMapping {
	return jamfpro.LDAPServerContainerMapping{
		UserMappings:                constructUserMappings(d),
		UserGroupMappings:           constructUserGroupMappings(d),
		UserGroupMembershipMappings: constructUserGroupMembershipMappings(d),
	}
}

func constructUserMappings(d *schema.ResourceData) jamfpro.LDAPServerSubsetMappingUsers {
	if v, ok := d.GetOk("user_mappings"); ok && len(v.([]interface{})) > 0 {
		mappingMap := v.([]interface{})[0].(map[string]interface{})
		return jamfpro.LDAPServerSubsetMappingUsers{
			MapObjectClassToAnyOrAll: mappingMap["map_object_class_to_any_or_all"].(string),
			ObjectClasses:            mappingMap["object_classes"].(string),
			SearchBase:               mappingMap["search_base"].(string),
			SearchScope:              mappingMap["search_scope"].(string),
			MapUserID:                mappingMap["map_user_id"].(string),
			MapUsername:              mappingMap["map_username"].(string),
			MapRealName:              mappingMap["map_realname"].(string),
			MapEmailAddress:          mappingMap["map_email_address"].(string),
			AppendToEmailResults:     mappingMap["append_to_email_results"].(string),
			MapDepartment:            mappingMap["map_department"].(string),
			MapBuilding:              mappingMap["map_building"].(string),
			MapRoom:                  mappingMap["map_room"].(string),
			MapPhone:                 mappingMap["map_phone"].(string),
			MapPosition:              mappingMap["map_position"].(string),
			MapUserUUID:              mappingMap["map_user_uuid"].(string),
		}
	}
	return jamfpro.LDAPServerSubsetMappingUsers{}
}

func constructUserGroupMappings(d *schema.ResourceData) jamfpro.LDAPServerSubsetMappingUserGroups {
	if v, ok := d.GetOk("user_group_mappings"); ok && len(v.([]interface{})) > 0 {
		mappingMap := v.([]interface{})[0].(map[string]interface{})
		return jamfpro.LDAPServerSubsetMappingUserGroups{
			MapObjectClassToAnyOrAll: mappingMap["map_object_class_to_any_or_all"].(string),
			ObjectClasses:            mappingMap["object_classes"].(string),
			SearchBase:               mappingMap["search_base"].(string),
			SearchScope:              mappingMap["search_scope"].(string),
			MapGroupID:               mappingMap["map_group_id"].(string),
			MapGroupName:             mappingMap["map_group_name"].(string),
			MapGroupUUID:             mappingMap["map_group_uuid"].(string),
		}
	}
	return jamfpro.LDAPServerSubsetMappingUserGroups{}
}

func constructUserGroupMembershipMappings(d *schema.ResourceData) jamfpro.LDAPServerSubsetMappingUserGroupMemberships {
	if v, ok := d.GetOk("user_group_membership_mappings"); ok && len(v.([]interface{})) > 0 {
		mappingMap := v.([]interface{})[0].(map[string]interface{})
		return jamfpro.LDAPServerSubsetMappingUserGroupMemberships{
			UserGroupMembershipStoredIn:   mappingMap["user_group_membership_stored_in"].(string),
			MapGroupMembershipToUserField: mappingMap["map_group_membership_to_user_field"].(string),
			AppendToUsername:              mappingMap["append_to_username"].(string),
			UseDN:                         mappingMap["use_dn"].(bool),
			RecursiveLookups:              mappingMap["recursive_lookups"].(bool),
			GroupMembershipEnabledWhenUserMembershipSelected: mappingMap["group_membership_enabled_when_user_membership_selected"].(bool),
			MapUserMembershipToGroupField:                    mappingMap["map_user_membership_to_group_field"].(string),
			MapUserMembershipUseDN:                           mappingMap["map_user_membership_use_dn"].(bool),
			MapObjectClassToAnyOrAll:                         mappingMap["map_object_class_to_any_or_all"].(string),
			ObjectClasses:                                    mappingMap["object_classes"].(string),
			SearchBase:                                       mappingMap["search_base"].(string),
			SearchScope:                                      mappingMap["search_scope"].(string),
			Username:                                         mappingMap["username"].(string),
			GroupID:                                          mappingMap["group_id"].(string),
			UserGroupMembershipUseLDAPCompare:                mappingMap["user_group_membership_use_ldap_compare"].(bool),
			MembershipScopingOptimization:                    mappingMap["membership_scoping_optimization"].(bool),
		}
	}
	return jamfpro.LDAPServerSubsetMappingUserGroupMemberships{}
}
